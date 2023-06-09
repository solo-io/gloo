ARG ENVOY_IMAGE
FROM $ENVOY_IMAGE
ARG GOARCH

RUN apk update;
RUN apk -U upgrade

COPY gloo-fed-apiserver-linux-${GOARCH} /usr/local/bin/gloo-fed-apiserver

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
