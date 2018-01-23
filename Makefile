SOURCES := $(shell find . -name *.go)

build: glue

fmt:
	gofmt -w module/ xds/ config/
	goimports -w module/ xds/ config/

glue: $(SOURCES)
	go build -o glue

run: glue
	./hack/run-local.sh

kill:
	killall envoy

clean:
	rm -f glue
