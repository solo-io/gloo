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
help: FIRST_COLUMN_WIDTH=35
help: ## Output the self-documenting make targets
	@grep -hE '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-$(FIRST_COLUMN_WIDTH)s\033[0m %s\n", $$1, $$2}'

#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

PACKAGE_PATH:=github.com/solo-io/solo-projects
# Kind of a hack to make sure _output exists
z := $(shell mkdir -p $(OUTPUT_DIR))
SOURCES := $(shell find . -name "*.go" | grep -v test.go)

GCS_BUCKET := glooctl-plugins
WASM_GCS_PATH := glooctl-wasm
FED_GCS_PATH := glooctl-fed

ENVOY_GLOO_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee:1.23.1-patch7
ENVOY_GLOO_DEBUG_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee-debug:1.23.1-patch7
ENVOY_GLOO_FIPS_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee-fips:1.23.1-patch7
ENVOY_GLOO_FIPS_DEBUG_IMAGE ?= gcr.io/gloo-ee/envoy-gloo-ee-fips-debug:1.23.1-patch7

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

LDFLAGS := "-X github.com/solo-io/solo-projects/pkg/version.Version=$(VERSION) -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore"
GCFLAGS := 'all=-N -l'

GO_BUILD_FLAGS := GO111MODULE=on CGO_ENABLED=0 GOARCH=$(DOCKER_GOARCH)
UNAME_M ?=$(shell uname -m)

# Passed by cloudbuild
GCLOUD_PROJECT_ID := $(GCLOUD_PROJECT_ID)
BUILD_ID := $(BUILD_ID)

TEST_IMAGE_TAG := test-$(BUILD_ID)
TEST_ASSET_DIR := $(ROOTDIR)/_test
GCR_REPO_PREFIX := gcr.io/$(GCLOUD_PROJECT_ID)

GINKGO_ENV := GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore ACK_GINKGO_RC=true ACK_GINKGO_DEPRECATIONS=1.16.5

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
update-all-deps: install-go-tools install-node-packages ## install-go-tools and install-node-packages

DEPSGOBIN=$(ROOTDIR)/.bin

# https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
.PHONY: install-go-tools
install-go-tools: mod-download install-test-tools ## Download and install Go dependencies
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
	GOBIN=$(DEPSGOBIN) go install github.com/solo-io/protoc-gen-openapi

.PHONY: install-test-tools
install-test-tools:
	mkdir -p $(DEPSGOBIN)
	GOBIN=$(DEPSGOBIN) go install github.com/onsi/ginkgo/ginkgo

.PHONY: mod-download
mod-download:
	go mod download all

