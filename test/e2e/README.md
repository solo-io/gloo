# In Memory End-to-End tests
This directory contains end-to-end tests that do not require Kubernetes, and persist configuration in memory.

*Note: All commands should be run from the root directory of the Gloo repository*

## Background
For a high level overview of these tests, see the [Gloo Edge Open Source e2e README](https://github.com/solo-io/gloo/blob/main/test/e2e/README.md#background)

### Example Test
We have an [example test](./example_test.go) which outlines how these tests work. It also provides some examples for basic testing operations. When writing a new e2e test, we recommend reviewing the example test first.

## CI
These tests are run by [build-bot](https://github.com/solo-io/build-bot) in Google Cloud as part of our CI pipeline.

Failing tests can be retried using the build-bot [comment directives](https://github.com/solo-io/build-bot#issue-comment-directives). Be sure to include a link to the failed logs in the PR conversation for debugging purposes.

## Local Development
### Setup
For these tests to run, we require that our gateway-proxy component be previously built as a docker image.

If you have there are local changes to the component, we can rely on a previously published image and no setup is required.

However, if there are changes to the component, refer to the [Envoyinit README](https://github.com/solo-io/gloo/blob/main/projects/envoyinit) for build instructions.

### Run Tests
The `test` make target runs ginkgo with a set of useful flags. The following environment variables can be configured for this target:

| Name              | Default | Description                                                                                                                                                                                                                                        |
|-------------------|---------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ENVOY_IMAGE_TAG   | ""      | The tag of the gloo-ee-envoy-wrapper-docker image built during setup                                                                                                                                                                               |
| TEST_PKG          | ""      | The path to the package of the test suite you want to run                                                                                                                                                                                          |
| WAIT_ON_FAIL      | 0       | Set to 1 to prevent Ginkgo from cleaning up the Gloo Edge installation in case of failure. Useful to exec into inspect resources created by the test. A command to resume the test run (and thus clean up resources) will be logged to the output. |
| INVALID_TEST_REQS | fail    | The behavior for tests which depend on environment conditions that aren't satisfied. Options are `skip`, `run`, `fail`                                                                                                                             |
| SERVICE_LOG_LEVEL | ""      | The log levels used for services. See "Controlling Log Verbosity of Services" below.                                                                                                                                                               |    
| ENVOY_BINARY      | ""      | The location of the Envoy binary to use for testing. By default, the tests will parse `ENVOY_GLOO_IMAGE_VERSION` from the Makefile, pull the docker image corresponding to that tag, and extract that version of Envoy out of the retrieved image (see [here](https://github.com/solo-io/solo-projects/blob/eaffe45805c5792f9702ad803fd2066f0c5d85e3/test/services/envoy/factory.go#L47-L72)). This environment variable can be used to specify a custom built version of Envoy residing somewhere on your local file system. |

#### Controlling Log Verbosity of Services
Multiple services (Gloo, Envoy, Discovery) are executed in parallel to run these tests. By default, these services log at the `info` level. To change the log level of a service, set the `SERVICE_LOG_LEVEL` environment variable to a comma separated list of `service:level` pairs.

Options for services are:
- gateway-proxy
- gloo
- uds
- fds
- ext-auth-service
- rate-limit-service

Options for log levels are:
- debug
- info
- warn
- error

For example, to set the log level of the Gloo service to `debug` and the Ext Auth service to `error`, set:

```bash
SERVICE_LOG_LEVEL=gloo:debug,ext-auth-service:error TEST_PKG=./test/e2e/... make test
```

#### Controlling Log Verbosity of Ginkgo Runner
Ginkgo has 4 verbosity settings, whose details can be found in the [Ginkgo docs](https://onsi.github.io/ginkgo/#controlling-verbosity)

To control these settings, pass the flags using the `GINKGO_USER_FLAGS` environment variable.

For example, to set the Ginkgo runner to `very verbose` mode, set:
```bash
GINKGO_USER_FLAGS=-vv TEST_PKG=./test/e2e/... make test
```

#### Using Recently Published Image (Most Common)
This is the most common pattern. If there are no changes to the `gateway-proxy` component, and no `ENVOY_IMAGE_TAG` is specified, our tests will identify the most recently published image (for the base LTS branch) and use that version.

```bash
TEST_PKG=./test/e2e/... make test
```

#### Using Previously Published Image
To specify a particular version that was previously published, set the `ENVOY_IMAGE_TAG` accordingly.

```bash
ENVOY_IMAGE_TAG=1.13.0 TEST_PKG=./test/e2e/... make test
```

#### Using Locally Built Image
If there are changes to the component, the image will need to be rebuilt locally (see [setup tests](#setup)). After rebuilding the image, supply the tag of that image when running the tests:

```bash
ENVOY_IMAGE_TAG=0.0.1-local TEST_PKG=./test/e2e/... make test
```

### Debugging Tests
#### Use WAIT_ON_FAIL
When Ginkgo encounters a [test failure](https://onsi.github.io/ginkgo/#mental-model-how-ginkgo-handles-failure) it will attempt to clean up relevant resources, which includes stopping the running instance of Envoy, preventing the developer from inspecting the state of the Envoy instance for useful clues.

To avoid this clean up, run the test(s) with `WAIT_ON_FAIL=1`. When the test fails, it will halt execution, allowing you to inspect the state of the Envoy instance.

Once halted, use `docker ps` to determine the admin port for the Envoy instance, and follow the recommendations for [debugging Envoy](https://github.com/solo-io/gloo/tree/main/projects/envoyinit#debug), specifically the parts around interacting with the Administration interface.

#### Use INVALID_TEST_REQS
Certain tests require environmental conditions to be true for them to succeed. For example, certain tests will only run on a Linux machine.

Setting `INVALID_TEST_REQS=skip` runs all tests locally, and any tests which will not run in the local environment will be skipped. The default behavior is that they fail.

## Additional Notes
### AWS Tests
We have a setup guide for configuring the AWS credentials needed for the tests in our [Gloo E2E README](https://github.com/solo-io/gloo/blob/main/test/e2e/README.md).

Solo's AWS security has been tightened, so it _may_ not be possible to generate personal AIM credentials anymore - at least without the proper permissions. 
You can configure your local credentials using the credentials found in our [AWS start page](https://soloio.awsapps.com/start#/) by
1. Selecting the `developers` AWS account
2. Click on "Command line or programmatic access" option
3. Use the credentials shown, _including_ the Session Token
    - The tests are set up to use the session token automatically when running locally through the `os.Getenv("GCLOUD_BUILD_ID")` check.
    - _Note: From experience, these credentials update every day, so you may need to update the credentials as necessary._

You will also need to set your `AWS_SHARED_CREDENTIALS_FILE` environment variable to the **absolute path** to your AWS credentials. 
The default location where AWS stores credentials is `~/.aws/credentials`.
