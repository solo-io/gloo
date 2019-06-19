FROM alpine

RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY gloo-grpcserver-linux-amd64 /usr/local/bin/gloo-grpcserver

ENTRYPOINT ["/usr/local/bin/gloo-grpcserver"]