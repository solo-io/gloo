FROM alpine:3.13.5
ARG EXTRA_PACKAGES

RUN apk update && apk add ca-certificates $EXTRA_PACKAGES && rm -rf /var/cache/apk/*
COPY extauth-linux-amd64 /usr/local/bin/extauth

USER 10101

ENTRYPOINT ["/usr/local/bin/extauth"]