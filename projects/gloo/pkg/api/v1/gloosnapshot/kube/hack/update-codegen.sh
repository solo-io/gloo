
#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Heavily inspired by <https://github.com/kubernetes-sigs/gateway-api/blob/main/hack/update-codegen.sh>.

readonly SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE}")"/.. && pwd)"
readonly OUTPUT_DIR="${SCRIPT_ROOT}/client"

readonly GOPATH="$(mktemp -d)"
mkdir -p $GOPATH/src/github.com/solo-io/gloo
ln -s $SCRIPT_ROOT $GOPATH/src/github.com/solo-io/gloo

readonly CODEGEN_PKG=$(go list -f '{{ .Dir }}' -m k8s.io/code-generator)
readonly GENGO_PKG=$(go list -f '{{ .Dir }}' -m k8s.io/gengo/v2)

readonly ROOT_PKG=github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot
readonly OUTPUT_PKG=${ROOT_PKG}/kube/client
readonly APIS_PKG=${ROOT_PKG}/kube/apis

readonly CLIENTSET_NAME=versioned
readonly CLIENTSET_PKG_NAME=clientset
readonly VERSIONS=( v1 )

PROJECT_INPUT_DIRS_SPACE=""
PROJECT_INPUT_DIRS_COMMA=""
for VERSION in "${VERSIONS[@]}"
do
    PROJECT_INPUT_DIRS_SPACE+="${APIS_PKG}/gloosnapshot.gloo.solo.io/gloosnapshot "
    PROJECT_INPUT_DIRS_COMMA+="${APIS_PKG}/gloosnapshot.gloo.solo.io/gloosnapshot,"
done
PROJECT_INPUT_DIRS_SPACE="${PROJECT_INPUT_DIRS_SPACE%,}" # drop trailing space
PROJECT_INPUT_DIRS_COMMA="${PROJECT_INPUT_DIRS_COMMA%,}" # drop trailing comma

readonly COMMON_FLAGS="--go-header-file ${GENGO_PKG}/boilerplate/boilerplate.go.txt"

echo "Generating clientset at ${OUTPUT_PKG}/${CLIENTSET_PKG_NAME}"
go run $CODEGEN_PKG/cmd/client-gen \
    --clientset-name "${CLIENTSET_NAME}" \
    --input-base "${APIS_PKG}" \
    --input "${PROJECT_INPUT_DIRS_COMMA//${APIS_PKG}/}" \
    --output-dir "$OUTPUT_DIR/${CLIENTSET_PKG_NAME}" \
    --output-pkg "${OUTPUT_PKG}/${CLIENTSET_PKG_NAME}" \
    ${COMMON_FLAGS}

echo "Generating listers at ${OUTPUT_PKG}/listers"
go run $CODEGEN_PKG/cmd/lister-gen \
    --output-dir "$OUTPUT_DIR/listers" \
    --output-pkg "${OUTPUT_PKG}/listers" \
    ${COMMON_FLAGS} \
    ${PROJECT_INPUT_DIRS_COMMA}

echo "Generating informers at ${OUTPUT_PKG}/informers"
go run $CODEGEN_PKG/cmd/informer-gen \
    --versioned-clientset-package "${OUTPUT_PKG}/${CLIENTSET_PKG_NAME}/${CLIENTSET_NAME}" \
    --listers-package "${OUTPUT_PKG}/listers" \
    --output-dir "$OUTPUT_DIR/informers" \
    --output-pkg "${OUTPUT_PKG}/informers" \
    ${COMMON_FLAGS} \
    ${PROJECT_INPUT_DIRS_COMMA}
