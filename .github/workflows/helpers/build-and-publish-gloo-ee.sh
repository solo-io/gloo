# !/bin/bash
# --------------------------------------
# This script builds and publishes a gloo-ee release...
#   * images:           https://quay.io/organization/solo-io
#   * helm chart:       https://console.cloud.google.com/storage/browser/gloo-ee-test-helm
# ...based on the current branch
# --------------------------------------
VERSION=$(go run .github/workflows/helpers/find-latest-local-version.go)-$(git rev-parse --abbrev-ref HEAD)-$(git rev-parse --short HEAD)
# for example: VERSION=$(1.14.0-beta2)-$(4052-one-click)-$(7c8df00ef)
HELM_BUCKET="gs://gloo-ee-test-helm"
HELM_SYNC_DIR="_output/helm"
HELM_DIR="install/helm"

# update local reference to gloo
GLOO_REPO_OVERRIDE=""
if [ -z "$GLOO_SHA" ]
then
      echo "GLOO_SHA is unset; skipping gloo update"
else
      echo "GLOO_SHA is set to $GLOO_SHA; updating go.mod"
      go get github.com/solo-io/gloo@$GLOO_SHA
      go mod tidy
      GLOO_REPO_OVERRIDE="https://storage.googleapis.com/gloo-ee-test-helm"
fi

# build and push images
VERSION=$VERSION make build-kind-images-non-fips -B
VERSION=$VERSION TAGGED_VERSION=$VERSION make docker-push-non-fips -B

# create appropriate Values.yaml and Chart.yaml files
VERSION=$VERSION GLOO_REPO_OVERRIDE=$GLOO_REPO_OVERRIDE make init-helm

# Complicated block ripped from gloo-ee Makefile.  Roughly, this block...
#   1. Grabs GENERATION id of helm repo
#   2. Downloads helm index.yaml
#   3. Packages local helm chart and merges it into the local index.yaml
#   4. Uploads the local helm chart and index.yaml to the helm repo
#   5. If the helm repo has been updated since the last download, the upload will fail
until $(GENERATION=$(gsutil ls -a $HELM_BUCKET/index.yaml | tail -1 | cut -f2 -d '#') && \
                gsutil cp -v $HELM_BUCKET/index.yaml $HELM_SYNC_DIR/index.yaml && \
                helm package --destination $HELM_SYNC_DIR/charts $HELM_DIR/gloo-ee >> /dev/null && \
                helm repo index $HELM_SYNC_DIR --merge $HELM_SYNC_DIR/index.yaml && \
                gsutil -m rsync $HELM_SYNC_DIR/charts $HELM_BUCKET/charts && \
                gsutil -h x-goog-if-generation-match:"$GENERATION" cp $HELM_SYNC_DIR/index.yaml $HELM_BUCKET/index.yaml); do \
    echo "Failed to upload new helm index (updated helm index since last download?). Trying again"; \
    sleep 2; \
done

# provide (hopefully) useful output
echo "Successfully published a test build of gloo-ee!" > published-gloo-ee.txt
echo "  Version:    $VERSION" >> published-gloo-ee.txt
echo "  Helm Repo:  https://console.cloud.google.com/storage/browser/gloo-ee-test-helm" >> published-gloo-ee.txt
echo "  Image Repo: https://quay.io/organization/solo-io" >> published-gloo-ee.txt
echo "Can Install Via:" >> published-gloo-ee.txt
echo "❯ helm repo add gloo-ee-test https://storage.googleapis.com/gloo-ee-test-helm" >> published-gloo-ee.txt
echo "❯ helm repo update" >> published-gloo-ee.txt
echo "❯ helm install -n gloo-system gloo-ee-test gloo-ee-test/gloo-ee --create-namespace --version $VERSION --set-string license_key=\$GLOO_LICENSE_KEY --set gloo-fed.enabled=false --set gloo-fed.glooFedApiserver.enable=false" >> published-gloo-ee.txt
cat published-gloo-ee.txt
