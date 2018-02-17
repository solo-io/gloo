SOURCES := $(shell find . -name *.go)
PKGDIRS := config/ module/ pkg/ xds/

PACKAGE_PATH:=github.com/solo-io/gloo/pkg/platform/kube

proto:
	cd api/v1/ && \
	protoc \
	-I=. \
	-I=$(GOPATH)/src \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/ \
	--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:\
	../../pkg/api/types/v1 \
	./*.proto

build-debug: glue-debug

build: glue

fmt:
	gofmt -w $(PKGDIRS)
	goimports -w $(PKGDIRS)

glue: $(SOURCES)
	go build -i -o glue cmd/glue/*.go

glue-debug: $(SOURCES)
	go build -i -gcflags "-N -l" -o glue-debug cmd/glue/*.go

hackrun: glue
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit

clean:
	rm -f glue

