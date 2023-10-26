# Using an unreleased Gloo OSS branch
**This setup step is only necessary if you are attempting to import Gloo OSS changes that have not yet been released**

Solo-Projects does not yet build all the resources that it depends upon. We have [an issue](https://github.com/solo-io/gloo/issues/5243) to track this technical debt. In the meantime, you must execute the following steps to run regression tests locally (or in CI):

## 1. Import Gloo OSS changes
```bash
go get github.com/solo-io/gloo@BRANCH-NAME
```

## 2. Configure Gloo OSS Helm Chart
As part of CI in Gloo, we publish the helm chart. These assets are published to a slightly different location for Gloo PRs as opposed to Gloo releases. You can see in the `build-bot` logs the path to the published helm chart:

```text
Step #14 - "release-chart": Uploading helm chart to gs://solo-public-tagged-helm with name gloo-1.12.0-beta5-6341.tgz
```

We need to update the [gloo subchart dependency](https://github.com/solo-io/solo-projects/blob/393276665446d69fcfde7d5e65cc9c678ab3a100/install/helm/gloo-ee/requirements-template.yaml#L2) to point to this path. An example is below, please note that the `repository` needs to change as well as the `version`:

```text
dependencies:
  - name: gloo
    repository: https://storage.googleapis.com/solo-public-tagged-helm
    version: 1.12.0-beta5-6341
```
## 3. Propagate and import API changes

Navigate the to [Push API Changes to solo-apis](https://github.com/solo-io/gloo/actions/workflows/push-solo-apis-branch.yaml) action on OSS

Run the workflow, specifying the appropriate feature branch and target LTS branch

Once the action is complete, there will be a solo-apis branch than can be imported:
```bash
go get github.com/solo-io/solo-apis@sync-apis/gloo-<LTS-BRANCH>/gloo-<FEATURE-BRANCH>
```

Checkout the [solo-apis branches page](https://github.com/solo-io/solo-apis/branches) to find the branch for your changes