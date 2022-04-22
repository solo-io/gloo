FROM alpine:3.15.4
ARG EXTRA_PACKAGES

RUN apk -U upgrade && apk add ca-certificates $EXTRA_PACKAGES && rm -rf /var/cache/apk/*
COPY rate-limit-linux-amd64 /usr/local/bin/rate-limit

USER 10101

ENTRYPOINT ["/usr/local/bin/rate-limit"]
