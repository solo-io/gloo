FROM alpine:3.13.11
ARG GOARCH

RUN apk -U upgrade
COPY gloo-fed-rbac-validating-webhook-linux-${GOARCH} /usr/local/bin/gloo-fed-rbac-validating-webhook

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-rbac-validating-webhook"]
