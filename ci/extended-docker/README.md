# Extended Docker Script

## Running Locally

If needed for testing, run from the root of this repo:
VERSION=<your-version> TAGGED_VERSION=<your-version> IMAGE_REPO=localhost:5000 ci/extended-docker/extended-docker.sh

Do not set CREATE_ASSETS to true as the make target will run without a VERSION which is unsafe.