FROM envoyproxy/envoy:v1.11.0
RUN apt-get update
COPY envoy_config_grpcserver.yaml /etc/envoy.yaml
CMD /usr/local/bin/envoy -c /etc/envoy.yaml