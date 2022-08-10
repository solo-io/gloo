FROM alpine:3.15.5
ARG GOARCH

RUN apk -U upgrade
COPY gloo-fed-linux-${GOARCH} /usr/local/bin/gloo-fed

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed"]
