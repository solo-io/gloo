FROM alpine:3.13.10
ARG EXTRA_PACKAGES
ARG DOCKER_GOARCH

RUN apk -U upgrade && apk add ca-certificates $EXTRA_PACKAGES && rm -rf /var/cache/apk/*
COPY rate-limit-linux-${DOCKER_GOARCH} /usr/local/bin/rate-limit

USER 10101

ENTRYPOINT ["/usr/local/bin/rate-limit"]
