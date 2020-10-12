FROM alpine

RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY grpcserver-linux-amd64 /usr/local/bin/grpcserver

USER 10101

ENTRYPOINT ["/usr/local/bin/grpcserver"]