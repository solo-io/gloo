FROM alpine:3.17.6

ARG DOCKER_GOARCH

# Needed for access to AWS
RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY certgen-linux-$DOCKER_GOARCH /usr/local/bin/certgen

USER 10101

ENTRYPOINT ["/usr/local/bin/certgen"]