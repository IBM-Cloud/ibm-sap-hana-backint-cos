package backint

import (
	"fmt"

	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/cos"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/global"
	"github.com/ibm-cloud/ibm-sap-hana-backint-cos/utils/logging"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

/*
Deleting the cloud objects specified in the input file
*/
func DeleteCloudObjects(
	s3Client *s3.S3,
) bool {
	success := true
	global.Logger.Debug("Function: delete")

	cosObjects := getCosObjectsForDelete(s3Client)
	deleteResults := cos.DeleteMultiple(s3Client, cosObjects)

	for _, r := range deleteResults {
		logging.BackintResultMsgs.AddKeyword(
			r.Status,
			[]string{r.ETag, r.Key},
		)
		if r.Status == "ERROR" {
			global.Logger.Error(
				fmt.Sprintf("Failed to delete object '%s' with ETag '%s'.",
					r.Key,
					r.ETag,
				),
			)
			success = false
		}
	}
	return success
}
