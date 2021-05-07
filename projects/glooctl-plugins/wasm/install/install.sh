#!/usr/bin/env sh

set -eu

VERSION=${GLOOCTL_WASM_VERSION:-latest}
if [ "$(uname -s)" = "Darwin" ]; then OS="darwin"; else OS="linux"; fi

(
  cd "${HOME}" || exit 1
  mkdir -p ".gloo/bin"
  echo "Downloading glooctl-wasm version: ${VERSION}"
  curl -o .gloo/bin/glooctl-wasm \
    -sL https://storage.googleapis.com/glooctl-plugins/glooctl-wasm/${VERSION}/glooctl-wasm-${OS}-amd64
  chmod +x ".gloo/bin/glooctl-wasm"
)

echo "The glooctl wasm plugin was successfully installed 🎉"
echo ""
echo "Ensure that you have glooctl installed, if you do not, run:"
echo "  curl -sL https://run.solo.io/glooctl/install | sh "
echo ""
echo "Add the Glooctl CLI to your path with:"
echo "  export PATH=\$HOME/.gloo/bin:\$PATH"
echo ""
echo "Now run:"
echo "  glooctl wasm --help     # see the commands available to you"
echo "Please see visit the Gloo Edge website for more info:  https://docs.solo.io/gloo-edge/latest/installation/advanced_configuration/wasm/"
exit 0
