SOURCES := $(shell find . -name *.go)

BINARY:=upstream-for-events

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v  -o $@ *.go

docker: $(BINARY)
	docker build -t solo-io/$(BINARY):v0.1 .

clean:
	rm -f $(BINARY)
