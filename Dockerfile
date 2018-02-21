FROM scratch
COPY gloo-ingress /
ENTRYPOINT ["/gloo-ingress"]