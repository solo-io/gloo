#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
PACKAGE_PATH:=github.com/solo-io/solo-projects
RELATIVE_OUTPUT_DIR ?= _output
OUTPUT_DIR ?= $(ROOTDIR)/$(RELATIVE_OUTPUT_DIR)
# Kind of a hack to make sure _output exists
z := $(shell mkdir -p $(OUTPUT_DIR))
SOURCES := $(shell find . -name "*.go" | grep -v test.go)
RELEASE := "true"

# If you just put your username, then that refers to your account at hub.docker.com
# To use quay images, set the IMAGE_REPO to "quay.io/solo-io"
# To use dockerhub images, set the IMAGE_REPO to "soloio"
# To use gcr images, set the IMAGE_REPO to "gcr.io/$PROJECT_NAME"
IMAGE_REPO := quay.io/solo-io

# TODO: use $(shell git describe --tags)
ifeq ($(TAGGED_VERSION),)
	TAGGED_VERSION := vdev
	RELEASE := "false"
endif

VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)

ENVOY_GLOO_IMAGE ?= $(IMAGE_REPO)/envoy-gloo-ee:1.15.0-rc1

LDFLAGS := "-X github.com/solo-io/solo-projects/pkg/version.Version=$(VERSION)"
GCFLAGS := 'all=-N -l'

GO_BUILD_FLAGS := GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64

# Passed by cloudbuild
GCLOUD_PROJECT_ID := $(GCLOUD_PROJECT_ID)
BUILD_ID := $(BUILD_ID)

TEST_IMAGE_TAG := test-$(BUILD_ID)
TEST_ASSET_DIR := $(ROOTDIR)/_test
GCR_REPO_PREFIX := gcr.io/$(GCLOUD_PROJECT_ID)


#----------------------------------------------------------------------------------
# Macros
#----------------------------------------------------------------------------------

# If both GCLOUD_PROJECT_ID and BUILD_ID are set, define a function that takes a docker image name
# and returns a '-t' flag that can be passed to 'docker build' to create a tag for a test image.
# If the function is not defined, any attempt at calling if will return nothing (it does not cause en error)
ifneq ($(GCLOUD_PROJECT_ID),)
ifneq ($(BUILD_ID),)
define get_test_tag_option
	-t $(GCR_REPO_PREFIX)/$(1):$(TEST_IMAGE_TAG)
endef
endif
endif

# Same as above, but returns only the tag name withouth the '-t' prefix
ifneq ($(GCLOUD_PROJECT_ID),)
ifneq ($(BUILD_ID),)
define get_test_tag
	$(GCR_REPO_PREFIX)/$(1):$(TEST_IMAGE_TAG)
endef
endif
endif

#----------------------------------------------------------------------------------
# Repo setup
#----------------------------------------------------------------------------------

# https://www.viget.com/articles/two-ways-to-share-git-hooks-with-your-team/
.PHONY: init
init:
	git config core.hooksPath .githooks

.PHONY: update-all-deps
update-all-deps: update-deps update-ui-deps

.PHONY: mod-download
mod-download:
	go mod download

.PHONY: update-deps
update-deps: mod-download
	go get -v golang.org/x/tools/cmd/goimports@v0.0.0-20200423205358-59e73619c742
	go get -v github.com/gogo/protobuf/gogoproto@v1.3.1
	go get -v github.com/gogo/protobuf/protoc-gen-gogo@v1.3.1
	go get -v github.com/solo-io/protoc-gen-ext@v0.0.7
	go get -v github.com/google/wire/cmd/wire@v0.4.0
	go get -v github.com/golang/mock/mockgen@v1.4.3

update-ui-deps:
	yarn --cwd=projects/gloo-ui install

.PHONY: fmt-changed
fmt-changed:
	git diff --name-only | grep '.*.go$$' | xargs goimports -w

.PHONY: check-format
check-format:
	NOT_FORMATTED=$$(gofmt -l ./projects/ ./pkg/ ./test/ ./install/) && if [ -n "$$NOT_FORMATTED" ]; then echo These files are not formatted: $$NOT_FORMATTED; exit 1; fi


