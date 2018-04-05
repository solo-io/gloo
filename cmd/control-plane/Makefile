SOURCES := $(shell find . -name "*.go")
BINARY:=gloo
VERSION:=$(shell cat version)
IMAGE_TAG?=v$(VERSION)

build: $(BINARY)
build-static: $(BINARY)-static
build-debug: control-plane-debug

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v  -o $@ *.go

$(BINARY)-static: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '-extldflags "-static"' -o $@ *.go

docker: $(BINARY)-static
	docker build -t soloio/$(BINARY):$(IMAGE_TAG) .

$(BINARY)-debug: $(SOURCES)
	go build -i -gcflags "-N -l" -o $(BINARY)-debug *.go

# build a container with debug symbols
docker-debug: $(BINARY)-debug
	docker build -t soloio/$(BINARY):$(IMAGE_TAG)-debug -f Dockerfile.debug .

clean:
	rm -f $(BINARY) $(BINARY)-debug $(BINARY)-static
