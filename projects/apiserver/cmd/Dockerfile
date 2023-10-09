ARG ENVOY_IMAGE
FROM $ENVOY_IMAGE
ARG GOARCH

# Update our deps to make cve toil lower
# install wget for our default probes
RUN apt-get update \
  && apt-get upgrade -y \
  && apt-get install wget -y \  
  && apt-get clean 

COPY gloo-fed-apiserver-linux-${GOARCH} /usr/local/bin/gloo-fed-apiserver

USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
