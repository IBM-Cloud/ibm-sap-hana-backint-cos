# IBM Backint Agent for SAP HANA with IBM Cloud Object Storage

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/v/release/IBM-Cloud/ibm-sap-hana-backint-cos)](https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos/releases)

A high-performance Backint agent for SAP HANA that enables backup and recovery operations using IBM Cloud Object Storage. This agent is specifically designed and optimized for SAP workloads running on IBM Power Servers.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
  - [Configuration File Structure](#configuration-file-structure)
  - [Configuration Parameters](#configuration-parameters)
  - [Key Prefixes](#key-prefixes)
  - [Complete Configuration Example](#complete-configuration-example)
  - [Validate Configuration](#validate-configuration)
- [SAP HANA Configuration](#sap-hana-configuration)
  - [Recommended Parameters](#recommended-parameters)
  - [Parameter Combinations](#parameter-combinations)
- [Architecture](#architecture)
- [Troubleshooting](#troubleshooting)
- [Performance Tuning](#performance-tuning)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Overview

The IBM Backint agent for SAP HANA provides a certified interface between SAP HANA databases and IBM Cloud Object Storage (COS). It implements the SAP Backint API to enable:

- **Data backups**: Full and incremental database backups.
- **Log backups**: Continuous transaction log backups.
- **Catalog backups**: Backup catalog management.
- **Recovery operations**: Point-in-time and full database recovery.

This implementation supports **regional endpoints only** and requires object versioning and object lock to be enabled on the IBM Cloud Object Storage bucket.

## Key Features

- ✅ **High Performance**: Parallel multistreaming with configurable concurrency
- ✅ **Data Protection**: Object versioning, retention policies, and legal hold support
- ✅ **Reliability**: Automatic retry logic with exponential backoff for network issues
- ✅ **Flexibility**: Configurable chunk sizes, prefixes, and tagging
- ✅ **Enterprise Ready**: Comprehensive logging and error handling
- ✅ **SAP Certified**: Implements SAP Backint API specification

## Prerequisites

### 1. System Requirements

- **Operating System**: IBM Power Servers (ppc64le architecture).
- **SAP HANA**: Compatible with SAP HANA 2.0 and later.
- **Permissions**: User must have read/write access to SAP HANA directories.

### 2. IBM Cloud Object Storage Requirements

- **Existing COS Instance**: You must have an active IBM Cloud Object Storage instance.
- **Existing Bucket**: A bucket must be created within the COS instance.
- **Regional Resiliency**: Bucket must use **Regional resiliency location only**.
  - ❌ Cross Region resiliency is **not supported**.
- **Object Versioning**: Must be **enabled** on the bucket.
- **Object Lock**: Must be **enabled** on the bucket.
- **Endpoint Type**: Use **direct endpoints only** for optimal performance.

### 3. Network Configuration (Recommended for Best Performance)

For optimal backup performance, configure a Virtual Private Endpoint (VPE) gateway:

1. **Create IBM Cloud VPC**: Set up a Virtual Private Cloud (VPC) in IBM Cloud.
2. **Create VPE Gateway**: Create a Virtual Private Endpoint gateway of type **Cloud Object Storage** within the VPC.
3. **Configure Transit Gateway**: Attach both the VPC and the Power Virtual Server Workspace to the same Transit Gateway to enable connectivity.
4. **Update Hosts File**: Add the virtual private endpoint IP address and hostname to `/etc/hosts` on the Power Virtual Server instance:
   ```bash
   # Example entry in /etc/hosts
   <VPE_IP_ADDRESS>  s3.direct.<region>.cloud-object-storage.appdomain.cloud
   ```

This configuration provides:
- **Better Performance**: Direct, private network path to Cloud Object Storage
- **Enhanced Security**: Traffic stays within IBM Cloud private network
- **Lower Latency**: Reduced network hops compared to public endpoints

### 4. API Key Permissions

The API key used for authentication must have the following IBM Cloud IAM permissions:

| Permission                                             |
|--------------------------------------------------------|
| cloud-object-storage.bucket.head                       |
| cloud-object-storage.bucket.get_lifecycle              |
| cloud-object-storage.bucket.get                        |
| cloud-object-storage.object.put                        |
| cloud-object-storage.object.post_complete_upload       |
| cloud-object-storage.object.post_initiate_upload       |
| cloud-object-storage.object.put_part                   |
| cloud-object-storage.object.put_object_lock_retention  |
| cloud-object-storage.object.head                       |
| cloud-object-storage.object.get                        |
| cloud-object-storage.bucket.get_versioning             |
| cloud-object-storage.object.put_object_lock_legal_hold |
| cloud-object-storage.object.head_version               |
| cloud-object-storage.bucket.get_versions               |
| cloud-object-storage.object.get_version                |

## Installation

### Step 1: Download the Agent

Download the latest release from the [GitHub releases page](https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos/releases).

The release is provided as a ZIP archive with the naming format:
```
ibm-sap-hana-backint-cos_<VERSION>_linux_ppc64le.zip
```

**Automatic download of latest version:**
```bash
LATEST_VERSION=$(curl -s https://api.github.com/repos/IBM-Cloud/ibm-sap-hana-backint-cos/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
wget https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos/releases/download/v${LATEST_VERSION}/ibm-sap-hana-backint-cos_${LATEST_VERSION}_linux_ppc64le.zip
```

**Or download a specific version:**
```bash
VERSION=1.1.0  # Replace with desired version
wget https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos/releases/download/v${VERSION}/ibm-sap-hana-backint-cos_${VERSION}_linux_ppc64le.zip
```

### Step 2: Extract the Archive

```bash
# If you used LATEST_VERSION variable
unzip ibm-sap-hana-backint-cos_${LATEST_VERSION}_linux_ppc64le.zip

# Or if you used VERSION variable
unzip ibm-sap-hana-backint-cos_${VERSION}_linux_ppc64le.zip
```

The extracted package contains:
- `hdbbackint`: The IBM Backint agent executable.
- `sample_hdbbackint.cfg`: Sample configuration file.
- `README.md`: Documentation.
- `CHANGELOG.md`: Version history.
- `LICENSE`: License information.

### Step 3: Install to SAP HANA Directory

SAP HANA expects the `hdbbackint` executable at:
```
/usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint
```

**Option 1: Copy the executable**
```bash
chmod +x hdbbackint
sudo cp hdbbackint /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint
```

**Option 2: Create a symbolic link**
```bash
chmod +x hdbbackint
sudo ln -s $(pwd)/hdbbackint /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint
```

## Configuration

### Configuration File Structure

The IBM Backint agent requires a configuration file in INI format with the following sections:

- **[cloud_storage]** (mandatory): IBM Cloud Object Storage connection settings.
- **[backint]** (optional): Backint agent performance settings.
- **[objects]** (optional): Object naming and retention settings.
- **[trace]** (optional): Logging configuration.

Place the configuration file in a directory accessible by SAP HANA, preferably:
```
/usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg
```

### Configuration Parameters

#### [cloud_storage] Section (Mandatory)

| Parameter            | Values                                                                                                                          | Required | Description                                                                                                                      |
|----------------------|---------------------------------------------------------------------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| `auth_mode`          | `apikey`                                                                                                                        | ✅ Yes   | Authentication method (currently only apikey supported).                                                                         |
| `auth_keypath`       | `<file_path>`                                                                                                                   | ✅ Yes   | Full path to file containing IBM Cloud API key.                                                                                  |
| `bucket`             | `<bucket_name>`                                                                                                                 | ✅ Yes   | Name of IBM Cloud Object Storage bucket.                                                                                         |
| `region`             | `au-syd`, `br-sao`, `ca-tor`, `eu-de`, `eu-es`, `eu-gb`, `in-che`, `jp-osa`, `jp-tok`, `us-east`, `us-south`                            | ✅ Yes   | Region of the COS bucket.                                                                                                        |
| `endpoint_url`       | `<endpoint_url>`                                                                                                                | ✅ Yes   | Regional **direct** endpoint URL for the COS bucket (e.g., `https://s3.direct.us-south.cloud-object-storage.appdomain.cloud`).  |
| `ibm_auth_endpoint`  | `https://private.iam.cloud.ibm.com/identity/token` or `https://iam.cloud.ibm.com/identity/token`                               | ⬜ No    | IAM authentication endpoint<br>**Default**: `https://private.iam.cloud.ibm.com/identity/token`                                   |

#### [objects] Section (Optional)

| Parameter                        | Values                    | Required | Description                                                                                                                                                |
|----------------------------------|---------------------------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `remove_key_prefix`              | `<prefix_string>`         | ⬜ No    | String to remove from the beginning of backup file paths when creating object keys.                                                                        |
| `additional_key_prefix`          | `<prefix_string>`         | ⬜ No    | String to prepend to object keys (useful for organizing backups by database).                                                                              |
| `object_tags`                    | `<Key1=Val1,Key2=Val2>`   | ⬜ No    | Tags to apply to COS objects (max 10 key-value pairs).                                                                                                     |
| `object_lock_retention_mode`     | `None`, `cmp`, `gov`      | ⬜ No    | Object retention mode. `cmp` = COMPLIANCE mode. `gov` = GOVERNANCE mode. <br>**Default**: `None`                                                                                     |
| `object_lock_retention_period`   | `<years,months,days>`     | ⬜ No    | Retention period as comma-separated integers (e.g., `1,6,15` for 1 year, 6 months, 15 days). All three values required, at least one must be > 0.          |
| `object_lock_legal_hold_status`  | `ON`, `OFF`               | ⬜ No    | Legal hold prevents deletion until disabled.<br>**Default**: `OFF`                                                                                         |

#### [backint] Section (Optional)

| Parameter                  | Values                     | Required | Description                                                                                                                                                  |
|----------------------------|----------------------------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `max_concurrency`          | `<integer>`                | ⬜ No    | Number of concurrent upload threads per backup stream.<br>**Default**: `4`                                                                                   |
| `recover_max_concurrency`  | `<integer>`                | ⬜ No    | Number of concurrent download threads during restore.<br>**Default**: `4`                                                                                    |
| `multipart_chunksize`      | `<size>` or `<size><unit>` | ⬜ No    | Chunk size for multipart uploads. Units: `KB`, `MB`, `GB` (case-insensitive).<br>Examples: `134000000`, `100MB`, `1GB`.<br>**Default**: `134000000` (128 MB) |

#### [trace] Section (Optional)

| Parameter          | Values                                                  | Required | Description                                     |
|--------------------|---------------------------------------------------------|----------|-------------------------------------------------|
| `agent_log_level`  | `debug`, `info`, `warning`, `error`, `critical`, `http` | ⬜ No    | Logging verbosity level.<br>**Default**: `info` |

### Key Prefixes

The agent uses the full backup pipe name as the storage key by default. You can manipulate storage keys using prefix parameters.

**Example Backup Command:**
```sql
BACKUP DATA FOR <dbname> USING BACKINT ('/usr/sap/<sid>/SYS/global/hdb/backint/DB_<dbname>/<identifier>')
```

**Default Storage Key:**
```
/usr/sap/<sid>/SYS/global/hdb/backint/DB_<dbname>/<identifier>_databackup<postfix>
```

#### Using remove_key_prefix

Remove a prefix from the storage key:

```ini
[objects]
remove_key_prefix = /usr/sap/<sid>/SYS/global/hdb/backint/
```

**Result:**
```
DB_<dbname>/<identifier>_databackup<postfix>
```

#### Using additional_key_prefix

Add a prefix to the storage key:

```ini
[objects]
additional_key_prefix = myDB/
```

**Result:**
```
myDB/DB_<dbname>/<identifier>_databackup<postfix>
```

#### Combined Example

```ini
[objects]
remove_key_prefix = /usr/sap/<sid>/SYS/global/hdb/backint/
additional_key_prefix = PROD_HANA/
```

**Final Storage Key:**
```
PROD_HANA/DB_<dbname>/<identifier>_databackup<postfix>
```

### Complete Configuration Example

```ini
[cloud_storage]
auth_mode = apikey
auth_keypath = /usr/sap/HDB/SYS/global/hdb/opt/apikey.txt
bucket = my-hana-backups
region = us-south
endpoint_url = https://s3.direct.us-south.cloud-object-storage.appdomain.cloud
ibm_auth_endpoint = https://private.iam.cloud.ibm.com/identity/token

[objects]
remove_key_prefix = /usr/sap/HDB/SYS/global/hdb/backint/
additional_key_prefix = PROD_HDB/
object_tags = Environment=Production,Application=SAP_HANA,Owner=DBA_Team
object_lock_retention_mode = cmp
object_lock_retention_period = 0,3,0
object_lock_legal_hold_status = OFF

[backint]
max_concurrency = 4
recover_max_concurrency = 4
multipart_chunksize = 100MB

[trace]
agent_log_level = info
```

### Validate Configuration

Before using the agent, validate your configuration file:

```bash
/usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint -p /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg -check
```

A successful validation will output:
```
Configuration file is valid.
```

## SAP HANA Configuration

Configure SAP HANA to use the Backint agent by setting parameters in the `[backup]` section of `global.ini`:

```
/usr/sap/<SAPSID>/SYS/global/hdb/custom/config/global.ini
```

### Required Parameters

```ini
[backup]
data_backup_parameter_file = /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg
log_backup_parameter_file = /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg
catalog_backup_parameter_file = /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg
catalog_backup_using_backint = true
log_backup_using_backint = true
parallel_data_backup_backint_channels = 4
data_backup_buffer_size = <refer table Buffer Size Recommendations below>
parallel_data_backup_backint_size_threshold = 400
backint_response_timeout = 1800
```

**Parameter Descriptions:**

- **data_backup_parameter_file**: Path to configuration file for data backups (mandatory)
- **log_backup_parameter_file**: Path to configuration file for log backups (required if using Backint for logs)
- **catalog_backup_parameter_file**: Path to configuration file for catalog backups (required if using Backint for catalog)
- **catalog_backup_using_backint**: Enable Backint for catalog backups
- **log_backup_using_backint**: Enable Backint for log backups
- **parallel_data_backup_backint_channels**: Number of parallel backup streams
- **data_backup_buffer_size**: Buffer size in MB (see table below for recommendations)
- **parallel_data_backup_backint_size_threshold**: Minimum backup size in GB to trigger parallel processing
- **backint_response_timeout**: Timeout in seconds for Backint operations

#### Buffer Size Recommendations

| System Memory      | Recommended data_backup_buffer_size |
|--------------------|-------------------------------------|
| < 1 TB             | 1024 MB                             |
| ≥ 1 TB and < 6 TB  | 2048 MB                             |
| ≥ 6 TB and < 24 TB | 4096 MB                             |
| ≥ 24 TB            | 4096 MB                             |

### Parameter Combinations

Recommended combinations of SAP HANA and Backint agent parameters for optimal performance:

| HANA: parallel_data_backup_backint_channels | Agent: max_concurrency | Total Concurrent Operations |
|---------------------------------------------|------------------------|-----------------------------|
| 8                                           | 4                      | 32                          |
| 8                                           | 2                      | 16                          |
| 4                                           | 4                      | 16                          |
| 4                                           | 2                      | 8                           |

**Note**: Values above these combinations generally do not provide further performance improvements and may lead to resource bottlenecks.

### Apply Configuration Changes

After updating `global.ini`, apply changes:

```bash
su - <sid>adm
hdbnsutil -reconfig
```

**⚠️ Important**: The `hdbnsutil -reconfig` command is required for changes to take effect.

## Architecture

The IBM Backint agent implements a multi-threaded architecture:

1. **SAP HANA** initiates backup/restore operations via Backint API
2. **Backint Agent** receives requests through named pipes
3. **Upload/Download Manager** coordinates parallel operations
4. **IBM COS SDK** handles multipart uploads and downloads with retry logic
5. **IBM Cloud Object Storage** stores backup data with versioning and retention

```
┌─────────────┐
│  SAP HANA   │
└──────┬──────┘
       │ Backint API
       ▼
┌─────────────────────┐
│  Backint Agent      │
│  - Config Parser    │
│  - Pipe Manager     │
│  - Concurrency Ctrl │
└──────┬──────────────┘
       │ IBM COS SDK
       ▼
┌─────────────────────┐
│ IBM Cloud Object    │
│ Storage (Regional)  │
│  - Versioning       │
│  - Object Lock      │
│  - Retention        │
└─────────────────────┘
```

## Troubleshooting

### Common Issues

**Issue**: Configuration validation fails
```bash
# Check file syntax
cat /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg

# Validate configuration
hdbbackint -p /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg -check
```

**Issue**: Permission denied errors
```bash
# Check file permissions
ls -la /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint*

# Ensure <sid>adm can read files
sudo chown <sid>adm:sapsys /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint*
sudo chmod 750 /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint
sudo chmod 640 /usr/sap/<SID>/SYS/global/hdb/opt/hdbbackint.cfg
```

**Issue**: API key authentication fails
```bash
# Verify API key file exists and is readable
cat /path/to/apikey.txt

# Test API key manually
curl -X POST "https://private.iam.cloud.ibm.com/identity/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=$(cat /path/to/apikey.txt)"
```

**Issue**: Backup hangs or times out
- Increase `backint_response_timeout` in SAP HANA configuration.
- Check network connectivity to IBM Cloud Object Storage endpoint.
- Review agent logs for detailed error messages.
- Reduce `max_concurrency` if system resources are constrained.

### Log Locations

- **Backint Agent Logs**: Check SAP HANA trace directory
  ```
  /usr/sap/<SID>/HDB<instance>/trace/DB_<SID>/backint.log
  ```
- **SAP HANA Backup Logs**:
  ```
  /usr/sap/<SID>/HDB<instance>/trace/DB_<SID>/backup.log
  ```

## Performance Tuning

### Upload Performance

1. **Adjust max_concurrency**: Start with 4, increase if CPU/network allows.
2. **Optimize multipart_chunksize**: Larger chunks (100-512MB) for high-bandwidth networks.
3. **Increase parallel_data_backup_backint_channels**: Use 4-8 channels for large databases.
4. **Tune data_backup_buffer_size**: Match to available system memory.

### Download Performance

- **Adjust recover_max_concurrency**: Start with 4, increase for faster restores.

### Network Optimization

- Use IBM Cloud Virtual private endpoint gateways for better performance and security.

## Examples

### Perform a Full Data Backup

```sql
-- Connect to SAP HANA
hdbsql -u SYSTEM -d <DBNAME>

-- Execute backup
BACKUP DATA USING BACKINT ('/usr/sap/<SID>/SYS/global/hdb/backint/DB_<DBNAME>/FULL_BACKUP');
```

### Perform a Differential Backup

```sql
BACKUP DATA DIFFERENTIAL USING BACKINT ('/usr/sap/<SID>/SYS/global/hdb/backint/DB_<DBNAME>/DIFF_BACKUP');
```

### Restore from Backup

```sql
-- Recover to most recent state
RECOVER DATABASE UNTIL TIMESTAMP '<timestamp>' USING BACKINT;

-- Recover to specific backup
RECOVER DATA USING BACKINT ('<backup_id>');
```

### List Available Backups

```sql
-- View backup catalog
SELECT * FROM M_BACKUP_CATALOG ORDER BY SYS_START_TIME DESC;
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code standards and formatting
- Development setup
- Pull request process
- Commit signing requirements (DCO)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Copyright 2026 IBM Corp.

## Support

### Documentation

- [Changelog](CHANGELOG.md) - Version history and release notes
- [Security Policy](SECURITY.md) - Security vulnerability reporting
- [SNAPPY Tools](docs/SNAPPYTOOLS.md) - Using the agent with SNAPPY tools

### Getting Help

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos/discussions)
- **Security**: Report security vulnerabilities privately via [SECURITY.md](SECURITY.md)

### Related Resources

- [IBM Cloud Object Storage Documentation](https://cloud.ibm.com/docs/cloud-object-storage)
- [SAP HANA Backup and Recovery Guide](https://help.sap.com/docs/SAP_HANA_PLATFORM)
- [SAP Note 3644731](https://me.sap.com/notes/3644731) - Official SAP Note for this agent

---

**Maintained by IBM Cloud**
