#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
PACKAGE_PATH:=github.com/solo-io/solo-projects
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
# Generated Code
#----------------------------------------------------------------------------------

.PHONY: generated-code
generated-code: $(OUTPUT_DIR)/.generated-code

SUBDIRS:=projects
$(OUTPUT_DIR)/.generated-code: $(OUTPUT_DIR)/protoc-gen-solo-kit
	PATH=$(OUTPUT_DIR):$${PATH} go generate ./...
	gofmt -w $(SUBDIRS)
	goimports -w $(SUBDIRS)
	touch $@


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

ENVOYINIT_DIR=projects/envoyinit/cmd
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
GH_REPO:=solo-projects

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


#----------------------------------------------------------------------------------
# protoc plugin binary
#----------------------------------------------------------------------------------

$(OUTPUT_DIR)/protoc-gen-solo-kit: $(SOURCES)
	go build -o $@ cmd/protoc-gen-solo-kit/main.go
