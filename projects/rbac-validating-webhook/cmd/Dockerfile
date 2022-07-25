FROM alpine:3.13.11
RUN apk -U upgrade
COPY gloo-fed-rbac-validating-webhook-linux-amd64 /usr/local/bin/gloo-fed-rbac-validating-webhook

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-rbac-validating-webhook"]
