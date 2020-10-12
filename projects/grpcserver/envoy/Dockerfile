FROM envoyproxy/envoy:v1.15.2

RUN apt-get update

USER 10101

ENTRYPOINT ["/usr/local/bin/envoy", "-c", "/etc/envoy/config.yaml"]
