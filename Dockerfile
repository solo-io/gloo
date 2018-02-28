FROM scratch
COPY gloo-function-discovery /
ENTRYPOINT ["/gloo-function-discovery"]