# Cloudbuild Configuration

We leverage [Google Cloud Build](https://cloud.google.com/build) to deploy a part of our CI/CD platform.

The main configuration is defined in the `cloudbuild.yaml` in the repository root. That build spins of sub-builds, which are defined below:


## [Publish Artifacts](publish-artifacts.yaml)
Build and publish docker images and helm chart. build-bot does not submit it; it runs from the release and PR image workflows below.


## [Run Tests](run-tests.yaml)
Run all tests which do not depend on Kubernetes

## Releases

Publishing a GitHub release starts [`release.yaml`](../../.github/workflows/release.yaml). It waits for approval in environment `release`, then submits [Publish Artifacts](publish-artifacts.yaml) as the release service account. Tag builds from build-bot no longer publish.

Before approving, check that the tag's commit is on its `v1.<minor>.x` branch. The workflow file runs from the tagged commit, so its own checks are not a guarantee.

To stop a mistaken release, cancel the Cloud Build: `gcloud builds cancel <id> --project solo-public`. The log of the "Submit and wait for Cloud Build" step prints the id and the command as soon as the build starts; the run summary shows them only after the job ends. Cancelling the workflow run does not stop the build.

Only a version with no suffix or `-patchN` is a regular release. Mark every other suffix (for example `-beta1`, `-rc1`, `-fips`) as a pre-release when creating the release. Validation fails otherwise, so edit the release and re-run the failed jobs.

For a new release line (for example `v1.23.x`), release-owners create the branch from the newest line, so it carries the workflows. On the new branch, add it to `on.pull_request_target.branches` in `.github/workflows/pr-images.yaml`. A repo admin adds the exact branch name to the deployment branches of environment `release`, and IT adds it to the PR-images branch list of the `gloo` workload identity provider. Branch names are listed exactly; patterns come later. `release.yaml` accepts tags up to `v1.29`.

## PR images

To get images for a PR, the PR author adds the label `publish-pr-images`. A label or push by anyone else doesn't start a build. [`pr-images.yaml`](../../.github/workflows/pr-images.yaml) then waits for a second person (not the author) to approve in environment `release` before it submits [Publish Artifacts](publish-artifacts.yaml) as the PR build service account. Fork PRs whose author can't add labels get no PR images.

The build runs the PR's code with the quay.io robot logged in, so approvers check the PR's Makefile, `ci/` and Dockerfile changes first. The run title shows the PR head commit; approve only if it is the commit you reviewed. Approvers must not be the PR author or have pushed commits to the PR (prevent self-review only blocks whoever started the run).

Images are tagged `<version>-<PR number>`. A push by the PR author asks for approval again. Pushes by anyone else, including the `keep pr updated` base merges, don't rebuild, so the images can lag the PR head. PR builds that don't push (`run-tests.yaml`, including the glooctl build) still run from build-bot without approval.
