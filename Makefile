SOURCES := $(shell find . -name *.go)
PKGDIRS := config/ module/ pkg/ xds/

PACKAGE_PATH:=github.com/solo-io/glue

# kubernetes custom clientsets
clientset:
	./vendor/k8s.io/code-generator/generate-groups.sh all \
		$(PACKAGE_PATH)/config/watcher/crd/client \
		$(PACKAGE_PATH)/config/watcher/crd \
		"solo.io:v1"

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
