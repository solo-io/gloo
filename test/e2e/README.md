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

Also, you can run a test to validate that we can handle sending large xds snapshots over GRPC with the flag: `RUN_XDS_SCALE_TESTS`. 
These tests are disabled by default for performance reasons but should be run on changes to the xds clients and added to any 
nightly suites. 

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
