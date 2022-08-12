FROM alpine:3.13.11
ARG GOARCH

RUN apk -U upgrade
COPY gloo-fed-apiserver-linux-${GOARCH} /usr/local/bin/gloo-fed-apiserver

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
