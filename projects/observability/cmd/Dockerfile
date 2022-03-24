FROM alpine:3.13.8

RUN apk -U upgrade
COPY observability-linux-amd64 /usr/local/bin/observability

USER 10101

ENTRYPOINT ["/usr/local/bin/observability"]