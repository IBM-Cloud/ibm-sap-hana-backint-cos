# [3644731](https://me.sap.com/notes/3644731) - Install and Configure IBM Backint agent for SAP HANA with IBM Cloud Object Storage

# Symptom

This SAP Note applies only to the IBM Backint agent for SAP HANA with IBM Cloud Object Storage and SAP workloads running on IBM Power Servers.

You can use IBM Backint agent for SAP HANA to backup and recover SAP HANA using Cloud Object Storage offered by IBM.

IBM Backint agent for SAP HANA with IBM Cloud Object Storage **supports regional endpoints only**.

This SAP Note describes how to install and configure the IBM Backint for SAP HANA with IBM Cloud Object Storage.

# Other Terms

SAP HANA, Backint, IBM, Cloud Object Storage, Backup and Recovery, Power Virtual Server

# Prerequisites

1. **IBM Cloud Object Storage Requirements**

* It is required to have an **existing** IBM Cloud Object Storage (COS) instance and an **existing** bucket within this instance.
* Bucket should be in **Regional resiliency location only**. Cross Region resiliency and Single data center options are not supported.
* **Object versioning** and **object lock** must be **enabled on the bucket**.

2. **API key permissions**

* To authenticate and upload/restore from IBM Cloud Object storage, current authentication mechanism is using an API KEY with following permissions are required:

|**Role**|
| - |
| cloud-object-storage.bucket.head |
| cloud-object-storage.bucket.get_lifecycle |
| cloud-object-storage.bucket.get |
| cloud-object-storage.object.put |
| cloud-object-storage.object.post_complete_upload |
| cloud-object-storage.object.post_initiate_upload |
| cloud-object-storage.object.put_part |
| cloud-object-storage.object.put_object_lock_retention |
| cloud-object-storage.object.head |
| cloud-object-storage.object.get |
| cloud-object-storage.bucket.get_versioning |
| cloud-object-storage.object.put_object_lock_legal_hold |
| cloud-object-storage.object.head_version |
| cloud-object-storage.bucket.get_versions |
| cloud-object-storage.object.get_version |


# Solution

## Install the agent

1. Download the latest release from [github repository](https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos/releases).

2. **Unpack the package** to a directory of your choice with the correct permissions in the system.

   The extracted package contains

   * `hdbbackint`: The IBM Backint agent **executable** for SAP HANA.
   * `sample_hdbbackint.cfg`: A **sample configuration** file for IBM SAP HANA backint agent.
   * `Readme.pdf`: Agent **manual**


3. **Create Symbolic Links**

   SAP HANA expects the `hdbbackint` executable to be in the following path:

   `/usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint`

   * Option 1: Copy the `hdbbackint` executable to the path `/usr/sap/<SID>/SYS/global/hdb/opt/`
   * Option 2: Create a symbolic link that points from `/usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint` to the `hdbbackint` executable that was extracted from the package.


## Configure the IBM Backint agent for SAP HANA with IBM Cloud Object Storage

   The IBM Backint agent for SAP HANA with IBM Cloud Object Storage requires a parameter file in the INI file format.

   An example parameter file `sample_hdbbackint.cfg` is part of the release package.

   You can use this example file, or create a new configuration file.

   Place the configuration file in a directory that has the permissions to allow SAP HANA to access it, preferably in path `/usr/sap/<SID>/SYS/global/hdb/opt/` where the executable exists.



   The configuration file consists of the following sections:

   * cloud_storage (mandatory)
   * backint (optional)
   * objects (optional)
   * trace (optional)

   To make sure that the `hdbbackint` agent runs without errors, first the configuration file is validated. Defaults are set if these parameters are not defined in the file. The configuration file is mandatory to execute the `hdbbackint` agent.


   The configuration file consists of the following sections which contain key-value pair settings:


