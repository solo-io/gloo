---
menuTitle: "Dev Environment"
title: "Setting up the Development Environment"
weight: 1
---

## Environment Setup

### Prerequisites

Developing on Gloo Gateway requires the following to be installed on your system:

- [`Make`](https://www.gnu.org/software/make/)
- [`Git`](https://git-scm.com/)
- [`Go`](https://golang.org/) (`solo-io` projects are built using version `1.23.3`)
- `Protoc` (`solo-io` projects are built using version `3.6.1`)
- Standard development tools like `gcc`

To install all the requirements, do the following:
```bash
# install build tools
# - macOS
xcode-select --install
# - other operating systems:
# distro meta-packages, like "build-essential", should have what is required

############################################
# install version of go in go.mod
############################################
# - macOS:
brew install go@1.23
# - other operating systems:
# follow directions at https://go.dev/doc/install

# install protoc
make install-protoc

# install go related tools
make install-go-tools
```

### Setting up the Gloo Gateway Repositories

Next, we'll clone the Gloo Gateway and Solo-Kit source code. Solo-Kit is required for code generation in Gloo Gateway. 

To clone the repository:
```bash
git clone https://github.com/solo-io/gloo
# or with SSH
git clone git@github.com:solo-io/gloo.git
```

To run the `main.go` files locally in your system make sure to have a [`Kubernetes Cluster`](https://kubernetes.io/docs/setup/) running.

You should now be able to run any `main.go` file in the Gloo Gateway repository using:

```bash
go run <path-to-cmd>/main.go
```

For example:
```bash
# run gloo locally
go run projects/gloo/cmd/main.go
# run discovery locally
go run projects/discovery/cmd/main.go
# run gateway locally
go run projects/gateway/cmd/main.go
```

Awesome! You're ready to start developing on Gloo Gateway! Check out the [Writing Upstream Plugins Guide]({{% versioned_link_path fromRoot="/guides/dev/writing-upstream-plugins" %}}) to see how to add plugins to gloo.

### Developing with a local K8s cluster (kind)

Developers without a K8s cluster are encouraged to use [Docker](https://docs.docker.com) and [Kind](https://kind.sigs.k8s.io).

Prerequisites:
* [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* [Docker](https://docs.docker.com/install/)
* [Helm](https://github.com/helm/helm)

Stand up Gloo Gateway running on Kind and begin making changes:
```bash
# prepare kind cluster, build images, and upload them
make kind-setup

# install Gloo
helm upgrade --install -n gloo-system --create-namespace gloo ./_test/gloo-1.0.1-dev.tgz --values ./test/kubernetes/e2e/tests/manifests/common-recommendations.yaml

############################################
# make changes to the code in the repo ...
############################################

# Update the Gloo Gateway installation in your cluster to apply your changes
make -B kind-build-and-load-gloo
```

### Code Generation

Follow these steps to make changes to Gloo Gateway's Protobuf-based API.

Confirm code generation works with Gloo Gateway:
```bash
make -B generated-code
echo $?
```

The `echo` output should be `0` if everything worked correctly.
