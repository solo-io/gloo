FROM busybox
COPY gloo-k8s-service-discovery /
ENTRYPOINT ["/gloo-k8s-service-discovery"]