ARG BASE_IMAGE

FROM $BASE_IMAGE


ARG GOARCH=amd64

# Needed for access to AWS
RUN apk upgrade --update-cache \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY certgen-linux-$GOARCH /usr/local/bin/certgen

USER 10101

ENTRYPOINT ["/usr/local/bin/certgen"]