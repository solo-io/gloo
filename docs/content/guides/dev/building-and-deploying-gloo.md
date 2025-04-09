---
title: "Building and Deploying Gloo Gateway from Source"
weight: 6
---

You can build and deploy Gloo Gateway Open Source from the source code.

## Before you begin

1. Follow the [setup guide]({{% versioned_link_path fromRoot="/guides/dev/setting-up-dev-environment/" %}}) to clone the Gloo Gateway repository and install the project dependencies.
2. In your terminal, navigate to the Gloo Gateway project.
3. Continue with this guide to learn how you can build Gloo Gateway from the top-level [Makefile](https://github.com/solo-io/gloo/blob/main/Makefile).

## Set up build dependencies and code generation {#setup}

From the Gloo Gateway project directory, run:
```sh
make -B install-go-tools generated-code
```

The `install-go-tools` target installs the dependencies to build Gloo Gateway in `_output/.bin`. 

The `generated-code` target does multiple things:

1. It generates the Go source code from the `.proto` files in the `api` directories. For example, `.proto` files in `projects/gloo/api` will be generated in `projects/gloo/pkg/api`.

2. It generates the custom resource definitions (CRDs) so that Gloo Gateway can interact with Kubernetes. These CRDs can be found in `install/helm/gloo/crds`.

3. It generates [`solo-kit`](https://github.com/solo-io/solo-kit) resources, event loops, emitters, snapshots, and resource clients. These tools are denoted by `.sk.go`, and are built using `solo-kit.json` configuration files.

## Build Gloo Gateway and Docker Images {#build}

The Gloo Gateway project has several products. The code for `gloo` itself can be found in `projects/gloo`. Each product has its own set of build targets. 

**Gloo Gateway binary**

To compile the `gloo` binary to the `_output/projects/gloo/` directory, run:

    make -B gloo

**Docker image**

To build the Docker image, run the following command. Review the following table to understand the command options.

    VERSION=0.0.1 make gloo-docker -B

| Option | Description                                                                                                                                  |
| ------ |----------------------------------------------------------------------------------------------------------------------------------------------|
| `VERSION` | An optional version number for the Docker image tag, such as `0.0.1`. The format *must* be valid [semantic versioning](https://semver.org/). |
| `IMAGE_REGISTRY` | The image registry for the image, such as the local host. This value is *required* on `arm64` or `m1` machines.                              |
| `gloo-docker` | The name for the Docker image.                                                                                                               |

Example output:

    Successfully tagged quay.io/solo-io/gloo:0.0.1

## Deploy Gloo Gateway {#deploy}

You can choose from [several Gloo Gateway installation options]({{% versioned_link_path fromRoot="/installation/preparation/#deployment-requirements" %}}). This guide assumes you deploy Gloo Gateway into a Kubernetes cluster that runs locally in [kind]({{% versioned_link_path fromRoot="/installation/platform_configuration/cluster_setup/#kind" %}}).

1. Install `gloo` into the cluster, if not already present.
   ```sh
   glooctl install gateway
   ```
2. Make the image accessible to the kind cluster. 
   {{< tabs >}} 
{{% tab name="x86_64" %}}
```sh
kind load docker-image quay.io/solo-io/gloo:0.0.1
```
{{% /tab %}} 
{{% tab name="arm64 or m1" %}}
```sh
kind load docker-image quay.io/solo-io/gloo:0.0.1
```
{{% /tab %}} 
   {{< /tabs >}}
4. Update the kind cluster to use the new image as the template. Note that the image tag varies depending on your machine and the image registry and tag version that you previously used.
   {{< tabs >}} 
{{% tab name="x86_64" %}}
```sh
kubectl -n gloo-system set image deployments/gloo gloo=quay.io/solo-io/gloo:0.0.1
```
{{% /tab %}} 
{{% tab name="arm64 or m1" %}}
```sh
kubectl -n gloo-system set image deployments/gloo gloo=quay.io/solo-io/gloo:0.0.1
```
{{% /tab %}} 
   {{< /tabs >}}
3. Verify that Kubernetes removes the old pod and spins up a new pod.{{% notice tip %}}Tip: You might find it easier to observe the update in real time by using a live monitoring tool like [k9s](https://k9scli.io/).{{% /notice %}}
   ```
   $ kubectl -n gloo-system get pods
   NAME                            READY   STATUS    RESTARTS   AGE
   discovery-77d7cd499d-vg2rl      1/1     Running   0          82s
   gateway-proxy-5b88d8995-hzt27   1/1     Running   0          82s
   gloo-6578f95696-4fdth           1/1     Running   0          82s
   gloo-bdd7ddb9d-l977c            0/1     Running   0          9s
   ```
   After a few seconds, the old pod is deleted:
   ```
   $ kubectl -n gloo-system get pods
   NAME                            READY   STATUS    RESTARTS   AGE
   discovery-77d7cd499d-vg2rl      1/1     Running   0          89s
   gateway-proxy-5b88d8995-hzt27   1/1     Running   0          89s
   gloo-bdd7ddb9d-l977c            1/1     Running   0          16s
   ```

## Make and debug changes {#changes}

Now that you can build and deploy your own builds, you can start editing Gloo Gateway itself. An in-depth discussion on making code changes is outside the scope of this guide. However, a brief overview of Gloo's software architecture will make things easier to understand.

Gloo Gateway is primarily composed of *plugins*. You can find the code for these plugins in `projects/gloo/pkg/plugins`. Gloo Gateway has a lot of them! Generally speaking, these plugins take Gloo Gateway configuration as input, and translate Envoy configuration as output. You can test any changes you make to plugins by following two main steps:
1. Unit tests, including dynamic and static breakpoints
2. Deploying the code

### Unit tests

The fastest way to test and debug changes is by using the unit tests in these directories. The unit tests are fast to recompile and re-execute, and integrate smoothly with the debugger of your choice. You can find the unit tests for a plugin in its directory.

1. Navigate to the directory of a plugin, such as cross-site request forgery (CSRF).
   ```sh
   cd projects/gloo/pkg/plugins/csrf/
   ```
2. Run the plugin unit tests. You might use `ginkgo` or `gdb`.
   {{< tabs >}} 
{{% tab name="ginkgo" %}}
```sh
ginkgo
```
{{% /tab %}} 
{{% tab name="gdb" %}}
1. Create a test file for the plugin by using Go.
   ```sh
   go test -c
   ```

   This command creates a a file called `csrf.test`.
2. Run `gdb` against the test file.
   ```sh
   gdb csrf.test
   ```
{{% /tab %}} 
   {{< /tabs >}}

### Dynamic breakpoints

Dynamic breakpoints can be added via gdb without having to rebuild the unit tests.

1. List the available source files. The `plugin.go` argument is optional but helpful to filter down all the source files available.
   ```sh
   info sources plugin.go
   ```
2. Specify the breakpoint by copying the *full path* of the filename, followed by the line on which to place it. For example:
   ```
   b /home/user/go/src/github.com/solo-io/gloo/projects/gloo/pkg/plugins/csrf/plugin.go:91
   ```

### Static breakpoints

Breakpoints can be added statically by editing the source code and adding `runtime.Breakpoint()` at the desired point in the code. Then rebuild and re-run in `gdb`. `gdb` breaks when it reaches that line. Then from that line, you can go up the stack to inspect variables.

Note that adding `runtime.Breakpoint()` in the code will likely cause the Go runtime to abort, because Go does not know how to handle the break signal. Therefore, you can only run the binary in `gbd`, not Go or `ginkgo`.

### Testing in deployment

Testing in deployment can take time to set up and run. Make sure that you get the unit tests updated and passing before deploying.

1. [Rebuild the Docker image](#build).
2. [Update the cluster](#deploy).
3. Check any output to stdout by viewing the logs.
   ```sh
   kubectl -n gloo-system logs deployments/gloo
   ```
