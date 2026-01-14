# Contributing to ibm-sap-hana-backint-cos

Thank you for your interest in contributing! This project follows IBM’s open‑source participation guidelines and secure‑by‑default practices. Contributions are welcome via Pull Requests (PRs).

---

## Prerequisites

- **Go:** Use the version declared in `go.mod` (Go ≥ 1.24 recommended).
- **Git:** Ability to sign commits (DCO).
- **pre-commit:** Hooks are defined in `.pre-commit-config.yaml`.

---

## Getting Started

```bash
# Clone the repo
git clone https://github.com/IBM-Cloud/ibm-sap-hana-backint-cos.git
cd ibm-sap-hana-backint-cos

# Install pre-commit
pip3 install pre-commit --break-system-packages || pip3 install --user pre-commit
pre-commit install
pre-commit run --all-files  # optional initial cleanup

# Ensure module files are tidy
go mod tidy
````

## Coding Standards

*   **Formatting:** `go fmt ./...` (enforced via CI and pre-commit).
*   **Imports:** Prefer simplified imports (`gofmt -s` or `goimports` if enabled).
*   **Linting:** `go vet ./...` (required). `golangci-lint run` if configured.
*   **Modules:** `go mod tidy` must produce no diff before committing.
*   **Errors:** Wrap errors using `fmt.Errorf("…: %w", err)` when appropriate.
*   **Security:** Never commit credentials or secrets. Use secure stores or env vars.


## Commit Rules (DCO)

All commits must be signed using the Developer Certificate of Origin:

```bash
git commit -s -m "feat: add COS signer for backint uploads"
```

This adds the required `Signed-off-by:` line automatically.


## Pre-commit Hooks

This repository uses `pre-commit` to enforce formatting, linting, and hygiene checks.

```bash
pre-commit install
pre-commit run --all-files
```

## Pull Request Workflow

1.  Create a topic branch:
    *   `feat/...`
    *   `fix/...`
    *   `docs/...`
    *   `ci/...`

2.  Ensure your changes pass local checks:

    ```bash
    go fmt ./...
    go vet ./...
    go test ./... -race
    go mod tidy
    pre-commit run --all-files
    ```

3.  Update documentation if behavior or parameters change.

4.  Link related issues and describe testing, coverage impact, and risk.

5.  PRs require:
    *   At least one maintainer review
    *   Passing CI checks


## Release Process

This project uses semantic versioning (`vX.Y.Z`). Maintainers handle:

1.  Updating CHANGELOG (if present)
2.  Tagging releases (`vX.Y.Z`)
3.  Publishing release notes and artifacts (if applicable)


## Security

Do **not** file public GitHub issues for security vulnerabilities.

Please follow the private reporting process outlined in `SECURITY.md`.


## Support

Open non-security questions or discussions via GitHub issues with relevant reproduction details, logs, and environment information.


## Thank You

Your contributions help keep IBM SAP HANA Backint → IBM Cloud Object Storage integration robust, secure, and enterprise‑ready.
