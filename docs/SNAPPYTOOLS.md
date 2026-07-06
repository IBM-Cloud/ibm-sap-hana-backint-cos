# Using IBM Backint Agent as Interface for SNAPPY Tool

IBM Backint agent for SAP HANA (IBM Backint agent) can be used to connect to the IBM Cloud Object Storage gathering the necessary information.

## Mandatory Arguments

| Argument | Description |
|:-------| :-----|
| -f | Function to be executed.<br>See [Supported Functions](#supported-functions) |
| -p | Parameter file.<br>See [Configuration Parameters](../README.md#configuration-parameters) |

## Supported Functions

### Verifying Bucket

`-f bucket-verify`

Checking if the given bucket exists and if versioning is enabled.

#### Additional mandatory arguments
No additional command line arguments needed.

#### Example

```
hdbbackint -p hdbbackint.cfg -f bucket-verify
```

### Verifying Lifecycle Settings of Bucket

`-f bucket-get-lifecycle`

Returning a list of all lifecycle rules defined for the given bucket.

The lifecycle rules are written to the file specified for argument `-r`.

Each line in the response file represents one rule. The output has the format:

`ID:<ID of the rule>;Expiration:<Days>`

#### Additional mandatory arguments
| Argument | Description |
|:-------| :-----|
| -r | Path of the response file containing the information provided by Cloud Object Storage |

#### Example

```
hdbbackint -p hdbbackint.cfg -f bucket-get-lifecycle -r ./info/bucket-get-lifecycle.out
```

### Getting Object List for Bucket

`-f bucket-get-list`

Returning a list of all objects in the given bucket.

The names of the objects are stored in the response file. Each line represents one object.

#### Additional mandatory arguments
| Argument | Description |
|:-------| :-----|
| -r | Path of the response file containing the information provided by Cloud Object Storage |

#### Example

```
hdbbackint -p hdbbackint.cfg -f bucket-get-list -r ./info/bucket-get-list.out
```

### Uploading a File

`-f file-upload`

Uploading a file to the IBM Cloud Object Storage.

#### Additional mandatory arguments
| Argument | Description |
|:-------| :-----|
| -source | Path to the file to be uploaded |
| -key | Key name for uploaded file |

#### Example

```
hdbbackint -p hdbbackint.cfg -f file-upload -source ./sources/data.tar -key data.tar
```
