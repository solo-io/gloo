FROM alpine:3.13.5

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY extauth-linux-amd64 /usr/local/bin/extauth

USER 10101

ENTRYPOINT ["/usr/local/bin/extauth"]