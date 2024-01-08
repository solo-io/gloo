FROM alpine:3.17.6

ARG DOCKER_GOARCH

RUN apk -U upgrade

COPY ingress-linux-$DOCKER_GOARCH /usr/local/bin/ingress

USER 10101

ENTRYPOINT ["/usr/local/bin/ingress"]