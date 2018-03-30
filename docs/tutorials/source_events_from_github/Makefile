BINARY:=image-pusher
IMAGE_TAG?=latest

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -i -v -o $@ $@.go

docker: $(BINARY)
	docker build -t soloio/$(BINARY):$(IMAGE_TAG) .
