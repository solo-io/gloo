ARG BASE_IMAGE

FROM rancher/kubectl:v1.34.3 as kubectl

FROM $BASE_IMAGE

RUN apk upgrade --update-cache

COPY --from=kubectl /bin/kubectl /usr/local/bin/

USER 10101