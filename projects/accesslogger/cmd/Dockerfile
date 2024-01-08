FROM alpine:3.17.6

ARG DOCKER_GOARCH
RUN apk -U upgrade && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY access-logger-linux-$DOCKER_GOARCH /usr/local/bin/access-logger

USER 10101

ENTRYPOINT ["/usr/local/bin/access-logger"]