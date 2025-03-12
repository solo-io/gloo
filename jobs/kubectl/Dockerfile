ARG BASE_IMAGE

FROM bitnami/kubectl:1.27.15 as kubectl

FROM $BASE_IMAGE

RUN apk upgrade --update-cache

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/