#----------------------------------------------------------------------------------
# Clean
#----------------------------------------------------------------------------------

# Important to clean before pushing new releases. Dockerfiles and binaries may not update properly
.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	rm -rf $(TEST_ASSET_DIR)
	rm -rf install/helm/gloo-os-with-ui/templates/
	rm -rf projects/gloo-ui/build
	git clean -xdf install

#----------------------------------------------------------------------------------
# Generated Code
#----------------------------------------------------------------------------------
PROTOC_IMPORT_PATH:=vendor_any

.PHONY: generate-all
generate-all: generated-code generated-ui


SUBDIRS:=projects install pkg test
.PHONY: generated-code
generated-code:
	GO111MODULE=on CGO_ENABLED=0 go generate ./...
	gofmt -w $(SUBDIRS)
	goimports -w $(SUBDIRS)
	go mod tidy

# Flags for all UI code generation
COMMON_UI_PROTOC_FLAGS=--plugin=protoc-gen-ts=projects/gloo-ui/node_modules/.bin/protoc-gen-ts \
		-I$(PROTOC_IMPORT_PATH)/github.com/solo-io \
		-I$(PROTOC_IMPORT_PATH)/github.com/envoyproxy/protoc-gen-validate \
		-I$(PROTOC_IMPORT_PATH)/github.com/gogo/protobuf \
		-I$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext \
		-I$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external \
		-I$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/external \
		--js_out=import_style=commonjs,binary:projects/gloo-ui/src/proto \

# Flags for UI code generation when we do not need to generate GRPC Web service code
UI_TYPES_PROTOC_FLAGS=$(COMMON_UI_PROTOC_FLAGS) \
		--ts_out=projects/gloo-ui/src/proto

# Flags for UI code generation when we need to generate GRPC Web service code
GRPC_WEB_SERVICE_PROTOC_FLAGS=$(COMMON_UI_PROTOC_FLAGS) \
		--ts_out=service=grpc-web:projects/gloo-ui/src/proto

