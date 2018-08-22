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
# Gateway
#----------------------------------------------------------------------------------
#
#---------
#--------- Graphql Stubs
#---------

GATEWAY_DIR=projects/gateway
GATEWAY_GRAPHQL_DIR=$(GATEWAY_DIR)/pkg/graphql
GATEWAY_GRAPHQL_GENERATED_FILES=$(GATEWAY_GRAPHQL_DIR)/models/generated.go $(GATEWAY_GRAPHQL_DIR)/graph/generated.go

.PHONY: gateway
gateway: $(GATEWAY_GRAPHQL_GENERATED_FILES)

$(GATEWAY_GRAPHQL_GENERATED_FILES): $(GATEWAY_GRAPHQL_DIR)/schema.graphql
	cd $(GATEWAY_GRAPHQL_DIR) && \
	gqlgen -v