include Makefile.docker

# https://www.gnu.org/software/make/manual/html_node/Special-Variables.html#Special-Variables
.DEFAULT_GOAL := help

#----------------------------------------------------------------------------------
# Help
#----------------------------------------------------------------------------------
# Our Makefile is quite large, and hard to reason through
# `make help` can be used to self-document targets
# To update a target to be self-documenting (and appear with the `help` command),
# place a comment after the target that is prefixed by `##`. For example:
#	custom-target: ## comment that will appear in the documentation when running `make help`
#
# **NOTE TO DEVELOPERS**
# As you encounter make targets that are frequently used, please make them self-documenting

.PHONY: help
help: NAME_COLUMN_WIDTH=35
help: LINE_COLUMN_WIDTH=5
help: ## Output the self-documenting make targets
	@grep -hnE '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "[:]|(## )"}; {printf "\033[36mL%-$(LINE_COLUMN_WIDTH)s%-$(NAME_COLUMN_WIDTH)s\033[0m %s\n", $$1, $$2, $$4}'

#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output
DEPSGOBIN := $(OUTPUT_DIR)/.bin
TEST_ASSET_DIR := $(ROOTDIR)/_test

# Important to use binaries built from module.
export PATH:=$(DEPSGOBIN):$(PATH)
export GOBIN:=$(DEPSGOBIN)

PACKAGE_PATH:=github.com/solo-io/solo-projects
# Kind of a hack to make sure _output exists
z := $(shell mkdir -p $(OUTPUT_DIR))
SOURCES := $(shell find . -name "*.go" | grep -v test.go)

GCS_BUCKET := glooctl-plugins
FED_GCS_PATH := glooctl-fed

ENVOY_GLOO_IMAGE_VERSION ?= 1.25.6-patch1
ENVOY_GLOO_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee:$(ENVOY_GLOO_IMAGE_VERSION)
ENVOY_GLOO_DEBUG_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee-debug:$(ENVOY_GLOO_IMAGE_VERSION)
ENVOY_GLOO_FIPS_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee-fips:$(ENVOY_GLOO_IMAGE_VERSION)
ENVOY_GLOO_FIPS_DEBUG_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee-fips-debug:$(ENVOY_GLOO_IMAGE_VERSION)

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

LDFLAGS :="-X github.com/solo-io/solo-projects/pkg/version.Version=$(VERSION) -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore"
LD_STATIC_LINKING_FLAGS:="-linkmode external -extldflags=-static -X github.com/solo-io/solo-projects/pkg/version.Version=$(VERSION) -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore"
GCFLAGS := 'all=-N -l'

GO_BUILD_FLAGS := GO111MODULE=on CGO_ENABLED=0 GOARCH=$(DOCKER_GOARCH)

# Passed by cloudbuild
GCLOUD_PROJECT_ID := $(GCLOUD_PROJECT_ID)
BUILD_ID := $(BUILD_ID)

TEST_IMAGE_TAG := test-$(BUILD_ID)
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

.PHONY: mod-download
mod-download: check-go-version
	go mod download all

# https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
.PHONY: install-go-tools
install-go-tools: mod-download ## Download and install Go dependencies
	go install istio.io/tools/cmd/protoc-gen-jsonshim
	go install istio.io/pkg/version
	go install github.com/solo-io/protoc-gen-ext
	go install golang.org/x/tools/cmd/goimports
	go install github.com/envoyproxy/protoc-gen-validate
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
	go install github.com/google/wire/cmd/wire
	go install github.com/solo-io/protoc-gen-openapi

.PHONY: update-all-deps
update-all-deps: install-go-tools install-node-packages ## install-go-tools and install-node-packages

.PHONY: install-node-packages
install-node-packages:
	npm install -g yarn --force
	npm install -g esbuild@0.16.14
	make install-graphql-js
	make update-ui-deps
	make build-stitching-bundles

.PHONY: update-ui-deps
update-ui-deps:
	yarn --cwd=$(APISERVER_UI_DIR) install

.PHONY: fmt-changed
fmt-changed:
	git diff --name-only | grep '.*.go$$' | xargs $(DEPSGOBIN)/goimports -w

.PHONY: fmt
fmt:
	$(DEPSGOBIN)/goimports -w $(shell ls -d */ | grep -E -v 'vendor|node_modules|_output|_test|changelog')

.PHONY: check-format
check-format:
	NOT_FORMATTED=$$(gofmt -l ./projects/ ./pkg/ ./test/ ./install/) && if [ -n "$$NOT_FORMATTED" ]; then echo These files are not formatted: $$NOT_FORMATTED; exit 1; fi

#----------------------------------------------------------------------------------
# Tests
#----------------------------------------------------------------------------------

GINKGO_VERSION ?= 2.8.1 # match our go.mod
GINKGO_ENV ?= GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore ACK_GINKGO_RC=true ACK_GINKGO_DEPRECATIONS=$(GINKGO_VERSION) VERSION=$(VERSION)
GINKGO_FLAGS ?= -fail-fast -trace -progress -compilers=4 -fail-on-pending
GINKGO_RANDOMIZE_FLAGS ?= -randomize-all # Not yet actively supported
GINKGO_REPORT_FLAGS ?= --json-report=test-report.json --junit-report=junit.xml -output-dir=$(OUTPUT_DIR)
GINKGO_COVERAGE_FLAGS ?= --cover --covermode=count --coverprofile=coverage.cov
TEST_PKG ?= ./... # Default to run all tests

# This is a way for a user executing `make test` to be able to provide flags which we do not include by default
# For example, you may want to run tests multiple times, or with various timeouts
GINKGO_USER_FLAGS ?=

.PHONY: install-test-tools
install-test-tools: check-go-version
install-test-tools:
	go install github.com/onsi/ginkgo/v2/ginkgo@v$(GINKGO_VERSION)

.PHONY: test
test: install-test-tools ## Run all tests, or only run the test package at {TEST_PKG} if it is specified
	$(GINKGO_ENV) ginkgo -ldflags=$(LDFLAGS) \
	$(GINKGO_FLAGS) $(GINKGO_REPORT_FLAGS) $(GINKGO_USER_FLAGS) \
	$(TEST_PKG)

.PHONY: test-with-coverage
test-with-coverage: GINKGO_FLAGS += $(GINKGO_COVERAGE_FLAGS)
test-with-coverage: test
	go tool cover -html $(OUTPUT_DIR)/coverage.cov

.PHONY: run-tests
run-tests: GINKGO_FLAGS += -skip-package=kube2e,federation-kube2e ## Run all tests, or only run the test package at {TEST_PKG} if it is specified
ifneq ($(RELEASE), "true")
run-tests: generate-extauth-test-plugins
run-tests: test
endif

# requires the environment variable KUBE2E_TESTS to be set to the test type you wish to run
.PHONY: run-ci-regression-tests
run-ci-regression-tests: TEST_PKG = ./test/kube2e/$(KUBE2E_TESTS) ## Run the Kubernetes E2E Tests in the {KUBE2E_TESTS} package
run-ci-regression-tests: test

# requires the environment variable KUBE2E_TESTS to be set to the test type you wish to run
.PHONY: run-ci-gloo-fed-regression-tests
run-ci-gloo-fed-regression-tests: TEST_PKG = ./test/federation-kube2e/$(KUBE2E_TESTS) ## Run the Federation Kubernetes E2E Tests in the {KUBE2E_TESTS} package
run-ci-gloo-fed-regression-tests: test

# command to run e2e tests
# requires the environment variable ENVOY_IMAGE_TAG to be set to the tag of the gloo-ee-envoy-wrapper Docker image you wish to run
.PHONY: run-e2e-tests
run-e2e-tests: TEST_PKG = ./test/e2e/ ## Run the in memory Envoy e2e tests (ENVOY_IMAGE_TAG)
run-e2e-tests: test

#----------------------------------------------------------------------------------
# Clean
#----------------------------------------------------------------------------------

# Important to clean before pushing new releases. Dockerfiles and binaries may not update properly
.PHONY: clean
clean: clean-artifacts
	rm -rf $(TEST_ASSET_DIR)
	rm -rf $(APISERVER_UI_DIR)/build
	rm -rf $(ROOTDIR)/vendor_any
	git clean -xdf install

.PHONY: clean-artifacts
clean-artifacts:
	rm -rf $(OUTPUT_DIR)

.PHONY: clean-vendor-any
clean-vendor-any:
	rm -rf vendor_any

