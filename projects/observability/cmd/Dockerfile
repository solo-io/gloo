FROM alpine:3.17.1
ARG DOCKER_GOARCH

RUN apk -U upgrade
COPY observability-linux-${DOCKER_GOARCH} /usr/local/bin/observability

USER 10101

ENTRYPOINT ["/usr/local/bin/observability"]