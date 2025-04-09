#!/bin/bash

set -euo pipefail

# The file containing CONFORMANCE_VERSION and CONFORMANCE_CHANNEL environment variables
SETUP_KIND_FILE="./ci/kind/setup-kind.sh"

# Function to source specific variables from the setup-kind.sh script
source_specific_vars() {
    local vars_to_source=("CONFORMANCE_VERSION" "CONFORMANCE_CHANNEL")

    for var in "${vars_to_source[@]}"; do
        eval $(grep "^$var=" "$SETUP_KIND_FILE")
    done
}

# Function to update the CONFORMANCE_VERSION variable in setup-kind.sh
update_conformance_version_in_setup_kind() {
    local key="CONFORMANCE_VERSION"
    local value="$1"

    if grep -q "^${key}=" "$SETUP_KIND_FILE"; then
        current_value=$(grep "^${key}=" "$SETUP_KIND_FILE" | sed -E 's/^.*:-([^}]+)}/\1/')
        if [ "$current_value" != "$value" ]; then
            echo "Updating $key in $SETUP_KIND_FILE from \"${current_value}\" to \"${value}\"..."
            sed -i.bak "s|^\(${key}=\"\${${key}:-\)[^\}]*\(}\"\)|\1${value}\2|" "$SETUP_KIND_FILE"
            rm "$SETUP_KIND_FILE.bak"
            echo "Updated $key to \"${value}\"."
        else
            echo "$key is already set to \"${value}\"."
        fi
    else
        echo "$key not found in $SETUP_KIND_FILE. Adding it..."
        echo "${key}=\"\${${key}:-${value}}\"" >> "$SETUP_KIND_FILE"
        echo "Added $key with value \"${value}\"."
    fi
}

# Source the required variables
source_specific_vars

# Update CONFORMANCE_VERSION in ./ci/kind/setup-kind.sh if needed
update_conformance_version_in_setup_kind "$CONFORMANCE_VERSION"

# Capitalize the first letter of CONFORMANCE_CHANNEL
CAPITALIZED_CHANNEL="$(echo "${CONFORMANCE_CHANNEL:0:1}" | tr '[:lower:]' '[:upper:]')${CONFORMANCE_CHANNEL:1}"

# Define output directory and filenames
OUT_DIR="${OUT_DIR:-projects/gateway2/crds}"
OUT_FILENAME="${OUT_FILENAME:-gateway-crds.yaml}"
TCPROUTE_FILENAME="${TCPROUTE_FILENAME:-"tcproute-crd.yaml"}"

# Create the output directory if it doesn't exist
mkdir -p "${OUT_DIR}"

# URLs for CRDs
GATEWAY_CRD_URL="https://github.com/kubernetes-sigs/gateway-api/releases/download/${CONFORMANCE_VERSION}/${CONFORMANCE_CHANNEL}-install.yaml"
TCPROUTE_CRD_URL="https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/refs/tags/${CONFORMANCE_VERSION}/config/crd/experimental/gateway.networking.k8s.io_tcproutes.yaml"

# Header to prepend to the TCPRoute CRD file
HEADER=$(cat <<EOF
# Copyright 2024 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# Gateway API ${CAPITALIZED_CHANNEL} channel install
#
---
EOF
)

# Function to compare files ignoring headers (for TCPRoute only)
compare_files_no_header() {
    local file1="$1"
    local file2="$2"

    # Strip the header from both files and compare the rest of the content
    tail -n +$(($(echo "$HEADER" | wc -l) + 1)) "$file1" > "$file1.stripped"
    tail -n +$(($(echo "$HEADER" | wc -l) + 1)) "$file2" > "$file2.stripped"
    cmp -s "$file1.stripped" "$file2.stripped"
    local result=$?
    rm -f "$file1.stripped" "$file2.stripped"
    return $result
}

download_gateway_crd() {
    local url="$1"
    local dest="$2"

    echo "Downloading Gateway CRD from $url to $dest..."
    curl -sLo "$dest.tmp" "$url"

    if [ -f "$dest" ] && cmp -s "$dest" "$dest.tmp"; then
        echo "No changes detected in $dest."
        rm "$dest.tmp"
    else
        mv "$dest.tmp" "$dest"
        echo "Updated $dest."
    fi
}

download_tcproute_crd() {
    local url="$1"
    local dest="$2"

    echo "Downloading TCPRoute CRD from $url to $dest..."
    curl -sLo "$dest.tmp" "$url"

    # Always create the temporary file with the header
    echo "$HEADER" > "$dest.tmp.full"
    cat "$dest.tmp" >> "$dest.tmp.full"

    if [ -f "$dest" ]; then
        # Compare files ignoring the header
        if compare_files_no_header "$dest" "$dest.tmp.full"; then
            echo "No changes detected in $dest."
            rm "$dest.tmp" "$dest.tmp.full"
            return
        fi
    fi

    # Update the file with the new content and header
    mv "$dest.tmp.full" "$dest"
    rm "$dest.tmp"
    echo "Updated $dest."
}

# Update sigs.k8s.io/gateway-api in go.mod
update_go_mod() {
    local module="sigs.k8s.io/gateway-api"
    if grep -q "$module" go.mod; then
        current_version=$(grep "$module" go.mod | awk '{print $2}')
        # Use pattern matching to check if the current version is a prefix of the conformance version as we're currently using an rc version
        if [[ "$current_version" != ${CONFORMANCE_VERSION}* ]]; then
            echo "Updating $module from $current_version to $CONFORMANCE_VERSION..."
            go get "$module@$CONFORMANCE_VERSION"
            go mod tidy
            echo "Updated $module to $CONFORMANCE_VERSION."
        else
            echo "$module is already at version $CONFORMANCE_VERSION."
        fi
    else
        echo "$module not found in go.mod."
    fi
}

# Update k8sgateway_api_version in nightly-tests/max_versions.env
update_max_versions_env() {
    local env_file=".github/workflows/.env/nightly-tests/max_versions.env"
    local key="k8sgateway_api_version="
    if [ -f "$env_file" ]; then
        if grep -q "$key" "$env_file"; then
            current_version=$(grep "$key" "$env_file" | cut -d '=' -f 2 | tr -d "'")
            if [ "$current_version" != "$CONFORMANCE_VERSION" ]; then
                echo "Updating $key in $env_file from '$current_version' to '$CONFORMANCE_VERSION'..."
                sed -i.bak "s|^$key.*|$key'$CONFORMANCE_VERSION'|" "$env_file"
                rm "$env_file.bak"
                echo "Updated $key to '$CONFORMANCE_VERSION'."
            else
                echo "$key is already set to '$CONFORMANCE_VERSION'."
            fi
        else
            echo "$key not found in $env_file. Adding it..."
            echo "$key'$CONFORMANCE_VERSION'" >> "$env_file"
            echo "Added $key with value '$CONFORMANCE_VERSION'."
        fi
    else
        echo "$env_file not found."
    fi
}

# Download Gateway API CRDs (leave as is)
download_gateway_crd "$GATEWAY_CRD_URL" "${OUT_DIR}/${OUT_FILENAME}"

# Download TCPRoute CRD (manage header)
download_tcproute_crd "$TCPROUTE_CRD_URL" "${OUT_DIR}/${TCPROUTE_FILENAME}"

# Update dependencies and environment
update_go_mod
update_max_versions_env

echo "Gateway API sync complete."
