SOURCES := $(shell find . -name *.go)
PKGDIRS := config/ module/ pkg/ xds/

PACKAGE_PATH:=github.com/solo-io/glue/internal/configwatcher/kube

# kubernetes custom clientsets
clientset:
	cd ${GOPATH}/src/k8s.io/code-generator && \
	./generate-groups.sh all \
		$(PACKAGE_PATH)/crd/client \
		$(PACKAGE_PATH)/crd \
		"solo.io:v1"

proto:
	cd pkg/api/types && \
	protoc -I. --go_out=. ./types.proto

build: glue

fmt:
	gofmt -w $(PKGDIRS)
	goimports -w $(PKGDIRS)

glue: $(SOURCES)
	go build -o glue

run: glue
	./hack/run-local.sh

unit:
	ginkgo -r -v config/ module/ pkg/ xds/

e2e:
	ginkgo -r -v test/e2e/

test: e2e unit

clean:
	rm -f glue
