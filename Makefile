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

build-debug: gloo-debug

build: gloo

fmt:
	gofmt -w $(PKGDIRS)
	goimports -w $(PKGDIRS)

gloo: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o $@ cmd/*.go

docker: gloo
	docker build -t solo-io/gloo:v1.0

gloo-debug: $(SOURCES)
	go build -i -gcflags "-N -l" -o gloo-debug cmd/*.go

hackrun: gloo
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit

clean:
	rm -f gloo