.PHONY: generated-ui
generated-ui:
	rm -rf projects/gloo-ui/src/proto
	mkdir -p projects/gloo-ui/src/proto
	ci/check-protoc.sh
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/gogo/protobuf/gogoproto/gogo.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/envoy/type/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/envoy/api/v2/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/http_uri.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/google/api/annotations.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/google/api/http.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/google/rpc/status.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/envoyproxy/protoc-gen-validate/validate/validate.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext/extproto/ext.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
	 	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/v1/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/external/envoy/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/external/envoy/*/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/external/envoy/*/*/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/external/envoy/*/*/*/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/external/envoy/*/*/*/*/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/core/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/options/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/options/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/options/*/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gateway/api/v1/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/*/*.proto
	protoc $(UI_TYPES_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/*/*/*.proto
	protoc $(GRPC_WEB_SERVICE_PROTOC_FLAGS) \
		$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/*.proto
	ci/fix-gen.sh
	gofmt -w $(SUBDIRS)
	goimports -w $(SUBDIRS)

#################
#     Build     #
#################

#----------------------------------------------------------------------------------
# allprojects
#----------------------------------------------------------------------------------
# helper for testing
.PHONY: allprojects
allprojects: grpcserver gloo extauth rate-limit observability

#----------------------------------------------------------------------------------
# grpcserver
#----------------------------------------------------------------------------------

GRPCSERVER_DIR=projects/grpcserver

$(OUTPUT_DIR)/grpcserver-linux-amd64:
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ projects/grpcserver/server/cmd/main.go

.PHONY: grpcserver
grpcserver: $(OUTPUT_DIR)/grpcserver-linux-amd64

$(OUTPUT_DIR)/Dockerfile.grpcserver: $(GRPCSERVER_DIR)/server/cmd/Dockerfile
	cp $< $@

.PHONY: grpcserver-docker
grpcserver-docker: grpcserver $(OUTPUT_DIR)/Dockerfile.grpcserver $(OUTPUT_DIR)/.grpcserver-docker

$(OUTPUT_DIR)/.grpcserver-docker: $(OUTPUT_DIR)/grpcserver-linux-amd64 $(OUTPUT_DIR)/Dockerfile.grpcserver
	docker build -t $(IMAGE_REPO)/grpcserver-ee:$(VERSION) $(call get_test_tag_option,grpcserver-ee) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.grpcserver
	touch $@

#----------------------------------------------------------------------------------
# grpcserver-envoy
#----------------------------------------------------------------------------------
CONFIG_YAML=envoy_config_grpcserver.yaml

.PHONY: grpcserver-envoy-docker
grpcserver-envoy-docker: $(OUTPUT_DIR)/Dockerfile.grpcserverenvoy
	docker build -t $(IMAGE_REPO)/grpcserver-envoy:$(VERSION) $(call get_test_tag_option,grpcserver-envoy) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.grpcserverenvoy

$(OUTPUT_DIR)/Dockerfile.grpcserverenvoy: $(GRPCSERVER_DIR)/envoy/Dockerfile
	cp $(GRPCSERVER_DIR)/envoy/$(CONFIG_YAML) $(OUTPUT_DIR)/$(CONFIG_YAML)
	cp $< $@

# helpers for local testing

CONFIG_DIR=/etc/config.yaml
GRPC_PORT=10101
GLOO_UI_PORT=20202

.PHONY: run-apiserver
run-apiserver:
	NO_AUTH=1 GRPC_PORT=$(GRPC_PORT) POD_NAMESPACE=gloo-system $(GO_BUILD_FLAGS) go run projects/grpcserver/server/cmd/main.go

.PHONY: run-envoy
run-envoy:
	envoy -c $(GRPCSERVER_DIR)/envoy/$(CONFIG_YAML) -l debug

.PHONY: run-ui
run-ui:
	yarn --cwd projects/gloo-ui install && \
	yarn --cwd projects/gloo-ui start
#----------------------------------------------------------------------------------
# UI
#----------------------------------------------------------------------------------

GRPCSERVER_UI_DIR=projects/gloo-ui
.PHONY: grpcserver-ui-build-local
# TODO rename this so the local build flag is not needed, infer from artifacts
# - should move the yarn build output to the _output dir
grpcserver-ui-build-local:
ifneq ($(LOCAL_BUILD),)
	yarn --cwd $(GRPCSERVER_UI_DIR) install && \
	yarn --cwd $(GRPCSERVER_UI_DIR) build
endif

# If building locally, set LOCAL_BUILD=true
.PHONY: grpcserver-ui-docker
grpcserver-ui-docker: grpcserver-ui-build-local
	docker build -t $(IMAGE_REPO)/grpcserver-ui:$(VERSION) $(call get_test_tag_option,grpcserver-ui) $(GRPCSERVER_UI_DIR) -f $(GRPCSERVER_UI_DIR)/Dockerfile


#----------------------------------------------------------------------------------
# RateLimit
#----------------------------------------------------------------------------------

RATELIMIT_DIR=projects/rate-limit
RATELIMIT_SOURCES=$(shell find $(RATELIMIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/rate-limit-linux-amd64: $(RATELIMIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(RATELIMIT_DIR)/cmd/main.go

.PHONY: rate-limit
rate-limit: $(OUTPUT_DIR)/rate-limit-linux-amd64

$(OUTPUT_DIR)/Dockerfile.rate-limit: $(RATELIMIT_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: rate-limit-docker
rate-limit-docker: $(OUTPUT_DIR)/.rate-limit-docker

$(OUTPUT_DIR)/.rate-limit-docker: $(OUTPUT_DIR)/rate-limit-linux-amd64 $(OUTPUT_DIR)/Dockerfile.rate-limit
	docker build -t $(IMAGE_REPO)/rate-limit-ee:$(VERSION) $(call get_test_tag_option,rate-limit-ee) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.rate-limit
	touch $@

#----------------------------------------------------------------------------------
# ExtAuth
#----------------------------------------------------------------------------------

EXTAUTH_DIR=projects/extauth
EXTAUTH_SOURCES=$(shell find $(EXTAUTH_DIR) -name "*.go" | grep -v test | grep -v generated.go)
EXTAUTH_GO_BUILD_IMAGE=golang:1.14.0-alpine

$(OUTPUT_DIR)/Dockerfile.extauth.build: $(EXTAUTH_DIR)/Dockerfile
	cp $< $@

$(OUTPUT_DIR)/Dockerfile.extauth: $(EXTAUTH_DIR)/cmd/Dockerfile
	cp $< $@

$(OUTPUT_DIR)/.extauth-docker-build: $(EXTAUTH_SOURCES) $(OUTPUT_DIR)/Dockerfile.extauth.build
	docker build -t $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) -f $(OUTPUT_DIR)/Dockerfile.extauth.build \
		--build-arg GO_BUILD_IMAGE=$(EXTAUTH_GO_BUILD_IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg GITHUB_TOKEN \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(OUTPUT_DIR)/extauth-linux-amd64: $(OUTPUT_DIR)/.extauth-docker-build
	docker create -ti --name extauth-temp-container $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) bash
	docker cp extauth-temp-container:/extauth-linux-amd64 $(OUTPUT_DIR)/extauth-linux-amd64
	docker rm -f extauth-temp-container

# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(OUTPUT_DIR)/verify-plugins-linux-amd64: $(OUTPUT_DIR)/.extauth-docker-build
	docker create -ti --name verify-plugins-temp-container $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) bash
	docker cp verify-plugins-temp-container:/verify-plugins-linux-amd64 $(OUTPUT_DIR)/verify-plugins-linux-amd64
	docker rm -f verify-plugins-temp-container

# Build extauth binaries
.PHONY: extauth
extauth: $(OUTPUT_DIR)/extauth-linux-amd64 $(OUTPUT_DIR)/verify-plugins-linux-amd64

# Build ext-auth-plugins docker image
.PHONY: auth-plugins
auth-plugins: $(OUTPUT_DIR)/verify-plugins-linux-amd64
	docker build --no-cache -t $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) -f projects/extauth/plugins/Dockerfile \
		--build-arg GO_BUILD_IMAGE=$(EXTAUTH_GO_BUILD_IMAGE) \
		--build-arg GC_FLAGS=$(GCFLAGS) \
		--build-arg VERIFY_SCRIPT=$(RELATIVE_OUTPUT_DIR)/verify-plugins-linux-amd64 \
		--build-arg GITHUB_TOKEN \
		.

# Build extauth server docker image
.PHONY: extauth-docker
extauth-docker: $(OUTPUT_DIR)/.extauth-docker

$(OUTPUT_DIR)/.extauth-docker: $(OUTPUT_DIR)/extauth-linux-amd64 $(OUTPUT_DIR)/verify-plugins-linux-amd64 $(OUTPUT_DIR)/Dockerfile.extauth
	docker build -t $(IMAGE_REPO)/extauth-ee:$(VERSION) $(call get_test_tag_option,extauth-ee) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.extauth
	touch $@


#----------------------------------------------------------------------------------
# Observability
#----------------------------------------------------------------------------------

OBSERVABILITY_DIR=projects/observability
OBSERVABILITY_SOURCES=$(shell find $(OBSERVABILITY_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/observability-linux-amd64: $(OBSERVABILITY_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(OBSERVABILITY_DIR)/cmd/main.go

.PHONY: observability
observability: $(OUTPUT_DIR)/observability-linux-amd64

$(OUTPUT_DIR)/Dockerfile.observability: $(OBSERVABILITY_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: observability-docker
observability-docker: $(OUTPUT_DIR)/.observability-docker

$(OUTPUT_DIR)/.observability-docker: $(OUTPUT_DIR)/observability-linux-amd64 $(OUTPUT_DIR)/Dockerfile.observability
	docker build -t $(IMAGE_REPO)/observability-ee:$(VERSION) $(call get_test_tag_option,observability-ee) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.observability
	touch $@

#----------------------------------------------------------------------------------
# Gloo
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gloo-linux-amd64: $(GLOO_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_DIR)/cmd/main.go


.PHONY: gloo
gloo: $(OUTPUT_DIR)/gloo-linux-amd64

$(OUTPUT_DIR)/Dockerfile.gloo: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@


.PHONY: gloo-docker
gloo-docker: $(OUTPUT_DIR)/.gloo-docker

$(OUTPUT_DIR)/.gloo-docker: $(OUTPUT_DIR)/gloo-linux-amd64 $(OUTPUT_DIR)/Dockerfile.gloo
	docker build $(call get_test_tag_option,gloo-ee) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.gloo \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) \
		-t $(IMAGE_REPO)/gloo-ee:$(VERSION)
	touch $@

gloo-docker-dev: $(OUTPUT_DIR)/gloo-linux-amd64 $(OUTPUT_DIR)/Dockerfile.gloo
	docker build -t $(IMAGE_REPO)/gloo-ee:$(VERSION) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.gloo --no-cache
	touch $@

#----------------------------------------------------------------------------------
# glooctl
#----------------------------------------------------------------------------------

CLI_DIR=projects/gloo/cli

$(OUTPUT_DIR)/glooctl: $(SOURCES)
	GO111MODULE=on go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

$(OUTPUT_DIR)/glooctl-linux-amd64: $(SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

$(OUTPUT_DIR)/glooctl-darwin-amd64: $(SOURCES)
	$(GO_BUILD_FLAGS) GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

$(OUTPUT_DIR)/glooctl-windows-amd64.exe: $(SOURCES)
	$(GO_BUILD_FLAGS) GOOS=windows go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

.PHONY: glooctl
glooctl: $(OUTPUT_DIR)/glooctl
.PHONY: glooctl-linux-amd64
glooctl-linux-amd64: $(OUTPUT_DIR)/glooctl-linux-amd64
.PHONY: glooctl-darwin-amd64
glooctl-darwin-amd64: $(OUTPUT_DIR)/glooctl-darwin-amd64
.PHONY: glooctl-windows-amd64
glooctl-windows-amd64: $(OUTPUT_DIR)/glooctl-windows-amd64.exe

.PHONY: build-cli
build-cli: glooctl-linux-amd64 glooctl-darwin-amd64 glooctl-windows-amd64

#----------------------------------------------------------------------------------
# Envoy init
#----------------------------------------------------------------------------------

ENVOYINIT_DIR=cmd/envoyinit
ENVOYINIT_SOURCES=$(shell find $(ENVOYINIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/envoyinit-linux-amd64: $(ENVOYINIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go

.PHONY: envoyinit
envoyinit: $(OUTPUT_DIR)/envoyinit-linux-amd64

$(OUTPUT_DIR)/Dockerfile.envoyinit: $(ENVOYINIT_DIR)/Dockerfile
	cp $< $@

$(OUTPUT_DIR)/docker-entrypoint.sh: $(ENVOYINIT_DIR)/docker-entrypoint.sh
	cp $< $@

.PHONY: gloo-ee-envoy-wrapper-docker
gloo-ee-envoy-wrapper-docker: $(OUTPUT_DIR)/.gloo-ee-envoy-wrapper-docker

$(OUTPUT_DIR)/.gloo-ee-envoy-wrapper-docker: $(OUTPUT_DIR)/envoyinit-linux-amd64 $(OUTPUT_DIR)/Dockerfile.envoyinit $(OUTPUT_DIR)/docker-entrypoint.sh
	docker build $(call get_test_tag_option,gloo-ee-envoy-wrapper) $(OUTPUT_DIR) -f $(OUTPUT_DIR)/Dockerfile.envoyinit \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) \
		-t $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION)
	touch $@


#----------------------------------------------------------------------------------
# Deployment Manifests / Helm
#----------------------------------------------------------------------------------
HELM_SYNC_DIR_FOR_GLOO_EE := $(OUTPUT_DIR)/helm
HELM_SYNC_DIR_RO_UI_GLOO := $(OUTPUT_DIR)/helm_gloo_os_ui
HELM_DIR := install/helm
MANIFEST_DIR := install/manifest
MANIFEST_FOR_RO_UI_GLOO := gloo-with-read-only-ui-release.yaml
MANIFEST_FOR_GLOO_EE := glooe-release.yaml
GLOOE_HELM_BUCKET := gs://gloo-ee-helm
GLOO_OS_UI_HELM_BUCKET := gs://gloo-os-ui-helm

.PHONY: manifest
manifest: helm-template init-helm produce-manifests

# creates Chart.yaml, values.yaml, and requirements.yaml
.PHONY: helm-template
helm-template:
	mkdir -p $(MANIFEST_DIR)
	mkdir -p $(HELM_SYNC_DIR_FOR_GLOO_EE)
	mkdir -p $(HELM_SYNC_DIR_RO_UI_GLOO)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION)

.PHONY: init-helm
init-helm: helm-template $(OUTPUT_DIR)/.helm-initialized

$(OUTPUT_DIR)/.helm-initialized:
	helm repo add helm-hub https://kubernetes-charts.storage.googleapis.com/
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm repo add dev-portal https://storage.googleapis.com/dev-portal-helm
	helm dependency update install/helm/gloo-ee
	# see install/helm/gloo-os-with-ui/README.md
	mkdir -p install/helm/gloo-os-with-ui/templates
	cp install/helm/gloo-ee/templates/_helpers.tpl install/helm/gloo-os-with-ui/templates/_helpers.tpl
	cp install/helm/gloo-ee/templates/*-apiserver-*.yaml install/helm/gloo-os-with-ui/templates/
	cp install/helm/gloo-ee/templates/40-settings.yaml install/helm/gloo-os-with-ui/templates/40-settings.yaml
	helm dependency update install/helm/gloo-os-with-ui
	touch $@

.PHONY: produce-manifests
produce-manifests: init-helm
	helm template glooe install/helm/gloo-ee --namespace gloo-system > $(MANIFEST_DIR)/$(MANIFEST_FOR_GLOO_EE)
	helm template gloo install/helm/gloo-os-with-ui --namespace gloo-system > $(MANIFEST_DIR)/$(MANIFEST_FOR_RO_UI_GLOO)

.PHONY: fetch-package-and-save-helm
fetch-package-and-save-helm: init-helm
ifeq ($(RELEASE),"true")
	until $$(GENERATION=$$(gsutil ls -a $(GLOOE_HELM_BUCKET)/index.yaml | tail -1 | cut -f2 -d '#') && \
					gsutil cp -v $(GLOOE_HELM_BUCKET)/index.yaml $(HELM_SYNC_DIR_FOR_GLOO_EE)/index.yaml && \
					helm package --destination $(HELM_SYNC_DIR_FOR_GLOO_EE)/charts $(HELM_DIR)/gloo-ee >> /dev/null && \
					helm repo index $(HELM_SYNC_DIR_FOR_GLOO_EE) --merge $(HELM_SYNC_DIR_FOR_GLOO_EE)/index.yaml && \
					gsutil -m rsync $(HELM_SYNC_DIR_FOR_GLOO_EE)/charts $(GLOOE_HELM_BUCKET)/charts && \
					gsutil -h x-goog-if-generation-match:"$$GENERATION" cp $(HELM_SYNC_DIR_FOR_GLOO_EE)/index.yaml $(GLOOE_HELM_BUCKET)/index.yaml); do \
		echo "Failed to upload new helm index (updated helm index since last download?). Trying again"; \
		sleep 2; \
	done
	until $$(GENERATION=$$(gsutil ls -a $(GLOO_OS_UI_HELM_BUCKET)/index.yaml | tail -1 | cut -f2 -d '#') && \
					gsutil cp -v $(GLOO_OS_UI_HELM_BUCKET)/index.yaml $(HELM_SYNC_DIR_RO_UI_GLOO)/index.yaml && \
					helm package --destination $(HELM_SYNC_DIR_RO_UI_GLOO)/charts $(HELM_DIR)/gloo-os-with-ui >> /dev/null && \
					helm repo index $(HELM_SYNC_DIR_RO_UI_GLOO) --merge $(HELM_SYNC_DIR_RO_UI_GLOO)/index.yaml && \
					gsutil -m rsync $(HELM_SYNC_DIR_RO_UI_GLOO)/charts $(GLOO_OS_UI_HELM_BUCKET)/charts && \
					gsutil -h x-goog-if-generation-match:"$$GENERATION" cp $(HELM_SYNC_DIR_RO_UI_GLOO)/index.yaml $(GLOO_OS_UI_HELM_BUCKET)/index.yaml); do \
		echo "Failed to upload new helm index (updated helm index since last download?). Trying again"; \
		sleep 2; \
	done
endif

#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------

CHANGELOGS_BUCKET=gloo-ee-changelogs

.PHONY: upload-changelog
upload-changelog:
ifeq ($(RELEASE),"true")
	gsutil -m cp -r './changelog' gs://$(CHANGELOGS_BUCKET)/$(VERSION)
endif

.PHONY: upload-github-release-assets
upload-github-release-assets: produce-manifests
	$(GO_BUILD_FLAGS) go run ci/upload_github_release_assets.go

DEPS_DIR=$(OUTPUT_DIR)/dependencies/$(VERSION)
DEPS_BUCKET=gloo-ee-dependencies

.PHONY: publish-dependencies
publish-dependencies: $(DEPS_DIR)/go.mod $(DEPS_DIR)/go.sum $(DEPS_DIR)/dependencies $(DEPS_DIR)/dependencies.json \
	$(DEPS_DIR)/build_env $(DEPS_DIR)/verify-plugins-linux-amd64
	gsutil cp -r $(DEPS_DIR) gs://$(DEPS_BUCKET)

$(DEPS_DIR):
	mkdir -p $(DEPS_DIR)

$(DEPS_DIR)/dependencies: $(DEPS_DIR) go.mod
	GO111MODULE=on go list -m all > $@

$(DEPS_DIR)/dependencies.json: $(DEPS_DIR) go.mod
	GO111MODULE=on go list -m --json all > $@

$(DEPS_DIR)/go.mod: $(DEPS_DIR) go.mod
	cp go.mod $(DEPS_DIR)

$(DEPS_DIR)/go.sum: $(DEPS_DIR) go.sum
	cp go.sum $(DEPS_DIR)

$(DEPS_DIR)/build_env: $(DEPS_DIR)
	echo "GO_BUILD_IMAGE=$(EXTAUTH_GO_BUILD_IMAGE)" > $@
	echo "GC_FLAGS=$(GCFLAGS)" >> $@

$(DEPS_DIR)/verify-plugins-linux-amd64: $(OUTPUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)
	cp $(OUTPUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)

#----------------------------------------------------------------------------------
# Docker push
#----------------------------------------------------------------------------------

DOCKER_IMAGES :=
ifeq ($(RELEASE),"true")
	DOCKER_IMAGES := docker
endif

.PHONY: docker docker-push
docker: grpcserver-ui-docker grpcserver-envoy-docker grpcserver-docker rate-limit-docker extauth-docker gloo-docker \
	gloo-ee-envoy-wrapper-docker observability-docker auth-plugins

# Depends on DOCKER_IMAGES, which is set to docker if RELEASE is "true", otherwise empty (making this a no-op).
# This prevents executing the dependent targets if RELEASE is not true, while still enabling `make docker`
# to be used for local testing.
# docker-push is intended to be run by CI
docker-push: $(DOCKER_IMAGES)
ifeq ($(RELEASE),"true")
	docker push $(IMAGE_REPO)/rate-limit-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/grpcserver-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/grpcserver-envoy:$(VERSION) && \
	docker push $(IMAGE_REPO)/grpcserver-ui:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) && \
	docker push $(IMAGE_REPO)/observability-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/extauth-ee:$(VERSION)
	docker push $(IMAGE_REPO)/ext-auth-plugins:$(VERSION)
endif

push-kind-images: docker
	kind load docker-image $(IMAGE_REPO)/rate-limit-ee:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/grpcserver-ee:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/grpcserver-envoy:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/grpcserver-ui:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/gloo-ee:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/observability-ee:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/extauth-ee:$(VERSION) --name $(CLUSTER_NAME)
	kind load docker-image $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) --name $(CLUSTER_NAME)

#----------------------------------------------------------------------------------
# Build assets for regression tests
#----------------------------------------------------------------------------------
#
# The following targets are used to generate the assets on which the regression tests rely upon. The following actions are performed:
#
#   1. Push the images to GCR (images have been tagged as $(GCR_REPO_PREFIX)/<image-name>:$(TEST_IMAGE_TAG)
#   2. Generate GlooE value files providing overrides to make the image elements point to GCR
#      - override the repository prefix for all repository names (e.g. quay.io/solo-io/gateway -> gcr.io/gloo-ee/gateway)
#      - set the tag for each image to TEST_IMAGE_TAG
#   3. Package the Gloo Helm chart to the _test directory (also generate an index file)
#
# The regression tests will use the generated Gloo Chart to install Gloo to the GKE test cluster.

.PHONY: build-test-assets
build-test-assets: docker push-test-images build-test-chart build-os-with-ui-test-chart

.PHONY: build-kind-assets
build-kind-assets: push-kind-images build-kind-chart

TEST_DOCKER_TARGETS := grpcserver-ui-docker-test grpcserver-envoy-docker-test grpcserver-docker-test rate-limit-docker-test extauth-docker-test observability-docker-test gloo-docker-test gloo-ee-envoy-wrapper-docker-test

.PHONY: push-test-images $(TEST_DOCKER_TARGETS)
push-test-images: $(TEST_DOCKER_TARGETS)

grpcserver-docker-test: $(OUTPUT_DIR)/grpcserver-linux-amd64 $(OUTPUT_DIR)/.grpcserver-docker
	docker push $(call get_test_tag,grpcserver-ee)

grpcserver-envoy-docker-test: grpcserver-envoy-docker $(OUTPUT_DIR)/Dockerfile.grpcserverenvoy
	docker push $(call get_test_tag,grpcserver-envoy)

grpcserver-ui-docker-test: grpcserver-ui-build-local grpcserver-ui-docker
	docker push $(call get_test_tag,grpcserver-ui)

rate-limit-docker-test: $(OUTPUT_DIR)/rate-limit-linux-amd64 $(OUTPUT_DIR)/Dockerfile.rate-limit
	docker push $(call get_test_tag,rate-limit-ee)

extauth-docker-test: $(OUTPUT_DIR)/extauth-linux-amd64 $(OUTPUT_DIR)/Dockerfile.extauth
	docker push $(call get_test_tag,extauth-ee)

observability-docker-test: $(OUTPUT_DIR)/observability-linux-amd64 $(OUTPUT_DIR)/Dockerfile.observability
	docker push $(call get_test_tag,observability-ee)

gloo-docker-test: gloo-docker
	docker push $(call get_test_tag,gloo-ee)

gloo-ee-envoy-wrapper-docker-test: $(OUTPUT_DIR)/envoyinit-linux-amd64 $(OUTPUT_DIR)/Dockerfile.envoyinit gloo-ee-envoy-wrapper-docker
	docker push $(call get_test_tag,gloo-ee-envoy-wrapper)

.PHONY: build-test-chart
build-test-chart:
	mkdir -p $(TEST_ASSET_DIR)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(TEST_IMAGE_TAG) $(GCR_REPO_PREFIX)
	helm repo add helm-hub https://kubernetes-charts.storage.googleapis.com/
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm dependency update install/helm/gloo-ee
	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-ee
	helm repo index $(TEST_ASSET_DIR)

.PHONY: build-os-with-ui-test-chart
build-os-with-ui-test-chart: init-helm
	mkdir -p $(TEST_ASSET_DIR)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(TEST_IMAGE_TAG) $(GCR_REPO_PREFIX)
	helm repo add helm-hub https://kubernetes-charts.storage.googleapis.com/
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm dependency update install/helm/gloo-os-with-ui
	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-os-with-ui
	helm repo index $(TEST_ASSET_DIR)

.PHONY: build-kind-chart
build-kind-chart:
	mkdir -p $(TEST_ASSET_DIR)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION)
	helm repo add helm-hub https://kubernetes-charts.storage.googleapis.com/
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm dependency update install/helm/gloo-ee
	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-ee
	helm repo index $(TEST_ASSET_DIR)


#----------------------------------------------------------------------------------
# Printing makefile variables utility
#----------------------------------------------------------------------------------

# use `make print-MAKEFILE_VAR` to print the value of MAKEFILE_VAR

print-%  : ; @echo $($*)
