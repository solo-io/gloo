SOURCES := $(shell find . -name *.go)
BINARY:=gloo-function-discovery
VERSION:=$(shell cat version)
IMAGE_TAG?=v$(VERSION)

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v -o $@ *.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):$(IMAGE_TAG) .

test:
	ginkgo -r -v .

clean:
	rm -f $(BINARY)
