# Release Process

This document is for maintainers who need to cut a release for `solo-io/gloo`.

## Overview

A pushed tag is not the full release process for this repo.

- The GitHub Release object must exist before `glooctl` assets can be uploaded.
- `glooctl` assets are uploaded by Build Bot / Google Cloud Build, not by GitHub Actions.
- The release-triggered GitHub Actions are follow-up jobs. They do not attach `glooctl` binaries to the GitHub release page.

Relevant implementation details:

- [`cloudbuild.yaml`](cloudbuild.yaml)
- [`ci/cloudbuild/publish-artifacts.yaml`](ci/cloudbuild/publish-artifacts.yaml)
- [`Makefile`](Makefile)
- [`ci/upload_github_release_assets.go`](ci/upload_github_release_assets.go)
- [`.github/workflows/push-docs.yaml`](.github/workflows/push-docs.yaml)
- [`.github/workflows/push-solo-apis-branch.yaml`](.github/workflows/push-solo-apis-branch.yaml)

## Release Checklist

1. Start from the branch you intend to release.

   For example, use `main` for the current development line or an LTS branch such as `v1.21.x` for a patch release.

2. Create and push an annotated tag.

   ```bash
   git checkout <release-branch>
   git pull --ff-only origin <release-branch>
   git tag -a vX.Y.Z[-betaN|-rcN] -m "vX.Y.Z[-betaN|-rcN]"
   git push origin vX.Y.Z[-betaN|-rcN]
   ```

3. Create the GitHub Release from that tag.

   Use the GitHub Releases UI or `gh release create`.
   Create it promptly after pushing the tag, and publish it rather than leaving it as a draft.

   ```bash
   gh release create vX.Y.Z \
     --repo solo-io/gloo \
     --verify-tag \
     --title vX.Y.Z
   ```

   Add `--prerelease` for beta and release candidate builds.

4. Monitor the release automation.

   There are three important pieces:

   - Build Bot / Cloud Build should run the `publish-artifacts` pipeline. This is the job that publishes images, the Helm chart, and `glooctl` release assets.
   - `push-docs` runs on `release.created`.
   - `Push API Changes to solo-apis` runs on `release.published`.

5. Verify the GitHub release assets.

   Expected assets from `publish-glooctl`:

   - `glooctl-linux-amd64`
   - `glooctl-linux-amd64.sha256`
   - `glooctl-linux-arm64`
   - `glooctl-linux-arm64.sha256`
   - `glooctl-darwin-amd64`
   - `glooctl-darwin-amd64.sha256`
   - `glooctl-darwin-arm64`
   - `glooctl-darwin-arm64.sha256`
   - `glooctl-windows-amd64.exe`
   - `glooctl-windows-amd64.exe.sha256`

   Quick check:

   ```bash
   gh release view vX.Y.Z --repo solo-io/gloo --json assets,url
   ```

6. Verify follow-up automation.

   - Confirm `push-docs` completed or capture the failure for follow-up.
   - Confirm `Push API Changes to solo-apis` completed and review the resulting PR in `solo-apis` if one was opened.

## Checking Cloud Build

If the GitHub release exists but assets are still missing, inspect Build Bot / Cloud Build.

The root Cloud Build config submits a child build named `publish-artifacts`, and the child build's `release-chart` step runs `publish-helm-chart` and `publish-glooctl`.

If you have `gcloud` access, start with:

```bash
gcloud builds list \
  --project=solo-public \
  --limit=20 \
  --sort-by=~createTime \
  --format='table(id,status,createTime,substitutions.TAG_NAME,substitutions.REPO_NAME,logUrl)'
```

Then inspect the matching build:

```bash
gcloud builds describe BUILD_ID \
  --project=solo-public \
  --format='yaml(id,status,createTime,finishTime,substitutions,steps.id,logUrl)'
```

And view logs:

```bash
gcloud builds log BUILD_ID --project=solo-public
```

Look for:

- `substitutions.TAG_NAME: vX.Y.Z`
- the `publish-artifacts` child build
- the `release-chart` step
- `publish-glooctl`
- `upload_github_release_assets.go`

If your team uses a different GCP project for Build Bot, substitute that project name accordingly.

## Common Failure Modes

### The release page exists but has no assets

This usually means the GitHub Release object did not exist when `publish-glooctl` tried to upload assets, or the Cloud Build release job failed.

Fix:

1. Confirm the GitHub Release object exists for the tag.
2. Inspect the Build Bot / Cloud Build history for that tag.
3. Re-run or manually trigger the `publish-artifacts` release build if needed.

### `push-docs` fails

This workflow is separate from asset publishing. A `push-docs` failure does not explain missing GitHub release assets by itself.

Check the workflow logs and fix the docs-copy issue independently.

### `Push API Changes to solo-apis` succeeds but no change lands automatically

That workflow pushes to `solo-apis`. Review and approve the resulting `solo-apis` PR if one is created.
