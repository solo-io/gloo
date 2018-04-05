ROOTDIR := $(shell pwd)

#----------------------------------------------------------------------------------
# Build
#----------------------------------------------------------------------------------

proto:
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

#----------------------------------------------------------------------------------
# Docs
#----------------------------------------------------------------------------------

doc: proto
	go run docs/gen_docs.go
#	godocdown pkg/api/types/v1/ > docs/go.md

site: doc
	mkdocs build

docker-docs: site
	docker build -t soloio/nginx-docs:v$(VERSION) -f Dockerfile.site .

#----------------------------------------------------------------------------------
# Test
#----------------------------------------------------------------------------------

hackrun: $(BINARY)
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit





# TODO: dependnencies
# binaries:
#  make
#  protoc
#  go
#  protoc-gen-doc
#  docker
#  mkdocs

# libs
#  libproto

# go packages#
#  github.com/gogo/protobuf