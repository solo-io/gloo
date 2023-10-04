ARG ENVOY_IMAGE
FROM $ENVOY_IMAGE
ARG GOARCH

RUN apt-get update \
  && apt-get upgrade -y \
  && apt-get clean 

COPY gloo-fed-apiserver-linux-${GOARCH} /usr/local/bin/gloo-fed-apiserver

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
