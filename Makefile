
proto:
	cd api/v1/ && \
	protoc \
	-I=. \
	-I=$(GOPATH)/src \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/ \
	--plugin=protoc-gen-doc=$(HOME)/go/bin/protoc-gen-doc \
    --doc_out=../../doc \
    --doc_opt=markdown,README.md \
	--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:\
	../../pkg/api/types/v1 \
	./*.proto
