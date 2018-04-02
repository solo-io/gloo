dirs:=image-pusher star-tweeter slack-bot
IMAGE_TAG?=latest

.PHONY: all $(dirs)
all: $(dirs)

$(dirs):
	cd $@ && \
	CGO_ENABLED=0 GOOS=linux go build -i -v -o $@ *.go && \
	docker build -t soloio/$@:$(IMAGE_TAG) .
