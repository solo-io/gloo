FROM alpine:3.15.4

RUN apk -U upgrade
COPY gloo-fed-linux-amd64 /usr/local/bin/gloo-fed

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed"]