| Section       | Key                           | Possible Values                                                                            |           | Description                                                                                                                                                                                                                                                                                                                      |
|---------------|-------------------------------|--------------------------------------------------------------------------------------------|-----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| cloud_storage | auth_mode                     | apikey                                                                                     | Mandatory | Possible authentication options                                                                                                                                                                                                                                                                                                  |
|               | auth_keypath                  | <api_key_file_path>                                                                        | Optional  | Full pathname to file containing the just the IBM Cloud api key.  Required if the auth_mode type is "apikey".                                                                                                                                                                                                                    |
|               | bucket                        | <bucket_name>                                                                              | Mandatory | Name of Cloud Object Storage bucket                                                                                                                                                                                                                                                                                              |
|               | region                        | au-syd, br-sao, ca-tor, eu-de, eu-es, eu-gb, jp-osa, jp-tok, us-east, us-south             | Mandatory | Region of Cloud Object Storage bucket                                                                                                                                                                                                                                                                                            |
|               | endpoint_url                  | <endpoint_url>                                                                             | Mandatory | Endpoint URL of Cloud Object Storage bucket                                                                                                                                                                                                                                                                                      |
|               | ibm_auth_enpoint              | https://private.iam.cloud.ibm.com/identity/token, https://iam.cloud.ibm.com/identity/token | Optional  | URL used for IAM authentication.  **Default**: https://private.iam.cloud.ibm.com/identity/token                                                                                                                                                                                                                                      |
| objects       | remove_key_prefix             | <prefix_string>                                                                            | Optional  | Backint uses the whole pipe name as the storage key for backups.  You can specify a string to be removed from the resulting storage key.                                                                                                                                                                                         |
|               | additional_key_prefix         | <prefix_string>                                                                            | Optional  | You can add database-specific prefix to the storage key for backups.                                                                                                                                                                                                                                                             |
|               | object_tags                   | <Key1=Val1,Key2=Val2>                                                                      | Optional  | Tags added to Cloud Object storage object. A maximum of 10 key value pairs is supported.                                                                                                                                                                                                                                         |
|               | object_lock_retention_mode    | None, cmp                                                                                  | Optional  | If set to "cmp", the Object Retention is switched on. For more information see [retention period](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-ol-overview#ol-terminology-retention-period) feature for IBM Cloud Object Storage.   **Default**: None                                              |
|               | object_lock_retention_period  | <object_lock_retention_period>                                                             | Optional  | If set to "cmp", the Object Retention is switched on. For more information see retention period feature for IBM Cloud Object Storage.   **Default**: None                                                                                                                                                                        |
|               | object_lock_legal_hold_status | ON, OFF                                                                                    | Optional  | A legal hold is like a retention period in that it prevents an object version from being overwritten or deleted. For more information see [legal hold](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-ol-overview#ol-terminology-legal-hold) feature for IBM Cloud object storage.  **Default**: OFF |
| backint       | max_concurrency               | <value_integer>                                                                                  | Optional  | Number of concurrent requests made to IBM Cloud object Storage. This value should be configured based on system resources.  **Default**: 10                                                                                                                                                                                      |
|               | multipart_chunksize           | <size_in_bytes> or `<size><unit>`, while `<unit>` can be one of the following: KB, MB or GB (not case sensitive), and `<size>` must not be 0.                                                                      | Optional  | Data transfer chunk size. This value should be configured based on system resources.  **Default**: 134000000                                                                                                                                                                                                                     |
| trace         | agent_log_level               | debug, info, warning, error,critical, http                                                                | Optional  | Trace level for the IBM SAP HANA Backint Agent for IBM Cloud Object Storage.  **Default**: info                                                                                                                                                                                                                                  |

### Key Prefixes

Assuming starting a backup using the following command: 

````
BACKUP DATA FOR <dbname> USING BACKINT ('/usr/sap/<sid>/SYS/global/hdb/backint/DB_<dbname>/<identifier>)
````

By default, the `hdbbackint` agent uses the whole backup name `/usr/sap/<sid>/SYS/global/hdb/backint/DB_<dbname>/<identifier>_databackup<post_fix>` as the storage key for backups.

To manipulate the storage keys, `hdbbackint` provides two parameters which can be defined in the `hdbbackint.cfg` file.

**remove_key_prefix**

The value of the key _remove_key_prefix_ specifies a string which will be removed from the beginning of the storage key name.

For example the following setting

```
remove_key_prefix = /usr/sap/<sid>/SYS/global/hdb/backint/
```

will result in the shorter storage key: `DB_<dbname>/<identifier>_databackup<post_fix>`

**additional_key_prefix**

The value of the key _additional_key_prefix_ specifies a string which will be added to the beginning of the storage key name.

For example:

```
additional_key_prefix = myDB/
````

will prepend `myDB/` to the resulting storage key.

For example, if the following options are used together:

```
remove_key_prefix = /usr/sap/<sid>/SYS/global/hdb/backint/
additional_key_prefix = myDB/
````

Then the final storage key will be:

`myDB/DB_<dbname>/<identifier>_databackup<post_fix>`


### Validate the hdbbackint configuration file 

The configuration file of the `hdbbackint` agent can be validated by executing the following command: 

```
hdbbackint -p <hdbbackint_configuration_file> -check
```


## Configure SAP HANA database to use the parameter File

SAP HANA database uses the following parameters to configure the usage of the configuration file `hdbbackint.cfg`. These parameters are set in the `backup` section of the `global.ini` file. This file is located in `/usr/sap/<SAPSID>/SYS/global/hdb/custom/config/`.

```
[backup]
data_backup_parameter_file = <path_to_hdbbackint.cfg>
log_backup_parameter_file = <path_to_hdbbackint.cfg>
catalog_backup_parameter_file =< path_to_hdbbackint.cfg>
catalog_backup_using_backint = true
log_backup_using_backint = true
parallel_data_backup_backint_channels = 8
data_backup_buffer_size = 1024
parallel_data_backup_backint_size_threshold = 400
backint_response_timeout = 1800
```

* **data_backup_parameter_file**

   Mandatory. Used or data backups.

* **log_backup_parameter_file**

   If log backups are written using Backint, this parameter must be configured.

* **catalog_backup_parameter_file**

   If catalog backups are done using Backint, this parameter must be configured.

* **parallel_data_backup_backint_channels**

   Specify the number of channels to be used for multistreaming.

**Caution: After updating the entries it is important to run the `hdbnsutil -reconfig` command as sidadm user for the changes to take effect.**

### Recommended Configuration Parameters

The value of the **data_backup_buffer_size** parameter should be set based on the total memory available on the VM. The following sizes are recommended:

| System Memory     | Recommended data_backup_buffer_size |
| ------------------| ------------------------------------|
| < 1TB             | 1024                                |
| ≥ 1 TB and < 6TB  | 2048                                |
| ≥ 6 TB and < 24TB | 4096                                |
| ≥ 24 TB           | 4096                                |

#### HANA & Backint Agent Recommended Parameter Combination
This section lists recommended combinations of:

* **HANA parameter:** 

   **parallel_data_backup_backint_channels** defines how many Backint channels SAP HANA starts in parallel during a data backup.
* **Backint Agent parameter:** 

   **max_concurrency** controls the maximum number of parallel processing threads used by the `hdbbackint` agent.

These combinations help ensure optimal backup throughput and resource utilization.


| HANA Parameter (parallel_data_backup_backint_channels) | Agent Parameter (max_concurrency) |
| -------------------------------------------------------| ----------------------------------|
| 8                                                      | 4                                 |
| 8                                                      | 2                                 |
| 4                                                      | 4                                 |
| 4                                                      | 2                                 |

**Note: Values above these combinations generally do not provide further performance improvements and may lead to resource bottlenecks.**
