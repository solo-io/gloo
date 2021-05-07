
#----------------------------------------------------------------------------------
# glooctl wasm plugin
#----------------------------------------------------------------------------------
CLI_PLUGINS_DIR=projects/glooctl-plugins
WASM_CLI_PLUGIN_DIR=$(CLI_PLUGINS_DIR)/wasm

.PHONY: glooctl-wasm-linux-amd64
glooctl-wasm-linux-amd64: $(OUTPUT_DIR)/glooctl-wasm-linux-amd64
$(OUTPUT_DIR)/glooctl-wasm-linux-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(WASM_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: glooctl-wasm-darwin-amd64
glooctl-wasm-darwin-amd64: $(OUTPUT_DIR)/glooctl-wasm-darwin-amd64
$(OUTPUT_DIR)/glooctl-wasm-darwin-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(WASM_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: glooctl-wasm-windows-amd64
glooctl-wasm-windows-amd64: $(OUTPUT_DIR)/glooctl-wasm-windows-amd64.exe
$(OUTPUT_DIR)/glooctl-wasm-windows-amd64.exe: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(WASM_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: build-wasm-cli
build-wasm-cli: install-go-tools glooctl-wasm-linux-amd64 glooctl-wasm-darwin-amd64 glooctl-wasm-windows-amd64

.PHONY: install-wasm-cli
install-wasm-cli:
	go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o ${GOPATH}/bin/glooctl-wasm $(WASM_CLI_PLUGIN_DIR)/cmd/main.go

#----------------------------------------------------------------------------------
# glooctl fed plugin
#----------------------------------------------------------------------------------
CLI_PLUGINS_DIR=projects/glooctl-plugins
FED_CLI_PLUGIN_DIR=$(CLI_PLUGINS_DIR)/fed

.PHONY: glooctl-fed-linux-amd64
glooctl-fed-linux-amd64: $(OUTPUT_DIR)/glooctl-fed-linux-amd64
$(OUTPUT_DIR)/glooctl-fed-linux-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(FED_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: glooctl-fed-darwin-amd64
glooctl-fed-darwin-amd64: $(OUTPUT_DIR)/glooctl-fed-darwin-amd64
$(OUTPUT_DIR)/glooctl-fed-darwin-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(FED_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: glooctl-fed-windows-amd64
glooctl-fed-windows-amd64: $(OUTPUT_DIR)/glooctl-fed-windows-amd64.exe
$(OUTPUT_DIR)/glooctl-fed-windows-amd64.exe: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(FED_CLI_PLUGIN_DIR)/cmd/main.go

.PHONY: build-fed-cli
build-fed-cli: install-go-tools glooctl-fed-linux-amd64 glooctl-fed-darwin-amd64 glooctl-fed-windows-amd64

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

build-and-upload-gcs-release-assets: check-gsutil build-wasm-cli build-fed-cli
# Only push assets if RELEASE is set to true
ifeq ($(RELEASE), "true")
	gsutil -m cp \
	$(OUTPUT_DIR)/glooctl-wasm-linux-amd64 \
	$(OUTPUT_DIR)/glooctl-wasm-darwin-amd64 \
	$(OUTPUT_DIR)/glooctl-wasm-windows-amd64.exe \
	gs://$(GCS_BUCKET)/$(WASM_GCS_PATH)/$(VERSION)/
	gsutil -m cp \
	$(OUTPUT_DIR)/glooctl-fed-linux-amd64 \
	$(OUTPUT_DIR)/glooctl-fed-darwin-amd64 \
	$(OUTPUT_DIR)/glooctl-fed-windows-amd64.exe \
	gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/$(VERSION)/
ifeq ($(ON_DEFAULT_BRANCH), "true")
	# We're on latest default git branch, so push /latest and updated install script
	gsutil -m cp -r gs://$(GCS_BUCKET)/$(WASM_GCS_PATH)/$(VERSION)/* gs://$(GCS_BUCKET)/$(WASM_GCS_PATH)/latest/
	gsutil -m cp -r gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/$(VERSION)/* gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/latest/
	gsutil cp projects/glooctl-plugins/wasm/install/install.sh gs://$(GCS_BUCKET)/$(WASM_GCS_PATH)/install.sh
	gsutil cp projects/glooctl-plugins/fed/install/install.sh gs://$(GCS_BUCKET)/$(FED_GCS_PATH)/install.sh
endif
endif
