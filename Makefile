#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
PACKAGE_PATH:=github.com/solo-io/solo-kit
OUTPUT_DIR ?= $(ROOTDIR)/_output
SOURCES := $(shell find . -name "*.go" | grep -v test.go)
VERSION ?= $(shell git describe --tags)

#----------------------------------------------------------------------------------
# Repo init
#----------------------------------------------------------------------------------

# https://www.viget.com/articles/two-ways-to-share-git-hooks-with-your-team/
.PHONY: init
init:
	git config core.hooksPath .githooks

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
.PHONY: install-plugin
install-plugin: ${GOPATH}/bin/protoc-gen-solo-kit

.PHONY: install-gqlgen
install-gqlgen: ${GOPATH}/bin/gqlgen
${GOPATH}/bin/gqlgen: ${GOPATH}/src/github.com/99designs/gqlgen/
	cd ${GOPATH}/src/github.com/99designs/gqlgen/ && \
	dep ensure && \
	go install

${GOPATH}/src/github.com/99designs/gqlgen/:
	mkdir -p ${GOPATH}/src/github.com/99designs && \
	cd ${GOPATH}/src/github.com/99designs && \
	git clone --branch v0.4.4 --depth 1 https://github.com/99designs/gqlgen/

$(OUTPUT_DIR)/protoc-gen-solo-kit: $(SOURCES)
	go build -o $@ cmd/generator/main.go

${GOPATH}/bin/protoc-gen-solo-kit: $(OUTPUT_DIR)/protoc-gen-solo-kit
	cp $(OUTPUT_DIR)/protoc-gen-solo-kit ${GOPATH}/bin/


GOGO_PROTO_VERSION=v$(shell grep -C 1 github.com/gogo/protobuf  Gopkg.toml|grep version |cut -d'"' -f 2)
GOLANG_PROTO_VERSION=v$(shell grep -C 1 github.com/golang/protobuf  Gopkg.toml|grep version |cut -d'"' -f 2)
.PHONY: install-gen-tools
install-gogo-proto: ${GOPATH}/bin/protoc-gen-solo-kit
	mkdir -p  ${GOPATH}/src/github.com/gogo/ 
	mkdir -p  ${GOPATH}/src/github.com/golang/ 
	cd  ${GOPATH}/src/github.com/gogo/ && if [ -d protobuf ]; then cd protobuf && git fetch && git checkout $(GOGO_PROTO_VERSION); \
		else  git clone --branch $(GOGO_PROTO_VERSION) http://github.com/gogo/protobuf; fi
	cd  ${GOPATH}/src/github.com/golang/ && if [ -d protobuf ]; then cd protobuf && git fetch && git checkout $(GOLANG_PROTO_VERSION); \
		else  git clone --branch $(GOLANG_PROTO_VERSION) http://github.com/golang/protobuf; fi
	go install github.com/gogo/protobuf/protoc-gen-gogo

.PHONY: install-gen-tools
install-gen-tools: install-gogo-proto install-gqlgen install-plugin

#################
#################
#               #
#     Build     #
#               #
#               #
#################
#################
#################

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
	CGO_ENABLED=0 go build -o $@ projects/apiserver/cmd/main.go

$(OUTPUT_DIR)/apiserver-linux-amd64: apiserver-dependencies $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $@ projects/apiserver/cmd/main.go

$(OUTPUT_DIR)/apiserver-darwin-amd64: apiserver-dependencies $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o $@ projects/apiserver/cmd/main.go


$(OUTPUT_DIR)/Dockerfile.apiserver: $(APISERVER_DIR)/cmd/Dockerfile
	cp $< $@

apiserver-docker: $(OUTPUT_DIR)/apiserver-linux-amd64 $(OUTPUT_DIR)/Dockerfile.apiserver
	docker build -t soloio/apiserver-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.apiserver


gloo-i-docker-update:
	cd projects/apiserver/ui && if [ -d gloo-i ]; then cd gloo-i && git pull && cd ..; else  git clone https://github.com/solo-io/gloo-i gloo-i/; fi
	cd projects/apiserver/ui && docker build -t soloio/gloo-i-ee:$(VERSION) .

gloo-i-docker:
	cd projects/apiserver/ui && if [ -d gloo-i ]; then cd gloo-i && git pull && cd ..; else  git clone https://github.com/solo-io/gloo-i gloo-i/; fi
	cd projects/apiserver/ui && docker build -t soloio/gloo-i-ee:$(VERSION) .

#----------------------------------------------------------------------------------
# Gateway
#----------------------------------------------------------------------------------

