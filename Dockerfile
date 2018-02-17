FROM busybox
COPY gloo /
ENTRYPOINT ["/gloo"]