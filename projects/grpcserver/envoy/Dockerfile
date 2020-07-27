FROM envoyproxy/envoy:v1.14.3
RUN apt-get update
COPY envoy_config_grpcserver.yaml /etc/envoy.yaml
CMD /usr/local/bin/envoy -c /etc/envoy.yaml