.PHONY: install-node-packages
install-node-packages:
	npm install -g yarn
	make install-graphql-js
	make update-ui-deps

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
	rm -rf $(ROOTDIR)/projects/gloo/pkg/api
	rm -rf $(ROOTDIR)/projects/gloo-fed/pkg/api
	rm -rf $(ROOTDIR)/projects/apiserver/pkg/api
	rm -rf $(ROOTDIR)/projects/glooctl-plugins/fed/pkg/api
	rm -rf $(ROOTDIR)/projects/apiserver/server/services/single_cluster_resource_handler/*

.PHONY: run-tests
run-tests: install-node-packages ## Run all tests, or only run the test package at {TEST_PKG} if it is specified
ifneq ($(RELEASE), "true")
	PATH=$(DEPSGOBIN):$$PATH go generate ./test/extauth/plugins/... ./projects/extauth/plugins/...
	$(GINKGO_ENV) VERSION=$(VERSION) $(DEPSGOBIN)/ginkgo -ldflags=$(LDFLAGS) -r -failFast -trace -progress -compilers=4 -failOnPending -noColor -skipPackage=kube2e,gloo-fed-e2e $(TEST_PKG)
endif

# command to run regression tests with guaranteed access to $(DEPSGOBIN)/ginkgo
# requires the environment variable KUBE2E_TESTS to be set to the test type you wish to run
.PHONY: run-ci-regression-tests
run-ci-regression-tests: install-test-tools ## Run the Kubernetes E2E Tests in the {KUBE2E_TESTS} package
run-ci-regression-tests:
	# We intentionally leave out the `-r` ginkgo flag, since we are specifying the exact package that we want run
	$(GINKGO_ENV) $(DEPSGOBIN)/ginkgo -failFast -trace -progress -race -failOnPending -noColor ./test/kube2e/$(KUBE2E_TESTS)

# command to run regression tests with guaranteed access to $(DEPSGOBIN)/ginkgo
# requires the environment variable KUBE2E_TESTS to be set to the test type you wish to run
.PHONY: run-ci-gloo-fed-regression-tests
run-ci-gloo-fed-regression-tests: install-test-tools ## Run the Kubernetes E2E Tests in the {gloo-fed-e2e} package
	# We intentionally leave out the `-r` ginkgo flag, since we are specifying the exact package that we want run
	$(GINKGO_ENV) $(DEPSGOBIN)/ginkgo -failFast -trace -progress -race -failOnPending -noColor ./test/gloo-fed-e2e

# command to run e2e tests
# requires the environment variable ENVOY_IMAGE_TAG to be set to the tag of the gloo-ee-envoy-wrapper Docker image you wish to run
.PHONY: run-e2e-tests
run-e2e-tests: ## Run the in memory Envoy e2e tests (ENVOY_IMAGE_TAG)
run-e2e-tests: install-test-tools
	# # We intentionally leave out the `-r` ginkgo flag, since we are specifying the exact package that we want run
	$(GINKGO_ENV) $(DEPSGOBIN)/ginkgo -failFast -trace -progress -race -compilers=4 -failOnPending ./test/e2e/

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
generated-code: update-licenses ## Evaluate go generate
	rm -rf $(ROOTDIR)/vendor_any
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
	PATH=$(DEPSGOBIN):$$PATH go run $(ROOTDIR)/install/helm/gloo-ee/generate.go $(VERSION) --generate-helm-docs $(USE_DIGESTS) # Generate Helm Documentation



#################
#     Build     #
#################

#----------------------------------------------------------------------------------
# allprojects
#----------------------------------------------------------------------------------
# helper for testing
.PHONY: allprojects
allprojects: gloo-fed-apiserver gloo extauth extauth-fips rate-limit rate-limit-fips observability caching

#----------------------------------------------------------------------------------
# Gloo Fed
#----------------------------------------------------------------------------------

GLOO_FED_DIR=$(ROOTDIR)/projects/gloo-fed
GLOO_FED_SOURCES=$(shell find $(GLOO_FED_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gloo-fed-linux-$(DOCKER_GOARCH): $(GLOO_FED_SOURCES)
	CGO_ENABLED=0 GOARCH=$(DOCKER_GOARCH) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_DIR)/cmd/main.go

.PHONY: gloo-fed
gloo-fed: $(OUTPUT_DIR)/gloo-fed-linux-$(DOCKER_GOARCH)

.PHONY: gloo-fed-docker
gloo-fed-docker: $(OUTPUT_DIR)/gloo-fed-linux-$(DOCKER_GOARCH)
	docker buildx build --load -t $(IMAGE_REPO)/gloo-fed:$(VERSION) $(DOCKER_BUILD_ARGS) $(OUTPUT_DIR) -f $(GLOO_FED_DIR)/cmd/Dockerfile --build-arg GOARCH=$(DOCKER_GOARCH);

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

# proto sources
APISERVER_DIR=$(ROOTDIR)/projects/apiserver/api/
APISERVER_SOURCES=$(shell find $(APISERVER_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(OUTPUT_DIR)/gloo-fed-apiserver-linux-$(DOCKER_GOARCH): $(APISERVER_SOURCES)
	cp -r projects/gloo/pkg/plugins/graphql/js $(OUTPUT_DIR)/js
	cp -r projects/ui/src/proto $(OUTPUT_DIR)/js
	CGO_ENABLED=0 GOARCH=$(DOCKER_GOARCH) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_APISERVER_DIR)/cmd/main.go

.PHONY: gloo-fed-apiserver
gloo-fed-apiserver: $(OUTPUT_DIR)/gloo-fed-apiserver-linux-$(DOCKER_GOARCH)

.PHONY: gloo-fed-apiserver-docker
gloo-fed-apiserver-docker: $(OUTPUT_DIR)/gloo-fed-apiserver-linux-$(DOCKER_GOARCH)
	docker buildx build --load -t $(IMAGE_REPO)/gloo-fed-apiserver:$(VERSION) $(DOCKER_BUILD_ARGS) $(OUTPUT_DIR) -f $(GLOO_FED_APISERVER_DIR)/cmd/Dockerfile --build-arg GOARCH=$(DOCKER_GOARCH);

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
	docker buildx build --load -t $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION) $(DOCKER_BUILD_ARGS) $(OUTPUT_DIR) -f $(GLOO_FED_APISERVER_ENVOY_DIR)/Dockerfile;

.PHONY: kind-load-gloo-fed-apiserver-envoy
kind-load-gloo-fed-apiserver-envoy: gloo-fed-apiserver-envoy-docker
	kind load docker-image --name $(CLUSTER_NAME) $(IMAGE_REPO)/gloo-fed-apiserver-envoy:$(VERSION)

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
	yarn install -d projects/gloo/pkg/plugins/graphql/js

.PHONY: run-apiserver
run-apiserver: checkprogram-protoc install-graphql-js checkenv-GLOO_LICENSE_KEY
# Todo: This should check that /etc/hosts includes the following line:
# 127.0.0.1 docker.internal
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=ignore \
	GRPC_PORT=$(GRPC_PORT) \
	POD_NAMESPACE=gloo-system \
	GRAPHQL_JS_ROOT="./projects/gloo/pkg/plugins/graphql/js/" \
	GRAPHQL_PROTO_ROOT="./projects/ui/src/proto/" \
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
	docker buildx build --load -t $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION) $(DOCKER_BUILD_ARGS) $(OUTPUT_DIR) -f $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/Dockerfile --build-arg GOARCH=$(DOCKER_GOARCH);

.PHONY: kind-load-gloo-fed-rbac-validating-webhook
kind-load-gloo-fed-rbac-validating-webhook: gloo-fed-rbac-validating-webhook-docker
	kind load docker-image --name $(CLUSTER_NAME) $(IMAGE_REPO)/gloo-fed-rbac-validating-webhook:$(VERSION)

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
	yarn --cwd projects/ui pbjs -t json -o src/Components/Features/Graphql/data/graphql.json \
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
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/*/*.proto

	$(PROTOC) -I$(APISERVER_DIR) \
	$(TS_OUT) \
	$(PROTOC_IMPORT_PATH)/github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/*/*/*.proto

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
build-ui: update-gloo-fed-ui-deps
ifneq ($(LOCAL_BUILD),)
	yarn --cwd $(APISERVER_UI_DIR) build
