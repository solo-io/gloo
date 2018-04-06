ROOTDIR := $(shell pwd)
PROTOS := $(shell find api/v1 -name "*.proto")
SOURCES := $(shell find . -name "*.go")
GENERATED_PROTO_FILES := $(shell find pkg/api/types/v1 -name "*.pb.go") docs/api.json

#----------------------------------------------------------------------------------
# Build
#----------------------------------------------------------------------------------

BINARIES ?= control-plane function-discovery kube-ingress-controller kube-upstream-discovery

DOCKER_USER=soloio

.PHONY: build
build: $(BINARIES)

docker: $(foreach BINARY,$(BINARIES),$(shell echo $(BINARY)-docker))
proto: $(GENERATED_PROTO_FILES)

$(GENERATED_PROTO_FILES): $(PROTOS)
	export DISABLE_SORT=1 && \
	cd api/v1/ && \
	mkdir -p $(ROOTDIR)/pkg/api/types/v1 && \
	protoc \
	-I=. \
	-I=$(GOPATH)/src \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/ \
	--plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc \
    --doc_out=$(ROOTDIR)/docs/ \
    --doc_opt=json,api.json \
	--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:\
	$(ROOTDIR)/pkg/api/types/v1 \
	./*.proto

define BINARY_TARGETS
$(eval VERSION := $(shell cat cmd/$(BINARY)/version))
$(eval IMAGE_TAG ?= v$(VERSION))

$(BINARY): $(PREREQUISITES)
	CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '-extldflags "-static"' -o $(BINARY) cmd/$(BINARY)/main.go
$(BINARY)-debug: $(PREREQUISITES)
	go build -i -gcflags "-N -l" -o $(BINARY)-debug cmd/$(BINARY)/main.go
$(BINARY)-docker: $(BINARY)
	docker build -t $(DOCKER_USER)/$(BINARY):$(IMAGE_TAG) -f cmd/$(BINARY)/Dockerfile .
$(BINARY)-docker-push: $(BINARY)
	docker push $(DOCKER_USER)/$(BINARY):$(IMAGE_TAG)
$(BINARY)-docker-debug: $(BINARY)-debug
	docker build -t $(DOCKER_USER)/$(BINARY)-debug:$(IMAGE_TAG) -f cmd/$(BINARY)/Dockerfile.debug .
endef

PREREQUISITES := $(SOURCES) $(GENERATED_PROTO_FILES)
$(foreach BINARY,$(BINARIES),$(eval $(BINARY_TARGETS)))


#----------------------------------------------------------------------------------
# Docs
#----------------------------------------------------------------------------------

doc: proto
	go run docs/gen_docs.go
#	godocdown pkg/api/types/v1/ > docs/go.md

site: doc
	mkdocs build

docker-docs: site
	docker build -t $(DOCKER_USER)/nginx-docs:v$(VERSION) -f Dockerfile.site .

#----------------------------------------------------------------------------------
# Test
#----------------------------------------------------------------------------------

hackrun: $(BINARY)
	./hack/run-local.sh

unit:
	ginkgo -r -v pkg/ xds/

e2e:
	ginkgo -r -v test/

test: e2e unit





# TODO: dependnencies
# binaries:
#  make
#  protoc
#  go
#  protoc-gen-doc ilackarms version
#  docker
#  mkdocs

# libs
#  libproto

# go packages#
#  github.com/gogo/protobuf