include ../../../Makefile.docker

VERSION := 0.0.1
ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output

#----------------------------------------------------------------------------------
# Docker
#----------------------------------------------------------------------------------

$(OUTPUT_DIR)/webhook-test-linux-$(DOCKER_GOARCH):
	CGO_ENABLED=0 GOARCH=$(DOCKER_GOARCH) GOOS=linux go build -o $@ internal/cmd/main.go

.PHONY: webhook-test-docker
webhook-test-docker: $(OUTPUT_DIR)/webhook-test-linux-$(DOCKER_GOARCH)
	docker build -t soloio/multicluster-admission-webhook-test:$(VERSION) $(OUTPUT_DIR) -f internal/cmd/Dockerfile;

#----------------------------------------------------------------------------------
# Helm
#----------------------------------------------------------------------------------
HELM_DIR := internal/helm/charts
HELM_ROOTDIR := $(shell pwd)/$(HELM_DIR)
HELM_OUTPUT_DIR ?= $(OUTPUT_DIR)/helm
CHART_DIR := $(HELM_DIR)/test
PACKAGED_CHART_DIR := $(HELM_OUTPUT_DIR)/test

# Generate Manifests from Helm Chart
.PHONY: chart-gen
chart-gen:
	go run internal/generate.go -chart
	goimports -w .
	go mod tidy

.PHONY: package-chart
package-chart: chart-gen
	helm package --destination $(PACKAGED_CHART_DIR) $(CHART_DIR)

.PHONY: deploy-test-chart
deploy-test-chart: package-chart
	helm install admission -n multicluster-admission $(PACKAGED_CHART_DIR)/multicluster-admission-webhook-test-$(VERSION).tgz

.PHONY: uninstall-test-chart
uninstall-test-chart:
	helm uninstall admission -n multicluster-admission

#----------------------------------------------------------------------------------
# Clean
#----------------------------------------------------------------------------------

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
