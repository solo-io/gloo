#!/bin/sh

set -eu

if [ -z "${GLOO_VERSION:-}" ]; then
  GLOO_VERSIONS=$(curl -sHL "Accept: application/vnd.github.v3+json" https://api.github.com/repos/solo-io/gloo/releases | jq -r '[.[] | select(.tag_name | test("-") | not) | select(.tag_name | startswith("v1.")) | .tag_name] | sort_by(.) | reverse | .[0]')
else
  GLOO_VERSIONS="${GLOO_VERSION}"
fi

if [ "$(uname -s)" = "Darwin" ]; then
  OS=darwin
else
  OS=linux
fi

arch=$(uname -m)
if [ "$arch" = "aarch64" ] || [ "$arch" = "arm64" ]; then
  GOARCH=arm64
else
  GOARCH=amd64
fi

for gloo_version in $GLOO_VERSIONS; do

tmp=$(mktemp -d /tmp/gloo.XXXXXX)
filename="glooctl-${OS}-${GOARCH}"
url="https://github.com/solo-io/gloo/releases/download/${gloo_version}/${filename}"

if curl -f ${url} >/dev/null 2>&1; then
  echo "Attempting to download glooctl version ${gloo_version}"
else
  continue
fi

(
  cd "$tmp"

  echo "Downloading ${filename}..."

  SHA=$(curl -sL "${url}.sha256" | cut -d' ' -f1)
  curl -sLO "${url}"
  echo "Download complete!, validating checksum..."
  checksum=$(openssl dgst -sha256 "${filename}" | awk '{ print $2 }')
  if [ "$checksum" != "$SHA" ]; then
    echo "Checksum validation failed." >&2
    exit 1
  fi
  echo "Checksum valid."
)

(
  cd "$HOME"
  mkdir -p ".gloo/bin"
  mv "${tmp}/${filename}" ".gloo/bin/glooctl"
  chmod +x ".gloo/bin/glooctl"
)

rm -r "$tmp"

echo "Gloo Gateway was successfully installed 🎉"
echo ""
echo "Add the gloo CLI to your path with:"
echo "  export PATH=\$HOME/.gloo/bin:\$PATH"
echo ""

echo "Now run:"
echo "  glooctl install gateway     # install gloo's function gateway functionality into the 'gloo-system' namespace"
echo "  glooctl install ingress     # install very basic Kubernetes Ingress support with Gloo into namespace gloo-system"
echo "  glooctl install knative     # install Knative serving with Gloo configured as the default cluster ingress"
echo "Please see visit the Gloo Installation guides for more:  https://docs.solo.io/gloo-edge/latest/installation/"
exit 0
done

echo "No versions of glooctl found."
exit 1
