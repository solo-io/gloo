ARG ENVOY_IMAGE
FROM node:16.14.2-alpine AS node
FROM $ENVOY_IMAGE
ARG DOCKER_GOARCH
RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*
COPY --from=node /usr/lib /usr/lib
COPY --from=node /usr/local/ /usr/local/

COPY js /usr/local/bin/js
COPY gloo-linux-${DOCKER_GOARCH} usr/local/bin/gloo
RUN apk update;
RUN npm install --prefix /usr/local/bin/js
USER 10101

ENTRYPOINT ["/usr/local/bin/gloo"]
