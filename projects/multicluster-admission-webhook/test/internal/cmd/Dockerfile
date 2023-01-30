FROM alpine:3.17.1

ARG DOCKER_GOARCH

# installs public root certs
# needed to establish trust with third party sources
RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY webhook-test-linux-${DOCKER_GOARCH} /usr/local/bin/webhook-test

ENTRYPOINT ["/usr/local/bin/webhook-test"]
