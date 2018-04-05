ROOTDIR := $(shell pwd)

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

#TODO: document these dependnencies
# binaries:
#  make
#  protoc
#  go
#  protoc-gen-doc

# libs
#  libproto

# go packages#
#  github.com/gogo/protobuf