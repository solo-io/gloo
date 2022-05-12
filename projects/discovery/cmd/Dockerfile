FROM alpine:3.15.4

ARG DOCKER_GOARCH

# Needed for access to AWS
RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY discovery-ee-linux-${DOCKER_GOARCH} /usr/local/bin/discovery

USER 10101

ENTRYPOINT ["/usr/local/bin/discovery"]