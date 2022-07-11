ARG ENVOY_IMAGE
FROM $ENVOY_IMAGE
ARG DOCKER_GOARCH
RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY gloo-linux-${DOCKER_GOARCH}  /usr/local/bin/gloo

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo"]