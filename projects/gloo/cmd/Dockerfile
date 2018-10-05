FROM alpine

RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY gloo-linux-amd64 /usr/local/bin/gloo

ENTRYPOINT ["/usr/local/bin/gloo"]