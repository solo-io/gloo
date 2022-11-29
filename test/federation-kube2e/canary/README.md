# Gloo Federation Canary Release Tests

## Background
These tests validate that Gloo Edge can perform a canary upgrade of Gloo Federation, and simultaneously have two Gloo Federation installs reconcile the same Federated resources.

## CI
These tests are run by a [GitHub action](https://github.com/solo-io/solo-projects/blob/master/.github/workflows/regression-tests.yaml) as part of our CI pipeline.

If a test fails, you can retry it from a [browser window](https://docs.github.com/en/actions/managing-workflow-runs/re-running-workflows-and-jobs#reviewing-previous-workflow-runs). If you do this, please make sure to comment on the Pull Request with a link to the failed logs for debugging purposes.

## Local Development

### Setup
For these tests to run, we require the following conditions:
- Gloo Edge Enterprise Helm chart archive be present in the `_test` folder,
- `glooctl` be built in the`_output` folder
- kind clusters (1 management, 2 remote) set up and loaded with the images to be installed by the helm chart

#### Use the CI Install Script
`./projects/gloo-fed/ci/setup-kind-canary.sh` gets run in CI to setup the test environment for the above requirements.
It accepts a number of environment variables, to control the creation of a kind cluster and deployment of Gloo resources to that kind cluster. For now, we can just rely on the default values.

Example:
```bash
./projects/gloo-fed/ci/setup-kind-canary.sh
```