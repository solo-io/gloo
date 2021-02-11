ROOTDIR ?= ../../..
RELATIVE_OUTPUT_DIR ?= _output
OUTPUT_DIR ?= $(ROOTDIR)/$(RELATIVE_OUTPUT_DIR)

## helper variables for when we're running from this file, and not from the root dir
VERSION ?= dev
CLUSTER_NAME ?= local

#----------------------------------------------------------------------------------
# Generated Code
#----------------------------------------------------------------------------------
PROTOC_IMPORT_PATH:=$(ROOTDIR)/vendor_any

.PHONY: generate-gloo-fed
generate-gloo-fed: generate-gloo-fed-code

DEPSGOBIN=$(ROOTDIR)/.bin

.PHONY: mod-download
mod-download:
	go mod download

.PHONY: update-deps
update-deps: mod-download
	mkdir -p $(DEPSGOBIN)
	GOBIN=$(DEPSGOBIN) go install istio.io/tools/cmd/protoc-gen-jsonshim
	GOBIN=$(DEPSGOBIN) go install github.com/solo-io/protoc-gen-ext
	GOBIN=$(DEPSGOBIN) go install golang.org/x/tools/cmd/goimports
	GOBIN=$(DEPSGOBIN) go install github.com/envoyproxy/protoc-gen-validate
	GOBIN=$(DEPSGOBIN) go install github.com/golang/protobuf/protoc-gen-go
	GOBIN=$(DEPSGOBIN) go install github.com/golang/mock/gomock
	GOBIN=$(DEPSGOBIN) go install github.com/golang/mock/mockgen
	GOBIN=$(DEPSGOBIN) go install github.com/google/wire/cmd/wire

.PHONY: clean-artifacts
clean-artifacts:
	rm -rf _output

.PHONY: clean-generated-protos
clean-generated-protos:
	rm -rf $(ROOTDIR)/projects/apiserver/api/fed.rpc/v1/*resources.proto

# Clean
.PHONY: clean-fed
clean-fed: clean-artifacts clean-generated-protos
	rm -rf $(ROOTDIR)/vendor_any
	rm -rf $(ROOTDIR)/projects/gloo-fed/pkg/api
	rm -rf $(ROOTDIR)/projects/apiserver/pkg/api
	rm -rf $(ROOTDIR)/projects/glooctl-extensions/fed/pkg/api

# Generated Code - Required to update Codgen Templates
.PHONY: generate-gloo-fed-code
generate-gloo-fed-code: clean-fed
	PATH=$(DEPSGOBIN):$$PATH go run $(ROOTDIR)/projects/gloo-fed/generate.go # Generates clients, controllers, etc
	PATH=$(DEPSGOBIN):$$PATH $(ROOTDIR)/projects/gloo-fed/ci/hack-fix-marshal.sh # TODO: figure out a more permanent way to deal with this
	PATH=$(DEPSGOBIN):$$PATH go run projects/gloo-fed/generate.go -apiserver # Generates apiserver protos into go code
	PATH=$(DEPSGOBIN):$$PATH go generate $(ROOTDIR)/projects/... # Generates mocks
	PATH=$(DEPSGOBIN):$$PATH goimports -w $(SUBDIRS)
	PATH=$(DEPSGOBIN):$$PATH go mod tidy
	#PATH=$(DEPSGOBIN):$$PATH make generated-ui

#----------------------------------------------------------------------------------
# Gloo Federation Projects
#----------------------------------------------------------------------------------

.PHONY: allgloofed
allgloofed: gloo-fed gloo-fed-rbac-validating-webhook gloo-fed-apiserver gloo-fed-apiserver-envoy

.PHONY: gloofed-docker
gloofed-docker: gloo-fed-docker gloo-fed-rbac-validating-webhook-docker gloo-fed-apiserver-docker gloo-fed-apiserver-envoy-docker

.PHONY: gloofed-load-kind-images
gloofed-load-kind-images: kind-load-gloo-fed kind-load-gloo-fed-rbac-validating-webhook kind-load-gloo-fed-apiserver kind-load-gloo-fed-apiserver-envoy

#----------------------------------------------------------------------------------
# Deployment Manifests / Helm
#----------------------------------------------------------------------------------

# creates Chart.yaml, values.yaml, and requirements.yaml
.PHONY: gloofed-helm-template
gloofed-helm-template:
	mkdir -p $(HELM_SYNC_DIR_GLOO_FED)
	sed -e 's/%version%/'$(VERSION)'/' $(GLOO_FED_CHART_DIR)/Chart-template.yaml > $(GLOO_FED_CHART_DIR)/Chart.yaml
	sed -e 's/%version%/'$(VERSION)'/' $(GLOO_FED_CHART_DIR)/values-template.yaml > $(GLOO_FED_CHART_DIR)/values.yaml

.PHONY: gloofedproduce-manifests
gloofedproduce-manifests:
	helm repo add gloo-fed https://storage.googleapis.com/gloo-fed-helm
	helm template gloo-fed install/helm/gloo-fed --namespace gloo-fed > $(MANIFEST_DIR)/$(MANIFEST_FOR_GLOO_FED)

.PHONY: package-gloo-fed-charts
package-gloo-fed-charts: gloofed-helm-template
	helm package --destination $(HELM_SYNC_DIR_GLOO_FED) $(GLOO_FED_CHART_DIR)

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
	docker build -t quay.io/solo-io/gloo-fed:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed
kind-load-gloo-fed: gloo-fed-docker
	kind load docker-image --name local quay.io/solo-io/gloo-fed:$(VERSION)


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
	docker build -t quay.io/solo-io/gloo-fed-rbac-validating-webhook:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed-rbac-validating-webhook
kind-load-gloo-fed-rbac-validating-webhook: gloo-fed-rbac-validating-webhook-docker
	kind load docker-image --name $(CLUSTER_NAME) quay.io/solo-io/gloo-fed-rbac-validating-webhook:$(VERSION)



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
	docker build -t quay.io/solo-io/gloo-fed-apiserver:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_APISERVER_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed-apiserver
kind-load-gloo-fed-apiserver: gloo-fed-apiserver-docker
	kind load docker-image --name $(CLUSTER_NAME) quay.io/solo-io/gloo-fed-apiserver:$(VERSION)

#----------------------------------------------------------------------------------
# apiserver-envoy
#----------------------------------------------------------------------------------
CONFIG_YAML=cfg.yaml

GLOO_FED_APISERVER_ENVOY_DIR=$(ROOTDIR)/projects/apiserver/apiserver-envoy

.PHONY: gloo-fed-apiserver-envoy-docker
gloo-fed-apiserver-envoy-docker:
	cp $(GLOO_FED_APISERVER_ENVOY_DIR)/$(CONFIG_YAML) $(OUTPUT_DIR)/$(CONFIG_YAML)
	docker build -t quay.io/solo-io/gloo-fed-apiserver-envoy:$(VERSION) $(OUTPUT_DIR) -f $(GLOO_FED_APISERVER_ENVOY_DIR)/Dockerfile;

.PHONY: kind-load-gloo-fed-apiserver-envoy
kind-load-gloo-fed-apiserver-envoy: gloo-fed-apiserver-envoy-docker
	kind load docker-image --name $(CLUSTER_NAME) quay.io/solo-io/gloo-fed-apiserver-envoy:$(VERSION)