endif

.PHONY: gloo-federation-console-docker
gloo-federation-console-docker: build-ui
	docker buildx build --load -t $(IMAGE_REPO)/gloo-federation-console:$(VERSION) $(DOCKER_BUILD_ARGS) $(APISERVER_UI_DIR) -f $(APISERVER_UI_DIR)/Dockerfile

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
	docker buildx build --load -t $(IMAGE_REPO)/rate-limit-ee-build-container:$(VERSION) \
		-f $(RATELIMIT_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_VERSION) \
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
	docker create -ti --name rate-limit-temp-container $(IMAGE_REPO)/rate-limit-ee-build-container:$(VERSION) bash
	docker cp rate-limit-temp-container:/rate-limit-linux-$(DOCKER_GOARCH) $(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH)
	docker rm -f rate-limit-temp-container

.PHONY: rate-limit
rate-limit: $(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH)

$(RATELIMIT_OUT_DIR)/Dockerfile: $(RATELIMIT_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: rate-limit-ee-docker
rate-limit-ee-docker: $(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker

$(RATELIMIT_OUT_DIR)/.rate-limit-ee-docker: $(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH) $(RATELIMIT_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REPO)/rate-limit-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(call get_test_tag_option,rate-limit-ee) $(RATELIMIT_OUT_DIR)
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
	docker buildx build --load -t $(IMAGE_REPO)/rate-limit-ee-build-container-fips:$(VERSION) \
		-f $(RATELIMIT_FIPS_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOBORING_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg GITHUB_TOKEN \
		$(DOCKER_GO_BORING_ARGS) \
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
	docker buildx build --load -t $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION) $(DOCKER_GO_BORING_ARGS) $(call get_test_tag_option,rate-limit-ee-fips) $(RATELIMIT_FIPS_OUT_DIR) \
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
	docker buildx build --load -t $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) \
		-f $(EXTAUTH_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_VERSION) \
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
	docker create -ti --name extauth-temp-container $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) bash
	docker cp extauth-temp-container:/extauth-linux-$(DOCKER_GOARCH) $(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH)
	docker rm -f extauth-temp-container

# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
# we need to be able to pass a variable to the extauth-ee-docker buildx build --load for Docker_CGO_ENABLED
$(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64: $(EXTAUTH_OUT_DIR)/.extauth-ee-docker-build
	docker create -ti --name verify-plugins-temp-container $(IMAGE_REPO)/extauth-ee-build-container:$(VERSION) bash
	docker cp verify-plugins-temp-container:/verify-plugins-linux-amd64 $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64
	docker rm -f verify-plugins-temp-container

# Build extauth binaries
.PHONY: extauth
extauth: $(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH) $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64

# Build ext-auth-plugins docker image (Cannot be built at all on Apple Silicon)
.PHONY: ext-auth-plugins-docker
ifeq  ($(IS_ARM_MACHINE), )
ext-auth-plugins-docker: $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64
	docker buildx build --load -t $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) -f projects/extauth/plugins/Dockerfile \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_VERSION) \
		--build-arg GC_FLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg VERIFY_SCRIPT=$(RELATIVE_EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 \
		--build-arg GITHUB_TOKEN \
		--build-arg USE_APK=true \
		$(DOCKER_BUILD_ARGS) \
		.
endif

# Build extauth server docker image
.PHONY: extauth-ee-docker
extauth-ee-docker: $(EXTAUTH_OUT_DIR)/.extauth-ee-docker

$(EXTAUTH_OUT_DIR)/.extauth-ee-docker: $(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH) $(EXTAUTH_OUT_DIR)/verify-plugins-linux-amd64 $(EXTAUTH_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REPO)/extauth-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(call get_test_tag_option,extauth-ee) $(EXTAUTH_OUT_DIR)
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

# GO_BORING_ARGS is set to amd64 currently this is because BORING will only work on amd64
# https://go.googlesource.com/go/+/refs/heads/dev.boringcrypto.go1.12/misc/boring/
$(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker-build: $(EXTAUTH_SOURCES) $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build
	docker buildx build --load -t $(IMAGE_REPO)/extauth-ee-build-container-fips:$(VERSION) \
		-f $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOBORING_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg GITHUB_TOKEN $(DOCKER_GO_BORING_ARGS) \
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
# NOTE: ext-auth-plugin will not build on arm64 machines
.PHONY: ext-auth-plugins-fips-docker
ext-auth-plugins-fips-docker: $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64
	docker buildx build --load -t $(IMAGE_REPO)/ext-auth-plugins-fips:$(VERSION) -f projects/extauth/plugins/Dockerfile \
		--build-arg GO_BUILD_IMAGE=$(GOBORING_VERSION) \
		--build-arg GC_FLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg VERIFY_SCRIPT=$(RELATIVE_EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 \
		--build-arg GITHUB_TOKEN $(DOCKER_GO_BORING_ARGS) \
		.

# Build extauth server docker image
.PHONY: extauth-ee-fips-docker
extauth-ee-fips-docker: $(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker

$(EXTAUTH_FIPS_OUT_DIR)/.extauth-ee-docker: $(EXTAUTH_FIPS_OUT_DIR)/extauth-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/verify-plugins-linux-amd64 $(EXTAUTH_FIPS_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REPO)/extauth-ee-fips:$(VERSION) $(DOCKER_GO_BORING_ARGS) $(call get_test_tag_option,extauth-ee-fips) $(EXTAUTH_FIPS_OUT_DIR) \
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
	docker buildx build --load -t $(IMAGE_REPO)/observability-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(call get_test_tag_option,observability-ee) $(OBS_OUT_DIR)
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
	docker buildx build --load -t $(IMAGE_REPO)/caching-ee:$(VERSION) $(call get_test_tag_option,caching-ee) $(CACHE_OUT_DIR) $(PLATFORM)
	touch $@

#----------------------------------------------------------------------------------
# Gloo
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)
GLOO_OUT_DIR=$(OUTPUT_DIR)/gloo

# the executable outputs as amd64 only because it is placed in an image that is amd64
$(GLOO_OUT_DIR)/gloo-linux-amd64: $(GLOO_SOURCES)
	GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_DIR)/cmd/main.go


.PHONY: gloo
gloo: $(GLOO_OUT_DIR)/gloo-linux-amd64

$(GLOO_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@


.PHONY: gloo-ee-docker
gloo-ee-docker: $(GLOO_OUT_DIR)/.gloo-ee-docker

$(GLOO_OUT_DIR)/.gloo-ee-docker: $(GLOO_OUT_DIR)/gloo-linux-amd64 $(GLOO_OUT_DIR)/Dockerfile
	cp -r projects/gloo/pkg/plugins/graphql/js $(GLOO_OUT_DIR)/js
	cp -r projects/ui/src/proto $(GLOO_OUT_DIR)/js
	docker buildx build --load $(call get_test_tag_option,gloo-ee) $(GLOO_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) $(DOCKER_GO_AMD_64_ARGS) \
		-t $(IMAGE_REPO)/gloo-ee:$(VERSION)
	touch $@

gloo-ee-docker-dev: $(GLOO_OUT_DIR)/gloo-linux-amd64 $(GLOO_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REPO)/gloo-ee:$(VERSION) $(DOCKER_BUILD_ARGS) $(GLOO_OUT_DIR) --no-cache
	touch $@

#----------------------------------------------------------------------------------
# Gloo with race detection enabled.
# This is intended to be used to aid in local debugging by swapping out this image in a running gloo instance
#----------------------------------------------------------------------------------
GLOO_RACE_OUT_DIR=$(OUTPUT_DIR)/gloo-race

$(GLOO_RACE_OUT_DIR)/Dockerfile.build: $(GLOO_DIR)/Dockerfile
	mkdir -p $(GLOO_RACE_OUT_DIR)
	cp $< $@

$(GLOO_RACE_OUT_DIR)/.gloo-race-ee-docker-build: $(GLOO_SOURCES) $(GLOO_RACE_OUT_DIR)/Dockerfile.build
	docker build -t $(IMAGE_REPO)/gloo-race-ee-build-container:$(VERSION) \
		-f $(GLOO_RACE_OUT_DIR)/Dockerfile.build \
		--build-arg GO_BUILD_IMAGE=$(GOLANG_VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GCFLAGS=$(GCFLAGS) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--build-arg USE_APK=true \
		--build-arg GITHUB_TOKEN \
		$(DOCKER_BUILD_ARGS) \
		.
	touch $@
# Build inside container as we need to target linux and must compile with CGO_ENABLED=1
# We may be running Docker in a VM (eg, minikube) so be careful about how we copy files out of the containers
$(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH): $(GLOO_RACE_OUT_DIR)/.gloo-race-ee-docker-build
	docker create -ti --name gloo-race-temp-container $(IMAGE_REPO)/gloo-race-ee-build-container:$(VERSION) bash
	docker cp gloo-race-temp-container:/gloo-linux-$(DOCKER_GOARCH) $(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH)
	docker rm -f gloo-race-temp-container

.PHONY: gloo-race
gloo-race: $(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH)

$(GLOO_RACE_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: gloo-race-ee-docker
gloo-race-ee-docker: $(GLOO_RACE_OUT_DIR)/.gloo-race-ee-docker
$(GLOO_RACE_OUT_DIR)/.gloo-race-ee-docker: $(GLOO_RACE_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH) $(GLOO_RACE_OUT_DIR)/Dockerfile
	cp -r projects/gloo/pkg/plugins/graphql/js $(GLOO_RACE_OUT_DIR)/js
	cp -r projects/ui/src/proto $(GLOO_RACE_OUT_DIR)/js
	docker build $(call get_test_tag_option,gloo-ee) $(GLOO_RACE_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_IMAGE) $(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_REPO)/gloo-ee:$(VERSION)-race
	touch $@
#----------------------------------------------------------------------------------
# Gloo with FIPS Envoy
#----------------------------------------------------------------------------------

GLOO_DIR=projects/gloo
GLOO_SOURCES=$(shell find $(GLOO_DIR) -name "*.go" | grep -v test | grep -v generated.go)
GLOO_FIPS_OUT_DIR=$(OUTPUT_DIR)/gloo-fips

$(GLOO_FIPS_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH): $(GLOO_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_DIR)/cmd/main.go


.PHONY: gloo-fips
gloo-fips: $(GLOO_FIPS_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH)

$(GLOO_FIPS_OUT_DIR)/Dockerfile: $(GLOO_DIR)/cmd/Dockerfile
	cp $< $@


.PHONY: gloo-fips-ee-docker
gloo-fips-ee-docker: $(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker

$(GLOO_FIPS_OUT_DIR)/.gloo-ee-docker: $(GLOO_FIPS_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH) $(GLOO_FIPS_OUT_DIR)/Dockerfile
	cp -r projects/gloo/pkg/plugins/graphql/js $(GLOO_FIPS_OUT_DIR)/js
	cp -r projects/ui/src/proto $(GLOO_FIPS_OUT_DIR)/js
	docker buildx build --load $(call get_test_tag_option,gloo-ee) $(GLOO_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_IMAGE) $(DOCKER_GO_BORING_ARGS) \
		-t $(IMAGE_REPO)/gloo-ee-fips:$(VERSION)
	touch $@

gloo-fips-ee-docker-dev: $(GLOO_FIPS_OUT_DIR)/gloo-linux-$(DOCKER_GOARCH) $(GLOO_FIPS_OUT_DIR)/Dockerfile
	docker buildx build --load -t $(IMAGE_REPO)/gloo-ee-fips:$(VERSION) $(DOCKER_BUILD_ARGS) $(GLOO_FIPS_OUT_DIR) --no-cache
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
discovery: $(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-$(DOCKER_GOARCH)
$(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery: $(DISCOVERY_DIR)/cmd/Dockerfile
	cp $< $@

.PHONY: discovery-ee-docker
discovery-ee-docker: $(DISCOVERY_OUTPUT_DIR)/discovery-ee-linux-$(DOCKER_GOARCH) $(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery
	docker buildx build --load $(DISCOVERY_OUTPUT_DIR) -f $(DISCOVERY_OUTPUT_DIR)/Dockerfile.discovery \
		$(DOCKER_BUILD_ARGS) -t $(IMAGE_REPO)/discovery-ee:$(VERSION) $(QUAY_EXPIRATION_LABEL)

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

# NOTE although it says amd64 it is determined by the architecture of the machine building it
# this is because of the dependency on github.com/solo-io/k8s-utils@v0.1.0/testutils/helper/install.go
.PHONY: glooctl
glooctl: $(OUTPUT_DIR)/glooctl
.PHONY: glooctl-linux
glooctl-linux: $(OUTPUT_DIR)/glooctl-linux-amd64
.PHONY: glooctl-darwin
glooctl-darwin: $(OUTPUT_DIR)/glooctl-darwin-amd64
.PHONY: glooctl-windows
glooctl-windows: $(OUTPUT_DIR)/glooctl-windows-amd64.exe

.PHONY: build-cli
build-cli: glooctl-linux glooctl-darwin glooctl-windows

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

$(ENVOYINIT_OUT_DIR)/envoyinit-linux-$(DOCKER_GOARCH): $(ENVOYINIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(ENVOYINIT_DIR)/main.go $(ENVOYINIT_DIR)/filter_types.gen.go

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
		-t $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) \
		-f $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit
	touch $@

.PHONY: gloo-ee-envoy-wrapper-debug-docker
gloo-ee-envoy-wrapper-debug-docker: $(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-debug-docker

$(ENVOYINIT_OUT_DIR)/.gloo-ee-envoy-wrapper-debug-docker: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-$(DOCKER_GOARCH) $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_OUT_DIR)/docker-entrypoint.sh
	docker buildx build --load $(call get_test_tag_option,gloo-ee-envoy-wrapper-debug) $(ENVOYINIT_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_DEBUG_IMAGE) $(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION)-debug \
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
	docker buildx build --load $(call get_test_tag_option,gloo-ee-envoy-wrapper-fips) $(ENVOYINIT_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_IMAGE) $(DOCKER_GO_BORING_ARGS) \
		-t $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) \
		-f $(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit
	touch $@

.PHONY: gloo-ee-envoy-wrapper-fips-debug-docker
gloo-ee-envoy-wrapper-fips-debug-docker: $(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-debug-docker

$(ENVOYINIT_FIPS_OUT_DIR)/.gloo-ee-envoy-wrapper-fips-debug-docker: $(ENVOYINIT_FIPS_OUT_DIR)/envoyinit-linux-amd64 $(ENVOYINIT_FIPS_OUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_FIPS_OUT_DIR)/docker-entrypoint.sh
	docker buildx build --load $(call get_test_tag_option,gloo-ee-envoy-wrapper-fips-debug) $(ENVOYINIT_FIPS_OUT_DIR) \
		--build-arg ENVOY_IMAGE=$(ENVOY_GLOO_FIPS_DEBUG_IMAGE) $(DOCKER_GO_BORING_ARGS) \
		-t $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION)-debug \
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
	PATH=$(DEPSGOBIN):$$PATH $(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION) --gloo-fed-repo-override="file://$(GLOO_FED_CHART_DIR)" $(USE_DIGESTS)

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
       extauth-ee-fips-docker gloo-ee-docker gloo-fips-ee-docker gloo-race-ee-docker gloo-ee-envoy-wrapper-docker gloo-ee-envoy-wrapper-debug-docker discovery-ee-docker\
       gloo-ee-envoy-wrapper-fips-docker gloo-ee-envoy-wrapper-fips-debug-docker observability-ee-docker caching-ee-docker ext-auth-plugins-docker ext-auth-plugins-fips-docker \
       gloo-fed-docker gloo-fed-apiserver-docker gloo-fed-apiserver-envoy-docker gloo-federation-console-docker gloo-fed-rbac-validating-webhook-docker

# Depends on DOCKER_IMAGES, which is set to docker if RELEASE is "true", otherwise empty (making this a no-op).
# This prevents executing the dependent targets if RELEASE is not true, while still enabling `make docker`
# to be used for local testing.
# docker-push is intended to be run by CI
docker-push: $(DOCKER_IMAGES)
.PHONY: docker-push
docker-push: docker-push-non-fips docker-push-fips docker-push-fed

docker-push-non-fips:
.PHONY: docker-push-non-fips
docker-push-non-fips:
ifeq ($(RELEASE), "true")
	docker push $(IMAGE_REPO)/rate-limit-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) && \
	docker push $(IMAGE_REPO)/observability-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/caching-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/extauth-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/discovery-ee:$(VERSION)
ifeq  ($(IS_ARM_MACHINE), )
	# these images are not built on ARM, so this is adding complexity.  There is no reason to add this as well to the normal build of gloo-ee.
	# so if pushing because of ARM, we will
	docker push $(IMAGE_REPO)/ext-auth-plugins:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee:$(VERSION)-race && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION)-debug
endif
#
endif

docker-push-fips:
.PHONY: docker-push-fips
docker-push-fips:
ifeq ($(RELEASE),"true")
	docker push $(IMAGE_REPO)/rate-limit-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION)-debug && \
	docker push $(IMAGE_REPO)/extauth-ee-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/ext-auth-plugins-fips:$(VERSION)
endif

docker-push-fed:
.PHONY: docker-push-fed
docker-push-fed:
ifeq ($(RELEASE),"true")
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
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee:$(VERSION)-race $(IMAGE_REPO)/gloo-ee:$(VERSION)-race && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper:$(VERSION) $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper:$(VERSION)-debug $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION)-debug && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper-fips:$(VERSION) $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/gloo-ee-envoy-wrapper-fips:$(VERSION)-debug $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION)-debug && \
	docker tag $(RETAG_IMAGE_REGISTRY)/observability-ee:$(VERSION) $(IMAGE_REPO)/observability-ee:$(VERSION) && \
	docker tag $(RETAG_IMAGE_REGISTRY)/caching-ee:$(VERSION) $(IMAGE_REPO)/caching-ee:$(VERSION) && \
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
	docker tag $(RETAG_IMAGE_REGISTRY)/caching-ee:$(VERSION)-extended $(IMAGE_REPO)/caching-ee:$(VERSION)-extended && \
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
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper:$(VERSION)-debug && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION) && \
	docker push $(IMAGE_REPO)/gloo-ee-envoy-wrapper-fips:$(VERSION)-debug && \
	docker push $(IMAGE_REPO)/observability-ee:$(VERSION) && \
	docker push $(IMAGE_REPO)/caching-ee:$(VERSION) && \
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
	docker push $(IMAGE_REPO)/caching-ee:$(VERSION)-extended && \
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

.PHONY: push-kind-images ## Build and load images into a kind cluster
push-kind-images:
ifeq ($(USE_FIPS),true)
push-kind-images: build-and-load-kind-images-fips
else
push-kind-images: build-and-load-kind-images-non-fips
endif

# Build and load images for a non-fips compliant (data plane) installation of Gloo Edge
# Used in CI during regression tests
.PHONY: build-and-load-kind-images-non-fips
build-and-load-kind-images-non-fips: build-kind-images-non-fips load-kind-images-non-fips

.PHONY: build-kind-images-non-fips
build-kind-images-non-fips: gloo-ee-docker
build-kind-images-non-fips: gloo-ee-envoy-wrapper-docker
build-kind-images-non-fips: rate-limit-ee-docker
build-kind-images-non-fips: extauth-ee-docker
# arm cannot build the ext-auth-plugin currently
ifeq ($(IS_ARM_MACHINE), )
build-kind-images-non-fips: ext-auth-plugins-docker
endif
build-kind-images-non-fips: observability-ee-docker
build-kind-images-non-fips: caching-ee-docker
build-kind-images-non-fips: discovery-ee-docker

.PHONY: load-kind-images-non-fips
load-kind-images-non-fips: kind-load-gloo-ee # gloo
load-kind-images-non-fips: kind-load-gloo-ee-envoy-wrapper # envoy
load-kind-images-non-fips: kind-load-rate-limit-ee # rate limit
load-kind-images-non-fips: kind-load-extauth-ee # ext auth
ifeq  ($(IS_ARM_MACHINE), )
load-kind-images-non-fips: kind-load-ext-auth-plugins # ext auth plugins
endif
load-kind-images-non-fips: kind-load-observability-ee # observability
load-kind-images-non-fips: kind-load-caching-ee # caching
load-kind-images-non-fips: kind-load-discovery-ee # discovery

# Build and load images for a fips compliant (data plane) installation of Gloo Edge
# Used in CI during regression tests
.PHONY: build-and-load-kind-images-fips
build-and-load-kind-images-fips: build-kind-images-fips load-kind-images-fips

.PHONY: build-kind-images-fips
build-kind-images-fips: gloo-fips-ee-docker # gloo
build-kind-images-fips: gloo-ee-envoy-wrapper-fips-docker # envoy
build-kind-images-fips: rate-limit-ee-fips-docker # rate limit
build-kind-images-fips: extauth-ee-fips-docker # ext auth
build-kind-images-fips: ext-auth-plugins-fips-docker # ext auth plugins
build-kind-images-fips: observability-ee-docker # observability
build-kind-images-fips: caching-ee-docker # caching
build-kind-images-fips: discovery-ee-docker # discovery

.PHONY: load-kind-images-fips
load-kind-images-fips: kind-load-gloo-ee-fips # gloo
load-kind-images-fips: kind-load-gloo-ee-envoy-wrapper-fips # envoy
load-kind-images-fips: kind-load-rate-limit-ee-fips # rate limit
load-kind-images-fips: kind-load-extauth-ee-fips # ext auth
load-kind-images-fips: kind-load-ext-auth-plugins-fips # ext auth plugins
load-kind-images-fips: kind-load-observability-ee # observability
load-kind-images-fips: kind-load-caching-ee # caching
load-kind-images-fips: kind-load-discovery-ee # discovery

# arm local development requires work around to deploy to docker registry instead of kind load docker-image
docker-push-local-arm:
.PHONY: docker-push-local-arm
# set release because we will be pushing docker images to image repo
docker-push-local-arm:
ifeq ($(USE_FIPS),true)
docker-push-local-arm: build-kind-images-fips docker-push-fips
else
docker-push-local-arm: build-kind-images-non-fips docker-push-non-fips
endif

.PHONY: build-kind-assets
build-kind-assets: push-kind-images build-test-chart

TEST_DOCKER_TARGETS := gloo-federation-console-docker-test apiserver-envoy-docker-test gloo-fed-apiserver-docker-test rate-limit-ee-docker-test extauth-ee-docker-test observability-ee-docker-test caching-ee-docker-test gloo-ee-docker-test gloo-ee-envoy-wrapper-docker-test

push-test-images: $(TEST_DOCKER_TARGETS)

gloo-fed-apiserver-docker-test: $(OUTPUT_DIR)/gloo-fed-apiserver-linux-$(DOCKER_GOARCH) $(OUTPUT_DIR)/.gloo-fed-apiserver-docker
	docker push $(call get_test_tag,gloo-fed-apiserver)

gloo-fed-apiserver-envoy-docker-test: gloo-fed-apiserver-envoy-docker $(OUTPUT_DIR)/Dockerfile
	docker push $(call get_test_tag,gloo-fed-apiserver-envoy)

gloo-federation-console-docker-test: build-ui gloo-federation-console-docker
	docker push $(call get_test_tag,gloo-federation-console)

rate-limit-ee-docker-test: $(RATELIMIT_OUT_DIR)/rate-limit-linux-$(DOCKER_GOARCH) $(RATELIMIT_OUT_DIR)/Dockerfile
	docker push $(call get_test_tag,rate-limit-ee)

extauth-ee-docker-test: $(EXTAUTH_OUT_DIR)/extauth-linux-$(DOCKER_GOARCH) $(EXTAUTH_OUT_DIR)/Dockerfile
	docker push $(call get_test_tag,extauth-ee)

observability-ee-docker-test: $(OBS_OUT_DIR)/observability-linux-$(DOCKER_GOARCH) $(OBS_OUT_DIR)/Dockerfile
	docker push $(call get_test_tag,observability-ee)

caching-ee-docker-test: $(OBS_OUT_DIR)/caching-linux-$(DOCKER_GOARCH) $(OBS_OUT_DIR)/Dockerfile
	docker push $(call get_test_tag,caching-ee)

gloo-ee-docker-test: gloo-ee-docker
	docker push $(call get_test_tag,gloo-ee)

gloo-ee-envoy-wrapper-docker-test: $(ENVOYINIT_OUT_DIR)/envoyinit-linux-$(DOCKER_GOARCH) $(ENVOYINIT_OUT_DIR)/Dockerfile.envoyinit gloo-ee-envoy-wrapper-docker
	docker push $(call get_test_tag,gloo-ee-envoy-wrapper)

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
	$(GO_BUILD_FLAGS) go run install/helm/gloo-ee/generate.go $(VERSION) $(USE_DIGESTS)
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
