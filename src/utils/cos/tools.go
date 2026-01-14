/*
Tools and type definitions for all IBM Cloud Object Storage functions
*/
package cos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hdbbackint/utils/config"
	"hdbbackint/utils/global"
	"os"
	"time"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
	"golang.org/x/sys/unix"
)

/*
Reader function for uploading data
*/
func (r *backintReader) Read(p []byte) (int, error) {
	// if r.compress {
	// 	global.Logger.Debug("compress")
	// 	readFromPipe, err := r.r.Read(global.Compress(p))
	// 	global.Logger.Debug(fmt.Sprintf("Read from pipe: %d", len(p)))
	// 	global.Logger.Debug(fmt.Sprintf("Compressed: %d", readFromPipe))
	// 	r.noOfbytes += int64(readFromPipe)
	// 	return readFromPipe, err
	// } else {
	readFromPipe, err := r.r.Read(p)
	r.noOfbytes += int64(readFromPipe)
	return readFromPipe, err
	// }
}

/*
Setting METADATA information
*/
func setMetaData() map[string]*string {
	metadata := make(map[string]*string)
	cmp := config.BackintConfig.CompressionString()
	metadata[global.METADATA_COMPRESSION_LABEL] = &cmp
	return metadata
}

/*
Setting up the information for uploading data to IBM Cloud Object Storage
*/
func setupUploadInputInfo(
	Key string,
	sourcePath string,
) (s3manager.UploadInput, *backintReader) {
	global.Logger.Debug("Opening the input pipe for reading.")
	rPipe, err := os.OpenFile(sourcePath, os.O_CREATE, os.ModeNamedPipe)
	global.CheckForError(
		err,
		fmt.Sprintf("Error opening named pipe '%s'", sourcePath),
		global.FAILURE,
	)

	// TODO compression -> Issue #7
	readerFromPipe := backintReader{
		r:         rPipe,
		noOfbytes: 0,
		// compress:  backintConfig.compression(),
		compress: false,
	}

	tags := config.BackintConfig.Tags()
	var pLockMode *string
	var pLockDate *time.Time
	if config.BackintConfig.ObjectLockRetentionMode() == "cmp" {
		lockMode := global.OBJECTLOCKMODE
		lockDate := config.BackintConfig.ObjectLockRetentionDate()
		pLockMode = &lockMode
		pLockDate = &lockDate
	}
	lockLegalHold := config.BackintConfig.ObjectLockLegalHoldStatus()

	input := s3manager.UploadInput{
		Bucket:                    aws.String(config.BackintConfig.BucketName()),
		Key:                       aws.String(Key),
		Body:                      &readerFromPipe,
		ObjectLockLegalHoldStatus: &lockLegalHold,
		ObjectLockMode:            pLockMode,
		ObjectLockRetainUntilDate: pLockDate,
		Tagging:                   &tags,
		Metadata:                  setMetaData(),
	}

	return input, &readerFromPipe
}

/*
Getting the number of parts from Cloud
*/
func calculateNumberOfParts(s3Client *s3.S3, size int64, Key string) (int64, int64) {
	noOfParts := getPartsCount(s3Client, Key)
	chunksize := size / noOfParts
	if size%noOfParts != 0 {
		chunksize++
	}

	global.Logger.Info(fmt.Sprintf(
		"Downloading '%s' with '%d' parts and a chunksize of '%d'.",
		Key,
		noOfParts,
		chunksize),
	)
	return noOfParts, chunksize
}

/*
Generating the infos for downloading the parts from COS

	Returns:

	Array of struct[]DownloadPart
	Number of parts to be downloaded
*/
func generateDownloadParts(
	s3Client *s3.S3,
	size int64,
	Key string,
) ([]DownloadPart, int64) {
	var downloadParts []DownloadPart
	noOfParts, chunksize := calculateNumberOfParts(s3Client, size, Key)

	for p := range noOfParts {
		start := p * chunksize
		end := start + chunksize - 1
		if end >= size {
			end = size - 1
		}

		byteRange := fmt.Sprintf("bytes %d-%d", start, end)

		dp := DownloadPart{
			Key:        Key,
			numParts:   noOfParts,
			partNumber: p + 1,
			byteRange:  byteRange,
		}
		downloadParts = append(downloadParts, dp)
	}
	return downloadParts, int64(len(downloadParts))
}

