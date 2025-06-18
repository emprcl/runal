BIN                    := bin
BIN_NAME               := runal
CLI_DIR                := cli
CLI_ENTRYPOINT         := main.go
GOLANG_BIN             := go
CGO_ENABLED            := 0
GOLANG_OS              := linux
GOLANG_ARCH            := amd64
GOLANG_BUILD_OPTS      += GOOS=$(GOLANG_OS)
GOLANG_BUILD_OPTS      += GOARCH=$(GOLANG_ARCH)
GOLANG_BUILD_OPTS      += CGO_ENABLED=$(CGO_ENABLED)
GOLANG_LINT            := $(BIN)/golangci-lint
GOLANG_LDFFLAGS        := -ldflags="-w -s"

$(BIN):
	mkdir -p $(BIN)

$(GOLANG_LINT): $(BIN)
	GOBIN=$$(pwd)/$(BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

build: $(BIN)
	$(GOLANG_BUILD_OPTS) $(GOLANG_BIN) build -C $(CLI_DIR) $(GOLANG_LDFFLAGS) -o ../$(BIN)/$(BIN_NAME) $(CLI_ENTRYPOINT)
	chmod +x $(BIN)/$(BIN_NAME)

checks: $(GOLANG_LINT)
	$(GOLANG_LINT) run ./...

test:
	$(GOLANG_BIN) test ./...

bench:
	$(GOLANG_BIN) test -bench=. -benchmem ./...
