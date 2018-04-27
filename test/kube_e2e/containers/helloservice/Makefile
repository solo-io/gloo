SOURCES := $(shell find . -name "*.go")

BINARY:=helloservice
IMAGE_TAG?=v0.1
build: $(BINARY)
OUTPUT_DIRECTORY ?= .

$(BINARY): $(SOURCES)
	GOOS=linux go build -v -i -o $(OUTPUT_DIRECTORY)/$@ *.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):$(IMAGE_TAG) .

clean:
	rm -f $(BINARY)
