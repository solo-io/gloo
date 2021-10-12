# Regression tests
This directory contains tests that install Gloo Edge Enterprise and runs regression tests.

*Note: All commands should be run from the root directory of the Solo-Projects repository*

## Setup
For these tests to run, we require the following conditions:
- Gloo Edge Enterprise Helm chart archive be present in the `_test` folder,
- `glooctl` be built in the`_output` folder
- kind cluster set up and loaded with the images to be installed by the helm chart

#### Use the CI Install Script (preferred)

[`ci/setup-kind.sh`](`https://github.com/solo-io/solo-projects/blob/master/ci/setup-kind.sh`) gets run in CI to setup the test environment for the above requirements.
It accepts a number of environment variables, to control the creation of a kind cluster and deployment of Gloo resources to that kind cluster.

| Name                  | Default    | Description |
| ---                   |   ---      |    ---      |
| CLUSTER_NAME          | kind       | The name of the cluster that will be generated |
| CLUSTER_NODE_VERSION  | v1.17.17@sha256:66f1d0d91a88b8a001811e2f1054af60eef3b669a9a74f9b6db871f2f1eeed00   | The version of the Node Docker image to use for booting the cluster |
| VERSION               | 0.0.0-kind | The version used to tag Gloo images that are deployed to the cluster |
| USE_FIPS              | false      | Whether to install fips compliant data plane images |

Example:
```bash
CLUSTER_NAME=solo-test-cluster CLUSTER_NODE_VERSION=v1.17.17@sha256:66f1d0d91a88b8a001811e2f1054af60eef3b669a9a74f9b6db871f2f1eeed00 VERSION=v1.0.0-solo-test USE_FIPS=false ci/setup-kind.sh
```



## Verify Your Setup
Before running your tests, it's worthwhile to verify that a cluster was created, and the proper images have been loaded.

```bash
CLUSTER_NAME=solo-test-cluster make kind-list-images
```

You should see the list of images in the cluster, including the ones you just uploaded

## Run Tests

To run the regression tests, your kubeconfig file must point to a running Kubernetes cluster.
`kubectl config current-context` should run `kind-<CLUSTER_NAME>`

#### (Option A) - Use the Make Target (preferred)

Use the same command that CI relies on:
```bash
KUBE2E_TESTS=<test-to-run> make run-ci-regression-tests
```

#### (Option B) - Use Ginkgo Directly

The make target just runs ginkgo with a set of useful flags. If you want to control the flags that are provided, you can run:
```bash
KUBE2E_TESTS=<test-to-run> ginkgo -r <other-flags>
```

#### Test Environment Variables
The below table contains the environment variables that can be used to configure the test execution.

| Name              | Default   | Description |
| ---               |   ---     |    ---      |
| KUBE2E_TESTS      | ""        | Name of the test suite to be run. Options: `'gateway', 'redis-clientside-sharding', 'gloo-mtls', 'wasm'` |
| DEBUG             | 0         | Set to 1 for debug log output |
| WAIT_ON_FAIL      | 0         | Set to 1 to prevent Ginkgo from cleaning up the Gloo Edge installation in case of failure. Useful to exec into inspect resources created by the test. A command to resume the test run (and thus clean up resources) will be logged to the output.
| TEAR_DOWN         | false     | Set to true to uninstall Gloo after the test suite completes |
| USE_FIPS          | false     | Set to true to use fips compliant data plane images |