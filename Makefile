SOURCES := $(shell find . -name *.go)
BINARY:=gloo
build-debug: gloo-debug

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v  -o $@ *.go

docker: $(BINARY)
	docker build -t solo-io/$(BINARY):v0.1 .

$(BINARY)-debug: $(SOURCES)
	go build -i -gcflags "-N -l" -o $(BINARY)-debug *.go

hackrun: $(BINARY)
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit

clean:
	rm -f $(BINARY) $(BINARY)-debug


api:
	export OUTDIR=$(PWD)/docs && make -C $(GOPATH)/src/github.com/solo-io/gloo-api proto

doc: api
	go run docs/gen_docs.go
#	godocdown pkg/api/types/v1/ > docs/go.md