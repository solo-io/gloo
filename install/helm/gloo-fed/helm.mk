ROOTDIR ?= ../../../

#----------------------------------------------------------------------------------
# Gloo Federation Projects
#----------------------------------------------------------------------------------

.PHONY: allgloofed
allgloofed: gloo-fed gloo-fed-rbac-validating-webhook

.PHONY: gloofed-docker
gloofed-docker: gloo-fed-docker gloo-fed-rbac-validating-webhook-docker

.PHONY: gloofed-load-kind-images
gloofed-load-kind-images: kind-load-gloo-fed kind-load-gloo-fed-rbac-validating-webhook

#----------------------------------------------------------------------------------
# Gloo Fed
#----------------------------------------------------------------------------------

GLOO_FED_DIR=$(ROOTDIR)/projects/gloo-fed
GLOO_FED_OUTPUT_DIR=$(GLOO_FED_DIR)/_output
GLOO_FED_SOURCES=$(shell find $(GLOO_FED_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(GLOO_FED_OUTPUT_DIR)/gloo-fed-linux-amd64: $(GLOO_FED_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_DIR)/cmd/main.go

.PHONY: gloo-fed
gloo-fed: $(GLOO_FED_OUTPUT_DIR)/gloo-fed-linux-amd64

.PHONY: gloo-fed-docker
gloo-fed-docker: $(GLOO_FED_OUTPUT_DIR)/gloo-fed-linux-amd64
	docker build -t quay.io/solo-io/gloo-fed:$(VERSION) $(GLOO_FED_OUTPUT_DIR) -f $(GLOO_FED_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed
kind-load-gloo-fed: gloo-fed-docker
	kind load docker-image --name local quay.io/solo-io/gloo-fed:$(VERSION)


#----------------------------------------------------------------------------------
# Gloo Fed Rbac Webhook
#----------------------------------------------------------------------------------
GLOO_FED_RBAC_WEBHOOK_DIR=$(ROOTDIR)/projects/rbac-validating-webhook
GLOO_FED_RBAC_WEBHOOK_OUTPUT_DIR=$(GLOO_FED_RBAC_WEBHOOK_DIR)/_output
GLOO_FED_RBAC_WEBHOOK_SOURCES=$(shell find $(GLOO_FED_RBAC_WEBHOOK_DIR) -name "*.go" | grep -v test | grep -v generated.go)

$(GLOO_FED_RBAC_WEBHOOK_OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-amd64: $(GLOO_FED_RBAC_WEBHOOK_SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/main.go

.PHONY: gloo-fed-rbac-validating-webhook
gloo-fed-rbac-validating-webhook: $(GLOO_FED_RBAC_WEBHOOK_OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-amd64

.PHONY: gloo-fed-rbac-validating-webhook-docker
gloo-fed-rbac-validating-webhook-docker: $(GLOO_FED_RBAC_WEBHOOK_OUTPUT_DIR)/gloo-fed-rbac-validating-webhook-linux-amd64
	docker build -t quay.io/solo-io/gloo-fed-rbac-validating-webhook:$(VERSION) $(GLOO_FED_RBAC_WEBHOOK_OUTPUT_DIR) -f $(GLOO_FED_RBAC_WEBHOOK_DIR)/cmd/Dockerfile;

.PHONY: kind-load-gloo-fed-rbac-validating-webhook
kind-load-gloo-fed-rbac-validating-webhook: gloo-fed-rbac-validating-webhook-docker
	kind load docker-image --name $(CLUSTER_NAME) quay.io/solo-io/gloo-fed-rbac-validating-webhook:$(VERSION)

