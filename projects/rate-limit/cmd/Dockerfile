FROM alpine:3.13.5

RUN apk -U upgrade && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY rate-limit-linux-amd64 /usr/local/bin/rate-limit

USER 10101

ENTRYPOINT ["/usr/local/bin/rate-limit"]
