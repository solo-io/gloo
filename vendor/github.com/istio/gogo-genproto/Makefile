SHA = c8c975543a134177cc41b64cbbf10b88fe66aa1d
GOOGLEAPIS_URL = https://raw.githubusercontent.com/googleapis/googleapis/$(SHA)

GOGO_PROTO_PKG := github.com/gogo/protobuf/gogoproto
GOGO_TYPES := github.com/gogo/protobuf/types
GOGO_DESCRIPTOR := github.com/gogo/protobuf/protoc-gen-gogo/descriptor
GOGO_GOOGLEAPIS := istio.io/gogo-genproto/googleapis

importmaps := \
	gogoproto/gogo.proto=$(GOGO_PROTO_PKG) \
	google/protobuf/any.proto=$(GOGO_TYPES) \
	google/protobuf/descriptor.proto=$(GOGO_DESCRIPTOR) \
	google/protobuf/duration.proto=$(GOGO_TYPES) \
	google/protobuf/timestamp.proto=$(GOGO_TYPES) \
	google/protobuf/wrappers.proto=$(GOGO_TYPES) \

comma := ,
empty :=
space := $(empty) $(empty)
mapping_with_spaces := $(foreach map,$(importmaps),M$(map),)
MAPPING := $(subst $(space),$(empty),$(mapping_with_spaces))
PLUGIN := --plugin=protoc-gen-gogoslick=gogoslick --gogoslick_out=$(MAPPING):googleapis
PROTOC = protoc

googleapis_protos = \
	google/api/http.proto \
	google/api/annotations.proto \
	google/rpc/status.proto \
	google/rpc/code.proto \
	google/rpc/error_details.proto \
	google/type/color.proto \
	google/type/date.proto \
	google/type/dayofweek.proto \
	google/type/latlng.proto \
	google/type/money.proto \
	google/type/postal_address.proto \
	google/type/timeofday.proto \

googleapis_packages = \
	google/api \
	google/rpc \
	google/type \

all: build

vendor:
	dep ensure --vendor-only

depend: vendor
	$(foreach var,$(googleapis_packages),mkdir -p googleapis/$(var);)

protoc.version:
	# Record protoc version
	@echo `protoc --version` > protoc.version

gogoslick: depend
	@go build -o gogoslick vendor/github.com/gogo/protobuf/protoc-gen-gogoslick/main.go

$(googleapis_protos): %:
	# Download $@ at $(SHA)
	@curl -sS $(GOOGLEAPIS_URL)/$@ -o googleapis/$@.tmp
	@sed -e '/^option go_package/d' googleapis/$@.tmp > googleapis/$@
	@rm googleapis/$@.tmp

$(googleapis_packages): %: gogoslick protoc.version $(googleapis_protos)
	# Generate $@
	@$(PROTOC) $(PLUGIN) -I googleapis  googleapis/$@/*.proto

generate: $(googleapis_packages)

format: generate
	# Format code
	@gofmt -l -s -w googleapis

build: format
	# Build code
	@go build ./...

clean:
	@rm gogoslick

.PHONY: all depend format build $(googleapis_protos) $(googleapis_packages) protoc.version clean
