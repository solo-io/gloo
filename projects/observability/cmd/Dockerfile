FROM alpine:3.13.10

RUN apk -U upgrade
COPY observability-linux-amd64 /usr/local/bin/observability

USER 10101

ENTRYPOINT ["/usr/local/bin/observability"]