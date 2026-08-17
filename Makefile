BINARY   := hdbbackint
GOOS     := linux
GOARCH   := ppc64le
BUILD_FLAGS := -trimpath
LDFLAGS  := -s -w

.PHONY: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
	  go build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) .
