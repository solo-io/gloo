SOURCES := $(shell find . -name *.go)
BINARY:=gloo
VERSION:=$(shell cat version)

build-debug: gloo-debug

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v  -o $@ *.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):v$(VERSION) .

$(BINARY)-debug: $(SOURCES)
	go build -i -gcflags "-N -l" -o $(BINARY)-debug *.go

# build a container with debug symbols
docker-debug: $(BINARY)-debug
	docker build -t soloio/$(BINARY):v$(VERSION)-debug -f Dockerfile.debug .

hackrun: $(BINARY)
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit

clean:
	rm -f $(BINARY) $(BINARY)-debug

proto:
	export OUTDIR=$(PWD)/docs && make -C $(GOPATH)/src/github.com/solo-io/gloo-api proto

doc: proto
	go run docs/gen_docs.go
#	godocdown pkg/api/types/v1/ > docs/go.md

site: doc
	mkdocs build && \
	docker build -t soloio/nginx-docs:v$(VERSION) -f Dockerfile.site .