ARG ENVOY_IMAGE
# envoy image is needed because we make calls to envoy in Gloo
# please check out https://github.com/solo-io/gloo/blob/main/projects/gloo/pkg/bootstrap/bootstrap_validation.go#L51
# Update our deps to make cve toil lower
# install wget for our default probes
FROM $ENVOY_IMAGE
ARG DOCKER_GOARCH
RUN apt-get update \
  && apt-get install --no-install-recommends -y ca-certificates \
  && apt-get upgrade -y \
  && apt-get install wget -y \
  && apt-get clean 

COPY gloo-linux-${DOCKER_GOARCH} usr/local/bin/gloo

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo"]
