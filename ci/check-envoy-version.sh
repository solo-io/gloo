#!/bin/bash
# helper script to check that the envoy version in solo-projects matches the envoy version in gloo
EE_ENVOY_VERSION=$1
GLOO_VERSION=$(go list -m github.com/solo-io/gloo | cut -d' ' -f2)

# Rip apart ENVOY_GLOO_IMAGE variable from a particular gloo tag's Makefile. Example line:
#   ENVOY_GLOO_IMAGE ?= quay.io/solo-io/envoy-gloo:1.25.4-patch1
# since we only care about the minor version, we can just grab the 3rd field after splitting on "."
gloo_oss_envoy_minor_version=$(\
    curl -s https://raw.githubusercontent.com/solo-io/gloo/$GLOO_VERSION/Makefile\
    | grep "ENVOY_GLOO_IMAGE ?= quay.io/solo-io/envoy-gloo:"\
    | cut -d. -f3)

gloo_ee_envoy_minor_version=$(\
    echo $EE_ENVOY_VERSION\
    | cut -d. -f2)
 
if [ "$gloo_oss_envoy_minor_version" == "$gloo_ee_envoy_minor_version" ]; then
    echo "gloo and solo-projects have matching envoy minor versions.  Continuing..."
else
    echo "gloo and solo-projects have mismatched envoy minor versions.  Exiting..."
    echo "gloo          envoy minor version: $gloo_oss_envoy_minor_version"
    echo "solo-projects envoy minor version: $gloo_ee_envoy_minor_version"
    exit 1
fi
