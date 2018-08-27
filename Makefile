#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
PACKAGE_PATH:=github.com/solo-io/solo-kit
OUTPUT_DIR ?= $(ROOTDIR)/_output
SOURCES := $(shell find . -name "*.go" | grep -v test)

#----------------------------------------------------------------------------------
# Protobufs
#----------------------------------------------------------------------------------

PROTOS := $(shell find api/v1 -name "*.proto")
GENERATED_PROTO_FILES := $(shell find pkg/api/v1/resources/core -name "*.pb.go")

.PHONY: all
all: build

.PHONY: proto
proto: $(GENERATED_PROTO_FILES)

$(GENERATED_PROTO_FILES): $(PROTOS)
	cd api/v1/ && \
	protoc \
	--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:$(GOPATH)/src/ \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/ \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf/ \
	-I=. \
	./*.proto

#----------------------------------------------------------------------------------
# Kubernetes Clientsets
#----------------------------------------------------------------------------------

$(OUTPUT_DIR):
	mkdir -p $@

.PHONY: clientset
clientset: $(OUTPUT_DIR) $(OUTPUT_DIR)/.clientset

$(OUTPUT_DIR)/.clientset: $(GENERATED_PROTO_FILES) $(SOURCES)
	cd ${GOPATH}/src/k8s.io/code-generator && \
	./generate-groups.sh all \
		$(PACKAGE_PATH)/pkg/api/v1/clients/kube/crd/client \
		$(PACKAGE_PATH)/pkg/api/v1/clients/kube/crd \
		"solo.io:v1"
	touch $@

#----------------------------------------------------------------------------------
# Generated Code
#----------------------------------------------------------------------------------

.PHONY: generated-code
generated-code: $(OUTPUT_DIR)/.generated-code

SUBDIRS:=pkg projects test
$(OUTPUT_DIR)/.generated-code:
	go generate ./...
	gofmt -w $(SUBDIRS)
	goimports -w $(SUBDIRS)
	touch $@

#----------------------------------------------------------------------------------
# protoc plugin binary
#----------------------------------------------------------------------------------

$(OUTPUT_DIR)/protoc-gen-solo-kit: $(SOURCES)
	go build -o $@ cmd/generator/main.go


.PHONY: install-plugin
install-plugin: $(OUTPUT_DIR)/protoc-gen-solo-kit
	cp $(OUTPUT_DIR)/protoc-gen-solo-kit ${GOPATH}/bin/


#----------------------------------------------------------------------------------
# Apiserver
#----------------------------------------------------------------------------------
#
#---------
#--------- Graphql Stubs
#---------

APISERVER_DIR=projects/apiserver
APISERVER_GRAPHQL_DIR=$(APISERVER_DIR)/pkg/graphql
APISERVER_GRAPHQL_GENERATED_FILES=$(APISERVER_GRAPHQL_DIR)/models/generated.go $(APISERVER_GRAPHQL_DIR)/graph/generated.go
APISERVER_SOURCES=$(shell find $(APISERVER_GRAPHQL_DIR) -name "*.go" | grep -v test | grep -v generated.go)

.PHONY: apiserver-dependencies
apiserver-dependencies: $(APISERVER_GRAPHQL_GENERATED_FILES)
$(APISERVER_GRAPHQL_GENERATED_FILES): $(APISERVER_DIR)/schema.graphql $(APISERVER_DIR)/gqlgen.yaml $(APISERVER_SOURCES)
	cd $(APISERVER_DIR) && \
	gqlgen -v

.PHONY: apiserver
apiserver: $(OUTPUT_DIR)/apiserver

# TODO(ilackarms): put these inside of a loop or function of some kind
$(OUTPUT_DIR)/apiserver: apiserver-dependencies $(SOURCES)
	go build -o $@ projects/apiserver/cmd/main.go

$(OUTPUT_DIR)/apiserver-linux-amd64:  $(SOURCES)
	GOOS=linux go build -o $@ projects/apiserver/cmd/main.go

$(OUTPUT_DIR)/apiserver-darwin-amd64:  $(SOURCES)
	GOOS=darwin go build -o $@ projects/apiserver/cmd/main.go

#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------

VERSION:=1
GH_ORG:=solo-io
GH_REPO:=solo-kit

RELEASE_BINARIES := \
	$(OUTPUT_DIR)/apiserver-linux-amd64 \
	$(OUTPUT_DIR)/apiserver-darwin-amd64

.PHONY: release-binaries
release-binaries: $(RELEASE_BINARIES)

.PHONY: release
release: release-binaries
	hack/create-release.sh github_api_token=$(GITHUB_TOKEN) owner=$(GH_ORG) repo=$(GH_REPO) tag=v$(VERSION)
	@$(foreach BINARY,$(RELEASE_BINARIES),hack/upload-github-release-asset.sh github_api_token=$(GITHUB_TOKEN) owner=solo-io repo=gloo tag=v$(VERSION) filename=$(BINARY);)

