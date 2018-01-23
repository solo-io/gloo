SOURCES := $(shell find . -name *.go)

build: glue

fmt:
	gofmt -w module/ xds/ config/
	goimports -w module/ xds/ config/

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
