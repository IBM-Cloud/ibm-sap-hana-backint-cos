# 0.0.3 (January 14, 2026)

## **New Features**

* **Configuration Validation (`--check`)**
  Added a new command-line option `--check` that validates the `hdbbackint.cfg` file.
  Use this flag to verify configuration correctness before running backups or restores.

* **Flexible Multipart Chunk Size Units**
  The `multipart_chunksize` parameter now supports multiple units — specify sizes in **KB**, **MB**, **GB** not case sensitive or as plain integers for greater flexibility and clarity.


## **Fixes & Improvements**

* Improved error handling for `hdbbackint.cfg`
  The agent no longer crashes when encountering invalid configuration content. Instead, detailed error messages are displayed to help identify and correct the issue quickly.
