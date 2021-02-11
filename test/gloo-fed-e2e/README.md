# gloo-fed E2e Tests
Gloo Fed E2E Tests

## Prerequisites
Before running any of these tests, run:
```
LICENSE_KEY=$LICENSE_KEY ./projects/gloo-fed/ci/setup-kind.sh local remote
```

Ensure that your glooctl version is >= v1.6.0.

This suite of tests will only run if both REMOTE_CLUSTER_CONTEXT
and LOCAL_CLUSTER_CONTEXT are set.

This suite of tests is mainly intended to be run in CI.