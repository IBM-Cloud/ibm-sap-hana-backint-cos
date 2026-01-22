# Using IBM Backint Agent as Interface for SNAPPY Tool

IBM Backint agent for SAP HANA (IBM Backint agent) can be used to connect to the IBM Cloud Object Storage gathering the necessary information.

## Arguments

| Argument | Mandatory | Description |
|:-------|:-----| :-----|
| -f | yes | Function to be executed, see [Supported Functions](#supported-functions) |
| -keypath | yes | Path of the file containing the APIKEY |
| -authendpoint | no | URL used for IAM authentication |
| -region | yes | Region of the IBM Cloud Object Storage Bucket |
| -endpoint | yes | Bucket endpoint URL |
| -bucket | yes | Name of the IBM Cloud Object Storage bucket |
| -source | yes for function file-upload | Path to the file to be uploaded |
| -key | yes for function file-upload | Key name for uploaded file |
| -r | yes for functions bucket-get-list and bucket-get-lifecycle | Path of the response file containing the information provided by Cloud Object Storage |


## Supported Functions

### Verifying Bucket

`-f bucket-verify`

Checking if the given bucket exists and if versioning is enabled.

Required arguments:
- keypath
- region
- endpoint
- bucket

### Verifying Lifecycle Settings of Bucket

`-f bucket-get-lifecycle`

Returning a list of all lifecycle rules defined for the given bucket.

The lifecycle rules are written to the file specified for argument `-r`.

Each line in the response file represents one rule. The output has the format:

`ID:<ID of the rule>;Expiration:<Days>`

Required arguments:
- keypath
- region
- endpoint
- bucket
- r

### Getting Object List for Bucket

`-f bucket-get-list`

Returning a list of all objects in the given bucket.

The names of the objects are stored in the response file. Each line represents one object.

Required arguments:
- keypath
- region
- endpoint
- bucket
- r

### Uploading a File

`-f file-upload`

Uploading a file to the IBM Cloud Object Storage.

Required arguments:
- keypath
- region
- endpoint
- bucket
- source
- key
