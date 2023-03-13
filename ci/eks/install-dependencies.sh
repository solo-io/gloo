#!/bin/bash
set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

git config --global url."https://${GITHUB_TOKEN}@github.com/solo-io/".insteadOf "https://github.com/solo-io/"

$SCRIPT_DIR/install-kubectl.sh
$SCRIPT_DIR/install-helm.sh
$SCRIPT_DIR/install-glooctl.sh