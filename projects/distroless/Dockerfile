ARG PACKAGE_DONOR_IMAGE
ARG BASE_IMAGE

FROM --platform=linux/amd64 $BASE_IMAGE AS distroless-amd64
FROM --platform=linux/arm64 $BASE_IMAGE AS distroless-arm64

# Collect the distroless packages that have newer versions in the Debian package donor, the openssl CLI
# and the certificates that distroless strips from ca-certificates. See sync-packages.sh for details.
FROM --platform=linux/amd64 $PACKAGE_DONOR_IMAGE AS delta-amd64
ENV DEBIAN_FRONTEND=noninteractive
COPY --from=distroless-amd64 / /base
COPY sync-packages.sh /sync-packages.sh
RUN bash /sync-packages.sh /base /delta openssl ca-certificates

FROM --platform=linux/arm64 $PACKAGE_DONOR_IMAGE AS delta-arm64
ENV DEBIAN_FRONTEND=noninteractive
COPY --from=distroless-arm64 / /base
COPY sync-packages.sh /sync-packages.sh
RUN bash /sync-packages.sh /base /delta openssl ca-certificates

FROM delta-${TARGETARCH} AS delta

FROM $BASE_IMAGE
# COPY can only add and overwrite files, so files of an old package version that the new version no
# longer ships are left in the image as orphans (not tracked by status.d). sync-packages.sh logs them.
COPY --from=delta /delta/ /
RUN ["/usr/bin/openssl", "version"]
