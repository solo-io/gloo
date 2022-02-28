include Makefile.docker

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

GCS_BUCKET := glooctl-plugins
WASM_GCS_PATH := glooctl-wasm
FED_GCS_PATH := glooctl-fed

# If you just put your username, then that refers to your account at hub.docker.com
# To use quay images, set the IMAGE_REPO to "quay.io/solo-io" (or leave unset)
# To use dockerhub images, set the IMAGE_REPO to "soloio"
# To use gcr images, set the IMAGE_REPO to "gcr.io/$PROJECT_NAME"
IMAGE_REPO ?= quay.io/solo-io

ifeq ($(TAGGED_VERSION),)
	TAGGED_VERSION := "v$(shell ./git-semver.sh)"
	RELEASE := "false"
endif

VERSION ?= $(shell echo $(TAGGED_VERSION) | sed -e "s/^refs\/tags\///" | cut -c 2-)

ENVOY_GLOO_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee:1.20.0-patch5
ENVOY_GLOO_FIPS_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee-fips:1.20.0-patch5

# The full SHA of the currently checked out commit
CHECKED_OUT_SHA := $(shell git rev-parse HEAD)
# Returns the name of the default branch in the remote `origin` repository, e.g. `master`
DEFAULT_BRANCH_NAME := $(shell git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@')
# Print the branches that contain the current commit and keep only the one that
# EXACTLY matches the name of the default branch (avoid matching e.g. `master-2`).
# If we get back a result, it mean we are on the default branch.
EMPTY_IF_NOT_DEFAULT := $(shell git branch --contains $(CHECKED_OUT_SHA) | grep -ow $(DEFAULT_BRANCH_NAME))

ON_DEFAULT_BRANCH := "false"
ifneq ($(EMPTY_IF_NOT_DEFAULT),)
    ON_DEFAULT_BRANCH = "true"
endif

LDFLAGS := "-X github.com/solo-io/solo-projects/pkg/version.Version=$(VERSION) -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn"
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

# If both GCLOUD_PROJECT_ID and BUILD_ID are set, define functions that take a docker image name
# and returns a tag name either with or without (_option) the '-t' flag that can be passed to 'docker build'
# to create a tag for a test image. If the function is not defined, any attempt at calling if will
# return nothing (it does not cause en error)
ifneq ($(GCLOUD_PROJECT_ID),)
ifneq ($(BUILD_ID),)
define get_test_tag_option
	-t $(GCR_REPO_PREFIX)/$(1):$(TEST_IMAGE_TAG)
endef
# Same as above, but returns only the tag name without the '-t' prefix
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
update-all-deps: install-go-tools update-ui-deps

DEPSGOBIN=$(ROOTDIR)/.bin

# https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
.PHONY: install-go-tools
install-go-tools: mod-download
	mkdir -p $(DEPSGOBIN)
	GOBIN=$(DEPSGOBIN) go install istio.io/tools/cmd/protoc-gen-jsonshim
	GOBIN=$(DEPSGOBIN) go install istio.io/pkg/version
	GOBIN=$(DEPSGOBIN) go install github.com/solo-io/protoc-gen-ext
	GOBIN=$(DEPSGOBIN) go install golang.org/x/tools/cmd/goimports
	GOBIN=$(DEPSGOBIN) go install github.com/envoyproxy/protoc-gen-validate
	GOBIN=$(DEPSGOBIN) go install github.com/golang/protobuf/protoc-gen-go
	GOBIN=$(DEPSGOBIN) go install github.com/golang/mock/gomock
	GOBIN=$(DEPSGOBIN) go install github.com/golang/mock/mockgen
	GOBIN=$(DEPSGOBIN) go install github.com/google/wire/cmd/wire
	GOBIN=$(DEPSGOBIN) go install github.com/onsi/ginkgo/ginkgo
	GOBIN=$(DEPSGOBIN) go install github.com/solo-io/protoc-gen-openapi

.PHONY: mod-download
mod-download:
	go mod download

.PHONY: clean-artifacts
clean-artifacts:
	rm -rf _output

.PHONY: clean-generated-protos
clean-generated-protos:
	rm -rf $(ROOTDIR)/projects/apiserver/api/fed.rpc/v1/*resources.proto
	rm -rf $(ROOTDIR)/projects/apiserver/api/rpc.edge.gloo/v1/*resources.proto

# Clean
.PHONY: clean-fed
clean-fed: clean-artifacts clean-generated-protos
	rm -rf $(ROOTDIR)/vendor_any
	rm -rf $(ROOTDIR)/projects/gloo-fed/pkg/api
	rm -rf $(ROOTDIR)/projects/apiserver/pkg/api
	rm -rf $(ROOTDIR)/projects/glooctl-plugins/fed/pkg/api
	rm -rf $(ROOTDIR)/projects/apiserver/server/services/single_cluster_resource_handler/*

# command to run regression tests with guaranteed access to $(DEPSGOBIN)/ginkgo
# requires the environment variable KUBE2E_TESTS to be set to the test type you wish to run
.PHONY: run-ci-regression-tests
run-ci-regression-tests: install-go-tools
	go env -w GOPRIVATE=github.com/solo-io
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn $(DEPSGOBIN)/ginkgo -r -failFast -trace -progress -race -compilers=4 -failOnPending -noColor ./test/regressions/$(KUBE2E_TESTS)/...

.PHONE: run-ci-gloo-fed-regression-tests
run-ci-gloo-fed-regression-tests: install-go-tools
	go env -w GOPRIVATE=github.com/solo-io
	REMOTE_CLUSTER_CONTEXT=kind-remote LOCAL_CLUSTER_CONTEXT=kind-local $(DEPSGOBIN)/ginkgo -r ./test/gloo-fed-e2e/...

# command to run e2e tests
# requires the environment variable ENVOY_IMAGE_TAG to be set to the tag of the gloo-ee-envoy-wrapper Docker image you wish to run
.PHONY: run-e2e-tests
run-e2e-tests: install-go-tools
	ginkgo -r -failFast -trace -progress -race -compilers=4 -failOnPending ./test/e2e/

.PHONY: update-ui-deps
update-ui-deps:
	yarn --cwd=$(APISERVER_UI_DIR) install

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
	rm -rf $(APISERVER_UI_DIR)/build
	rm -rf $(ROOTDIR)/vendor_any
	git clean -xdf install

#----------------------------------------------------------------------------------
# Generated Code
#----------------------------------------------------------------------------------
PROTOC_IMPORT_PATH:=$(ROOTDIR)/vendor_any

.PHONY: generate-all
generate-all: check-solo-apis generated-code generate-gloo-fed generate-helm-docs

GLOO_VERSION=$(shell echo $(shell go list -m github.com/solo-io/gloo) | cut -d' ' -f2)

.PHONY: check-solo-apis
check-solo-apis:
ifeq ($(GLOO_BRANCH_BUILD),)
	# Ensure that the gloo and solo-apis dependencies are in lockstep
	go get github.com/solo-io/solo-apis@gloo-$(GLOO_VERSION)
endif

SUBDIRS:=projects install pkg test
.PHONY: generated-code
generated-code: update-licenses
	rm -rf $(ROOTDIR)/vendor_any
	go mod tidy
	PATH=$(DEPSGOBIN):$$PATH GO111MODULE=on CGO_ENABLED=1 go generate ./...
	PATH=$(DEPSGOBIN):$$PATH goimports -w $(SUBDIRS)
	PATH=$(DEPSGOBIN):$$PATH go mod tidy
	ci/check-protoc.sh

.PHONY: generate-gloo-fed
generate-gloo-fed: generate-gloo-fed-code generated-gloo-fed-ui

# Generated Code - Required to update Codgen Templates
.PHONY: generate-gloo-fed-code
generate-gloo-fed-code: clean-fed
	PATH=$(DEPSGOBIN):$$PATH go run $(ROOTDIR)/projects/gloo-fed/generate.go # Generates clients, controllers, etc
	PATH=$(DEPSGOBIN):$$PATH $(ROOTDIR)/projects/gloo-fed/ci/hack-fix-marshal.sh # TODO: figure out a more permanent way to deal with this
	PATH=$(DEPSGOBIN):$$PATH go run projects/gloo-fed/generate.go -apiserver # Generates apiserver protos into go code
	PATH=$(DEPSGOBIN):$$PATH go generate $(ROOTDIR)/projects/... # Generates mocks
	PATH=$(DEPSGOBIN):$$PATH goimports -w $(SUBDIRS)
	PATH=$(DEPSGOBIN):$$PATH go mod tidy

.PHONY: generate-helm-docs
generate-helm-docs:
	PATH=$(DEPSGOBIN):$$PATH go run $(ROOTDIR)/install/helm/gloo-ee/generate.go $(VERSION) --generate-helm-docs # Generate Helm Documentation

#################
#     Build     #
#################

#----------------------------------------------------------------------------------
# allprojects
#----------------------------------------------------------------------------------
# helper for testing
.PHONY: allprojects
allprojects: gloo-fed-apiserver gloo extauth extauth-fips rate-limit rate-limit-fips observability

#----------------------------------------------------------------------------------
# Gloo Fed
#----------------------------------------------------------------------------------

GLOO_FED_DIR=$(ROOTDIR)/projects/gloo-fed
GLOO_FED_SOURCES=$(shell find $(GLOO_FED_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gloo-fed-linux-amd64: $(GLOO_FED_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_DIR)/cmd/main.go

.PHONY: gloo-fed
gloo-fed: $(OUTPUT_DIR)/gloo-fed-linux-amd64

.PHONY: gloo-fed-docker
gloo-fed-docker: $(OUTPUT_DIR)/gloo-fed-linux-amd64
	docker build -t $(IMAGE_REPO)/gloo-fed:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed
kind-load-gloo-fed: gloo-fed-docker
	kind load docker-image --name $(CLUSTER_NAME) $(IMAGE_REPO)/gloo-fed:$(VERSION)

#----------------------------------------------------------------------------------
# Gloo Federation Projects
#----------------------------------------------------------------------------------

.PHONY: gloofed-docker
gloofed-docker: gloo-fed-docker gloo-fed-rbac-validating-webhook-docker gloo-fed-apiserver-docker gloo-fed-apiserver-envoy-docker gloo-federation-console-docker

.PHONY: gloofed-load-kind-images
gloofed-load-kind-images: kind-load-gloo-fed kind-load-gloo-fed-rbac-validating-webhook kind-load-gloo-fed-apiserver kind-load-gloo-fed-apiserver-envoy kind-load-ui

#----------------------------------------------------------------------------------
# Gloo Fed Apiserver
#----------------------------------------------------------------------------------
GLOO_FED_APISERVER_DIR=$(ROOTDIR)/projects/apiserver
APISERVER_SOURCES=$(shell find $(APISERVER_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gloo-fed-apiserver-linux-amd64: $(APISERVER_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_APISERVER_DIR)/cmd/main.go

.PHONY: gloo-fed-apiserver
gloo-fed-apiserver: $(OUTPUT_DIR)/gloo-fed-apiserver-linux-amd64

.PHONY: gloo-fed-apiserver-docker
gloo-fed-apiserver-docker: $(OUTPUT_DIR)/gloo-fed-apiserver-linux-amd64
	docker build -t $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_APISERVER_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed-apiserver
kind-load-gloo-fed-apiserver: gloo-fed-apiserver-docker
	kind load docker-image --name $(CLUSTER_NAME) $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION)

#----------------------------------------------------------------------------------
# apiserver-envoy
#----------------------------------------------------------------------------------
CONFIG_YAML=cfg.yaml

GLOO_FED_APISERVER_ENVOY_DIR=$(ROOTDIR)/projects/apiserver/apiserver-envoy

.PHONY: gloo-fed-apiserver-envoy-docker
gloo-fed-apiserver-envoy-docker:
	cp $(GLOO_FED_APISERVER_ENVOY_DIR)/$(CONFIG_YAML) $(OUTPUT_DIR)/$(CONFIG_YAML)
	docker build -t $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_APISERVER_ENVOY_DIR)/Dockerfile;

.PHONY: kind-load-gloo-fed-apiserver-envoy
kind-load-gloo-fed-apiserver-envoy: gloo-fed-apiserver-envoy-docker
	kind load docker-image --name $(CLUSTER_NAME) $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION)

#----------------------------------------------------------------------------------
# helpers for local testing
#----------------------------------------------------------------------------------
GRPC_PORT=10101
CONFIG_YAML=cfg.yaml

.PHONY: run-apiserver
run-apiserver:
	GRPC_PORT=$(GRPC_PORT) POD_NAMESPACE=gloo-system $(GO_BUILD_FLAGS) go run projects/apiserver/cmd/main.go

.PHONY: run-envoy
run-envoy:
	envoy -c projects/apiserver/apiserver-envoy/$(CONFIG_YAML) -l debug

.PHONY: run-fed-ui
run-fed-ui:
	./hack/check-gloo-fed.sh && \
 	yarn --cwd $(APISERVER_UI_DIR) install && \
	yarn --cwd $(APISERVER_UI_DIR) start

.PHONY: run-single-cluster-ui
run-single-cluster-ui:
	yarn --cwd $(APISERVER_UI_DIR) install && \
	yarn --cwd $(APISERVER_UI_DIR) start

#----------------------------------------------------------------------------------
# Gloo Fed Rbac Webhook
#----------------------------------------------------------------------------------
GLOO_FED_RBAC_WEBHOOK_DIR=$(ROOTDIR)/projects/rbac-validating-webhook
GLOO_FED_RBAC_WEBHOOK_SOURCES=$(shell find $(GLOO_FED_RBAC_WEBHOOK_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-amd64: $(GLOO_FED_RBAC_WEBHOOK_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/main.go

.PHONY: gloo-fed-rbac-validating-webhook
gloo-fed-rbac-validating-webhook: $(OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-amd64

.PHONY: gloo-fed-rbac-validating-webhook-docker
gloo-fed-rbac-validating-webhook-docker: $(OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-amd64
	docker build -t $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed-rbac-validating-webhook
kind-load-gloo-fed-rbac-validating-webhook: gloo-fed-rbac-validating-webhook-docker
	kind load docker-image --name $(CLUSTER_NAME) $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION)

#----------------------------------------------------------------------------------
# ApiServer gRPC Code Generation
#----------------------------------------------------------------------------------

# proto sources
APISERVER_DIR=$(ROOTDIR)/projects/apiserver/api

COMMON_PROTOC_FLAGS=-I$(PROTOC_IMPORT_PATH)/github.com/envoyproxy/protoc-gen-validate \
	-I$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext \
	-I$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext/external \
	-I$(PROTOC_IMPORT_PATH)/ \
	-I$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/external \
	-I$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external \

GENERATED_TS_DIR=$(APISERVER_UI_DIR)/src/proto

TS_OUT=--plugin=protoc-gen-ts=$(APISERVER_UI_DIR)/node_modules/.bin/protoc-gen-ts \
			--ts_out=service=grpc-web:$(GENERATED_TS_DIR) \
			--js_out=import_style=commonjs,binary:$(GENERATED_TS_DIR)

PROTOC=protoc $(COMMON_PROTOC_FLAGS)

.PHONY: generated-gloo-fed-ui
generated-gloo-fed-ui: update-gloo-fed-ui-deps generated-gloo-fed-ui-deps
	mkdir -p $(APISERVER_UI_DIR)/pkg/api/fed.rpc/v1
	mkdir -p $(APISERVER_UI_DIR)/pkg/api/rpc.edge.gloo/v1
	./ci/fix-ui-gen.sh

.PHONY: generated-gloo-fed-ui-deps
generated-gloo-fed-ui-deps:
	rm -rf $(GENERATED_TS_DIR)
	mkdir -p $(GENERATED_TS_DIR)

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext/extproto/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/envoy/type/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/envoy/api/v2/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/google/api/annotations.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/google/api/http.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/external/google/rpc/status.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/envoyproxy/protoc-gen-validate/validate/validate.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext/extproto/ext.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	 $(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	 $(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/skv2/api/core/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	 $(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/skv2/api/multicluster/v1alpha1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/*/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/*/*/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/*/*/*/*/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/external/udpa/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/*/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gateway/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gateway/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo//gloo/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.enterprise.gloo/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.ratelimit/v1alpha1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/apiserver/api/*/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/skv2-enterprise/multicluster-admission-webhook/api/multicluster/v1alpha1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/rate-limiter/*/*.proto

#----------------------------------------------------------------------------------
# UI
#----------------------------------------------------------------------------------

APISERVER_UI_DIR=projects/ui
GLOO_UI_OUT_DIR=$(OUTPUT_DIR)/ui

.PHONY: update-gloo-fed-ui-deps
update-gloo-fed-ui-deps:
# TODO rename this so the local build flag is not needed, infer from artifacts
ifneq ($(LOCAL_BUILD),)
	yarn --cwd $(APISERVER_UI_DIR) install
endif

.PHONY: build-ui
build-ui: update-gloo-fed-ui-deps
ifneq ($(LOCAL_BUILD),)
	yarn --cwd $(APISERVER_UI_DIR) build
endif

.PHONY: gloo-federation-console-docker
gloo-federation-console-docker: build-ui
	docker build -t $(IMAGE_REPO)/gloo-federation-console:$(VERSION) $(APISERVER_UI_DIR) -f $(APISERVER_UI_DIR)/Dockerfile

.PHONY: kind-load-ui
kind-load-ui: gloo-federation-console-docker
	kind load docker-image --name $(CLUSTER_NAME) $(IMAGE_REPO)/gloo-federation-console:$(VERSION)

.PHONY: cleanup-node-modules
cleanup-node-modules:
	# Remove node_modules to save disk-space (Eg in CI)
	rm -rf $(APISERVER_UI_DIR)/node_modules

.PHONY: cleanup-local-docker-images
cleanup-local-docker-images:
	# Remove all of the kind images
	docker images | grep solo-io | grep -v envoy-gloo-ee |xargs -L1 echo | cut -d ' ' -f 1 | xargs -I{} docker image rm {}:$(VERSION)
	# Remove the downloaded envoy-gloo-ee image
	docker image rm $(ENVOY_GLOO_IMAGE)
	docker image rm $(ENVOY_GLOO_FIPS_IMAGE)

#----------------------------------------------------------------------------------
# RateLimit
#----------------------------------------------------------------------------------

RATELIMIT_DIR=projects/rate-limit
RATELIMIT_SOURCES=$(shell find $(RATELIMIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)
RATELIMIT_OUT_DIR=$(OUTPUT_DIR)/rate-limit
_ := $(shell mkdir -p $(RATELIMIT_OUT_DIR))

$(RATELIMIT_OUT_DIR)/Dockerfile.build: $(RATELIMIT_DIR)/Dockerfile
	cp $< $@

$(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker-build: $(RATELIMIT_SOURCES) $(RATELIMIT_OUT_DIR)/Dockerfile.build
	docker build -t $(IMAGE_REPO)/rate-limit-ee-build-container:$(VERSION) \
		-f $(RATELIMIT_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg USE_APK=true \
		--build-arg GITHUB_TOKEN \
		.
	touch $@


# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(RATELIMIT_OUT_DIR)/rate-limit-linux-amd64: $(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker-build
	docker create -ti --name rate-limit-temp-container $(IMAGE_REPO)/rate-limit-ee-build-container:$(VERSION) bash
	docker cp rate-limit-temp-container:/rate-limit-linux-amd64 $(RATELIMIT_OUT_DIR)/rate-limit-linux-amd64
	docker rm -f rate-limit-temp-container

.PHONY: rate-limit
rate-limit: $(RATELIMIT_OUT_DIR)/rate-limit-linux-amd64

$(RATELIMIT_OUT_DIR)/Dockerfile: $(RATELIMIT_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: rate-limit-ee-docker
rate-limit-ee-docker: $(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker

$(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker: $(RATELIMIT_OUT_DIR)/rate-limit-linux-amd64 $(RATELIMIT_OUT_DIR)/Dockerfile
	docker build -t $(IMAGE_REPO)/rate-limit-ee:$(VERSION) $(call get_test_tag_option,rate-limit-ee) $(RATELIMIT_OUT_DIR)
	touch $@

#----------------------------------------------------------------------------------
# RateLimit-fips
#----------------------------------------------------------------------------------

RATELIMIT_FIPS_OUT_DIR=$(OUTPUT_DIR)/rate-limit-fips
_ := $(shell mkdir -p $(RATELIMIT_FIPS_OUT_DIR))

$(RATELIMIT_FIPS_OUT_DIR)/Dockerfile.build: $(RATELIMIT_DIR)/Dockerfile
	cp $< $@

$(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker-build: $(RATELIMIT_SOURCES) $(RATELIMIT_FIPS_OUT_DIR)/Dockerfile.build
	docker build -t $(IMAGE_REPO)/rate-limit-ee-build-container-fips:$(VERSION) \
		-f $(RATELIMIT_FIPS_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOBORING_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg GITHUB_TOKEN \
		.
	touch $@


# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64: $(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker-build
	docker create -ti --name rate-limit-temp-container $(IMAGE_REPO)/rate-limit-ee-build-container-fips:$(VERSION) bash
	docker cp rate-limit-temp-container:/rate-limit-linux-amd64 $(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64
	docker rm -f rate-limit-temp-container

.PHONY: rate-limit-fips
rate-limit-fips: $(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64

$(RATELIMIT_FIPS_OUT_DIR)/Dockerfile: $(RATELIMIT_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: rate-limit-ee-fips-docker
rate-limit-ee-fips-docker: $(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker

$(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker: $(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64 $(RATELIMIT_FIPS_OUT_DIR)/Dockerfile
	docker build -t $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION) $(call get_test_tag_option,rate-limit-ee-fips) $(RATELIMIT_FIPS_OUT_DIR) \
       --build-arg EXTRA_PACKAGES=libc6-compat
	touch $@

#----------------------------------------------------------------------------------
# ExtAuth
#----------------------------------------------------------------------------------

EXTAUTH_DIR=projects/extauth
EXTAUTH_SOURCES=$(shell find $(EXTAUTH_DIR) -name "*.go" | grep -v test | grep -v generated.go)
EXTAUTH_OUT_DIR=$(OUTPUT_DIR)/extauth
RELATIVE_EXTAUTH_OUT_DIR=$(RELATIVE_OUTPUT_DIR)/extauth
_ := $(shell mkdir -p $(EXTAUTH_OUT_DIR))

$(EXTAUTH_OUT_DIR)/Dockerfile.build: $(EXTAUTH_DIR)/Dockerfile
	cp $< $@

$(EXTAUTH_OUT_DIR)/Dockerfile: $(EXTAUTH_DIR)/cmd/Dockerfile
	cp $< $@

$(EXTAUTH_OUT_DIR)/.extauth-ee-docker-build: $(EXTAUTH_SOURCES) $(EXTAUTH_OUT_DIR)/Dockerfile.build
	docker build -t $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) \
		-f $(EXTAUTH_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg USE_APK=true \
		--build-arg GITHUB_TOKEN \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
$(EXTAUTH_OUT_DIR)/extauth-linux-amd64: $(EXTAUTH_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name extauth-temp-container $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) bash
	docker cp extauth-temp-container:/extauth-linux-amd64 $(EXTAUTH_OUT_DIR)/extauth-linux-amd64
	docker rm -f extauth-temp-container

# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64: $(EXTAUTH_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name verify-plugins-temp-container $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) bash
	docker cp verify-plugins-temp-container:/verify-plugins-linux-amd64 $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64
	docker rm -f verify-plugins-temp-container

# Build extauth binaries
.PHONY: extauth
extauth: $(EXTAUTH_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64

# Build ext-auth-plugins docker image
.PHONY: ext-auth-plugins-docker
ext-auth-plugins-docker: $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64
	docker build -t $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) -f projects/extauth/plugins/Dockerfile \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_VERSION) \
		--build-arg GC_FLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg VERIFY_SCRIPT=$(RELATIVE_EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 \
		--build-arg GITHUB_TOKEN \
		--build-arg USE_APK=true \
		.

# Build extauth server docker image
.PHONY: extauth-ee-docker
extauth-ee-docker: $(EXTAUTH_OUT_DIR)/.extauth-ee-docker

$(EXTAUTH_OUT_DIR)/.extauth-ee-docker: $(EXTAUTH_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 $(EXTAUTH_OUT_DIR)/Dockerfile
	docker build -t $(IMAGE_REPO)/extauth-ee:$(VERSION) $(call get_test_tag_option,extauth-ee) $(EXTAUTH_OUT_DIR)
	touch $@

#----------------------------------------------------------------------------------
# ExtAuth-fips
#----------------------------------------------------------------------------------

EXTAUTH_FIPS_OUT_DIR=$(OUTPUT_DIR)/extauth_fips
RELATIVE_EXTAUTH_FIPS_OUT_DIR=$(RELATIVE_OUTPUT_DIR)/extauth_fips
_ := $(shell mkdir -p $(EXTAUTH_FIPS_OUT_DIR))

$(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build: $(EXTAUTH_DIR)/Dockerfile
	cp $< $@

$(EXTAUTH_FIPS_OUT_DIR)/Dockerfile: $(EXTAUTH_DIR)/cmd/Dockerfile
	cp $< $@

$(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker-build: $(EXTAUTH_SOURCES) $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build
	docker build -t $(IMAGE_REPO)/extauth-ee-build-container-fips:$(VERSION) \
		-f $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOBORING_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg GITHUB_TOKEN \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64: $(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name extauth-temp-container $(IMAGE_REPO)/extauth-ee-build-container-fips:$(VERSION) bash
	docker cp extauth-temp-container:/extauth-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64
	docker rm -f extauth-temp-container

# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64: $(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name verify-plugins-temp-container $(IMAGE_REPO)/extauth-ee-build-container-fips:$(VERSION) bash
	docker cp verify-plugins-temp-container:/verify-plugins-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64
	docker rm -f verify-plugins-temp-container

# Build extauth binaries
.PHONY: extauth-fips
extauth-fips: $(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64

# Build ext-auth-plugins docker image
.PHONY: ext-auth-plugins-fips-docker
ext-auth-plugins-fips-docker: $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64
	docker build -t $(IMAGE_REPO)/ext-auth-plugins-fips:$(VERSION) -f projects/extauth/plugins/Dockerfile \
		--build-arg GO_BUILD_IMAGE=$(GOBORING_VERSION) \
		--build-arg GC_FLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg VERIFY_SCRIPT=$(RELATIVE_EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 \
		--build-arg GITHUB_TOKEN \
		.

# Build extauth server docker image
.PHONY: extauth-ee-fips-docker
extauth-ee-fips-docker: $(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker

$(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker: $(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile
	docker build -t $(IMAGE_REPO)/extauth-ee-fips:$(VERSION) $(call get_test_tag_option,extauth-ee-fips) $(EXTAUTH_FIPS_OUT_DIR) \
		--build-arg EXTRA_PACKAGES=libc6-compat
	touch $@

#----------------------------------------------------------------------------------
# Observability
#----------------------------------------------------------------------------------

OBSERVABILITY_DIR=projects/observability
OBSERVABILITY_SOURCES=$(shell find $(OBSERVABILITY_DIR) -name "*.go" | grep -v test | grep -v generated.go)
OBS_OUT_DIR=$(OUTPUT_DIR)/observability

$(OBS_OUT_DIR)/observability-linux-amd64: $(OBSERVABILITY_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(OBSERVABILITY_DIR)/cmd/main.go

.PHONY: observability
observability: $(OBS_OUT_DIR)/observability-linux-amd64

$(OBS_OUT_DIR)/Dockerfile: $(OBSERVABILITY_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: observability-ee-docker
observability-ee-docker: $(OBS_OUT_DIR)/.observability-ee-docker

$(OBS_OUT_DIR)/.observability-ee-docker: $(OBS_OUT_DIR)/observability-linux-amd64 $(OBS_OUT_DIR)/Dockerfile
	docker build -t $(IMAGE_REPO)/observability-ee:$(VERSION) $(call get_test_tag_option,observability-ee) $(OBS_OUT_DIR)
	touch $@

#----------------------------------------------------------------------------------
# Gloo
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)
GLOO_OUT_DIR=$(OUTPUT_DIR)/gloo

$(GLOO_OUT_DIR)/gloo-linux-amd64: $(GLOO_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_DIR)/cmd/main.go


.PHONY: gloo
gloo: $(GLOO_OUT_DIR)/gloo-linux-amd64

$(GLOO_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@


.PHONY: gloo-ee-docker
gloo-ee-docker: $(GLOO_OUT_DIR)/.gloo-ee-docker

$(GLOO_OUT_DIR)/.gloo-ee-docker: $(GLOO_OUT_DIR)/gloo-linux-amd64 $(GLOO_OUT_DIR)/Dockerfile
	docker build $(call get_test_tag_option,gloo-ee) $(GLOO_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) \
		-t $(IMAGE_REPO)/gloo-ee:$(VERSION)
	touch $@

gloo-ee-docker-dev: $(GLOO_OUT_DIR)/gloo-linux-amd64 $(GLOO_OUT_DIR)/Dockerfile
	docker build -t $(IMAGE_REPO)/gloo-ee:$(VERSION) $(GLOO_OUT_DIR) --no-cache
	touch $@

#----------------------------------------------------------------------------------
# Gloo with FIPS Envoy
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)
GLOO_FIPS_OUT_DIR=$(OUTPUT_DIR)/gloo-fips

$(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64: $(GLOO_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_DIR)/cmd/main.go


.PHONY: gloo-fips
gloo-fips: $(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64

$(GLOO_FIPS_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@


.PHONY: gloo-fips-ee-docker
gloo-fips-ee-docker: $(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker

$(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker: $(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64 $(GLOO_FIPS_OUT_DIR)/Dockerfile
	docker build $(call get_test_tag_option,gloo-ee) $(GLOO_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_IMAGE) \
		-t $(IMAGE_REPO)/gloo-ee-fips:$(VERSION)
	touch $@

gloo-fips-ee-docker-dev: $(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64 $(GLOO_FIPS_OUT_DIR)/Dockerfile
	docker build -t $(IMAGE_REPO)/gloo-ee-fips:$(VERSION) $(GLOO_FIPS_OUT_DIR) --no-cache
	touch $@
#----------------------------------------------------------------------------------
# discovery (enterprise)
#----------------------------------------------------------------------------------

DISCOVERY_DIR=projects/discovery
DISCOVERY_SOURCES=$(shell find $(DISCOVERY_DIR) -name "*.go" | grep -v test | grep -v generated.go)
DISCOVERY_OUTPUT_DIR=$(OUTPUT_DIR)/$(DISCOVERY_DIR)

$(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-amd64: $(DISCOVERY_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(DISCOVERY_DIR)/cmd/main.go

.PHONY: discovery-ee
discovery: $(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-amd64

$(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery: $(DISCOVERY_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: discovery-ee-docker
discovery-ee-docker: $(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-amd64 $(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery
	docker build $(DISCOVERY_OUTPUT_DIR) -f $(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery \
		--build-arg GOARCH=amd64 \
		-t $(IMAGE_REPO)/discovery-ee:$(VERSION) $(QUAY_EXPIRATION_LABEL)

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
# Glooctl Plugins
#----------------------------------------------------------------------------------

# Include helm makefile so its targets can be ran from the root of this repo
include $(ROOTDIR)/projects/glooctl-plugins/plugins.mk

#----------------------------------------------------------------------------------
# Envoy init (BASE/SIDECAR)
#----------------------------------------------------------------------------------

ENVOYINIT_DIR=cmd/envoyinit
ENVOYINIT_SOURCES=$(shell find $(ENVOYINIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)
ENVOYINIT_OUT_DIR=$(OUTPUT_DIR)/envoyinit

$(ENVOYINIT_OUT_DIR)/envoyinit-linux-amd64: $(ENVOYINIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go $(ENVOYINIT_DIR)/filter_types.gen.go

.PHONY: envoyinit
envoyinit: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-amd64

$(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit: $(ENVOYINIT_DIR)/Dockerfile.envoyinit
	cp $< $@

$(ENVOYINIT_OUT_DIR)/docker-entrypoint.sh: $(ENVOYINIT_DIR)/docker-entrypoint.sh
	cp $< $@

.PHONY: gloo-ee-envoy-wrapper-docker
gloo-ee-envoy-wrapper-docker: $(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-docker

$(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-docker: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-amd64 $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_OUT_DIR)/docker-entrypoint.sh
	docker build $(call get_test_tag_option,gloo-ee-envoy-wrapper) $(ENVOYINIT_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) \
		-t $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) \
		-f $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit
	touch $@

#----------------------------------------------------------------------------------
# Fips Envoy init (BASE/SIDECAR)
#----------------------------------------------------------------------------------

ENVOYINIT_FIPS_OUT_DIR=$(OUTPUT_DIR)/envoyinit_fips

$(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64: $(ENVOYINIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go $(ENVOYINIT_DIR)/filter_types.gen.go

.PHONY: envoyinit-fips
envoyinit-fips: $(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64

$(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit: $(ENVOYINIT_DIR)/Dockerfile.envoyinit
	cp $< $@

$(ENVOYINIT_FIPS_OUT_DIR)/docker-entrypoint.sh: $(ENVOYINIT_DIR)/docker-entrypoint.sh
	cp $< $@

.PHONY: gloo-ee-envoy-wrapper-fips-docker
gloo-ee-envoy-wrapper-fips-docker: $(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-docker

$(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-docker: $(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64 $(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_FIPS_OUT_DIR)/docker-entrypoint.sh
	docker build $(call get_test_tag_option,gloo-ee-envoy-wrapper-fips) $(ENVOYINIT_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_IMAGE) \
		-t $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) \
		-f $(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit
	touch $@
#----------------------------------------------------------------------------------
# Deployment Manifests / Helm
#----------------------------------------------------------------------------------
HELM_SYNC_DIR_FOR_GLOO_EE := $(OUTPUT_DIR)/helm
HELM_SYNC_DIR_GLOO_FED := $(OUTPUT_DIR)/helm_gloo_fed
HELM_DIR := $(ROOTDIR)/install/helm
GLOOE_CHART_DIR := $(HELM_DIR)/gloo-ee
GLOO_FED_CHART_DIR := $(HELM_DIR)/gloo-fed
MANIFEST_DIR := install/manifest
MANIFEST_FOR_GLOO_EE := glooe-release.yaml
MANIFEST_FOR_GLOO_FED := gloo-fed-release.yaml
GLOOE_HELM_BUCKET := gs://gloo-ee-helm
GLOO_FED_HELM_BUCKET := gs://gloo-fed-helm

.PHONY: manifest
manifest: init-helm produce-manifests

# creates Chart.yaml, values.yaml, and requirements.yaml
.PHONY: helm-template
helm-template:
	mkdir -p $(MANIFEST_DIR)
	mkdir -p $(HELM_SYNC_DIR_FOR_GLOO_EE)
	PATH=$(DEPSGOBIN):$$PATH $(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION) --gloo-fed-repo-override="file://$(GLOO_FED_CHART_DIR)"

.PHONY: init-helm
init-helm: helm-template gloofed-helm-template $(OUTPUT_DIR)/.helm-initialized

$(OUTPUT_DIR)/.helm-initialized:
	helm repo add helm-hub https://charts.helm.sh/stable
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm repo add gloo-fed https://storage.googleapis.com/gloo-fed-helm
	helm dependency update install/helm/gloo-ee
	touch $@

.PHONY: produce-manifests
produce-manifests: gloofed-produce-manifests
	helm template glooe install/helm/gloo-ee --namespace gloo-system > $(MANIFEST_DIR)/$(MANIFEST_FOR_GLOO_EE)

.PHONY: package-gloo-edge-chart
package-gloo-edge-chart: init-helm
	helm package --destination $(HELM_SYNC_DIR_FOR_GLOO_EE) $(GLOOE_CHART_DIR)

.PHONY: fetch-package-and-save-helm
fetch-package-and-save-helm: init-helm package-gloo-fed-chart
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
	until $$(GENERATION=$$(gsutil ls -a $(GLOO_FED_HELM_BUCKET)/index.yaml | tail -1 | cut -f2 -d '#') && \
		  gsutil cp -v $(GLOO_FED_HELM_BUCKET)/index.yaml $(HELM_SYNC_DIR_GLOO_FED)/index.yaml && \
		  helm package --destination $(HELM_SYNC_DIR_GLOO_FED)/charts $(HELM_DIR)/gloo-fed >> /dev/null && \
		  helm repo index $(HELM_SYNC_DIR_GLOO_FED) --merge $(HELM_SYNC_DIR_GLOO_FED)/index.yaml && \
		  gsutil -m rsync $(HELM_SYNC_DIR_GLOO_FED)/charts $(GLOO_FED_HELM_BUCKET) && \
		  gsutil -h x-goog-if-generation-match:"$$GENERATION" cp $(HELM_SYNC_DIR_GLOO_FED)/index.yaml $(GLOO_FED_HELM_BUCKET)/index.yaml); do \
	echo "Failed to upload new helm index (updated helm index since last download?). Trying again"; \
	sleep 2; \
	done
endif

#----------------------------------------------------------------------------------
# Gloo Fed Deployment Manifests / Helm
#----------------------------------------------------------------------------------

# creates Chart.yaml, values.yaml, and requirements.yaml
.PHONY: gloofed-helm-template
gloofed-helm-template:
	mkdir -p $(HELM_SYNC_DIR_GLOO_FED)
	sed -e 's/%version%/'$(VERSION)'/' $(GLOO_FED_CHART_DIR)/Chart-template.yaml > $(GLOO_FED_CHART_DIR)/Chart.yaml
	sed -e 's/%version%/'$(VERSION)'/' $(GLOO_FED_CHART_DIR)/values-template.yaml > $(GLOO_FED_CHART_DIR)/values.yaml

.PHONY: gloofed-produce-manifests
gloofed-produce-manifests: gloofed-helm-template
	helm template gloo-fed install/helm/gloo-fed --namespace gloo-system > $(MANIFEST_DIR)/$(MANIFEST_FOR_GLOO_FED)

.PHONY: package-gloo-fed-chart
package-gloo-fed-chart: gloofed-helm-template
	helm package --destination $(HELM_SYNC_DIR_GLOO_FED) $(GLOO_FED_CHART_DIR)

#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------

.PHONY: upload-github-release-assets
upload-github-release-assets: produce-manifests
	$(GO_BUILD_FLAGS) go run ci/upload_github_release_assets.go

DEPS_DIR=$(OUTPUT_DIR)/dependencies/$(VERSION)
DEPS_BUCKET=gloo-ee-dependencies

.PHONY: publish-dependencies
publish-dependencies: $(DEPS_DIR)/go.mod $(DEPS_DIR)/go.sum $(DEPS_DIR)/dependencies $(DEPS_DIR)/dependencies.json \
	$(DEPS_DIR)/build_env $(DEPS_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)/fips-verify-plugins-linux-amd64
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
	echo "GO_BUILD_IMAGE=$(GOLANG_VERSION)" > $@
	echo "FIPS_GO_BUILD_IMAGE=$(GOBORING_VERSION)" >> $@
	echo "GC_FLAGS=$(GCFLAGS)" >> $@

$(DEPS_DIR)/verify-plugins-linux-amd64: $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)
	cp $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)
$(DEPS_DIR)/fips-verify-plugins-linux-amd64: $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)
	cp $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)/fips-verify-plugins-linux-amd64

#----------------------------------------------------------------------------------
# Docker push
#----------------------------------------------------------------------------------

DOCKER_IMAGES :=
ifeq ($(RELEASE),"true")
	DOCKER_IMAGES := docker
endif

.PHONY: docker docker-push
 docker: rate-limit-ee-docker rate-limit-ee-fips-docker extauth-ee-docker \
       extauth-ee-fips-docker gloo-ee-docker gloo-fips-ee-docker gloo-ee-envoy-wrapper-docker discovery-ee-docker\
       gloo-ee-envoy-wrapper-fips-docker observability-ee-docker ext-auth-plugins-docker ext-auth-plugins-fips-docker \
       gloo-fed-docker gloo-fed-apiserver-docker gloo-fed-apiserver-envoy-docker gloo-federation-console-docker gloo-fed-rbac-validating-webhook-docker

# Depends on DOCKER_IMAGES, which is set to docker if RELEASE is "true", otherwise empty (making this a no-op).
# This prevents executing the dependent targets if RELEASE is not true, while still enabling `make docker`
# to be used for local testing.
# docker-push is intended to be run by CI
docker-push: $(DOCKER_IMAGES)
ifeq ($(RELEASE),"true")
	docker push $(IMAGE_REPO)/rate-limit-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/observability-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/extauth-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/extauth-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/discovery-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) && \
	docker push $(IMAGE_REPO)/ext-auth-plugins-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-federation-console:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION)
endif

.PHONY: docker-push-extended
docker-push-extended:
ifeq ($(RELEASE),"true")
	ci/extended-docker/extended-docker.sh
endif

# Retag and push images that are already built
.PHONY: docker-retag-images
docker-retag-images:
ifeq ($(RELEASE),"true")
	docker tag $(RETAG_IMAGE_REGISTRY)/rate-limit-ee:$(VERSION) $(IMAGE_REPO)/rate-limit-ee:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/rate-limit-ee-fips:$(VERSION) $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee:$(VERSION) $(IMAGE_REPO)/gloo-ee:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-fips:$(VERSION) $(IMAGE_REPO)/gloo-ee-fips:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper:$(VERSION) $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper-fips:$(VERSION) $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/observability-ee:$(VERSION) $(IMAGE_REPO)/observability-ee:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/extauth-ee:$(VERSION) $(IMAGE_REPO)/extauth-ee:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/extauth-ee-fips:$(VERSION) $(IMAGE_REPO)/extauth-ee-fips:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/discovery-ee:$(VERSION) $(IMAGE_REPO)/discovery-ee:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/ext-auth-plugins:$(VERSION) $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/ext-auth-plugins-fips:$(VERSION) $(IMAGE_REPO)/ext-auth-plugins-fips:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed:$(VERSION) $(IMAGE_REPO)/gloo-fed:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed-apiserver:$(VERSION) $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed-apiserver-envoy:$(VERSION) $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-federation-console:$(VERSION) $(IMAGE_REPO)/gloo-federation-console:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed-rbac-validating-webhook:$(VERSION) $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/rate-limit-ee:$(VERSION)-extended $(IMAGE_REPO)/rate-limit-ee:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/rate-limit-ee-fips:$(VERSION)-extended $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee:$(VERSION)-extended $(IMAGE_REPO)/gloo-ee:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-fips:$(VERSION)-extended $(IMAGE_REPO)/gloo-ee-fips:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper:$(VERSION)-extended $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper-fips:$(VERSION)-extended $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/observability-ee:$(VERSION)-extended $(IMAGE_REPO)/observability-ee:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/extauth-ee:$(VERSION)-extended $(IMAGE_REPO)/extauth-ee:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/extauth-ee-fips:$(VERSION)-extended $(IMAGE_REPO)/extauth-ee-fips:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed:$(VERSION)-extended $(IMAGE_REPO)/gloo-fed:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed-apiserver:$(VERSION)-extended $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed-apiserver-envoy:$(VERSION)-extended $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-federation-console:$(VERSION)-extended $(IMAGE_REPO)/gloo-federation-console:$(VERSION)-extended && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-fed-rbac-validating-webhook:$(VERSION)-extended $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION)-extended

	docker push $(IMAGE_REPO)/rate-limit-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/observability-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/extauth-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/discovery-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/extauth-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) && \
	docker push $(IMAGE_REPO)/ext-auth-plugins-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-federation-console:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION) && \
	docker push $(IMAGE_REPO)/rate-limit-ee:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-ee:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-ee-fips:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/observability-ee:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/extauth-ee:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/extauth-ee-fips:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-fed:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-federation-console:$(VERSION)-extended && \
	docker push $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION)-extended
endif

# Helper targets for CI
CLUSTER_NAME?=kind

.PHONY: push-kind-images
push-kind-images:
ifeq ($(USE_FIPS),true)
push-kind-images: build-and-load-kind-images-fips
else
push-kind-images: build-and-load-kind-images-non-fips
endif

# Build and load images for a non-fips compliant (data plane) installation of Gloo Edge
# Used in CI during regression tests
.PHONY: build-and-load-kind-images-non-fips
build-and-load-kind-images-non-fips: gloo-ee-docker kind-load-gloo-ee # gloo
build-and-load-kind-images-non-fips: gloo-ee-envoy-wrapper-docker kind-load-gloo-ee-envoy-wrapper # envoy
build-and-load-kind-images-non-fips: rate-limit-ee-docker kind-load-rate-limit-ee # rate limit
build-and-load-kind-images-non-fips: extauth-ee-docker kind-load-extauth-ee # ext auth
build-and-load-kind-images-non-fips: ext-auth-plugins-docker kind-load-ext-auth-plugins # ext auth plugins
build-and-load-kind-images-non-fips: observability-ee-docker kind-load-observability-ee # observability
build-and-load-kind-images-non-fips: discovery-ee-docker kind-load-discovery-ee # discovery

# Build and load images for a fips compliant (data plane) installation of Gloo Edge
# Used in CI during regression tests
.PHONY: build-and-load-kind-images-fips
build-and-load-kind-images-fips: gloo-fips-ee-docker kind-load-gloo-ee-fips # gloo
build-and-load-kind-images-fips: gloo-ee-envoy-wrapper-fips-docker kind-load-gloo-ee-envoy-wrapper-fips # envoy
build-and-load-kind-images-fips: rate-limit-ee-fips-docker kind-load-rate-limit-ee-fips # rate limit
build-and-load-kind-images-fips: extauth-ee-fips-docker kind-load-extauth-ee-fips # ext auth
build-and-load-kind-images-fips: ext-auth-plugins-fips-docker kind-load-ext-auth-plugins-fips # ext auth plugins
build-and-load-kind-images-fips: observability-ee-docker kind-load-observability-ee # observability
build-and-load-kind-images-fips: discovery-ee-docker kind-load-discovery-ee # discovery

.PHONY: build-kind-assets
build-kind-assets: push-kind-images build-test-chart

TEST_DOCKER_TARGETS := gloo-federation-console-docker-test apiserver-envoy-docker-test gloo-fed-apiserver-docker-test rate-limit-ee-docker-test extauth-ee-docker-test observability-ee-docker-test gloo-ee-docker-test gloo-ee-envoy-wrapper-docker-test

push-test-images: $(TEST_DOCKER_TARGETS)

gloo-fed-apiserver-docker-test: $(OUTPUT_DIR)/gloo-fed-apiserver-linux-amd64 $(OUTPUT_DIR)/.gloo-fed-apiserver-docker
	docker push $(call get_test_tag,gloo-fed-apiserver)

gloo-fed-apiserver-envoy-docker-test: gloo-fed-apiserver-envoy-docker $(OUTPUT_DIR)/Dockerfile
	docker push $(call get_test_tag,gloo-fed-apiserver-envoy)

gloo-federation-console-docker-test: build-ui gloo-federation-console-docker
	docker push $(call get_test_tag,gloo-federation-console)

rate-limit-ee-docker-test: $(RATELIMIT_OUT_DIR)/rate-limit-linux-amd64 $(RATELIMIT_OUT_DIR)/Dockerfile
	docker push $(call get_test_tag,rate-limit-ee)

extauth-ee-docker-test: $(EXTAUTH_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_OUT_DIR)/Dockerfile
	docker push $(call get_test_tag,extauth-ee)

observability-ee-docker-test: $(OBS_OUT_DIR)/observability-linux-amd64 $(OBS_OUT_DIR)/Dockerfile
	docker push $(call get_test_tag,observability-ee)

gloo-ee-docker-test: gloo-ee-docker
	docker push $(call get_test_tag,gloo-ee)

gloo-ee-envoy-wrapper-docker-test: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-amd64 $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit gloo-ee-envoy-wrapper-docker
	docker push $(call get_test_tag,gloo-ee-envoy-wrapper)

.PHONY: build-test-chart
build-test-chart: build-test-chart-fed
	mkdir -p $(TEST_ASSET_DIR)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION) --gloo-fed-repo-override="file://$(GLOO_FED_CHART_DIR)"
	helm repo add helm-hub https://charts.helm.sh/stable
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm dependency update install/helm/gloo-ee
	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-ee
	helm repo index $(TEST_ASSET_DIR)

	# unarchive all files in the test archive directory in-place
	ls $(TEST_ASSET_DIR)/*.tgz | xargs -n1 tar -C $(TEST_ASSET_DIR) -xf

	echo "\n\
	\x1B[32mSuccessfully built 'gloo-ee' helm chart!  now render templates locally! \n\n\
	Local rendering Docs: https://helm.sh/docs/helm/helm_template/ \n\
	Suggested first command: \n\
	\t\033[1m$$ helm template ./_test/gloo-ee\033[0m\n"


.PHONY: build-test-chart-fed
build-test-chart-fed: gloofed-helm-template
	if [ $(FORCE_CLEAN_TEST_ASSET_DIR) = true ] ; then\
		rm -rf $(TEST_ASSET_DIR);\
	fi
	mkdir -p $(TEST_ASSET_DIR)
	helm repo add helm-hub https://charts.helm.sh/stable
	helm repo add gloo-fed https://storage.googleapis.com/gloo-fed-helm
	helm dependency update install/helm/gloo-fed
	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-fed
	helm repo index $(TEST_ASSET_DIR)

# Exclusively useful for testing with locally modified gloo-edge-OS builds.
# Assumes that a gloo-edge-OS chart is located at ../gloo/_test/gloo-dev.tgz, which
# points towards whatever modified build is being tested.
.PHONY: build-chart-with-local-gloo-dev
build-chart-with-local-gloo-dev:
	mkdir -p $(TEST_ASSET_DIR)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION)
	helm repo add helm-hub https://charts.helm.sh/stable
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm dependency update install/helm/gloo-ee
	echo replacing gloo chart $(ls install/helm/gloo-ee/charts/gloo*) with ../gloo/_test/gloo-dev.tgz
	rm install/helm/gloo-ee/charts/gloo*
	cp ../gloo/_test/gloo-dev.tgz install/helm/gloo-ee/charts/
	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-ee
	helm repo index $(TEST_ASSET_DIR)

#----------------------------------------------------------------------------------
# CI: Escrow
#----------------------------------------------------------------------------------

DEPOSITOR_NAME ?= 'Janice Morales'
DEPOSITOR_EMAIL ?= 'janice.morales@solo.io'
DEPOSITOR_PHONE ?= '(617)-893-7557'
DEPOSIT_DATE ?= '$(shell date)'

.PHONY: tar-repo
tar-repo:
	$(eval GIT_BRANCH=$(shell git rev-parse HEAD))
	git checkout v$(VERSION)
	go mod vendor
	tar -cvf gloo-ee-$(VERSION).tar ./
	mkdir -p _gloo-ee-source/
	mv gloo-ee-$(VERSION).tar _gloo-ee-source/
	git checkout $(GIT_BRANCH)

.PHONY: generate-escrow-pdf
generate-escrow-pdf: tar-repo
	$(eval DEPOSIT_BYTES=$(shell wc -c < _gloo-ee-source/gloo-ee-$(VERSION).tar))
	deno run --allow-read --allow-write --allow-net ci/escrow/modify-pdf.ts $(VERSION) $(DEPOSIT_BYTES) $(DEPOSITOR_NAME) $(DEPOSIT_DATE) $(DEPOSITOR_EMAIL) $(DEPOSITOR_PHONE)

#----------------------------------------------------------------------------------
# Third Party License Management
#----------------------------------------------------------------------------------
.PHONY: update-licenses
update-licenses:
	# check for GPL licenses, if there are any, this will fail
	GO111MODULE=on go run hack/oss_compliance/oss_compliance.go osagen -c "GNU General Public License v2.0,GNU General Public License v3.0,GNU Affero General Public License v3.0"

	GO111MODULE=on go run hack/oss_compliance/oss_compliance.go osagen -s "Mozilla Public License 2.0,GNU General Public License v2.0,GNU General Public License v3.0,GNU Lesser General Public License v2.1,GNU Lesser General Public License v3.0,GNU Affero General Public License v3.0" > ci/licenses/osa_provided.md
	GO111MODULE=on go run hack/oss_compliance/oss_compliance.go osagen -i "Mozilla Public License 2.0,GNU Lesser General Public License v2.1,GNU Lesser General Public License v3.0"> ci/licenses/osa_included.md

#----------------------------------------------------------------------------------
# Printing makefile variables utility
#----------------------------------------------------------------------------------

# use `make print-MAKEFILE_VAR` to print the value of MAKEFILE_VAR

print-%  : ; @echo $($*)

SCAN_DIR ?= $(OUTPUT_DIR)/scans
SCAN_BUCKET ?= solo-gloo-security-scans/glooe

.PHONY: publish-security-scan
publish-security-scan:
ifeq ($(RELEASE),"true")
	gsutil cp -r $(SCAN_DIR)/$(VERSION)/$(SCAN_FILE) gs://$(SCAN_BUCKET)/$(VERSION)/$(SCAN_FILE)
endif