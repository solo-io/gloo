# Regression tests
This directory contains regression tests for GlooE.

# How to run tests

- `docker`: builds all images
  - for local builds, set `LOCAL_BUILD` to `true`
  - when running locally, should set `LOCAL_BUILD=1` in order to build the ui resources
  - may want to set `VERSION` env var to `kind`
- `push-kind-images`: pushes images built by `make docker` target to your kind cluster
  - requires `CLUSTER_NAME` env var set. default kind cluster is named `kind`
- `build-test-chart` and `build-os-with-ui-test-chart`: zipped helm chart saved in the `_test` dir
  - may want to set `VERSION` env var to `kind`

Then run the tests from the appropriate directory with the right env vars, e.g.:
```shell script
	KUBE2E_TESTS=$KUBE_TEST_VAL ginkgo -r -failFast -trace -progress -race -compilers=4 -failOnPending -noColor ./test/regressions/...
```
where `$KUBE_TEST_VAL`, is set for the test suite you want to run, i.e. one of:
- 'gateway'
- 'gloomtls'
- 'redis-clientside-sharding'
- 'wasm'