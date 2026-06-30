# 2.0.0 (June 30, 2026)

## **New Features**

* **PowerVS OAuth Authentication (`auth_mode = oauth`)**
  Added a new keyless authentication mode for SAP HANA backup agents running on IBM Power Virtual Server.
  Set `auth_mode = oauth` to authenticate via the VPC Instance Metadata Service — no API key required.
  Tokens are automatically obtained and refreshed; the provider proactively renews 5 minutes before expiry
  so no in-flight S3 request ever sees a 401. The IAM `refresh_token` flow is deliberately bypassed as
  PowerVS-sourced tokens carry no refresh token. Existing `auth_mode = apikey` configurations are
  **fully compatible and require no changes**.

## **Bug Fixes**

* **Fixed Backup Upload Failure — TCP Write Timeout on PowerVS**
  On PowerVS, a stateful firewall can silently drop a long-running PUT TCP connection mid-part, causing
  a `write: connection timed out` error. The default SDK retryer treated this as non-retryable.
  A custom `cosRetryer` now wraps the SDK's `DefaultRetryer` and explicitly retries on write timeouts,
  re-dialling a fresh socket for each attempt (max 5 retries).

## **Dependency Updates**

* **IBM Cloud Go SDK**: added `github.com/IBM/go-sdk-core/v5` for `VpcInstanceAuthenticator`
* **Go**: 1.25.0 → 1.26.1

---

# 1.2.0 (April 9, 2026)

## **New Features**

* **in-mum Region Support**
  Added support for IBM Cloud Object Storage Mumbai (in-mum) datacenter region

# 1.1.0 (April 8, 2026)

## **New Features**

* **in-che Region Support**
  Added support for IBM Cloud Object Storage Chennai (in-che) datacenter region

# 1.0.0 (March 31, 2026)

## **Critical Fixes**

* **Fixed Deadlock During Parallel Restore Operations**
  Resolved a critical deadlock issue that caused restore operations to hang at ~96% completion when multiple files were being restored in parallel. The issue was caused by a global lock being held during blocking pipe write operations, preventing other download threads from storing data in the buffer. Implemented per-pipe mutexes to allow independent progress across parallel restore operations while maintaining data integrity.

* **Fixed Context Leak in Pipe Write Operations**
  Corrected a context leak where `defer cancel()` was called inside a loop, causing multiple uncancelled contexts to accumulate. This led to resource exhaustion and timeout errors. Now properly cancels contexts immediately after each write operation completes or times out.

* **Fixed Closure Variable Capture Bug**
  Resolved a bug where goroutines were capturing loop variables by reference, causing incorrect data slices to be written to pipes. Now properly captures loop iteration values before passing them to goroutines.

* **Fixed Premature Pipe Closure**
  Corrected an issue where pipes were being closed before all download operations completed, resulting in incomplete data transfers. Changed from asynchronous to synchronous waiting for download completion to ensure all data is written before pipes are closed.

## **Improvements**

* **Increased Pipe Write Timeout**
  Increased the default pipe write timeout from 30 seconds to 300 seconds (5 minutes) to accommodate slower HANA read operations during intensive recovery processing.

* **Added Network Retry Logic**
  Implemented automatic retry mechanism (3 attempts with exponential backoff) for S3 download operations to handle transient network issues such as connection timeouts, temporary endpoint slowness, and TCP connection drops.

* **Improved Data Safety**
  Added data copying before buffer deletion to prevent use-after-free bugs that could cause data corruption when the buffer lock is released before pipe writes complete.

* **Concurrency for Download**
  Added recover_max_concurrency to backint configuration file.

* **New GOVERNANCE option for object_lock_retention_mode**
  Added gov for object_lock_retention_mode parameter.

## **Dependency Updates**

* **Go**: 1.24.1 → 1.25.0
* **IBM COS SDK**: v1.13.0 → v1.14.0
* **github.com/gabriel-vasile/mimetype**: v1.4.11 → v1.4.13
* **github.com/go-openapi/errors**: v0.22.4 → v0.22.7
* **github.com/go-openapi/strfmt**: v0.25.0 → v0.26.1
* **github.com/go-playground/validator/v10**: v10.28.0 → v10.30.1
* **github.com/go-viper/mapstructure/v2**: v2.4.0 → v2.5.0
* **go.mongodb.org/mongo-driver**: v1.17.6 → v1.17.9
* **go.yaml.in/yaml/v2**: v2.4.3 → v2.4.4
* **golang.org/x/crypto**: v0.47.0 → v0.49.0
* **golang.org/x/net**: v0.49.0 → v0.52.0
* **golang.org/x/sys**: v0.40.0 → v0.42.0
* **golang.org/x/text**: v0.33.0 → v0.35.0
* **Added**: github.com/oklog/ulid/v2 v2.1.1

## **Technical Details**

* Implemented two-level locking strategy:
  - Buffer lock: Short-lived, protects shared buffer access across all pipes
  - Per-pipe mutex: Held during blocking pipe writes, allows parallel writes to different pipes
* Removed unnecessary `time.Sleep()` calls after pipe writes
* Enhanced error logging for retry attempts and timeout scenarios

# 0.0.3 (January 14, 2026)

## **New Features**

* **Configuration Validation (`--check`)**
  Added a new command-line option `--check` that validates the `hdbbackint.cfg` file.
  Use this flag to verify configuration correctness before running backups or restores.

* **Flexible Multipart Chunk Size Units**
  The `multipart_chunksize` parameter now supports multiple units — specify sizes in **KB**, **MB**, **GB** (not case sensitive) or as plain integers for greater flexibility and clarity.


## **Fixes & Improvements**

* Improved error handling for `hdbbackint.cfg`.
  The agent no longer crashes when encountering invalid configuration content. Instead, detailed error messages are displayed to help identify and correct the issue quickly.
