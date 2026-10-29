# Cloudbuild Configuration

We leverage [Google Cloud Build](https://cloud.google.com/build) to deploy a part of our CI/CD platform.

The main configuration is defined in the `cloudbuild.yaml` in the repository root. That build spins of sub-builds, which are defined below:


## [Publish Artifacts](publish-artifacts.yaml)
Build and publish docker images and helm chart


## [Run Tests](run-tests.yaml)
Run all tests which do not depend on Kubernetes

## Releases

Publishing a GitHub release starts [`release.yaml`](../../.github/workflows/release.yaml). It waits for approval in environment `release`, then submits [Publish Artifacts](publish-artifacts.yaml) as the release service account. Tag builds from build-bot no longer publish.

Before approving, check that the tag's commit is on its `v1.<minor>.x` branch. The workflow file runs from the tagged commit, so its own checks are not a guarantee.

To stop a mistaken release, cancel the Cloud Build: `gcloud builds cancel <id> --project solo-public`. The run summary prints the id. Cancelling the workflow run does not stop the build.
