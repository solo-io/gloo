# Gloo Federation

Gloo Edge Federation allows users to manage the configuration for all of their Gloo Edge instances from one place, no matter what platform they run on.

## Background
Gloo Edge Federation elevates Gloo Edge’s powerful routing features beyond the environment they live in, allowing users to create all new global routing features between different Gloo Edge instances. Gloo Edge Federation provides:
- Consistent configuration 
- Service failover 
- Unified debugging
- Automated Gloo Edge discovery across all of your Gloo Edge instances

## Build
*All make targets are currently defined in the [Makefile](./../../Makefile) at the root of the repository.*

The `VERSION` env variable determines the name of the tag for the image.

You may either inject the version yourself:
```bash
VERSION=<version name> make gloo-fed-docker
```

Or rely on the auto-generated version:
```shell
make gloo-fed-docker
```

## Release
During a Gloo Edge Enterprise release, the `gloo-fed` image is published to the [Google Cloud Registry](https://console.cloud.google.com/gcr/images/gloo-edge/GLOBAL) and the [Quay repository](https://quay.io/repository/solo-io/gloo-fed).


## Local Development
### Prerequisites
Before setting up Gloo Fed, make sure you have performed the following steps:

* Install [Docker](https://docs.docker.com/docker-for-mac/install/)
* Install [Homebrew](https://brew.sh/) and the following packages:

      brew install go
      brew install kubectl
      brew install kind
      brew install helm
      brew install findutils
      brew install gnu-sed
      brew install yarn

* Update the PATH in your .zshrc or .bashrc file:

      export GOPATH="${HOME}/go"
      export GOROOT="$(brew --prefix golang)/libexec"
      export PATH=${GOPATH}/bin:${GOROOT}/bin:/usr/local/opt/findutils/libexec/gnubin:/usr/local/bin:$PATH
      PATH="$(brew --prefix)/opt/gnu-sed/libexec/gnubin:$PATH"

* Install protoc version 3.6.1:

On a mac:

      brew install protoc@3.6

On a linux:

      curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-osx-x86_64.zip
      unzip protoc-3.6.1-osx-x86_64.zip -d protoc3
      mv protoc3/bin/* /usr/local/bin/
      mv protoc3/include/* /usr/local/include/
      rm -r protoc-3.6.1-osx-x86_64.zip protoc3

* Add to `/etc/hosts`:

      # kind setup
      127.0.0.1       host.docker.internal

### Running and testing Gloo Fed locally
Note that since this is an Enterprise feature, you will need a Gloo License Key below:

```shell script
# Run setup-kind in order to provision two kind clusters.
# This will create two clusters, "local" and "remote".
# "local" will run the gloo-fed control plane, and "remote" will run gloo.

GLOO_LICENSE_KEY=<your license key here> ./projects/gloo-fed/ci/setup-kind.sh local remote

# Register the clusters using glooctl.
glooctl cluster register --cluster-name kind-local --remote-context kind-local --local-cluster-domain-override host.docker.internal --federation-namespace gloo-system
glooctl cluster register --cluster-name kind-remote --remote-context kind-remote --local-cluster-domain-override host.docker.internal --federation-namespace gloo-system

# (Optional) Apply some test data (as of this writing, a FederatedUpstream and FederatedVirtualService) with:

kubectl apply -f projects/gloo-fed/example/resources/
# Verify that the upstream was federated to the remote cluster
kubectl get upstream -n gloo-system --context kind-remote i-was-federated -oyaml


# Now we're all set up!
# Let's view our remote Gloo installation.

kubectl get glooinstance -n gloo-system -oyaml

# To teardown kind clusters, run

./ci/teardown-kind.sh local remote
```
