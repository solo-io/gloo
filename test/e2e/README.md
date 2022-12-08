# End-to-end tests
This directory contains end-to-end tests that do not require kubernetes

*Note: All commands should be run from the root directory of the Gloo repository*

## (Option A) Run Envoy in a Docker Container
### Setup
For these tests to run, we require Envoy be built in a docker container. The `VERSION` env variable determines the name of the tag for that image.

```bash
VERSION=<version-name> make gloo-ee-envoy-wrapper-docker -B
```

### Run Tests
The `run-e2e-tests` make target runs ginkgo with a set of useful flags. The `ENVOY_IMAGE_TAG` environment variable must be set to the tag of the `gloo-ee-envoy-wrapper` Docker image you wish to run for this target:


Example:
```bash
ENVOY_IMAGE_TAG=<version-name> make run-e2e-tests
```

### Flags
You can use the following flags for the `ratelimit_test.go` to turn off key features.  This is done because the Focus will be ignored and run all `BeforeEach` functions, thus creating the containers, but not tearing them down. Set to `1`.
- DO_NOT_RUN_AEROSPIKE=1
- DO_NOT_RUN_REDIS=1
- DO_NOT_RUN_DYNAMO=1




