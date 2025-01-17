BASE_DIR ?= $(CURDIR)
SILENT ?= @

# Protocol Buffers

PROTOC_VERSION := 23.1

PROTOC_DIR := $(BASE_DIR)/.proto
$(PROTOC_DIR):
	$(SILENT)mkdir -p "$@"

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
PROTOC_OS = linux
endif
ifeq ($(UNAME_S),Darwin)
PROTOC_OS = osx
endif
PROTOC_ARCH=$(shell case $$(uname -m) in (arm64) echo aarch_64 ;; (*) uname -m ;; esac)

DOWNLOAD_DIR := $(PROTOC_DIR)/.downloads
$(DOWNLOAD_DIR):
	$(SILENT)mkdir -p "$@"

PROTOC_ZIP := protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip
PROTOC_FILE := $(DOWNLOAD_DIR)/$(PROTOC_ZIP)

.PRECIOUS: $(PROTOC_FILE)
$(PROTOC_FILE): $(DOWNLOAD_DIR)
	curl --output-dir $(DOWNLOAD_DIR) -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)"

PROTO_BIN := $(PROTOC_DIR)/bin
$(PROTO_BIN):
	$(SILENT)mkdir -p "$@"

PROTOC := $(PROTO_BIN)/protoc
$(PROTOC): $(PROTOC_FILE)
	$(SILENT)unzip -q -o -d "$(PROTOC_DIR)" "$(PROTOC_FILE)"
	$(SILENT)test -x "$@"

PROTOC_GEN_GO_BIN := $(PROTO_BIN)/protoc-gen-go
$(PROTOC_GEN_GO_BIN): $(PROTO_BIN)
	GOBIN=$(PROTO_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go

PROTOC_GEN_GO_GRPC_BIN := $(PROTO_BIN)/protoc-gen-go-grpc
$(PROTOC_GEN_GO_GRPC_BIN): $(PROTO_BIN)
	GOBIN=$(PROTO_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: proto-install
proto-install: $(PROTOC) $(PROTOC_GEN_GO_BIN) $(PROTOC_GEN_GO_GRPC_BIN)
