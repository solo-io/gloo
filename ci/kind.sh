#!/bin/bash

# make all the docker images
# write the output to a temp file so that we can grab the image names out of it
# also ensure we clean up the file once we're done
TEMP_FILE=$(mktemp)
VERSION=kind LOCAL_BUILD=true make docker --jobs=2 | tee ${TEMP_FILE}
make cleanup-node-modules

cleanup() {
    echo ">> Removing ${TEMP_FILE}"
    rm ${TEMP_FILE}
}
trap "cleanup" EXIT SIGINT

echo ">> Temporary output file ${TEMP_FILE}"

# grab the image names out of the `make docker` output
sed -nE 's|Successfully tagged (.*$)|\1|p' ${TEMP_FILE} | while read f; do kind load docker-image --name kind $f; done

# Now that the images are loaded into kind, we can delete them locally to save some disk space
make cleanup-local-docker-images

make VERSION=kind build-test-chart build-os-with-ui-test-chart glooctl-linux-amd64
