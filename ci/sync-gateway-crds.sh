#!/bin/bash

set -euo pipefail

# Function to source specific variables from the setup-kind.sh script
source_specific_vars() {
    local file="./ci/kind/setup-kind.sh"
    local vars_to_source=("CONFORMANCE_VERSION" "CONFORMANCE_CHANNEL")
    
    for var in "${vars_to_source[@]}"; do
        eval $(grep "^$var=" "$file")
    done
}

# Source the required variables
source_specific_vars

# Capitalize the first letter of CONFORMANCE_CHANNEL
CAPITALIZED_CHANNEL="$(echo "${CONFORMANCE_CHANNEL:0:1}" | tr '[:lower:]' '[:upper:]')${CONFORMANCE_CHANNEL:1}"

# Define output directory and filenames
OUT_DIR="${OUT_DIR:-projects/gateway2/crds}"
OUT_FILENAME="${OUT_FILENAME:-gateway-crds.yaml}"
TCPROUTE_FILENAME="tcproute-crd.yaml"

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
    
    if [ -f "$dest" ]; then
        # Prepend header to the temporary file
        echo "$HEADER" > "$dest.tmp.full"
        cat "$dest.tmp" >> "$dest.tmp.full"

        # Compare files ignoring the header
        if compare_files_no_header "$dest" "$dest.tmp.full"; then
            echo "No changes detected in $dest."
            rm "$dest.tmp" "$dest.tmp.full"
            return
        fi
    fi

    # Update the file with the new content and header
    echo "$HEADER" > "$dest"
    cat "$dest.tmp" >> "$dest"
    rm "$dest.tmp" "$dest.tmp.full"
    echo "Updated $dest."
}

# Download Gateway API CRDs (leave as is)
download_gateway_crd "$GATEWAY_CRD_URL" "${OUT_DIR}/${OUT_FILENAME}"

# Download TCPRoute CRD (manage header)
download_tcproute_crd "$TCPROUTE_CRD_URL" "${OUT_DIR}/${TCPROUTE_FILENAME}"

echo "CRD sync complete."
