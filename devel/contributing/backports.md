# Backports

## What is a backport?
A backport is a change that is introduced on the main branch, and then applied to a previous version of Gloo Gateway.

## When is backporting appropriate?
For a backport to be appropriate it must fit the following criteria:
- The change must have a clear rationale for why it is needed on a previous version of Gloo Gateway.
- The change must be a low-risk change, this typically means it is a bug fix or a non-breaking change.
- The proposed change is targeted to a [supported, stable release branch](https://docs.solo.io/gloo-edge/latest/reference/support/).
- If the change is a feature request, you must have explicit approval from the product and engineering teams. This approval can also be solicited on the backport prs themselves
  - Generally a backport for a feature should have been requested by at least one of these teams to be considered in the first place

## How to identify a backport
On the issue that tracks the desired functionality, apply a `release/1.N` label to indicate the version of Gloo Gateway you wish the request to be supported on.

For example, if there is a `release/1.15` label, that means the issue is targeted to be introduced first on the stable main branch, and then backported to Gloo Gateway 1.15.x.

## How to create a backport
First, create a PR to introduce the change on the main branch. Doing so ensures that changes are tested and reviewed before being applied to a previous version of Gloo Gateway. Also, given that we use [protocol buffers](https://developers.google.com/protocol-buffers) for our API definitions, introducing API changes to our main branch first ensures we will not have API compatibility issues when backporting.

Once the change has been merged into the main branch, create a PR to backport the change to the desired release branch. The PR should be titled `[BRANCH NAME]: <PR title>` (ie `[1.14]: Original PR Title`). To create a backport branch we recommend the following:
- Use [cherry-pick](https://git-scm.com/docs/git-cherry-pick) to apply changes to a previous version of Gloo Gateway.
  - Resolve any conflicts that have arisen due to drift between LTS branches
  - If there is significant drift that causes the cherry-pick to be non-trivial, consider re-implementing the change from scratch rather than "backporting"
- Modify the changelog to be in the proper directory
  - Note: In the rare case that a new feature (indicated by `NEW_FEATURE` in the [changelog file](https://github.com/solo-io/go-utils/blob/main/changelogutils/README.md)) is being backported, changelog-bot
    will complain with an error similar to `Only one version greater than the latest release <release> valid, found...`. This is because our changelog-bot does not allow `NEW_FEATURE` in a backport/patch release. To get around this error, change the changelog type from `NEW_FEATURE` to `FIX` for backports.
- Validate that Proto fields have the same numbers as in main
- Title your PR to start with the major.minor version that you are backporting to (e.g. 1.13 for 1.13.x branch)