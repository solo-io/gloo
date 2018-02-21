FROM scratch
COPY gloo /
ENTRYPOINT ["/gloo"]
