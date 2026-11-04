# Cloudbuild Configuration

We leverage [Google Cloud Build](https://cloud.google.com/build) to deploy a part of our CI/CD platform.

The main configuration is defined in the `cloudbuild.yaml` in the repository root. That build spins of sub-builds, which are defined below:


## [Publish Artifacts](publish-artifacts.yaml)
Build and publish docker images and helm chart. build-bot does not submit it; it runs from the release and PR image workflows below.


## [Run Tests](run-tests.yaml)
Run all tests which do not depend on Kubernetes

## Releases

Publishing a GitHub release starts [`release.yaml`](../../.github/workflows/release.yaml). It waits for approval in environment `release`, then submits [Publish Artifacts](publish-artifacts.yaml) as the release service account. Tag builds from build-bot no longer publish.

The second-person rule depends on environment `release` settings: required reviewers (release-owners), "Prevent self-reviews" on, and "Allow administrators to bypass configured protection rules" off. With the bypass on, an admin can publish alone. Prevent self-review blocks only whoever published the release, so a person must publish it: the workflow rejects releases published by an App, and no shared machine-user token may be able to create releases.

Before approving, check that the tag's commit is on its `v1.<minor>.x` branch. The workflow file runs from the tagged commit, so its own checks are not a guarantee.

To stop a mistaken release, cancel the Cloud Build: `gcloud builds cancel <id> --project solo-public`. The log of the "Submit and wait for Cloud Build" step prints the id and the command as soon as the build is submitted, so you can cancel it while it is queued; the run summary shows them only after the job ends. Cancelling the workflow run does not stop the build.

Only a version with no suffix or `-patchN` is a regular release. Mark every other suffix (for example `-beta1`, `-rc1`, `-fips`) as a pre-release when creating the release. Validation fails otherwise, so edit the release and re-run the failed jobs.

For a new release line (for example `v1.23.x`), release-owners create the branch from the newest line, so it carries the workflows. On the new branch, add it to the base-branch list in `.github/workflows/pr-images-comment.yaml` on `main` (a repo admin lands it past `close-main-prs.yaml`). A repo admin adds the exact branch name to the deployment branches of environment `release`, and IT adds it to the PR-images branch list of the `gloo` workload identity provider. Branch names are listed exactly; patterns come later. `release.yaml` accepts tags up to `v1.29`.

## PR images

To get images for a PR, comment `/build-pr-images` on it (the PR author or someone with write access), or run "PR images" from the Actions tab on the PR's base branch with the PR number, the head SHA and your login. [`pr-images.yaml`](../../.github/workflows/pr-images.yaml) checks that the PR is open, the branch is its base and the SHA is its head. It then waits for a second person to approve in environment `release` before it submits [Publish Artifacts](publish-artifacts.yaml) as the PR build service account. The requester and the PR author can't approve; the build fails if they do. After approval the build checks again that the PR is open and its head is still the requested commit, and fails otherwise. A direct dispatch leaves `comment_id` empty.

The build runs the PR's code with the quay.io robot logged in, so approvers check the PR's Makefile, `ci/` and Dockerfile changes first. The run title shows the PR head commit; approve only if it is the commit you reviewed. Approvers must not have pushed commits to the PR.

Each new push needs a new request. Images are tagged `<version>-<PR number>` and can lag the PR head. For a fork PR, someone with write access can comment on the author's behalf. PR builds that don't push (`run-tests.yaml`, including the glooctl build) still run from build-bot without approval.
