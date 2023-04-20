# Usage

## Setup

The following programs must be installed in order to follow the subsequent guidelines:
 - Make 
   - Install on Mac OS using `brew install make`
   - Install on Ubuntu using `sudo apt-get install -y make`
 - Golang v1.20
   - Follow the instructions for your operating system at https://golang.org/doc/install to install
 - Docker v20.10.x
   - Follow the instructions for your operating system at https://docs.docker.com/get-docker/ to install
 - Git
   - Follow the instructions for your operating system at https://git-scm.com/book/en/v2/Getting-Started-Installing-Git to install
 - Helm
   - Follow the instructions for your operating system at https://helm.sh/docs/intro/install/ to install

Run the following command to prepare your local environment for development once you have the dependencies listed above installed
```bash
# install build dependencies
make install-go-tools
# set up git hooks
make init
```

 - Note: If you do not have access to private Solo.io Github repositories for whatever reason, you will not be able to execute the above initialization commands.

## Building Docker Images

You can use the following command to build all Gloo Edge Enterprise Docker images:
```bash
TAGGED_VERSION=<VERSION> IMAGE_REG='quay.io/solo-io' make docker
```

The following variables can be set to affect the output:
 |Variable Name|Default Value|Description|
 |-------------|-------------|-----------|
 |`TAGGED_VERSION`|The current Semver, as determined by the repo's `HEAD` commit|The version of Gloo Edge Enterprise that a deposit is being made for. Must follow [Semver format](https://semver.org/) and be prefixed with a literal "v", i.e., v1.8.1. Included in the tag of all built images.|
 |`IMAGE_REG`|quay.io/solo-io|Docker image repository that built images are tagged for. Just put your usernameto refer to your account at hub.docker.com. To use gcr images, set the `IMAGE_REG` to `gcr.io/<PROJECT_NAME>`|

The `docker` make target will build the following images:
|Image Name|Command to Build Individual Image|Dependencies|
|----------|---------------------------------|------------|
|rate-limit-ee|`make rate-limit-ee-docker`||
|rate-limit-ee-fips|`make rate-limit-ee-fips-docker`||
|extauth-ee|`make extauth-ee-docker`||
|extauth-ee-fips|`make extauth-ee-fips-docker`||
|gloo-ee|`make gloo-ee-docker`|gcr.io/gloo-ee/envoy-gloo-ee|
|gloo-ee-fips|`make gloo-fips-ee-docker`|gcr.io/gloo-ee/envoy-gloo-ee-fips|
|gloo-ee-envoy-wrapper|`make gloo-ee-envoy-wrapper-docker`|gcr.io/gloo-ee/envoy-gloo-ee|
|gloo-ee-envoy-wrapper-fips|`make gloo-ee-envoy-wrapper-fips-docker`|gcr.io/gloo-ee/envoy-gloo-ee-fips|
|observability-ee|`make observability-ee-docker`||
|ext-auth-plugins|`make ext-auth-plugins-docker`||
|ext-auth-plugins-fips|`make ext-auth-plugins-fips-docker`||
|gloo-fed|`make gloo-fed-docker`||
|gloo-fed-apiserver|`make gloo-fed-apiserver-docker`||
|gloo-fed-apiserver-envoy|`make gloo-fed-apiserver-envoy-docker`||
|gloo-federation-console|`make gloo-federation-console-docker`||
|gloo-fed-rbac-validating-webhook|`make gloo-fed-rbac-validating-webhook-docker`||

### Notes:
 - You will need the gcr.io/gloo-ee/envoy-gloo-ee Docker image in order to build the gloo-ee and gloo-ee-envoy-wrapper images
   - You can set the `ENVOY_GLOO_IMAGE` environment variable to `gcr.io/gloo-ee/envoy-gloo-ee:<tag>` to point to a locally tagged version of this image
   - Ex: `TAGGED_VERSION=<VERSION> IMAGE_REG='quay.io/solo-io' ENVOY_GLOO_IMAGE=gcr.io/gloo-ee/envoy-gloo-ee:1.19.0-patch4 make docker`
 - You will need the gcr.io/gloo-ee/envoy-gloo-ee-fips Docker image in order to build the gloo-ee-fips and gloo-ee-envoy-wrapper-fips images
   - You can set the `ENVOY_GLOO_FIPS_IMAGE` environment variable to `gcr.io/gloo-ee/envoy-gloo-ee-fips:<tag>` to point to a locally tagged version of this image
   - Ex: `TAGGED_VERSION=<VERSION> IMAGE_REG='quay.io/solo-io' ENVOY_GLOO_FIPS_IMAGE=gcr.io/gloo-ee/envoy-gloo-ee-fips:1.19.0-patch4 make docker`
 - `make docker` attempts to build all of the images listed above, and will stop execution if any of the images fails to build. You can prevent this behavior by calling `make` with the `-i` flag, which will build as many images as possible, even if some fail.

