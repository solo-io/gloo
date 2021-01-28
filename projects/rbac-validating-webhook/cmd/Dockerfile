FROM alpine:3.11

COPY gloo-fed-rbac-validating-webhook-linux-amd64 /usr/local/bin/gloo-fed-rbac-validating-webhook

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-rbac-validating-webhook"]