/*
Trying to write downloaded part to pipe.
If the next part in the row is already stored in downloadedParts, write it to pipe,
increase the nextIndex counter and
delete the part from downloadedParts

!!Caution:!!
This function must be handled carefully.
All writes and reads to the map (downloadedParts) and to the pipe
must be locked directly before and after the action!
*/
func sendDataToHANA(fifo *os.File,
	nextIndex *int64,
	downloadedParts *ByteMap,
	index int64,
	pipeBufferSize int,
	buffer *bytes.Buffer,
) bool {

	writeToPipeLock.Lock()

	// Storing data into buffer
	(*downloadedParts)[index] = buffer.Bytes()

	global.Logger.Debug(fmt.Sprintf(
		"'%s': Writing part #%d to buffer, nextIndex = %d.",
		fifo.Name(),
		index,
		*nextIndex,
	))

	global.Logger.Debug(fmt.Sprintf(
		"'%s': downloadedParts length: %d",
		fifo.Name(),
		len(*downloadedParts),
	))

	writeToPipeLock.Unlock()

	for {
		writeToPipeLock.Lock()

		data := getBufferedDataForIndex(downloadedParts, nextIndex, fifo.Name())
		if data == nil {
			writeToPipeLock.Unlock()
			break
		}

		global.Logger.Debug(fmt.Sprintf(
			"'%s': Writing part #%d to pipe",
			fifo.Name(),
			*nextIndex,
		))

		if !writeDataToPipe(fifo, data, nextIndex, pipeBufferSize) {
			writeToPipeLock.Unlock()
			return false
		}

		deleteBufferedDataForIndex(downloadedParts, nextIndex, fifo.Name())

		*nextIndex++

		global.Logger.Debug(fmt.Sprintf(
			"'%s': Increased nextIndex to #%d",
			fifo.Name(),
			*nextIndex,
		))
		writeToPipeLock.Unlock()
	}
	return true
}

/*
Writing the data to pipe
Due to hang problems (looks like it is caused by HANA processing itself) the
data must be splitted into smaller portions so that HANA can process the data.
The size of the portion is set by config parameter pipe_chunksize_KB.
In addition, writing to pipe stops after 30 seconds, if not successfull
*/
func writeDataToPipe(fifo *os.File, data []byte, nextIndex *int64, pipeBufferSize int) bool {

	// Processing the data portions
	for i := 0; i < len(data); i += pipeBufferSize {
		end := min(i+pipeBufferSize, len(data))

		// Setting the timeout for the Write statement
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(30)*time.Second,
		)
		defer cancel()

		written := make(chan error, 1)

		go func() {
			_, err := fifo.Write(data[i:end])
			time.Sleep(
				time.Duration(config.BackintConfig.Timeout()) * time.Microsecond,
			)
			written <- err
		}()

		select {
		case <-ctx.Done():
			// Timout writing to pipe
			global.Logger.Error(fmt.Sprintf(
				"Error writing part #%d to pipe '%s': %s",
				*nextIndex,
				fifo.Name(),
				errors.New("Timeout")),
			)
		case err := <-written:
			if err != nil {
				global.Logger.Error(fmt.Sprintf(
					"'%s': Error writing part #%d to pipe: %s",
					fifo.Name(),
					*nextIndex,
					err,
				))
				return false
			}
		}
	}
	return true
}

func getBufferedDataForIndex(
	downloadedParts *ByteMap,
	index *int64,
	pipeName string) []byte {
	data, available := (*downloadedParts)[*index]

	if !available {
		global.Logger.Debug(fmt.Sprintf(
			"'%s': part with %d not found in buffer.",
			pipeName,
			*index,
		))
		return nil
	}
	return data
}

func deleteBufferedDataForIndex(
	downloadedParts *ByteMap,
	index *int64,
	pipeName string) {
	global.Logger.Debug(fmt.Sprintf(
		"'%s': Deleting part #%d from buffer",
		pipeName,
		*index,
	))
	delete(*downloadedParts, *index)
}

func openPipeForWriting(pipeName string) *os.File {
	// Opening destination pipe for writing
	fifo, err := os.OpenFile(pipeName, os.O_WRONLY, os.ModeNamedPipe)
	global.CheckForError(err,
		fmt.Sprintf("Error opening named pipe '%s'", pipeName),
		global.FAILURE,
	)
	return fifo
}

func getPipeBufferSize(fifo *os.File) int {
	const F_GETPIPE_SZ = 1032
	size, err := unix.FcntlInt(fifo.Fd(), F_GETPIPE_SZ, 0)
	if err != nil {
		global.Logger.Info(fmt.Sprintf(
			"Setting the chunksize for piping to default %d",
			global.PIPE_BUFFER_SIZE,
		))
		return global.PIPE_BUFFER_SIZE
	}
	global.Logger.Info(fmt.Sprintf(
		"Setting the chunksize for piping to pipe buffer size %d",
		size,
	))
	return size
}