## Pushing Docker Images

You can use the following command to push all Gloo Edge Enterprise Docker images to the repository specified by `IMAGE_REG`:
```bash
TAGGED_VERSION=<VERSION> IMAGE_REG='quay.io/solo-io' make docker-push
```
 - Please note that the above command will have no effect if `TAGGED_VERSION` is not explicitly set
 - You should run the prior `make docker` step with the same values of `TAGGED_VERSION` and `IMAGE_REG` set in order to ensure that the proper images are pushed
 - You can push an individual image using the following command:
    ```bash
    docker push $IMAGE_REG/<IMAGE NAME>:$VERSION
    ```
    Where `<IMAGE NAME>` is one of the values listed in the prior `make docker` step.

### Development Tip: Loading Docker Images into Kind Cluster for Local Development

If using a [Kind](https://kind.sigs.k8s.io/) Cluster for development or local testing, you can sidestep pushing your Docker images to a remote repository. Instead, use the following command to load local images directly into the repository:
```bash
CLUSTER_NAME=<CLUSTER NAME> TAGGED_VERSION=<VERSION> IMAGE_REG='quay.io/solo-io' make push-kind-images
```

Where `CLUSTER_NAME` is set to the name of the Kind cluster, and `TAGGED_VERSION` and `IMAGE_REG` are set to the values used during the prior `make docker` step.

To load Gloo Fed images to kind use the `make gloofed-load-kind-images` target:
```bash
CLUSTER_NAME=<CLUSTER NAME> TAGGED_VERSION=<VERSION> IMAGE_REG='quay.io/solo-io' make gloofed-load-kind-images
```

## Building Helm Charts

You can use the following command to build helm charts for Gloo Edge Enterprise and Gloo Fed which reference the images built in the previous steps:
```bash
TAGGED_VERSION=<VERSION> IMAGE_REG='quay.io/solo-io' make build-test-chart
```

Where `TAGGED_VERSION` and `IMAGE_REG` are set to the values used during the prior `make docker` and `make docker-push` steps.

The above command will create two compressed helm charts in the `_test` directory, `_test/gloo-ee-<TAGGED_VERSION>.tgz`, which points to the Gloo Edge Enterprise Chart, and `_test/gloo-fed-<TAGGED_VERSION>.tgz`, which points to the Gloo Fed chart. The uncompressed versions of these charts can be found in `install/helm/gloo-ee` and `install/helm/gloo-fed`

You can install these helm charts using the following command:
```bash
helm install <NAME> <PATH_TO_CHART>
```

Where `<NAME>` is the name of the chart to install, typically set to either `gloo` or `gloo-fed`

For example, the following command can be used to install the generated Gloo Edge Enterprise Chart:
```bash
# In this case, Docker images and Helm Charts have been built with
# TAGGED_VERSION=v1.9.0-beta10
helm install gloo _test/gloo-ee-1.9.0-beta10.tgz

# Install the uncompressed chart. Has the same affect as the above command
helm install gloo install/helm/gloo-ee
```

## Updated instructions for the grpcserver

### prep

- get the right version of protoc (3.6.1)
  - the make target below will warn you if you need to update

### build

```bash
make update-ui-deps
make generated-ui
make run-apiserver
```

## Building `extauth` components locally
We build the `extauth` binaries inside a [docker container](projects/extauth/cmd/Dockerfile) for reproducibility. 
Since it needs to access private git repositories, the container relies on a GitHub token to be provided via the 
`GITHUB_TOKEN` environment variable. If you need to build the `extauth` binaries locally, you have to generate a token. 
You can do that by opening the settings page for your GitHub account and navigating to 
`Developer Settings > Personal access tokens > Generate new token`. Once you have a token, you can export it to your 
environment and run the desired `make` target, e.g.:

``` 
export GITHUB_TOKEN=<your token> 
make extauth
```

## Noteworthy make targets

- `docker`: builds all images
  - for local builds, set `LOCAL_BUILD` to `true`
  - when running locally, should set `LOCAL_BUILD=1` in order to build the ui resources
  - may want to set `VERSION` env var to `kind`
- `push-kind-images`: pushes images built by `make docker` target to your kind cluster
  - requires `CLUSTER_NAME` env var set. default kind cluster is named `kind`
- `build-test-chart` and `build-test-chart-fed`: zipped helm chart saved in the `_test` dir
  - may want to set `VERSION` env var to `kind`

## Additional Notes

- Shared projects across Solo.io.
- This repo contains the git history for Gloo and Solo-Kit.
- Make sure to follow the [developer guide](https://github.com/solo-io/dev-docs/blob/master/new_hire_guide.md#dev-environment-setup-guides) for IDE and git config setup.

## Helm Repositories
- [GlooE](https://console.cloud.google.com/storage/browser/gloo-ee-helm)
- [Gloo with read-only UI](https://console.cloud.google.com/storage/browser/gloo-os-ui-helm)
