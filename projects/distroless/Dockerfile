# Since we only publish amd64 images, we use the amd64 variant.

ARG PACKAGE_DONOR_IMAGE
ARG BASE_IMAGE

FROM --platform=linux/amd64 $PACKAGE_DONOR_IMAGE as donor
ENV DEBIAN_FRONTEND=noninteractive
# only get the things we say we want
RUN echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf.d/00-keepclean
RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf.d/00-keepclean

# retry apt-get update and install to work around transient network issues
RUN echo 'Acquire::Retries "3";' >> /etc/apt/apt.conf.d/01-bad-sources
RUN echo 'Acquire::http::timeout "80";' >> /etc/apt/apt.conf.d/01-bad-sources
RUN echo 'Acquire::https::timeout "80";' >> /etc/apt/apt.conf.d/01-bad-sources
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install --no-install-recommends -y ca-certificates

FROM --platform=linux/amd64 $BASE_IMAGE

# ca-certificates depends on openssl
COPY --from=donor /etc/ssl /etc/ssl
COPY --from=donor /usr/lib/ssl /usr/lib/ssl
COPY --from=donor /usr/bin/c_rehash usr/bin/c_rehash
COPY --from=donor /usr/bin/openssl /usr/bin/openssl

# openssl depends on libssl1.1
COPY --from=donor /usr/lib/x86_64-linux-gnu/engines-1.1 /usr/lib/x86_64-linux-gnu/engines-1.1
COPY --from=donor /usr/lib/x86_64-linux-gnu/libcrypto.so.1.1 /usr/lib/x86_64-linux-gnu/libcrypto.so.1.1
COPY --from=donor /usr/lib/x86_64-linux-gnu/libssl.so.1.1 /usr/lib/x86_64-linux-gnu/libssl.so.1.1

# Copy over the certs
COPY --from=donor /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=donor /usr/share/doc/ca-certificates/ /usr/share/doc/ca-certificates/
COPY --from=donor /usr/sbin/update-ca-certificates /usr/sbin/update-ca-certificates