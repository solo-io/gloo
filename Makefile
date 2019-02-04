#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
PACKAGE_PATH:=github.com/solo-io/solo-projects
OUTPUT_DIR ?= $(ROOTDIR)/_output
SOURCES := $(shell find . -name "*.go" | grep -v test.go | grep -v '\.\#*')
RELEASE := "true"
ifeq ($(TAGGED_VERSION),)
	# TAGGED_VERSION := $(shell git describe --tags)
	# This doesn't work in CI, need to find another way...
	TAGGED_VERSION := vdev
	RELEASE := "false"
endif
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)
LDFLAGS := "-X github.com/solo-io/solo-projects/pkg/version.Version=$(VERSION)"

#----------------------------------------------------------------------------------
# Repo setup
#----------------------------------------------------------------------------------

# https://www.viget.com/articles/two-ways-to-share-git-hooks-with-your-team/
.PHONY: init
init:
	git config core.hooksPath .githooks

.PHONY: update-deps
update-deps:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/gogo/protobuf/gogoproto
	go get -u github.com/gogo/protobuf/protoc-gen-gogo
	go get -u github.com/lyft/protoc-gen-validate
	go get -u github.com/paulvollmer/2gobytes

.PHONY: pin-repos
pin-repos:
	go run pin_repos.go

.PHONY: check-format
check-format:
	NOT_FORMATTED=$$(gofmt -l ./projects/ ./pkg/ ./test/) && if [ -n "$$NOT_FORMATTED" ]; then echo These files are not formatted: $$NOT_FORMATTED; exit 1; fi

#----------------------------------------------------------------------------------
# Clean
#----------------------------------------------------------------------------------

# Important to clean before pushing new releases. Dockerfiles and binaries may not update properly
.PHONY: clean
clean:
	rm -rf _output
	git clean -xdf install

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
	mkdir -p $(OUTPUT_DIR)
	touch $@


#################
#     Build     #
#################

#----------------------------------------------------------------------------------
# allprojects
#----------------------------------------------------------------------------------
# helper for testing
.PHONY: allprojects
allprojects: apiserver gloo glooctl rate-limit sqoop

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
APISERVER_GQL_SCHEMAS=$(shell find $(APISERVER_DIR)/gql_schemas -name "*.graphql")

.PHONY: apiserver-dependencies
apiserver-dependencies: $(APISERVER_GRAPHQL_GENERATED_FILES)
$(APISERVER_GRAPHQL_GENERATED_FILES): $(APISERVER_GQL_SCHEMAS) $(APISERVER_DIR)/gqlgen.yaml $(APISERVER_SOURCES)
	cd $(APISERVER_DIR) && \
	go run gqlgen.go -v

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
# Observability
#----------------------------------------------------------------------------------

