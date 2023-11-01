ARG ENVOY_IMAGE
FROM $ENVOY_IMAGE as base


FROM frolvlad/alpine-glibc:alpine-3.17_glibc-2.34
ENV loglevel=info

RUN apk update;
RUN apk -U upgrade

RUN mkdir -p /etc/envoy
COPY --from=base /usr/local/bin/envoy /usr/local/bin/envoy
COPY --from=base /usr/local/bin/libsaxon-solo.so /usr/local/bin/libsaxon-solo.so

RUN apk update;
RUN apk -U upgrade

COPY gloo-fed-apiserver-linux-amd64 /usr/local/bin/gloo-fed-apiserver

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
