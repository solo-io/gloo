SOURCES := $(shell find . -name *.go)
BINARY:=gloo-function-discovery
VERSION:=$(shell cat version)

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v -o $@ *.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):v$(VERSION) .

test:
	ginkgo -r -v .

clean:
	rm -f $(BINARY)
