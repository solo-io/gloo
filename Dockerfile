FROM alpine:3.7
COPY gloo /
ENTRYPOINT ["/gloo"]
