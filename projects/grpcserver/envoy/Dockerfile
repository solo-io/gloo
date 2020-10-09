FROM envoyproxy/envoy:v1.15.2

RUN apt-get update

ENTRYPOINT ["/usr/local/bin/envoy", "-c", "/etc/envoy/config.yaml"]