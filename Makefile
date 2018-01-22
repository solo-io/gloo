SOURCES := $(shell find . -name *.go)

build: mockway

fmt:
	gofmt -w module/ xds/ config/
	goimports -w module/ xds/ config/

mockway: $(SOURCES)
	go build -o mockway

run: mockway
	./mockway -f module/example/example_config.yml

clean:
	rm -f mockway
