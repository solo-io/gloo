proto:
	export DISABLE_SORT=1 && \
	cd api/v1/ && \
	protoc \
	-I=. \
	-I=$(GOPATH)/src \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/ \
	--plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc \
    --doc_out=$(OUTDIR) \
    --doc_opt=json,api.json \
	--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:\
	../../pkg/api/types/v1 \
	./*.proto
