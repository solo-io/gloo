SOURCES := $(shell find . -name *.go)

BINARY:=upstream-for-events
IMAGE_TAG?=v0.1

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v  -o $@ *.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):$(IMAGE_TAG) .

clean:
	rm -f $(BINARY)
