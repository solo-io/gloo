FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY extauth-linux-amd64 /usr/local/bin/extauth

ENTRYPOINT ["/usr/local/bin/extauth"]