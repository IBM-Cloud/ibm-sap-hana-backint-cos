# IBM Backint Agent golang version

## Building executable

```
go mod init hdbbackint
go mod tidy
env GOOS=linux GOARCH=ppc64le go build -trimpath hdbbackint
```
