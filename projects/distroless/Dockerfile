ARG PACKAGE_DONOR_IMAGE
ARG BASE_IMAGE
ARG TARGETARCH

FROM $PACKAGE_DONOR_IMAGE AS donor

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install --no-install-recommends -y ca-certificates openssl

FROM $BASE_IMAGE

# The distroless base ships the OpenSSL 3 runtime libraries (libssl.so.3,
# libcrypto.so.3 and the provider modules) but not the openssl CLI, so the
# binary and its support files are taken from the donor.
COPY --from=donor /etc/ssl /etc/ssl
COPY --from=donor /usr/lib/ssl /usr/lib/ssl
COPY --from=donor /usr/bin/c_rehash usr/bin/c_rehash
COPY --from=donor /usr/bin/openssl /usr/bin/openssl

# Copy over the certs
COPY --from=donor /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=donor /usr/share/doc/ca-certificates/ /usr/share/doc/ca-certificates/
COPY --from=donor /usr/sbin/update-ca-certificates /usr/sbin/update-ca-certificates