.PHONY: clean-generated-protos
clean-generated-protos:
	rm -rf $(ROOTDIR)/projects/apiserver/api/fed.rpc/v1/*resources.proto
	rm -rf $(ROOTDIR)/projects/apiserver/api/rpc.edge.gloo/v1/*resources.proto

# Clean
.PHONY: clean-fed
clean-fed: clean-generated-protos
	rm -rf $(ROOTDIR)/vendor_any
	rm -rf $(ROOTDIR)/projects/gloo/pkg/api
	rm -rf $(ROOTDIR)/projects/gloo-fed/pkg/api
	rm -rf $(ROOTDIR)/projects/apiserver/pkg/api
	rm -rf $(ROOTDIR)/projects/glooctl-plugins/fed/pkg/api
	rm -rf $(ROOTDIR)/projects/apiserver/server/services/single_cluster_resource_handler/*


#----------------------------------------------------------------------------------
# Generated Code
#----------------------------------------------------------------------------------
PROTOC_IMPORT_PATH:=$(ROOTDIR)/vendor_any

# https://github.com/solo-io/skv2/blob/2c57f51fd5f459d4895c523134cd4e90bf54eb15/codegen/collector/compiler.go#L86
# MAX_CONCURRENT_PROTOCS sets the upper limit for the number of concurrent `protoc` processes executed during codegen
# When this is not set, the host machine executing codegen may run out of available file descriptors
MAX_CONCURRENT_PROTOCS ?= 10

.PHONY: generate-all
generate-all: generated-code generate-gloo-fed generate-helm-docs build-stitching-bundles
generate-all: fmt
generate-all:
	go mod tidy

GLOO_VERSION=$(shell echo $(shell go list -m github.com/solo-io/gloo) | cut -d' ' -f2)

# makes sure you are running codegen with the correct Go version
.PHONY: check-go-version
check-go-version:
	./ci/check-go-version.sh

.PHONY: check-solo-apis
check-solo-apis:
ifeq ($(GLOO_BRANCH_BUILD),)
	# Ensure that the gloo and solo-apis dependencies are in lockstep
	go get github.com/solo-io/solo-apis@gloo-$(GLOO_VERSION)
endif

.PHONY: check-envoy-version
check-envoy-version:
	./ci/check-envoy-version.sh $(ENVOY_GLOO_IMAGE_VERSION)

.PHONY: generated-code
generated-code: check-go-version clean-vendor-any update-licenses ## Evaluate go generate
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore GO111MODULE=on CGO_ENABLED=1 go generate ./...
	ci/check-protoc.sh

.PHONY: generate-gloo-fed
generate-gloo-fed: generate-gloo-fed-code generated-gloo-fed-ui

# Generated Code - Required to update Codgen Templates
.PHONY: generate-gloo-fed-code
generate-gloo-fed-code: clean-fed
	go run $(ROOTDIR)/projects/gloo-fed/generate.go # Generates clients, controllers, etc
	$(ROOTDIR)/projects/gloo-fed/ci/hack-fix-marshal.sh # TODO: figure out a more permanent way to deal with this
	go run projects/gloo-fed/generate.go -apiserver # Generates apiserver protos into go code
	go generate $(ROOTDIR)/projects/... # Generates mocks

.PHONY: generate-helm-docs
generate-helm-docs:
	go run $(ROOTDIR)/install/helm/gloo-ee/generate.go $(VERSION) --generate-helm-docs $(USE_DIGESTS) # Generate Helm Documentation


.PHONY: generate-extauth-test-plugins
generate-extauth-test-plugins:
	go generate $(ROOTDIR)/test/extauth/plugins/... $(ROOTDIR)/projects/extauth/plugins/...

#################
#     Build     #
#################

#----------------------------------------------------------------------------------
# Gloo Fed
#----------------------------------------------------------------------------------

GLOO_FED_DIR=$(ROOTDIR)/projects/gloo-fed
GLOO_FED_OUT_DIR=$(OUTPUT_DIR)/gloo-fed
GLOO_FED_SOURCES=$(shell find $(GLOO_FED_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(GLOO_FED_OUT_DIR)/Dockerfile.build: $(GLOO_FED_DIR)/Dockerfile
	mkdir -p $(GLOO_FED_OUT_DIR)
	cp $< $@

$(GLOO_FED_OUT_DIR)/.gloo-fed-ee-docker-build: $(GLOO_FED_SOURCES) $(GLOO_FED_OUT_DIR)/Dockerfile.build
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-fed-ee-build-container:$(VERSION) \
		-f $(GLOO_FED_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg GITHUB_TOKEN \
		--build-arg DOCKER_GOARCH=$(DOCKER_GOARCH) \
		.
	touch $@

$(GLOO_FED_OUT_DIR)/gloo-fed-linux-$(DOCKER_GOARCH): $(GLOO_FED_OUT_DIR)/.gloo-fed-ee-docker-build
	docker create -ti --name gloo-fed-temp-container $(IMAGE_REGISTRY)/gloo-fed-ee-build-container:$(VERSION) bash
	docker cp gloo-fed-temp-container:/gloo-fed-linux-$(DOCKER_GOARCH) $(GLOO_FED_OUT_DIR)/gloo-fed-linux-$(DOCKER_GOARCH)
	docker rm -f gloo-fed-temp-container


.PHONY: gloo-fed
gloo-fed: $(OUTPUT_DIR)/gloo-fed-linux-$(DOCKER_GOARCH)

.PHONY: gloo-fed-docker
gloo-fed-docker: $(GLOO_FED_OUT_DIR)/gloo-fed-linux-$(DOCKER_GOARCH)
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-fed:$(VERSION) $(DOCKER_BUILD_ARGS) $(GLOO_FED_OUT_DIR) -f $(GLOO_FED_DIR)/cmd/Dockerfile --build-arg GOARCH=$(DOCKER_GOARCH);

#----------------------------------------------------------------------------------
# Gloo Federation Projects
#----------------------------------------------------------------------------------

.PHONY: gloofed-docker
gloofed-docker: gloo-fed-docker gloo-fed-rbac-validating-webhook-docker gloo-fed-apiserver-docker gloo-fed-apiserver-envoy-docker gloo-federation-console-docker

.PHONY: gloofed-load-kind-images
gloofed-load-kind-images: kind-load-gloo-fed kind-load-gloo-fed-rbac-validating-webhook kind-load-gloo-fed-apiserver kind-load-gloo-fed-apiserver-envoy kind-load-ui

.PHONY: remove-all-gloofed-images
remove-all-gloofed-images: remove-gloofed-ui-images remove-gloofed-controller-images

.PHONY: remove-gloofed-ui-images
remove-gloofed-ui-images:
	docker image rm $(IMAGE_REGISTRY)/gloo-fed-apiserver:$(VERSION)
	docker image rm $(IMAGE_REGISTRY)/gloo-fed-apiserver-envoy:$(VERSION)
	docker image rm $(IMAGE_REGISTRY)/gloo-federation-console:$(VERSION)

.PHONY: remove-gloofed-controller-images
remove-gloofed-controller-images:
	docker image rm $(IMAGE_REGISTRY)/gloo-fed:$(VERSION)
	docker image rm $(IMAGE_REGISTRY)/gloo-fed-rbac-validating-webhook:$(VERSION)

#----------------------------------------------------------------------------------
# Gloo Fed Apiserver
#----------------------------------------------------------------------------------
GLOO_FED_APISERVER_DIR=$(ROOTDIR)/projects/apiserver
GLOO_FED_APISERVER_OUT_DIR=$(OUTPUT_DIR)/apiserver

# proto sources
APISERVER_DIR=$(ROOTDIR)/projects/apiserver/api/
APISERVER_SOURCES=$(shell find $(APISERVER_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(GLOO_FED_APISERVER_OUT_DIR)/Dockerfile.build: $(GLOO_FED_APISERVER_DIR)/Dockerfile
	mkdir -p $(GLOO_FED_APISERVER_OUT_DIR)
	cp $< $@

# the executable outputs as amd64 only because it is placed in an image that is amd64
$(GLOO_FED_APISERVER_OUT_DIR)/.gloo-fed-apiserver-docker-build: $(GLOO_FED_SOURCES) $(GLOO_FED_APISERVER_OUT_DIR)/Dockerfile.build
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-fed-apiserver-build-container:$(VERSION) \
		-f $(GLOO_FED_APISERVER_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LD_STATIC_LINKING_FLAGS) \
		--build-arg GITHUB_TOKEN \
		--build-arg DOCKER_GOARCH=amd64 \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(GLOO_FED_APISERVER_OUT_DIR)/gloo-fed-apiserver-linux-amd64: $(GLOO_FED_APISERVER_OUT_DIR)/.gloo-fed-apiserver-docker-build
	docker create -ti --name gloo-fed-apiserver-temp-container $(IMAGE_REGISTRY)/gloo-fed-apiserver-build-container:$(VERSION) bash
	docker cp gloo-fed-apiserver-temp-container:/gloo-fed-apiserver-linux-amd64 $(GLOO_FED_APISERVER_OUT_DIR)/gloo-fed-apiserver-linux-amd64
	docker rm -f gloo-fed-apiserver-temp-container

.PHONY: gloo-fed-apiserver
gloo-fed-apiserver: $(GLOO_FED_APISERVER_OUT_DIR)/gloo-fed-apiserver-linux-amd64

.PHONY: gloo-fed-apiserver-docker
gloo-fed-apiserver-docker: $(GLOO_FED_APISERVER_OUT_DIR)/gloo-fed-apiserver-linux-amd64
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-fed-apiserver:$(VERSION) \
		$(DOCKER_GO_AMD_64_ARGS)  \
		$(GLOO_FED_APISERVER_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) \
		-f $(GLOO_FED_APISERVER_DIR)/cmd/Dockerfile --build-arg GOARCH=amd64;

#----------------------------------------------------------------------------------
# apiserver-envoy
#----------------------------------------------------------------------------------
CONFIG_YAML=cfg.yaml

GLOO_FED_APISERVER_ENVOY_DIR=$(ROOTDIR)/projects/apiserver/apiserver-envoy

.PHONY: gloo-fed-apiserver-envoy-docker
gloo-fed-apiserver-envoy-docker:
	cp $(GLOO_FED_APISERVER_ENVOY_DIR)/$(CONFIG_YAML) $(OUTPUT_DIR)/$(CONFIG_YAML)
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-fed-apiserver-envoy:$(VERSION) $(DOCKER_BUILD_ARGS) $(OUTPUT_DIR) -f $(GLOO_FED_APISERVER_ENVOY_DIR)/Dockerfile;

#----------------------------------------------------------------------------------
# helpers for local testing
#----------------------------------------------------------------------------------
GRPC_PORT=10101
CONFIG_YAML=cfg.yaml

# use this target to run the UI locally
.PHONY: run-ui-local
run-ui-local:
	make -j 3 run-apiserver run-envoy run-single-cluster-ui

.PHONY: install-graphql-js
install-graphql-js:
	yarn --cwd projects/gloo/pkg/plugins/graphql/v8go install

.PHONY: build-stitching-bundles
build-stitching-bundles:
	esbuild projects/gloo/pkg/plugins/graphql/v8go/schema-diff.js --bundle --outfile=projects/gloo/pkg/plugins/graphql/v8go/schema-diff_bundled.js --log-limit=1 --minify
	esbuild projects/gloo/pkg/plugins/graphql/v8go/stitching.js --bundle --outfile=projects/gloo/pkg/plugins/graphql/v8go/stitching_bundled.js --log-limit=1 --minify

.PHONY: run-apiserver
run-apiserver: checkprogram-protoc install-graphql-js build-stitching-bundles checkenv-GLOO_LICENSE_KEY
# Todo: This should check that /etc/hosts includes the following line:
# 127.0.0.1 docker.internal
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore \
	GRPC_PORT=$(GRPC_PORT) \
	POD_NAMESPACE=gloo-system \
	$(GO_BUILD_FLAGS) go run projects/apiserver/cmd/main.go

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

$(OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-$(DOCKER_GOARCH): $(GLOO_FED_RBAC_WEBHOOK_SOURCES)
	CGO_ENABLED=0 GOARCH=$(DOCKER_GOARCH) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/main.go

.PHONY: gloo-fed-rbac-validating-webhook
gloo-fed-rbac-validating-webhook: $(OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-$(DOCKER_GOARCH)

.PHONY: gloo-fed-rbac-validating-webhook-docker
gloo-fed-rbac-validating-webhook-docker: $(OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-$(DOCKER_GOARCH)
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-fed-rbac-validating-webhook:$(VERSION) $(DOCKER_BUILD_ARGS) $(OUTPUT_DIR) -f $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/Dockerfile --build-arg GOARCH=$(DOCKER_GOARCH);

#----------------------------------------------------------------------------------
# ApiServer gRPC Code Generation
#----------------------------------------------------------------------------------

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
generated-gloo-fed-ui: update-gloo-fed-ui-deps generated-gloo-fed-ui-deps generated-graphqlschema-json-descriptor
	mkdir -p $(APISERVER_UI_DIR)/pkg/api/fed.rpc/v1
	mkdir -p $(APISERVER_UI_DIR)/pkg/api/rpc.edge.gloo/v1
	./ci/fix-ui-gen.sh

# Generate json descriptor used to convert between graphql protobuf messages and json
.PHONY: generated-graphqlschema-json-descriptor
generated-graphqlschema-json-descriptor:
	yarn --cwd projects/ui pbjs -t json -o src/Components/Features/Graphql/source-data/graphql.json \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql.proto \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext/external/google/protobuf/struct.proto \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext/external/google/protobuf/wrappers.proto \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/protoc-gen-ext/external/google/protobuf/duration.proto \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-kit/api/v1/ref.proto

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
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl/*.proto

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
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/*.proto

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
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/*.proto

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
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-projects/projects/gloo/api/enterprise/*/*/*.proto

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
build-ui: update-ui-deps
	yarn --cwd $(APISERVER_UI_DIR) build

.PHONY: gloo-federation-console-docker
gloo-federation-console-docker: build-ui
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-federation-console:$(VERSION) $(DOCKER_BUILD_ARGS) $(APISERVER_UI_DIR) -f $(APISERVER_UI_DIR)/Dockerfile

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
	docker buildx build --load -t $(IMAGE_REGISTRY)/rate-limit-ee-build-container:$(VERSION) \
		-f $(RATELIMIT_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_ALPINE_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg USE_APK=true \
		--build-arg GITHUB_TOKEN \
		$(DOCKER_BUILD_ARGS) \
		--build-arg DOCKER_CGO_ENABLED=$(GO_ENABLE_CGO) \
		.
	touch $@


# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH): $(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker-build
	docker create -ti --name rate-limit-temp-container $(IMAGE_REGISTRY)/rate-limit-ee-build-container:$(VERSION) bash
	docker cp rate-limit-temp-container:/rate-limit-linux-$(DOCKER_GOARCH) $(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH)
	docker rm -f rate-limit-temp-container

.PHONY: rate-limit
rate-limit: $(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH)

$(RATELIMIT_OUT_DIR)/Dockerfile: $(RATELIMIT_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: rate-limit-ee-docker
rate-limit-ee-docker: $(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker

$(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker: $(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH) $(RATELIMIT_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/rate-limit-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(call get_test_tag_option,rate-limit-ee) $(RATELIMIT_OUT_DIR)
	touch $@

#----------------------------------------------------------------------------------
# RateLimit-fips
#----------------------------------------------------------------------------------

RATELIMIT_FIPS_OUT_DIR=$(OUTPUT_DIR)/rate-limit-fips
_ := $(shell mkdir -p $(RATELIMIT_FIPS_OUT_DIR))

$(RATELIMIT_FIPS_OUT_DIR)/Dockerfile.build: $(RATELIMIT_DIR)/Dockerfile
	cp $< $@

# GO_BORING_ARGS is set to amd64 currently this is because BORING will only work on amd64
# https://go.googlesource.com/go/+/refs/heads/dev.boringcrypto.go1.12/misc/boring/
$(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker-build: $(RATELIMIT_SOURCES) $(RATELIMIT_FIPS_OUT_DIR)/Dockerfile.build
	docker buildx build --load -t $(IMAGE_REGISTRY)/rate-limit-ee-build-container-fips:$(VERSION) \
		-f $(RATELIMIT_FIPS_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_ALPINE_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg GITHUB_TOKEN \
		--build-arg USE_APK=true \
		$(DOCKER_GO_BORING_ARGS) \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64: $(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker-build
	docker create -ti --name rate-limit-temp-container $(IMAGE_REGISTRY)/rate-limit-ee-build-container-fips:$(VERSION) bash
	docker cp rate-limit-temp-container:/rate-limit-linux-amd64 $(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64
	docker rm -f rate-limit-temp-container

.PHONY: rate-limit-fips
rate-limit-fips: $(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64

$(RATELIMIT_FIPS_OUT_DIR)/Dockerfile: $(RATELIMIT_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: rate-limit-ee-fips-docker
rate-limit-ee-fips-docker: $(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker

$(RATELIMIT_FIPS_OUT_DIR)/.rate-limit-ee-docker: $(RATELIMIT_FIPS_OUT_DIR)/rate-limit-linux-amd64 $(RATELIMIT_FIPS_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/rate-limit-ee-fips:$(VERSION) $(DOCKER_GO_BORING_ARGS) $(call get_test_tag_option,rate-limit-ee-fips) $(RATELIMIT_FIPS_OUT_DIR) \
       --build-arg EXTRA_PACKAGES=libc6-compat
	touch $@

#----------------------------------------------------------------------------------
# ExtAuth
#----------------------------------------------------------------------------------
# When referencing files within a Docker build context, we use relative paths
DOCKER_OUTPUT_DIR := ./_output

EXTAUTH_DIR=projects/extauth
EXTAUTH_SOURCES=$(shell find $(EXTAUTH_DIR) -name "*.go" | grep -v test | grep -v generated.go)
EXTAUTH_OUT_DIR=$(OUTPUT_DIR)/extauth
DOCKER_EXTAUTH_OUT_DIR=$(DOCKER_OUTPUT_DIR)/extauth
_ := $(shell mkdir -p $(EXTAUTH_OUT_DIR))

$(EXTAUTH_OUT_DIR)/Dockerfile.build: $(EXTAUTH_DIR)/Dockerfile
	cp $< $@

$(EXTAUTH_OUT_DIR)/Dockerfile: $(EXTAUTH_DIR)/cmd/Dockerfile
	cp $< $@

$(EXTAUTH_OUT_DIR)/.extauth-ee-docker-build: $(EXTAUTH_SOURCES) $(EXTAUTH_OUT_DIR)/Dockerfile.build
	docker buildx build --load -t $(IMAGE_REGISTRY)/extauth-ee-build-container:$(VERSION) \
		-f $(EXTAUTH_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_ALPINE_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg USE_APK=true \
		--build-arg GITHUB_TOKEN \
		$(DOCKER_BUILD_ARGS) \
		--build-arg DOCKER_CGO_ENABLED=$(GO_ENABLE_CGO) \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
$(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH): $(EXTAUTH_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name extauth-temp-container $(IMAGE_REGISTRY)/extauth-ee-build-container:$(VERSION) bash
	docker cp extauth-temp-container:/extauth-linux-$(DOCKER_GOARCH) $(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH)
	docker rm -f extauth-temp-container

# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
# we need to be able to pass a variable to the extauth-ee-docker buildx build --load for Docker_CGO_ENABLED
$(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64: $(EXTAUTH_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name verify-plugins-temp-container $(IMAGE_REGISTRY)/extauth-ee-build-container:$(VERSION) bash
	docker cp verify-plugins-temp-container:/verify-plugins-linux-amd64 $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64
	docker rm -f verify-plugins-temp-container

# Build extauth binaries
.PHONY: extauth
extauth: $(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH) $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64

# Build ext-auth-plugins docker image (Cannot be built at all on Apple Silicon)
.PHONY: ext-auth-plugins-docker
ext-auth-plugins-docker: $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64
	docker buildx build --load -t $(IMAGE_REGISTRY)/ext-auth-plugins:$(VERSION) -f projects/extauth/plugins/Dockerfile \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_ALPINE_IMAGE_NAME) \
		--build-arg GC_FLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg VERIFY_SCRIPT=$(DOCKER_EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 \
		--build-arg GITHUB_TOKEN \
		--build-arg USE_APK=true \
		$(DOCKER_BUILD_ARGS) \
		.

# Build extauth server docker image
.PHONY: extauth-ee-docker
extauth-ee-docker: $(EXTAUTH_OUT_DIR)/.extauth-ee-docker

$(EXTAUTH_OUT_DIR)/.extauth-ee-docker: $(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH) $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 $(EXTAUTH_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/extauth-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(call get_test_tag_option,extauth-ee) $(EXTAUTH_OUT_DIR)
	touch $@

#----------------------------------------------------------------------------------
# ExtAuth-fips
#----------------------------------------------------------------------------------

EXTAUTH_FIPS_OUT_DIR=$(OUTPUT_DIR)/extauth_fips
DOCKER_EXTAUTH_FIPS_OUT_DIR=$(DOCKER_OUTPUT_DIR)/extauth_fips
_ := $(shell mkdir -p $(EXTAUTH_FIPS_OUT_DIR))

$(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build: $(EXTAUTH_DIR)/Dockerfile
	cp $< $@

$(EXTAUTH_FIPS_OUT_DIR)/Dockerfile: $(EXTAUTH_DIR)/cmd/Dockerfile
	cp $< $@

# GO_BORING_ARGS is set to amd64 currently this is because BORING will only work on amd64
# https://go.googlesource.com/go/+/refs/heads/dev.boringcrypto.go1.12/misc/boring/
$(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker-build: $(EXTAUTH_SOURCES) $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build
	docker buildx build --load -t $(IMAGE_REGISTRY)/extauth-ee-build-container-fips:$(VERSION) \
		-f $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_ALPINE_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg GITHUB_TOKEN \
		--build-arg USE_APK=true \
		$(DOCKER_GO_BORING_ARGS) \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64: $(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name extauth-temp-container $(IMAGE_REGISTRY)/extauth-ee-build-container-fips:$(VERSION) bash
	docker cp extauth-temp-container:/extauth-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64
	docker rm -f extauth-temp-container

# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64: $(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name verify-plugins-temp-container $(IMAGE_REGISTRY)/extauth-ee-build-container-fips:$(VERSION) bash
	docker cp verify-plugins-temp-container:/verify-plugins-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64
	docker rm -f verify-plugins-temp-container

# Build extauth binaries
.PHONY: extauth-fips
extauth-fips: $(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64

# Build ext-auth-plugins docker image
# NOTE: ext-auth-plugin will not build on arm64 machines
.PHONY: ext-auth-plugins-fips-docker
ext-auth-plugins-fips-docker: $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64
	docker buildx build --load -t $(IMAGE_REGISTRY)/ext-auth-plugins-fips:$(VERSION) -f projects/extauth/plugins/Dockerfile \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_ALPINE_IMAGE_NAME) \
		--build-arg GC_FLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg VERIFY_SCRIPT=$(DOCKER_EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 \
		--build-arg GITHUB_TOKEN \
		--build-arg USE_APK=true \
		 $(DOCKER_GO_BORING_ARGS) \
		.

# Build extauth server docker image
.PHONY: extauth-ee-fips-docker
extauth-ee-fips-docker: $(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker

$(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker: $(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/extauth-ee-fips:$(VERSION) $(DOCKER_GO_BORING_ARGS) $(call get_test_tag_option,extauth-ee-fips) $(EXTAUTH_FIPS_OUT_DIR) \
		--build-arg EXTRA_PACKAGES=libc6-compat
	touch $@

#----------------------------------------------------------------------------------
# Observability
#----------------------------------------------------------------------------------

OBSERVABILITY_DIR=projects/observability
OBSERVABILITY_SOURCES=$(shell find $(OBSERVABILITY_DIR) -name "*.go" | grep -v test | grep -v generated.go)
OBS_OUT_DIR=$(OUTPUT_DIR)/observability

$(OBS_OUT_DIR)/observability-linux-$(DOCKER_GOARCH): $(OBSERVABILITY_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(OBSERVABILITY_DIR)/cmd/main.go

.PHONY: observability
observability: $(OBS_OUT_DIR)/observability-linux-$(DOCKER_GOARCH)

$(OBS_OUT_DIR)/Dockerfile: $(OBSERVABILITY_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: observability-ee-docker
observability-ee-docker: $(OBS_OUT_DIR)/.observability-ee-docker

$(OBS_OUT_DIR)/.observability-ee-docker: $(OBS_OUT_DIR)/observability-linux-$(DOCKER_GOARCH) $(OBS_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/observability-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(call get_test_tag_option,observability-ee) $(OBS_OUT_DIR)
	touch $@

#----------------------------------------------------------------------------------
# Caching
#----------------------------------------------------------------------------------

CACHING_DIR=projects/caching
CACHING_SOURCES=$(shell find $(CACHING_DIR) -name "*.go" | grep -v test | grep -v generated.go)
CACHE_OUT_DIR=$(OUTPUT_DIR)/caching

$(CACHE_OUT_DIR)/caching-linux-amd64: $(CACHING_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CACHING_DIR)/cmd/main.go

.PHONY: caching
caching: $(CACHE_OUT_DIR)/caching-linux-amd64

$(CACHE_OUT_DIR)/Dockerfile: $(CACHING_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: caching-ee-docker
caching-ee-docker: $(CACHE_OUT_DIR)/.caching-ee-docker

$(CACHE_OUT_DIR)/.caching-ee-docker: $(CACHE_OUT_DIR)/caching-linux-amd64 $(CACHE_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/caching-ee:$(VERSION) $(call get_test_tag_option,caching-ee) $(CACHE_OUT_DIR) $(PLATFORM)
	touch $@

#----------------------------------------------------------------------------------
# Gloo
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)
GLOO_OUT_DIR=$(OUTPUT_DIR)/gloo




$(GLOO_OUT_DIR)/Dockerfile.build: $(GLOO_DIR)/Dockerfile
	mkdir -p $(GLOO_OUT_DIR)
	cp $< $@

# the executable outputs as amd64 only because it is placed in an image that is amd64
$(GLOO_OUT_DIR)/.gloo-ee-docker-build: install-node-packages $(GLOO_SOURCES) $(GLOO_OUT_DIR)/Dockerfile.build
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-ee-build-container:$(VERSION) \
		-f $(GLOO_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LD_STATIC_LINKING_FLAGS) \
		--build-arg USE_APK=true \
		--build-arg GITHUB_TOKEN \
		${DOCKER_GO_AMD_64_ARGS} \
		.
	touch $@

# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(GLOO_OUT_DIR)/gloo-linux-amd64: $(GLOO_OUT_DIR)/.gloo-ee-docker-build
	docker create -ti --name gloo-temp-container $(IMAGE_REGISTRY)/gloo-ee-build-container:$(VERSION) bash
	docker cp gloo-temp-container:/gloo-linux-amd64 $(GLOO_OUT_DIR)/gloo-linux-amd64
	docker rm -f gloo-temp-container

.PHONY: gloo
gloo: $(GLOO_OUT_DIR)/gloo-linux-amd64

$(GLOO_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@


.PHONY: gloo-ee-docker
gloo-ee-docker: $(GLOO_OUT_DIR)/.gloo-ee-docker

$(GLOO_OUT_DIR)/.gloo-ee-docker: $(GLOO_OUT_DIR)/gloo-linux-amd64 $(GLOO_OUT_DIR)/Dockerfile
	docker buildx build --load $(call get_test_tag_option,gloo-ee) $(GLOO_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) \
		$(DOCKER_GO_AMD_64_ARGS) \
		-t $(IMAGE_REGISTRY)/gloo-ee:$(VERSION)
	touch $@

gloo-ee-docker-dev: $(GLOO_OUT_DIR)/gloo-linux-amd64 $(GLOO_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(GLOO_OUT_DIR) --no-cache
	touch $@

#----------------------------------------------------------------------------------
# Gloo with race detection enabled.
# This is intended to be used to aid in local debugging by swapping out this image in a running gloo instance
#----------------------------------------------------------------------------------
GLOO_RACE_OUT_DIR=$(OUTPUT_DIR)/gloo-race

$(GLOO_RACE_OUT_DIR)/Dockerfile.build: $(GLOO_DIR)/Dockerfile
	mkdir -p $(GLOO_RACE_OUT_DIR)
	cp $< $@

$(GLOO_RACE_OUT_DIR)/.gloo-ee-race-docker-build: $(GLOO_SOURCES) $(GLOO_RACE_OUT_DIR)/Dockerfile.build
	docker build -t $(IMAGE_REGISTRY)/gloo-race-ee-build-container:$(VERSION) \
		-f $(GLOO_RACE_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LD_STATIC_LINKING_FLAGS) \
		--build-arg USE_APK=true \
		--build-arg GITHUB_TOKEN \
		--build-arg RACE=-race \
		$(DOCKER_BUILD_ARGS) \
		.
	touch $@
# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH): $(GLOO_RACE_OUT_DIR)/.gloo-ee-race-docker-build
	docker create -ti --name gloo-race-temp-container $(IMAGE_REGISTRY)/gloo-race-ee-build-container:$(VERSION) bash
	docker cp gloo-race-temp-container:/gloo-linux-$(DOCKER_GOARCH) $(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH)
	docker rm -f gloo-race-temp-container

.PHONY: gloo-race
gloo-race: $(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH)

$(GLOO_RACE_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: gloo-ee-race-docker
gloo-ee-race-docker: $(GLOO_RACE_OUT_DIR)/.gloo-ee-race-docker
$(GLOO_RACE_OUT_DIR)/.gloo-ee-race-docker: $(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH) $(GLOO_RACE_OUT_DIR)/Dockerfile
	docker build $(call get_test_tag_option,gloo-ee) $(GLOO_RACE_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) $(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_REGISTRY)/gloo-ee:$(VERSION)-race
	touch $@
#----------------------------------------------------------------------------------
# Gloo with FIPS Envoy
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)
GLOO_FIPS_OUT_DIR=$(OUTPUT_DIR)/gloo-fips

$(GLOO_FIPS_OUT_DIR)/Dockerfile.build: $(GLOO_DIR)/Dockerfile
	mkdir -p $(GLOO_FIPS_OUT_DIR)
	cp $< $@

$(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker-build: $(GLOO_SOURCES) $(GLOO_FIPS_OUT_DIR)/Dockerfile.build
	docker build -t $(IMAGE_REGISTRY)/gloo-ee-fips-build-container:$(VERSION) \
		-f $(GLOO_FIPS_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_IMAGE_NAME) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LD_STATIC_LINKING_FLAGS) \
		--build-arg GITHUB_TOKEN \
		--build-arg USE_APK=true \
		$(DOCKER_GO_BORING_ARGS) \
		.
	touch $@


$(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64: $(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker-build
	docker create -ti --name gloo-fips-temp-container $(IMAGE_REGISTRY)/gloo-ee-fips-build-container:$(VERSION) bash
	docker cp gloo-fips-temp-container:/gloo-linux-amd64 $(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64
	docker rm -f gloo-fips-temp-container

.PHONY: gloo-fips
gloo-fips: $(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64

$(GLOO_FIPS_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@


.PHONY: gloo-ee-fips-docker
gloo-ee-fips-docker: $(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker

$(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker: $(GLOO_FIPS_OUT_DIR)/gloo-linux-amd64 $(GLOO_FIPS_OUT_DIR)/Dockerfile
	docker buildx build --load $(call get_test_tag_option,gloo-ee) $(GLOO_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_IMAGE) \
		$(DOCKER_GO_BORING_ARGS) \
		-t $(IMAGE_REGISTRY)/gloo-ee-fips:$(VERSION)
	touch $@

gloo-ee-fips-docker-dev: $(GLOO_FIPS_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH) $(GLOO_FIPS_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REGISTRY)/gloo-ee-fips:$(VERSION) $(DOCKER_BUILD_ARGS) $(GLOO_FIPS_OUT_DIR) --no-cache
	touch $@
#----------------------------------------------------------------------------------
# discovery (enterprise)
#----------------------------------------------------------------------------------

DISCOVERY_DIR=projects/discovery
DISCOVERY_SOURCES=$(shell find $(DISCOVERY_DIR) -name "*.go" | grep -v test | grep -v generated.go)
DISCOVERY_OUTPUT_DIR=$(OUTPUT_DIR)/$(DISCOVERY_DIR)

$(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-$(DOCKER_GOARCH): $(DISCOVERY_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(DISCOVERY_DIR)/cmd/main.go

.PHONY: discovery-ee
discovery-ee: $(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-$(DOCKER_GOARCH)
$(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery: $(DISCOVERY_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: discovery-ee-docker
discovery-ee-docker: $(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-$(DOCKER_GOARCH) $(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery
	docker buildx build --load $(DISCOVERY_OUTPUT_DIR) -f $(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery \
		$(DOCKER_BUILD_ARGS) -t $(IMAGE_REGISTRY)/discovery-ee:$(VERSION) $(QUAY_EXPIRATION_LABEL)

#----------------------------------------------------------------------------------
# glooctl
#----------------------------------------------------------------------------------

CLI_DIR=projects/gloo/cli
$(OUTPUT_DIR)/glooctl: $(SOURCES)
	GO111MODULE=on go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

$(OUTPUT_DIR)/glooctl-linux-$(GOARCH): $(SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

$(OUTPUT_DIR)/glooctl-darwin-$(GOARCH): $(SOURCES)
	$(GO_BUILD_FLAGS) GOOS=darwin GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

$(OUTPUT_DIR)/glooctl-windows-amd64.exe: $(SOURCES)
	$(GO_BUILD_FLAGS) GOOS=windows go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/main.go

# NOTE although it says amd64 it is determined by the architecture of the machine building it
# this is because of the dependency on github.com/solo-io/k8s-utils@v0.1.0/testutils/helper/install.go
.PHONY: glooctl
glooctl: $(OUTPUT_DIR)/glooctl
.PHONY: glooctl-linux-$(GOARCH)
glooctl-linux-$(GOARCH): $(OUTPUT_DIR)/glooctl-linux-$(GOARCH)
.PHONY: glooctl-darwin-$(GOARCH)
glooctl-darwin-$(GOARCH): $(OUTPUT_DIR)/glooctl-darwin-$(GOARCH)
.PHONY: glooctl-windows
glooctl-windows: $(OUTPUT_DIR)/glooctl-windows-amd64.exe

.PHONY: build-cli-local
build-cli-local: glooctl-$(GOOS)-$(GOARCH) ## Build the CLI according to your local GOOS and GOARCH

.PHONY: build-cli
build-cli: glooctl-linux-$(GOARCH) glooctl-darwin-$(GOARCH) glooctl-windows

#----------------------------------------------------------------------------------
# Glooctl Plugins
#----------------------------------------------------------------------------------

# Include helm makefile so its targets can be ran from the root of this repo
include $(ROOTDIR)/projects/glooctl-plugins/plugins.mk

#----------------------------------------------------------------------------------
# Envoy init (BASE/SIDECAR)
#----------------------------------------------------------------------------------

ENVOYINIT_DIR=projects/envoyinit
ENVOYINIT_SOURCES=$(shell find $(ENVOYINIT_DIR) -name "*.go" | grep -v test | grep -v generated.go)
ENVOYINIT_OUT_DIR=$(OUTPUT_DIR)/$(ENVOYINIT_DIR)

$(ENVOYINIT_OUT_DIR)/envoyinit-linux-$(DOCKER_GOARCH): $(ENVOYINIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go

.PHONY: envoyinit
envoyinit: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-$(DOCKER_GOARCH)

$(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit: $(ENVOYINIT_DIR)/Dockerfile.envoyinit
	cp $< $@

$(ENVOYINIT_OUT_DIR)/docker-entrypoint.sh: $(ENVOYINIT_DIR)/docker-entrypoint.sh
	cp $< $@

.PHONY: gloo-ee-envoy-wrapper-docker
gloo-ee-envoy-wrapper-docker: $(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-docker
$(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-docker: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-$(DOCKER_GOARCH) $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_OUT_DIR)/docker-entrypoint.sh
	docker buildx build --load $(call get_test_tag_option,gloo-ee-envoy-wrapper) $(ENVOYINIT_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) $(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_REGISTRY)/gloo-ee-envoy-wrapper:$(VERSION) \
		-f $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit
	touch $@

.PHONY: gloo-ee-envoy-wrapper-debug-docker
gloo-ee-envoy-wrapper-debug-docker: $(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-debug-docker

$(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-debug-docker: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-$(DOCKER_GOARCH) $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_OUT_DIR)/docker-entrypoint.sh
	docker buildx build --load $(call get_test_tag_option,gloo-ee-envoy-wrapper-debug) $(ENVOYINIT_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_DEBUG_IMAGE) $(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_REGISTRY)/gloo-ee-envoy-wrapper:$(VERSION)-debug \
		-f $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit
	touch $@

#----------------------------------------------------------------------------------
# Fips Envoy init (BASE/SIDECAR)
#----------------------------------------------------------------------------------

ENVOYINIT_FIPS_OUT_DIR=$(OUTPUT_DIR)/envoyinit_fips

$(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64: $(ENVOYINIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go

.PHONY: envoyinit-fips
envoyinit-fips: $(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64
$(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit: $(ENVOYINIT_DIR)/Dockerfile.envoyinit
	cp $< $@

$(ENVOYINIT_FIPS_OUT_DIR)/docker-entrypoint.sh: $(ENVOYINIT_DIR)/docker-entrypoint.sh
	cp $< $@

.PHONY: gloo-ee-envoy-wrapper-fips-docker
gloo-ee-envoy-wrapper-fips-docker: $(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-docker

$(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-docker: $(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64 $(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_FIPS_OUT_DIR)/docker-entrypoint.sh
	docker buildx build --load $(call get_test_tag_option,gloo-ee-envoy-wrapper-fips) $(ENVOYINIT_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_IMAGE) $(DOCKER_GO_BORING_ARGS) \
		-t $(IMAGE_REGISTRY)/gloo-ee-envoy-wrapper-fips:$(VERSION) \
		-f $(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit
	touch $@

.PHONY: gloo-ee-envoy-wrapper-fips-debug-docker
gloo-ee-envoy-wrapper-fips-debug-docker: $(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-debug-docker

$(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-debug-docker: $(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64 $(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_FIPS_OUT_DIR)/docker-entrypoint.sh
	docker buildx build --load $(call get_test_tag_option,gloo-ee-envoy-wrapper-fips-debug) $(ENVOYINIT_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_DEBUG_IMAGE) $(DOCKER_GO_BORING_ARGS) \
		-t $(IMAGE_REGISTRY)/gloo-ee-envoy-wrapper-fips:$(VERSION)-debug \
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
GLOOE_HELM_BUCKET := gs://gloo-ee-helm
GLOO_FED_HELM_BUCKET := gs://gloo-fed-helm

# creates Chart.yaml, values.yaml, and requirements.yaml
USE_DIGESTS:=""
ifeq ($(RELEASE), "true")
		USE_DIGESTS="--use-digests"
endif

.PHONY: helm-template
helm-template:
	mkdir -p $(HELM_SYNC_DIR_FOR_GLOO_EE)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION) --gloo-fed-repo-override="file://$(GLOO_FED_CHART_DIR)" $(USE_DIGESTS) --gloo-repo-override=$(GLOO_REPO_OVERRIDE)

.PHONY: init-helm
init-helm: helm-template gloofed-helm-template $(OUTPUT_DIR)/.helm-initialized

$(OUTPUT_DIR)/.helm-initialized:
	helm repo add helm-hub https://charts.helm.sh/stable
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	helm repo add gloo-fed https://storage.googleapis.com/gloo-fed-helm
	helm dependency update install/helm/gloo-ee
	touch $@

.PHONY: package-gloo-edge-chart
package-gloo-edge-chart: init-helm
	helm package --destination $(HELM_SYNC_DIR_FOR_GLOO_EE) $(GLOOE_CHART_DIR)

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
	until $$(GENERATION=$$(gsutil ls -a $(GLOO_FED_HELM_BUCKET)/index.yaml | tail -1 | cut -f2 -d '#') && \
		  gsutil cp -v $(GLOO_FED_HELM_BUCKET)/index.yaml $(HELM_SYNC_DIR_GLOO_FED)/index.yaml && \
		  helm package --destination $(HELM_SYNC_DIR_GLOO_FED)/charts $(HELM_DIR)/gloo-fed >> /dev/null && \
		  helm repo index $(HELM_SYNC_DIR_GLOO_FED) --merge $(HELM_SYNC_DIR_GLOO_FED)/index.yaml && \
		  gsutil -m rsync $(HELM_SYNC_DIR_GLOO_FED)/charts $(GLOO_FED_HELM_BUCKET)/charts && \
		  gsutil -h x-goog-if-generation-match:"$$GENERATION" cp $(HELM_SYNC_DIR_GLOO_FED)/index.yaml $(GLOO_FED_HELM_BUCKET)/index.yaml); do \
	echo "Failed to upload new helm index (updated helm index since last download?). Trying again"; \
	sleep 2; \
	done
endif

#----------------------------------------------------------------------------------
# Gloo Fed Deployment Manifests / Helm
#----------------------------------------------------------------------------------
GLOO_FED_VERSION=$(VERSION)
GLOO_FED_APISERVER_VERSION=$(VERSION)
GLOO_FED_APISERVER_ENVOY_VERSION=$(VERSION)
GLOO_FEDERATION_CONSOLE_VERSION=$(VERSION)
GLOO_FED_RBAC_VALIDATING_WEBHOOK_VERSION=$(VERSION)
ifeq ($(RELEASE), "true")
		GLOO_FED_VERSION=$(VERSION)@$(shell docker manifest inspect "quay.io/solo-io/gloo-fed:$(VERSION)" -v | jq ".Descriptor.digest")
		GLOO_FED_APISERVER_VERSION=$(VERSION)@$(shell docker manifest inspect "quay.io/solo-io/gloo-fed-apiserver:$(VERSION)" -v | jq ".Descriptor.digest")
		GLOO_FED_APISERVER_ENVOY_VERSION=$(VERSION)@$(shell docker manifest inspect "quay.io/solo-io/gloo-fed-apiserver-envoy:$(VERSION)" -v | jq ".Descriptor.digest")
		GLOO_FEDERATION_CONSOLE_VERSION=$(VERSION)@$(shell docker manifest inspect "quay.io/solo-io/gloo-federation-console:$(VERSION)" -v | jq ".Descriptor.digest")
		GLOO_FED_RBAC_VALIDATING_WEBHOOK_VERSION=$(VERSION)@$(shell docker manifest inspect "quay.io/solo-io/gloo-fed-rbac-validating-webhook:$(VERSION)" -v | jq ".Descriptor.digest")
endif

# creates Chart.yaml, values.yaml, and requirements.yaml
.PHONY: gloofed-helm-template
gloofed-helm-template:
	mkdir -p $(HELM_SYNC_DIR_GLOO_FED)
	sed -e 's/%version%/'$(VERSION)'/' $(GLOO_FED_CHART_DIR)/Chart-template.yaml > $(GLOO_FED_CHART_DIR)/Chart.yaml
	sed -e 's/%gloo-fed-version%/'$(GLOO_FED_VERSION)'/'\
		-e 's/%gloo-fed-apiserver-version%/'$(GLOO_FED_APISERVER_VERSION)'/'\
		-e 's/%gloo-fed-apiserver-envoy-version%/'$(GLOO_FED_APISERVER_ENVOY_VERSION)'/'\
		-e 's/%gloo-federation-console-version%/'$(GLOO_FEDERATION_CONSOLE_VERSION)'/'\
		-e 's/%gloo-fed-rbac-validating-webhook-version%/'$(GLOO_FED_RBAC_VALIDATING_WEBHOOK_VERSION)'/'\
		$(GLOO_FED_CHART_DIR)/values-template.yaml > $(GLOO_FED_CHART_DIR)/values.yaml

.PHONY: package-gloo-fed-chart
package-gloo-fed-chart: gloofed-helm-template
	helm package --destination $(HELM_SYNC_DIR_GLOO_FED) $(GLOO_FED_CHART_DIR)

#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------

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
	echo "GO_BUILD_IMAGE=$(GOLANG_ALPINE_IMAGE_NAME)" > $@
	echo "FIPS_GO_BUILD_IMAGE=$(GOLANG_IMAGE_NAME)" >> $@
	echo "GC_FLAGS=$(GCFLAGS)" >> $@

$(DEPS_DIR)/verify-plugins-linux-amd64: $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)
	cp $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)
$(DEPS_DIR)/fips-verify-plugins-linux-amd64: $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)
	cp $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 $(DEPS_DIR)/fips-verify-plugins-linux-amd64

# Intended only to be run by CI
# Build and push docker images to the defined IMAGE_REGISTRY
.PHONY: publish-docker
ifeq ($(RELEASE), "true")
publish-docker: docker
publish-docker: docker-push
publish-docker: docker-fed
publish-docker: docker-fed-push
endif

# Intended only to be run by CI
# Re-tag docker images previously pushed to the ORIGINAL_IMAGE_REGISTRY,
# and push them to a secondary repository, defined at IMAGE_REGISTRY
.PHONY: publish-docker-retag
ifeq ($(RELEASE),"true")
publish-docker-retag: docker-retag
publish-docker-retag: docker-push
publish-docker-retag: docker-fed-retag
publish-docker-retag: docker-fed-push
endif

#----------------------------------------------------------------------------------
# Docker
#----------------------------------------------------------------------------------

# ARM64 Support in Gloo Edge Enterprise is not feature complete: https://github.com/solo-io/gloo/issues/5471
# There are certain images which cannot be built locally on a machine using an ARM processor architecture.
# We have some recipes below which attempt to build (or push) all the images, so we add the following guard:
#
# 	ifeq ($(IS_ARM_MACHINE), )
#	...
# 	endif
#
# To our understanding, the following images will not work if they are built and run locally:
#	gloo-ee race image
#	gloo-ee-envoy-wrapper debug image
#	gloo-ee-envoy-wrapper debug fips image
#	ext-auth-plugins image
#	ext-auth-plugins fips image
#
# It is our long-term goal to support these, but since they are not critical to everyday development,
# we have not completed the support. In the meantime, we are tracking these in the above list.

# Note: Order matters. We want the matcher to match the most specific (ex. `race` and `debug`) first.
docker-push-%-race:
	docker push $(IMAGE_REGISTRY)/$*:$(VERSION)-race

docker-retag-%-race:
	docker tag $(ORIGINAL_IMAGE_REGISTRY)/$*:$(VERSION)-race $(IMAGE_REGISTRY)/$*:$(VERSION)-race

docker-push-%-debug:
	docker push $(IMAGE_REGISTRY)/$*:$(VERSION)-debug

docker-retag-%-debug:
	docker tag $(ORIGINAL_IMAGE_REGISTRY)/$*:$(VERSION)-debug $(IMAGE_REGISTRY)/$*:$(VERSION)-debug

docker-push-%:
	docker push $(IMAGE_REGISTRY)/$*:$(VERSION)

docker-retag-%:
	docker tag $(ORIGINAL_IMAGE_REGISTRY)/$*:$(VERSION) $(IMAGE_REGISTRY)/$*:$(VERSION)

.PHONY: docker-info
	echo "------------- Docker Info --------------"
	docker info
	# docker system df **this can take a long time to run**
	docker image ls
	echo "----------------------------------------"

# Build Gloo Enterprise docker images using the defined IMAGE_REGISTRY, VERSION
.PHONY: docker
docker: check-go-version
docker: # Build Control Plane images
docker: gloo-ee-docker
docker: gloo-ee-fips-docker
docker: discovery-ee-docker
docker: observability-ee-docker
docker: # Build Data Plane images
docker: gloo-ee-envoy-wrapper-docker
docker: gloo-ee-envoy-wrapper-fips-docker
docker: extauth-ee-docker
docker: extauth-ee-fips-docker
docker: rate-limit-ee-docker
docker: rate-limit-ee-fips-docker
docker: caching-ee-docker
docker: # Build AMD64-only supported images
ifeq ($(IS_ARM_MACHINE), )
docker: gloo-ee-race-docker
docker: gloo-ee-envoy-wrapper-debug-docker
docker: gloo-ee-envoy-wrapper-fips-debug-docker
docker: ext-auth-plugins-docker
docker: ext-auth-plugins-fips-docker
endif

# Push docker images to the defined IMAGE_REGISTRY
.PHONY: docker-push
docker-push: # Push Control Plane images
docker-push: docker-push-gloo-ee
docker-push: docker-push-gloo-ee-fips
docker-push: docker-push-discovery-ee
docker-push: docker-push-observability-ee
docker-push: # Push Data Plane images
docker-push: docker-push-gloo-ee-envoy-wrapper
docker-push: docker-push-gloo-ee-envoy-wrapper-fips
docker-push: docker-push-extauth-ee
docker-push: docker-push-extauth-ee-fips
docker-push: docker-push-rate-limit-ee
docker-push: docker-push-rate-limit-ee-fips
docker-push: docker-push-caching-ee
docker-push: # Push AMD64-only supported images
ifeq ($(IS_ARM_MACHINE), )
docker-push: docker-push-gloo-ee-race
docker-push: docker-push-gloo-ee-envoy-wrapper-debug
docker-push: docker-push-gloo-ee-envoy-wrapper-fips-debug
docker-push: docker-push-ext-auth-plugins
docker-push: docker-push-ext-auth-plugins-fips
endif

# Re-tag docker images previously pushed to the ORIGINAL_IMAGE_REGISTRY,
# and tag them with a secondary repository, defined at IMAGE_REGISTRY
.PHONY: docker-retag
docker-retag: # Re-tag Control Plane images
docker-retag: docker-retag-gloo-ee
docker-retag: docker-retag-gloo-ee-fips
docker-retag: docker-retag-discovery-ee
docker-retag: docker-retag-observability-ee
docker-retag: # Re-tag Data Plane images
docker-retag: docker-retag-gloo-ee-envoy-wrapper
docker-retag: docker-retag-gloo-ee-envoy-wrapper-fips
docker-retag: docker-retag-extauth-ee
docker-retag: docker-retag-extauth-ee-fips
docker-retag: docker-retag-rate-limit-ee
docker-retag: docker-retag-rate-limit-ee-fips
docker-retag: docker-retag-caching-ee
docker-retag: # Re-tag AMD64-only supported images
ifeq ($(IS_ARM_MACHINE), )
docker-retag: docker-retag-gloo-ee-race
docker-retag: docker-retag-gloo-ee-envoy-wrapper-debug
docker-retag: docker-retag-gloo-ee-envoy-wrapper-fips-debug
docker-retag: docker-retag-ext-auth-plugins
docker-retag: docker-retag-ext-auth-plugins-fips
endif


# Federation Assets
# We use separate targets for Federation resources because at times we may publish these independently
# of Edge resources (when testing)
# In the future, it would be preferred if these were all treated as a single project

# Build Gloo Federation docker images using the defined IMAGE_REGISTRY, VERSION
.PHONY: docker-fed
docker-fed: gloo-fed-docker
docker-fed: gloo-fed-apiserver-docker
docker-fed: gloo-fed-apiserver-envoy-docker
docker-fed: gloo-federation-console-docker
docker-fed: gloo-fed-rbac-validating-webhook-docker

# Push Gloo Federation docker images to the defined IMAGE_REGISTRY
.PHONY: docker-fed-push
docker-fed-push: docker-push-gloo-fed
docker-fed-push: docker-push-gloo-fed-apiserver
docker-fed-push: docker-push-gloo-fed-apiserver-envoy
docker-fed-push: docker-push-gloo-federation-console
docker-fed-push: docker-push-gloo-fed-rbac-validating-webhook

# Re-tag Gloo Federation docker images previously pushed to the ORIGINAL_IMAGE_REGISTRY,
# and tag them with a secondary repository, defined at IMAGE_REGISTRY
.PHONY: docker-fed-retag
docker-fed-retag: docker-retag-gloo-fed
docker-fed-retag: docker-retag-gloo-fed-apiserver
docker-fed-retag: docker-retag-gloo-fed-apiserver-envoy
docker-fed-retag: docker-retag-gloo-federation-console
docker-fed-retag: docker-retag-gloo-fed-rbac-validating-webhook

#----------------------------------------------------------------------------------
# Build assets for Kube2e tests
#----------------------------------------------------------------------------------

CLUSTER_NAME ?= kind
INSTALL_NAMESPACE ?= gloo-system

kind-load-%:
	kind load docker-image $(IMAGE_REGISTRY)/$*:$(VERSION) --name $(CLUSTER_NAME)

kind-load-%-debug:
	kind load docker-image $(IMAGE_REGISTRY)/$*:$(VERSION)-debug --name $(CLUSTER_NAME)

# Build an image and load it into the KinD cluster
# Depends on: IMAGE_REGISTRY, VERSION, CLUSTER_NAME
# Envoy image may be specified via ENVOY_GLOO_IMAGE on the command line or at the top of this file
kind-build-and-load-%: %-docker kind-load-% ; ## Use to build specified image and load it into kind

# Reload an image in KinD
# This is useful to developers when changing a single component
# You can reload an image, which means it will be rebuilt and reloaded into the kind cluster
# using the same tag so that tests can be re-run
# Depends on: IMAGE_REGISTRY, VERSION, INSTALL_NAMESPACE , CLUSTER_NAME
# Envoy image may be specified via ENVOY_GLOO_IMAGE on the command line or at the top of this file
kind-reload-%: kind-build-and-load-% ## Use to build specified image, load it into kind, and restart its deployment
	kubectl rollout restart deployment/$* -n $(INSTALL_NAMESPACE)

# Useful utility for listing images loaded into the kind cluster
.PHONY: kind-list-images
kind-list-images:
	docker exec -ti $(CLUSTER_NAME)-control-plane crictl images | grep "solo-io"

# Useful utility for pruning images that were previously loaded into the kind cluster
.PHONY: kind-prune-images
kind-prune-images:
	docker exec -ti $(CLUSTER_NAME)-control-plane crictl rmi --prune

.PHONY: kind-build-and-load
kind-build-and-load: ## Use to build all images and load them into kind
ifeq ($(USE_FIPS),true)
kind-build-and-load: kind-build-and-load-gloo-ee-fips
kind-build-and-load: kind-build-and-load-gloo-ee-envoy-wrapper-fips
kind-build-and-load: kind-build-and-load-rate-limit-ee-fips
kind-build-and-load: kind-build-and-load-extauth-ee-fips
kind-build-and-load: kind-build-and-load-observability-ee
kind-build-and-load: kind-build-and-load-caching-ee
kind-build-and-load: kind-build-and-load-discovery-ee
ifeq  ($(IS_ARM_MACHINE), )
kind-build-and-load: kind-build-and-load-ext-auth-plugins-fips
endif # ARM support
else # non-fips support
kind-build-and-load: kind-build-and-load-gloo-ee
kind-build-and-load: kind-build-and-load-gloo-ee-envoy-wrapper
kind-build-and-load: kind-build-and-load-rate-limit-ee
kind-build-and-load: kind-build-and-load-extauth-ee
kind-build-and-load: kind-build-and-load-observability-ee
kind-build-and-load: kind-build-and-load-caching-ee
kind-build-and-load: kind-build-and-load-discovery-ee
ifeq  ($(IS_ARM_MACHINE), )
kind-build-and-load: kind-build-and-load-ext-auth-plugins
endif # ARM support
endif # non-fips support

define kind_reload_msg
The kind-reload-% targets exist in order to assist developers with the work cycle of
build->test->change->build->test. To that end, rebuilding/reloading every image, then
restarting every deployment is seldom necessary. Consider using kind-reload-% to do so
for a specific component, or kind-build-and-load to push new images for every component.
endef
export kind_reload_msg
.PHONY: kind-reload
kind-reload:
	@echo "$$kind_reload_msg"

.PHONY: build-test-chart
build-test-chart: build-test-chart-fed
	mkdir -p $(TEST_ASSET_DIR)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION) --gloo-fed-repo-override="file://$(GLOO_FED_CHART_DIR)" $(USE_DIGESTS)
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

# Exclusively useful for testing with locally modified gloo-edge-OS builds.
# Assumes that a gloo-edge-OS chart is located at ../gloo/_test/gloo-v0.0.0-dev.tgz, which
# one must have also run build-test-chart
# points towards whatever modified build is being tested.
.PHONY: build-chart-with-local-gloo-dev ## Helm chart that relies on oss chart for cross cutting deps
build-chart-with-local-gloo-dev: build-test-chart-fed
	mkdir -p $(TEST_ASSET_DIR)
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION) $(USE_DIGESTS)
	helm repo add helm-hub https://charts.helm.sh/stable
	helm repo add gloo https://storage.googleapis.com/solo-public-helm
	echo replacing gloo chart $(ls install/helm/gloo-ee/charts/gloo*) with ../gloo/_test/gloo-dev.tgz
	# TODO: If edge hits 2.0 then this needs to change
	#  we specify it this way to not delete the gloo-fed chart
	rm install/helm/gloo-ee/charts/gloo-1*
	cp ../gloo/_test/gloo-v0.0.0-dev.tgz install/helm/gloo-ee/charts/

	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-ee
	helm repo index $(TEST_ASSET_DIR)

#----------------------------------------------------------------------------------
# Build assets for Federation Kube2e tests
#----------------------------------------------------------------------------------

# Build all of the images used in the Federatation Kube2e tests
# We separate the building and loading of images into separate steps to ensure we can build
# images once and load them into multiple clusters (local and remote)
.PHONY: kind-build-federation-images
kind-build-federation-images: gloo-ee-docker
kind-build-federation-images: gloo-ee-envoy-wrapper-docker
kind-build-federation-images: discovery-ee-docker
kind-build-federation-images: gloo-fed-docker
kind-build-federation-images: gloo-fed-apiserver-envoy-docker
kind-build-federation-images: gloo-fed-rbac-validating-webhook-docker
kind-build-federation-images: gloo-fed-apiserver-docker
kind-build-federation-images: gloo-federation-console-docker

.PHONY: kind-load-federation-control-plane-images
kind-load-federation-control-plane-images: kind-load-gloo-ee
kind-load-federation-control-plane-images: kind-load-gloo-ee-envoy-wrapper
kind-load-federation-control-plane-images: kind-load-discovery-ee

.PHONY: kind-load-federation-management-plane-images
kind-load-federation-management-plane-images: kind-load-gloo-fed
kind-load-federation-management-plane-images: kind-load-gloo-fed-rbac-validating-webhook
kind-load-federation-management-plane-images: kind-load-gloo-fed-apiserver
kind-load-federation-management-plane-images: kind-load-gloo-fed-apiserver-envoy
kind-load-federation-management-plane-images: kind-load-gloo-federation-console

.PHONY: build-test-chart-fed
build-test-chart-fed: gloofed-helm-template
	mkdir -p $(TEST_ASSET_DIR)
	helm repo add helm-hub https://charts.helm.sh/stable
	helm repo add gloo-fed https://storage.googleapis.com/gloo-fed-helm
	helm dependency update install/helm/gloo-fed
	helm package --destination $(TEST_ASSET_DIR) $(HELM_DIR)/gloo-fed
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

# use `make checkprogram-PROGRAM_NAME` to check if a program exists in PATH
checkprogram-%:
	 @which $* > /dev/null || (echo "program '$*' not found in path. Do you need to install it?" && exit 1)

# use `make checkenv-ENV_VAR` to make sure ENV_VAR is set
checkenv-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

_local/redis.key:
	mkdir -p _local
	openssl req -new -x509 -newkey rsa:2048 -sha256 -nodes -keyout _local/redis.key -days 3560 -out _local/redis.crt -config - < _local/cert.conf

#----------------------------------------------------------------------------------
# Security Scanning utility
#	Scanning for solo-projects is performed by the open source Gloo Edge repo
#	These utilities make it easier for developers to perform equivalent scans
#----------------------------------------------------------------------------------

.PHONY: scan-version
scan-version:
	PATH=$(DEPSGOBIN):$$PATH go run ./hack/trivy/cli/main.go scan -v "$(VERSION)"
