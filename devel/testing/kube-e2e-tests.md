# Kubernetes End-to-End tests
These end-to-end tests install each of the 3 Gloo Gateway flavors (`gateway`, `ingress`, and `knative`) in a Kubernetes cluster, and run
end-to-end tests against them.

- [Background](#background)
  - [How do the tests work?](#how-do-the-tests-work)
- [CI](#ci)
- [Local Development](#local-development)

## Background
Kubernetes may be relied on for scheduling, persistence or security. These tests validate that Gloo Gateway can successfully operate within a Kubernetes cluster.

### How do the tests work?
Kubernetes end-to-end tests:
1. Install Gloo Gateway in Kubernetes cluster [using Helm](https://github.com/solo-io/gloo/blob/1f457f4ef5f32aedabc58ef164aeea92acbf481e/test/kube2e/gateway/gateway_suite_test.go#L84)
1. Apply Gloo resources using Kubernetes resource clients
1. Execute requests against the Envoy proxy and confirm the expected response. This validates that the Gloo resources have been picked up by the controllers, were been translated correctly into Envoy configuration, the configuration was sent to the Envoy proxy, and the proxy behaves appropriately.

## CI
These tests are run by a [GitHub action](/.github/workflows/regression-tests.yaml) as part of our CI pipeline.

If a test fails, you can retry it from a [browser window](https://docs.github.com/en/actions/managing-workflow-runs/re-running-workflows-and-jobs#reviewing-previous-workflow-runs). If you do this, please make sure to comment on the Pull Request with a link to the failed logs for debugging purposes.

## Local Development
See the [kubernetes e2e test README](/test/kubernetes/e2e/README.md) and [debugging guide](/test/kubernetes/e2e/debugging.md) for more details about running these tests.

See the [kube2e test README](/test/kube2e/README.md) for more details about running the legacy tests.