OBSERVABILITY_DIR=projects/observability
OBSERVABILITY_SOURCES=$(shell find $(OBSERVABILITY_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/observability-linux-amd64: $(OBSERVABILITY_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(OBSERVABILITY_DIR)/cmd/main.go

.PHONY: observability
observability: $(OUTPUT_DIR)/observability-linux-amd64

$(OUTPUT_DIR)/Dockerfile.observability: $(OBSERVABILITY_DIR)/cmd/Dockerfile
	cp $< $@

observability-docker: $(OUTPUT_DIR)/observability-linux-amd64 $(OUTPUT_DIR)/Dockerfile.observability
	docker build -t soloio/observability-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.observability

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

gloo-docker-dev: $(OUTPUT_DIR)/gloo-linux-amd64 $(OUTPUT_DIR)/Dockerfile.gloo
	docker build -t soloio/gloo-ee:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.gloo --no-cache

#----------------------------------------------------------------------------------
# Envoy init
#----------------------------------------------------------------------------------

ENVOYINIT_DIR=cmd/envoyinit
ENVOYINIT_SOURCES=$(shell find $(ENVOYINIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/envoyinit-linux-amd64: $(ENVOYINIT_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go

.PHONY: envoyinit
envoyinit: $(OUTPUT_DIR)/envoyinit-linux-amd64


$(OUTPUT_DIR)/Dockerfile.envoyinit: $(ENVOYINIT_DIR)/Dockerfile
	cp $< $@

.PHONY: gloo-ee-envoy-wrapper-docker
gloo-ee-envoy-wrapper-docker: $(OUTPUT_DIR)/envoyinit-linux-amd64 $(OUTPUT_DIR)/Dockerfile.envoyinit
	docker build -t soloio/gloo-ee-envoy-wrapper:$(VERSION)  $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.envoyinit


#----------------------------------------------------------------------------------
# Deployment Manifests / Helm
#----------------------------------------------------------------------------------


HELM_SYNC_DIR := $(OUTPUT_DIR)/helm
HELM_DIR := install/helm

.PHONY: manifest
manifest: helm-template init-helm install/gloo-ee.yaml install/distribution/glooe.yaml update-helm-chart

# creates Chart.yaml, values.yaml, and requirements.yaml
helm-template:
	go run install/helm/gloo-ee/generate.go $(VERSION)

update-helm-chart:
ifeq ($(RELEASE),"true")
	mkdir -p $(HELM_SYNC_DIR)/charts
	helm package --destination $(HELM_SYNC_DIR)/charts $(HELM_DIR)/gloo-ee
	helm repo index $(HELM_SYNC_DIR)
endif

install/gloo-ee.yaml:
	helm template install/helm/gloo-ee --namespace gloo-system --name=gloo-ee > $@

install/distribution/glooe.yaml:
	helm template install/helm/gloo-ee -f install/distribution/values.yaml --namespace gloo-system --name=gloo-ee > $@

init-helm:
	helm repo add helm-hub  https://kubernetes-charts.storage.googleapis.com/
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm dependency update install/helm/gloo-ee

#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------
GH_ORG:=solo-io
GH_REPO:=solo-projects

# For now, expecting people using the release to start from a glooctl CLI we provide, not
# installing the binaries locally / directly. So only uploading the CLI binaries to Github.
# The other binaries can be built manually and used, and docker images for everything will
# be published on release.
RELEASE_BINARIES :=
ifeq ($(RELEASE),"true")
	RELEASE_BINARIES := \
		$(OUTPUT_DIR)/glooctl-linux-amd64 \
		$(OUTPUT_DIR)/glooctl-darwin-amd64
endif

RELEASE_YAMLS :=
ifeq ($(RELEASE),"true")
	RELEASE_YAMLS := \
		install/gloo-ee.yaml
endif

.PHONY: release-binaries
release-binaries: $(RELEASE_BINARIES)


.PHONY: release-yamls
release-yamls: $(RELEASE_YAMLS)

# This is invoked by cloudbuild. When the bot gets a release notification, it kicks of a build with and provides a tag
# variable that gets passed through to here as $TAGGED_VERSION. If no tag is provided, this is a no-op. If a tagged
# version is provided, all the release binaries are uploaded to github.
# Create new releases by clicking "Draft a new release" from https://github.com/solo-io/solo-projects/releases
.PHONY: release
release: release-binaries release-yamls
ifeq ($(RELEASE),"true")
	@$(foreach BINARY,$(RELEASE_BINARIES),ci/upload-github-release-asset.sh owner=solo-io repo=solo-projects tag=$(TAGGED_VERSION) filename=$(BINARY) sha=TRUE;)
	@$(foreach YAML,$(RELEASE_YAMLS),ci/upload-github-release-asset.sh owner=solo-io repo=solo-projects tag=$(TAGGED_VERSION) filename=$(YAML);)
endif

#----------------------------------------------------------------------------------
# Docker
#----------------------------------------------------------------------------------
#
#---------
#--------- Push
#---------

DOCKER_IMAGES :=
ifeq ($(RELEASE),"true")
	DOCKER_IMAGES := docker
endif

.PHONY: docker docker-push
docker: apiserver-docker rate-limit-docker gloo-docker gloo-ee-envoy-wrapper-docker sqoop-docker observability-docker

# Depends on DOCKER_IMAGES, which is set to docker if RELEASE is "true", otherwise empty (making this a no-op).
# This prevents executing the dependent targets if RELEASE is not true, while still enabling `make docker`
# to be used for local testing.
# docker-push is intended to be run by CI
docker-push: $(DOCKER_IMAGES)
ifeq ($(RELEASE),"true")
	docker push soloio/sqoop-ee:$(VERSION) && \
	docker push soloio/rate-limit-ee:$(VERSION) && \
	docker push soloio/apiserver-ee:$(VERSION) && \
	docker push soloio/gloo-ee:$(VERSION) && \
	docker push soloio/gloo-ee-envoy-wrapper:$(VERSION) && \
	docker push soloio/observability-ee:$(VERSION)
endif


#----------------------------------------------------------------------------------
# Distribution
#----------------------------------------------------------------------------------

DISTRIBUTION_DIR=install/distribution

distribution:
ifeq ($(RELEASE),"true")
	go run $(ROOTDIR)/$(DISTRIBUTION_DIR) $(VERSION)
endif