GATEWAY_DIR=projects/gateway
GATEWAY_SOURCES=$(shell find $(GATEWAY_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gateway-linux-amd64: $(GATEWAY_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $@ $(GATEWAY_DIR)/cmd/main.go


.PHONY: gateway
gateway: $(OUTPUT_DIR)/gateway-linux-amd64

$(OUTPUT_DIR)/Dockerfile.gateway: $(GATEWAY_DIR)/cmd/Dockerfile
	cp $< $@

gateway-docker: $(OUTPUT_DIR)/gateway-linux-amd64 $(OUTPUT_DIR)/Dockerfile.gateway
	docker build -t soloio/gateway-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.gateway

#----------------------------------------------------------------------------------
# Sqoop
#----------------------------------------------------------------------------------

SQOOP_DIR=projects/sqoop
SQOOP_SOURCES=$(shell find $(SQOOP_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/sqoop-linux-amd64: $(SQOOP_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $@ $(SQOOP_DIR)/cmd/main.go


.PHONY: sqoop
sqoop: $(OUTPUT_DIR)/sqoop-linux-amd64

$(OUTPUT_DIR)/Dockerfile.sqoop: $(SQOOP_DIR)/cmd/Dockerfile
	cp $< $@

sqoop-docker: $(OUTPUT_DIR)/sqoop-linux-amd64 $(OUTPUT_DIR)/Dockerfile.sqoop
	docker build -t soloio/sqoop-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.sqoop

#----------------------------------------------------------------------------------
# Discovery
#----------------------------------------------------------------------------------

DISCOVERY_DIR=projects/discovery
DISCOVERY_SOURCES=$(shell find $(DISCOVERY_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/discovery-linux-amd64: $(DISCOVERY_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $@ $(DISCOVERY_DIR)/cmd/main.go


.PHONY: discovery
discovery: $(OUTPUT_DIR)/discovery-linux-amd64

$(OUTPUT_DIR)/Dockerfile.discovery: $(DISCOVERY_DIR)/cmd/Dockerfile
	cp $< $@

discovery-docker: $(OUTPUT_DIR)/discovery-linux-amd64 $(OUTPUT_DIR)/Dockerfile.discovery
	docker build -t soloio/discovery-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.discovery

#----------------------------------------------------------------------------------
# Gloo
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gloo-linux-amd64: $(GLOO_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $@ $(GLOO_DIR)/cmd/main.go


.PHONY: gloo
gloo: $(OUTPUT_DIR)/gloo-linux-amd64

$(OUTPUT_DIR)/Dockerfile.gloo: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@

gloo-docker: $(OUTPUT_DIR)/gloo-linux-amd64 $(OUTPUT_DIR)/Dockerfile.gloo
	docker build -t soloio/gloo-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.gloo

#----------------------------------------------------------------------------------
# Envot init
#----------------------------------------------------------------------------------

ENVOYINIT_DIR=cmd/envoyinit
ENVOYINIT_SOURCES=$(shell find $(ENVOYINIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/envoyinit-linux-amd64: $(ENVOYINIT_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $@ $(ENVOYINIT_DIR)/main.go

.PHONY: envoyinit
envoyinit: $(OUTPUT_DIR)/envoyinit-linux-amd64


$(OUTPUT_DIR)/Dockerfile.envoyinit: $(ENVOYINIT_DIR)/Dockerfile
	cp $< $@

data-plane-docker: $(OUTPUT_DIR)/envoyinit-linux-amd64 $(OUTPUT_DIR)/Dockerfile.envoyinit
	docker build -t soloio/data-plane-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.envoyinit


#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------
GH_ORG:=solo-io
GH_REPO:=solo-kit

RELEASE_BINARIES := \
	$(OUTPUT_DIR)/apiserver-linux-amd64 \
	$(OUTPUT_DIR)/apiserver-darwin-amd64 \
	$(OUTPUT_DIR)/gateway-linux-amd64 \
	$(OUTPUT_DIR)/gloo-linux-amd64 \
	$(OUTPUT_DIR)/discovery-linux-amd64 \
	$(OUTPUT_DIR)/envoyinit-linux-amd64

.PHONY: release-binaries
release-binaries: $(RELEASE_BINARIES)

.PHONY: release
release: release-binaries
	hack/create-release.sh github_api_token=$(GITHUB_TOKEN) owner=$(GH_ORG) repo=$(GH_REPO) tag=v$(VERSION)
	@$(foreach BINARY,$(RELEASE_BINARIES),hack/upload-github-release-asset.sh github_api_token=$(GITHUB_TOKEN) owner=solo-io repo=gloo tag=v$(VERSION) filename=$(BINARY);)

#----------------------------------------------------------------------------------
# Docker
#----------------------------------------------------------------------------------
#
#---------
#--------- Push
#---------

.PHONY: docker docker-push
docker: apiserver-docker discovery-docker gateway-docker gloo-docker sqoop-docker
docker-push:
	docker push soloio/sqoop-ee:$(VERSION) && \
	docker push soloio/gateway-ee:$(VERSION) && \
	docker push soloio/apiserver-ee:$(VERSION) && \
	docker push soloio/discovery-ee:$(VERSION) && \
	docker push soloio/gloo-ee:$(VERSION) && \
	docker push soloio/gloo-i-ee:$(VERSION)
