TAGGED_VERSION ?= vlatest
ifneq (,$(TAGGED_VERSION))
VERSION := $(shell echo $(TAGGED_VERSION) | cut -c 2-)
endif

LDFLAGS := "-X github.com/solo-io/gloo/pkg/version.Version=$(VERSION)"

# TODO (Law[revert]): need to add CLI tests here, see CI pr.yaml workflow
.PHONY: test
test:
	go test -ldflags=$(LDFLAGS) ./...

.PHONY: test-full
test-full:
	go test -ldflags=$(LDFLAGS) -count=1 ./...

# internal target used by controller_suite_test.go
envtest:
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)"

#----------------------------------------------------------------------------------
# Env test
#----------------------------------------------------------------------------------
ENVTEST_K8S_VERSION = 1.31

.PHONY: envtest-path
envtest-path:
	@$(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --arch=amd64

# Populates the shared setup-envtest store ahead of the test run.
# The test packages under projects/gateway2 run in parallel and each shell out to
# envtest-path, so against a cold store they race to unpack the same
# version-platform directory: setup-envtest wipes that directory before extracting,
# so one process deletes the binaries out from under the other. Once the version is
# on disk, envtest-path resolves it locally and no longer writes to the store.
.PHONY: envtest-download
envtest-download:
	@$(ENVTEST) use $(ENVTEST_K8S_VERSION) --arch=amd64 -p path > /dev/null

ENVTEST = $(DEPSGOBIN)/setup-envtest

DEPSGOBIN:=$(shell pwd)/.bin
export PATH:=$(DEPSGOBIN):$(PATH)
export GOBIN:=$(DEPSGOBIN)

install-go-tools:
	mkdir -p $(DEPSGOBIN)
	test -s $(ENVTEST) || go install sigs.k8s.io/controller-runtime/tools/setup-envtest@e4c1c38bcbdb
	$(MAKE) envtest-download
