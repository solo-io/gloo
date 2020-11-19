FROM alpine:3.12.1

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY rate-limit-linux-amd64 /usr/local/bin/rate-limit

USER 10101

ENTRYPOINT ["/usr/local/bin/rate-limit"]