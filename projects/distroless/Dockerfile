ARG PACKAGE_DONOR_IMAGE
ARG BASE_IMAGE
ARG TARGETARCH

FROM --platform=linux/amd64 $PACKAGE_DONOR_IMAGE AS donor-amd64
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install --no-install-recommends -y openssl

FROM --platform=linux/arm64 $PACKAGE_DONOR_IMAGE AS donor-arm64
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install --no-install-recommends -y openssl

FROM $PACKAGE_DONOR_IMAGE AS donor

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install --no-install-recommends -y ca-certificates openssl

FROM $BASE_IMAGE

# ca-certificates depends on openssl
COPY --from=donor /etc/ssl /etc/ssl
COPY --from=donor /usr/lib/ssl /usr/lib/ssl
COPY --from=donor /usr/bin/c_rehash usr/bin/c_rehash
COPY --from=donor /usr/bin/openssl /usr/bin/openssl

# openssl depends on libssl3 (OpenSSL 3) and provider modules
COPY --from=donor-amd64 /usr/lib/x86_64-linux-gnu/ossl-modules /usr/lib/x86_64-linux-gnu/ossl-modules
COPY --from=donor-amd64 /usr/lib/x86_64-linux-gnu/libcrypto.so.3 /usr/lib/x86_64-linux-gnu/libcrypto.so.3
COPY --from=donor-amd64 /usr/lib/x86_64-linux-gnu/libssl.so.3 /usr/lib/x86_64-linux-gnu/libssl.so.3

COPY --from=donor-arm64 /usr/lib/aarch64-linux-gnu/ossl-modules /usr/lib/aarch64-linux-gnu/ossl-modules
COPY --from=donor-arm64 /usr/lib/aarch64-linux-gnu/libcrypto.so.3 /usr/lib/aarch64-linux-gnu/libcrypto.so.3
COPY --from=donor-arm64 /usr/lib/aarch64-linux-gnu/libssl.so.3 /usr/lib/aarch64-linux-gnu/libssl.so.3

# Copy over the certs
COPY --from=donor /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=donor /usr/share/doc/ca-certificates/ /usr/share/doc/ca-certificates/
COPY --from=donor /usr/sbin/update-ca-certificates /usr/sbin/update-ca-certificates