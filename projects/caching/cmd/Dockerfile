FROM alpine:3.17.3
ARG EXTRA_PACKAGES

RUN apk -U upgrade && apk add ca-certificates $EXTRA_PACKAGES && rm -rf /var/cache/apk/*
COPY caching-linux-amd64 /usr/local/bin/caching

USER 10101

ENTRYPOINT ["/usr/local/bin/caching"]
