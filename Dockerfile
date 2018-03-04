FROM scratch
COPY gloo-ingress-controller /
ENTRYPOINT ["/gloo-ingress-controller"]