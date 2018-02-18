FROM busybox
COPY gloo-ingress /
ENTRYPOINT ["/gloo-ingress"]