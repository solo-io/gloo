SOURCES := $(shell find . -name *.go)
PKGDIRS := config/ module/ pkg/ xds/

PACKAGE_PATH:=github.com/solo-io/gloo/pkg/platform/kube

BINARY:=gloo-k8s-service-discovery

build: $(BINARY)

fmt:
	gofmt -w $(PKGDIRS)
	goimports -w $(PKGDIRS)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v -ldflags '-extldflags "-static"' -o $@ *.go

docker: $(BINARY)
	docker build -t solo-io/$(BINARY):v1.0 .

test:
	ginkgo -r -v .

clean:
	rm -f $(BINARY)
