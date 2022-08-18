FROM alpine:3.15.6
ARG GOARCH

RUN apk -U upgrade
COPY gloo-fed-apiserver-linux-${GOARCH} /usr/local/bin/gloo-fed-apiserver

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
