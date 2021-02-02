FROM envoyproxy/envoy:v1.13.0
RUN apt-get update
COPY ./cfg.yaml /etc/envoy.yaml
USER 10101
CMD /usr/local/bin/envoy -c /etc/envoy.yaml