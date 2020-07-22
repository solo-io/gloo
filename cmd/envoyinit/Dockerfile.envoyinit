ARG ENVOY_IMAGE
FROM $ENVOY_IMAGE

COPY envoyinit-linux-amd64 /usr/local/bin/envoyinit
COPY docker-entrypoint.sh /

ARG COMMIT=v3.2.0
ARG REPO=SpiderLabs/owasp-modsecurity-crs

RUN apk upgrade --update-cache \
    && apk add git \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /modsecurity \
    && cd /modsecurity \
    && git clone https://github.com/${REPO}.git owasp-modsecurity-crs \
    && cd owasp-modsecurity-crs \
    && git checkout -qf ${COMMIT} \
    && cd /modsecurity \
    && mkdir -p rules/crs && cp owasp-modsecurity-crs/rules/*.conf rules/crs \
    && mkdir data && cp owasp-modsecurity-crs/rules/*.data data \
    && rm -rf owasp-modsecurity-crs

WORKDIR /modsecurity/data

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/docker-entrypoint.sh"]
CMD []