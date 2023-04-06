#----------------------------------------------------------------------------------
# glooctl fed plugin
#----------------------------------------------------------------------------------
CLI_PLUGINS_DIR=projects/glooctl-plugins
FED_CLI_PLUGIN_DIR=$(CLI_PLUGINS_DIR)/fed

.PHONY: glooctl-fed-linux-$(GOARCH)
glooctl-fed-linux-$(GOARCH): $(OUTPUT_DIR)/glooctl-fed-linux-$(GOARCH)
$(OUTPUT_DIR)/glooctl-fed-linux-$(GOARCH): $(SOURCES)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(FED_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: glooctl-fed-darwin-$(GOARCH)
glooctl-fed-darwin-$(GOARCH): $(OUTPUT_DIR)/glooctl-fed-darwin-$(GOARCH)
$(OUTPUT_DIR)/glooctl-fed-darwin-$(GOARCH): $(SOURCES)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(FED_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: glooctl-fed-windows-$(GOARCH)
glooctl-fed-windows-$(GOARCH): $(OUTPUT_DIR)/glooctl-fed-windows-$(GOARCH).exe
$(OUTPUT_DIR)/glooctl-fed-windows-$(GOARCH).exe: $(SOURCES)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=windows go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(FED_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: build-fed-cli
build-fed-cli: install-go-tools glooctl-fed-linux-$(GOARCH) glooctl-fed-darwin-$(GOARCH) glooctl-fed-windows-$(GOARCH)

.PHONY: install-fed-cli
install-fed-cli:
	go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o ${GOPATH}/bin/glooctl-fed $(FED_CLI_PLUGIN_DIR)/cmd/main.go



##----------------------------------------------------------------------------------
## Release glooctl plugins
##----------------------------------------------------------------------------------

.PHONY: check-gsutil
check-gsutil:
ifeq (, $(shell which gsutil))
	$(error "No gsutil in $(PATH), follow the instructions at https://cloud.google.com/sdk/docs/install to install")
endif

build-and-upload-gcs-release-assets: check-gsutil build-fed-cli
# Only push assets if RELEASE is set to true
ifeq ($(RELEASE), "true")
	gsutil -m cp \
	$(OUTPUT_DIR)/glooctl-fed-linux-$(GOARCH) \
	$(OUTPUT_DIR)/glooctl-fed-darwin-$(GOARCH) \
	$(OUTPUT_DIR)/glooctl-fed-windows-$(GOARCH).exe \
	gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/$(VERSION)/
ifeq ($(ON_DEFAULT_BRANCH), "true")
	# We're on latest default git branch, so push /latest and updated install script
	gsutil -m cp -r gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/$(VERSION)/* gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/latest/
	gsutil cp projects/glooctl-plugins/fed/install/install.sh gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/install.sh
endif
endif
