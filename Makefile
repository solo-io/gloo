
proto:
	cd api/v1/ && \
	protoc \
	-I=. \
	-I=$(GOPATH)/src \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/ \
	--plugin=protoc-gen-doc=$(HOME)/go/bin/protoc-gen-doc \
    --doc_out=../../docs \
    --doc_opt=markdown,api.md \
	--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:\
	../../pkg/api/types/v1 \
	./*.proto

docs/go.md:
	godocdown pkg/api/types/v1/ > docs/go.md

doc: docs/go.md proto