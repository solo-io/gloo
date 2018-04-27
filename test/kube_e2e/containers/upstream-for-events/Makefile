SOURCES := $(shell find . -name "*.go")

BINARY:=upstream-for-events
IMAGE_TAG?=v0.1
OUTPUT_DIRECTORY ?= .

build: $(BINARY)

$(BINARY): $(SOURCES)
	GOOS=linux go build -v -i -o $(OUTPUT_DIRECTORY)/$@ *.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):$(IMAGE_TAG) .

clean:
	rm -f $(BINARY)
