ARG BASE_IMAGE

FROM cgr.dev/chainguard/kubectl:latest AS kubectl

FROM $BASE_IMAGE

RUN apk upgrade --update-cache

COPY --from=kubectl /usr/bin/kubectl /usr/local/bin/kubectl

USER 10101
