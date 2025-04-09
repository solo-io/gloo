#!/bin/sh

set -eu

if [ -z "${GLOO_VERSION:-}" ]; then
  GLOO_VERSIONS=$(curl -sHL "Accept: application/vnd.github.v3+json" https://api.github.com/repos/solo-io/gloo/releases | jq -r '[.[] | select(.tag_name | test("-") | not) | select(.tag_name | startswith("v1.")) | .tag_name] | sort_by(.) | reverse | .[0]')
else
  GLOO_VERSIONS="${GLOO_VERSION}"
fi

for gloo_version in $GLOO_VERSIONS; do
echo "${gloo_version}"
done
