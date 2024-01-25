FROM alpine:3.17.6

ARG DOCKER_GOARCH

RUN apk -U upgrade

COPY sds-linux-$DOCKER_GOARCH /usr/local/bin/sds

USER 10101

ENTRYPOINT ["/usr/local/bin/sds"]
