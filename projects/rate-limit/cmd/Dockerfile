FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY rate-limit-linux-amd64 /usr/local/bin/rate-limit

ENTRYPOINT ["/usr/local/bin/rate-limit"]