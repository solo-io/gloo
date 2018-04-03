SOURCES := $(shell find . -name "*.go")
BINARY:=gloo
VERSION:=$(shell cat version)
IMAGE_TAG?=v$(VERSION)

build: $(BINARY)
build-static: $(BINARY)-static
build-debug: gloo-debug


$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v  -o $@ *.go

$(BINARY)-static: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '-extldflags "-static"' -o $@ *.go

docker: $(BINARY)-static
	docker build -t soloio/$(BINARY):$(IMAGE_TAG) .

$(BINARY)-debug: $(SOURCES)
	go build -i -gcflags "-N -l" -o $(BINARY)-debug *.go

# build a container with debug symbols
docker-debug: $(BINARY)-debug
	docker build -t soloio/$(BINARY):$(IMAGE_TAG)-debug -f Dockerfile.debug .

hackrun: $(BINARY)
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit

clean:
	rm -f $(BINARY) $(BINARY)-debug $(BINARY)-static

proto:
	export OUTDIR=$(PWD)/docs && make -C $(GOPATH)/src/github.com/solo-io/gloo-api proto

doc: proto
	go run docs/gen_docs.go
#	godocdown pkg/api/types/v1/ > docs/go.md

site: doc
	mkdocs build 

docker-docs: site
	docker build -t soloio/nginx-docs:v$(VERSION) -f Dockerfile.site .
