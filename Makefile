SOURCES := $(shell find . -name *.go)

build: glue

fmt:
	gofmt -w module/ xds/ config/
	goimports -w module/ xds/ config/

glue: $(SOURCES)
	go build -o glue

run: glue
	./glue -f module/example/example_config.yml

clean:
	rm -f glue
