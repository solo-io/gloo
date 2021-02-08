#!/usr/bin/env sh

set -eu

VERSION=${GLOOCTL_FED_VERSION:-latest}
if [ "$(uname -s)" = "Darwin" ]; then OS="darwin"; else OS="linux"; fi

(
  cd "${HOME}" || exit 1
  mkdir -p ".gloo/bin"
  echo "Downloading glooctl-fed version: ${VERSION}"
  curl -o .gloo/bin/glooctl-fed \
    -sL https://storage.googleapis.com/glooctl-plugins/glooctl-fed/${VERSION}/glooctl-fed-${OS}-amd64
  chmod +x ".gloo/bin/glooctl-fed"
)

echo "The glooctl fed plugin was successfully installed 🎉"
echo ""
echo "Ensure that you have glooctl installed, if you do not, run:"
echo "  curl -sL https://run.solo.io/glooctl/install | sh "
echo ""
echo "Add the Glooctl CLI to your path with:"
echo "  export PATH=\$HOME/.gloo/bin:\$PATH"
echo ""
echo "Now run:"
echo "  glooctl fed --help     # see the commands available to you"
echo "Please see visit the Gloo Edge website for more info:  https://docs.solo.io/gloo-edge/master/introduction/gloo_federation/"
exit 0
