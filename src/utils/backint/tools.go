package backint

import (
	"strings"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/config"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/cos"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Generating the object Key name
*/
func generateCosObjectKeyname(pipeName string) string {
	Key, _ := strings.CutPrefix(pipeName, config.BackintConfig.RemoveKeyPrefix())
	Key = config.BackintConfig.AdditionalKeyPrefix() + Key

	if global.Args.Function == global.BACKUP {
		global.Logger.Info("'" + pipeName + "' -> '" + Key + "'.")
	} else {
		global.Logger.Info("'" + Key + "' -> '" + pipeName + "'.")
	}
	return Key
}

/*
Getting the source paths from the input file for function = BACKUP
*/
func getSourcePathsForBackup() []string {
	var sourcePaths []string
	for _, element := range global.InputFileContent {
		if element.Keyword != "PIPE" {
			continue
		}
		sourcePaths = append(sourcePaths, element.Parameter)
	}
	return sourcePaths
}

/*
Getting the list of object names and the ETags for function = DELETE
*/
func getCosObjectsForDelete(
	s3Client *s3.S3,
) []cos.CosObject {
	var cosObjects []cos.CosObject

	for _, element := range global.InputFileContent {
		if element.Keyword != "EBID" {
			continue
		}

		// Checking if object exists with specified EBID
		splitted := strings.Split(element.Parameter, " ")
		ETag := splitted[0]
		Key := splitted[1]
		cos_object := cos.CosObject{
			ETag:  ETag,
			Key:   Key,
			Found: false,
		}

		for _, cos_element := range cos.ListObjectsOfBucket(s3Client) {
			if cos_element.ETag == &ETag && cos_element.Key == &Key {
				cos_object.Found = true
				break
			}
		}
		cosObjects = append(cosObjects, cos_object)
	}
	return cosObjects
}

/*
Getting the list of object names and the ETags for function = RESTORE
*/
func getCosObjectsForRestore() []cos.CosObject {
	var cosObjects []cos.CosObject
	for _, element := range global.InputFileContent {
		splitted := strings.Split(element.Parameter, " ")

		etag := ""
		sourcePath := ""
		Key := ""
		destination := ""

		switch element.Keyword {
		case "EBID":
			etag = splitted[0]
			sourcePath = splitted[1]
			destination = sourcePath

			if len(splitted) == 3 {
				destination = splitted[2]
			}
		case "NULL":
			etag = ""
			sourcePath = splitted[0]
			destination = sourcePath
			if len(splitted) == 2 {
				destination = splitted[1]
			}

		default:
			return nil
		}

		Key = generateCosObjectKeyname(sourcePath)

		nextIndex := int64(1)

		cosObject := cos.CosObject{
			ETag:        etag,
			Key:         Key,
			Destination: destination,
			Found:       false,
			NextIndex:   &nextIndex,
		}

		cosObjects = append(cosObjects, cosObject)
	}
	return cosObjects
}
