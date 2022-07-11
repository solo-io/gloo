FROM alpine:3.13.10
ARG EXTRA_PACKAGES
ARG DOCKER_GOARCH

RUN apk -U upgrade && apk add ca-certificates $EXTRA_PACKAGES && rm -rf /var/cache/apk/*
COPY extauth-linux-${DOCKER_GOARCH} /usr/local/bin/extauth

USER 10101

ENTRYPOINT ["/usr/local/bin/extauth"]