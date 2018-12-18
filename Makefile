#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
PACKAGE_PATH:=github.com/solo-io/solo-projects
OUTPUT_DIR ?= $(ROOTDIR)/_output
SOURCES := $(shell find . -name "*.go" | grep -v test.go | grep -v '\.\#*')
VERSION ?= $(shell git describe --tags)
LDFLAGS := "-X github.com/solo-io/solo-projects/pkg/version.Version=$(VERSION)"

#----------------------------------------------------------------------------------
# Repo init
#----------------------------------------------------------------------------------

# https://www.viget.com/articles/two-ways-to-share-git-hooks-with-your-team/
.PHONY: init
init:
	git config core.hooksPath .githooks

#----------------------------------------------------------------------------------
# Clean
#----------------------------------------------------------------------------------

# Important to clean before pushing new releases. Dockerfiles and binaries may not update properly
.PHONY: clean
clean:
	rm -rf _output

#----------------------------------------------------------------------------------
# Generated Code
#----------------------------------------------------------------------------------

.PHONY: generated-code
generated-code: $(OUTPUT_DIR)/.generated-code

SUBDIRS:=projects
$(OUTPUT_DIR)/.generated-code:
	go generate ./...
	gofmt -w $(SUBDIRS)
	goimports -w $(SUBDIRS)
	touch $@


#################
#     Build     #
#################

#----------------------------------------------------------------------------------
# allprojects
#----------------------------------------------------------------------------------
# helper for testing
.PHONY: allprojects
allprojects: apiserver discovery gateway gloo glooctl rate-limit sqoop

#----------------------------------------------------------------------------------
# glooctl
#----------------------------------------------------------------------------------

CLI_DIR=projects/gloo/cli

$(OUTPUT_DIR)/glooctl: $(SOURCES)
	go build -ldflags=$(LDFLAGS) -o $@ $(CLI_DIR)/cmd/main.go


$(OUTPUT_DIR)/glooctl-linux-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(CLI_DIR)/cmd/main.go


$(OUTPUT_DIR)/glooctl-darwin-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags=$(LDFLAGS) -o $@ $(CLI_DIR)/cmd/main.go

.PHONY: glooctl
glooctl: $(OUTPUT_DIR)/glooctl
.PHONY: glooctl-linux-amd64
glooctl-linux-amd64: $(OUTPUT_DIR)/glooctl-linux-amd64
.PHONY: glooctl-darwin-amd64
glooctl-darwin-amd64: $(OUTPUT_DIR)/glooctl-darwin-amd64

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
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o $@ projects/apiserver/cmd/main.go

$(OUTPUT_DIR)/apiserver-linux-amd64: apiserver-dependencies $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ projects/apiserver/cmd/main.go

$(OUTPUT_DIR)/apiserver-darwin-amd64: apiserver-dependencies $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags=$(LDFLAGS) -o $@ projects/apiserver/cmd/main.go


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
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(GATEWAY_DIR)/cmd/main.go


.PHONY: gateway
gateway: $(OUTPUT_DIR)/gateway-linux-amd64

$(OUTPUT_DIR)/Dockerfile.gateway: $(GATEWAY_DIR)/cmd/Dockerfile
	cp $< $@

gateway-docker: $(OUTPUT_DIR)/gateway-linux-amd64 $(OUTPUT_DIR)/Dockerfile.gateway
	docker build -t soloio/gateway-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.gateway

#----------------------------------------------------------------------------------
# RateLimit
#----------------------------------------------------------------------------------

RATELIMIT_DIR=projects/rate-limit
RATELIMIT_SOURCES=$(shell find $(RATELIMIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/rate-limit-linux-amd64: $(RATELIMIT_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(RATELIMIT_DIR)/cmd/main.go

.PHONY: rate-limit
rate-limit: $(OUTPUT_DIR)/rate-limit-linux-amd64

$(OUTPUT_DIR)/Dockerfile.rate-limit: $(RATELIMIT_DIR)/cmd/Dockerfile
	cp $< $@

rate-limit-docker: $(OUTPUT_DIR)/rate-limit-linux-amd64 $(OUTPUT_DIR)/Dockerfile.rate-limit
	docker build -t soloio/rate-limit-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.rate-limit

#----------------------------------------------------------------------------------
# Sqoop
#----------------------------------------------------------------------------------

SQOOP_DIR=projects/sqoop
SQOOP_SOURCES=$(shell find $(SQOOP_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/sqoop-linux-amd64: $(SQOOP_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(SQOOP_DIR)/cmd/main.go


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
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(DISCOVERY_DIR)/cmd/main.go


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
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(GLOO_DIR)/cmd/main.go


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
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go

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
	$(OUTPUT_DIR)/rate-limit-linux-amd64 \
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
docker: apiserver-docker discovery-docker gateway-docker rate-limit-docker gloo-docker sqoop-docker
docker-push:
	docker push soloio/sqoop-ee:$(VERSION) && \
	docker push soloio/gateway-ee:$(VERSION) && \
	docker push soloio/rate-limit-ee:$(VERSION) && \
	docker push soloio/apiserver-ee:$(VERSION) && \
	docker push soloio/discovery-ee:$(VERSION) && \
	docker push soloio/gloo-ee:$(VERSION) && \
	docker push soloio/gloo-i-ee:$(VERSION)
