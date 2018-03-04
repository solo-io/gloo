SOURCES := $(shell find . -name *.go)

BINARY:=helloservice

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v  -o $@ *.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):v0.1 .

clean:
	rm -f $(BINARY)
