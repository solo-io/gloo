SOURCES := $(shell find . -name *.go)
PKGDIRS := config/ module/ pkg/ xds/

PACKAGE_PATH:=github.com/solo-io/gloo/pkg/platform/kube

build-debug: gloo-debug

build: gloo

fmt:
	gofmt -w $(PKGDIRS)
	goimports -w $(PKGDIRS)

gloo: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v -ldflags '-extldflags "-static"' -o $@ *.go

docker: gloo
	docker build -t solo-io/gloo:v1.0 .

gloo-debug: $(SOURCES)
	go build -i -gcflags "-N -l" -o gloo-debug *.go

hackrun: gloo
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit

clean:
	rm -f gloo

