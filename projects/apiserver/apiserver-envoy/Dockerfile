FROM envoyproxy/envoy:v1.16.2
RUN apt-get update && apt-get -y upgrade
USER 10101
CMD /usr/local/bin/envoy -c /etc/envoy/config.yaml