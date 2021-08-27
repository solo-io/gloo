FROM alpine:3.13.5

RUN apk -U upgrade
COPY gloo-fed-linux-amd64 /usr/local/bin/gloo-fed

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed"]
