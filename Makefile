SOURCES := $(shell find . -name *.go)
PKGDIRS := config/ module/ pkg/ xds/

PACKAGE_PATH:=github.com/solo-io/glue-storage

# kubernetes custom clientsets
clientset:
	cd ${GOPATH}/src/k8s.io/code-generator && \
	./generate-groups.sh all \
		$(PACKAGE_PATH)/crd/client \
		$(PACKAGE_PATH)/crd \
		"solo.io:v1"

clean:
	rm -rf $(PACKAGE_PATH)/crd/client
