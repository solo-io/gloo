SOURCES := $(shell find . -name *.go)
BINARY:=gloo-k8s-service-discovery

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v -o $@ *.go

docker: $(BINARY)
	docker build -t solo-io/$(BINARY):v0.1 .

test:
	ginkgo -r -v .

clean:
	rm -f $(BINARY)
