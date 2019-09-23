FROM soloio/envoy-gloo-ee:0.1.17

COPY envoyinit-linux-amd64 /usr/local/bin/envoyinit


ARG COMMIT=v3.1/dev
ARG REPO=SpiderLabs/owasp-modsecurity-crs

RUN apk upgrade --update-cache \
    && apk add git \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /modsecurity \
    && cd /modsecurity \
    && git clone https://github.com/${REPO}.git owasp-modsecurity-crs-3.1 \
    && cd owasp-modsecurity-crs-3.1 \
    && git checkout -qf ${COMMIT} \
    && cd /modsecurity \
    && mkdir -p rules/general \
    && mkdir -p rules/crs && cp owasp-modsecurity-crs-3.1/rules/*.conf rules/crs \
    && mkdir data && cp owasp-modsecurity-crs-3.1/rules/*.data data \
    && rm -rf  owasp-modsecurity-crs-3.1

WORKDIR /modsecurity/data

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/usr/local/bin/envoyinit"]
CMD []