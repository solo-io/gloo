---
title: Prepare to upgrade
weight: 10
description: Prepare your environment, review version changes, and review FAQs before you upgrade Gloo Gateway.
---

Before you upgrade Gloo Gateway, complete the following preparatory steps:
* [Prepare your environment](#prepare), such as upgrading your current version to the latest patch and upgrading any dependencies to the required supported versions. 
* [Review important changes](#review-changes) made to Gloo Gateway in version {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}}, including CRD, Helm, CLI, and feature changes.
* [Review frequently-asked questions](#faqs) about the upgrade process.

## Prepare your environment {#prepare}

Review the following preparatory steps that might be required for your environment.

### Upgrade your current minor version to the latest patch {#current-patch}

Before you upgrade your minor version, first upgrade your current version to the latest patch. For example, if you currently run Gloo Gateway Enterprise version {{< readfile file="static/content/version_gee_n-1_oldpatch.md" markdown="true">}}, first upgrade your installation to version {{< readfile file="static/content/version_gee_n-1.md" markdown="true">}}. This ensures that your current environment is up-to-date with any bug fixes or security patches before you begin the minor version upgrade process.

1. Find the latest patch of your minor version by checking the [Open Source changelog]({{% versioned_link_path fromRoot="/reference/changelog/open_source/" %}}) or [Enterprise changelog]({{% versioned_link_path fromRoot="/reference/changelog/enterprise/" %}}).
2. Go to the documentation set for your current minor version. For example, if you currently run Gloo Gateway Enterprise version {{< readfile file="static/content/version_gee_n-1_oldpatch.md" markdown="true">}}, use the drop-down menu in the header of this page to select **v{{< readfile file="static/content/version_geoss_n-1_minor.md" markdown="true">}}.x**.
3. Follow the upgrade guide, using the latest patch for your minor version.

### If required, perform incremental minor version updates {#minor-increment}

If you plan to upgrade to a version that is more than one minor version greater than your current version, such as to version {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}} from {{< readfile file="static/content/version_geoss_n-2_minor.md" markdown="true">}} or older, you must upgrade incrementally. For example, you must first use the upgrade guide in the v{{< readfile file="static/content/version_geoss_n-1_minor.md" markdown="true">}}.x documentation set to upgrade from {{< readfile file="static/content/version_geoss_n-2_minor.md" markdown="true">}} to {{< readfile file="static/content/version_geoss_n-1_minor.md" markdown="true">}}, and then follow the upgrade guide in the v{{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}}.x documentation set to upgrade from {{< readfile file="static/content/version_geoss_n-1_minor.md" markdown="true">}} to {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}}.

### Upgrade dependencies {#dependencies}

Check that your underlying infrastructure platform, such as Kubernetes, and other dependencies run a version that is supported for {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}}.

1. Review the [supported versions]({{% versioned_link_path fromRoot="/reference/support/#supported-versions" %}}) for dependencies such as Kubernetes, Istio, Helm, and more.
2. Compare the supported versions against the versions you currently use.
3. If necessary, upgrade your dependencies, such as consulting your cluster infrastructure provider to upgrade the version of Kubernetes that your cluster runs.

### Consider settings to avoid downtime {#downtime}

You might deploy Gloo Gateway in Kubernetes environments that use the Kubernetes load balancer, or in non-Kubernetes environments. Depending on your setup, you can take additional steps to avoid downtime during the upgrade process.

* **Kubernetes**: Enable [Envoy readiness and liveness probes]({{< versioned_link_path fromRoot="/operations/production_deployment/#enable-health-checks" >}}) during the upgrade. When these probes are set, Kubernetes sends requests only to the healthy Envoy proxy during the upgrade process, which helps to prevent potential downtime. The probes are not enabled in default installations because they can lead to timeouts or other poor getting started experiences. 
* **Non-Kubernetes**: Configure [health checks]({{< versioned_link_path fromRoot="/guides/traffic_management/request_processing/health_checks" >}}) on Envoy. Then, configure your load balancer to leverage these health checks, so that requests stop going to Envoy when it begins draining connections.

## Review version {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}} changes {#review-changes}

Review the following changes made to Gloo Gateway in version {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}}. For some changes, you might be required to complete additional steps during the upgrade process.

### Breaking changes

**AuthPlugin removed**: The `AuthPlugin` auth config type is removed in 1.20. If you use this config type, you must remove it from your Helm values file before upgrading. If you need to configure your own auth service, check out [Custom Auth server]({{< versioned_link_path fromRoot="/guides/security/auth/custom_auth/" >}}).

**Discovery disabled by default**: To improve performance, discovery is now disabled by default in 1.20. Upstreams are no longer automatically created for discovered Kubernetes services. If you want to continue using discovery, set `discovery.enabled=true` (or `gloo.discovery.enabled` in Gloo Gateway Enterprise) in your Helm values file before upgrading. For more information, see the [Discovery guide]({{< versioned_link_path fromRoot="/installation/advanced_configuration/fds_mode/" >}}).

**Envoy version upgrade**

The Envoy dependency in Gloo Gateway 1.20 was upgraded from 1.33.x to 1.35.x. This change includes the following upstream breaking changes. For more information about these changes, see the [Envoy v1.34 changelog documentation](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.34/v1.34) and [Envoy v1.35 changelog documentation](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.35/v1.35).

From Envoy v1.34:
* **Extproc - Tracing changes**: When the `modeOverride` fields of the header and trailer modes have the value `DEFAULT` or are unset, no change is made to the processing mode that is set in the filter configuration. Additionally, a bug was fixed in which local replies were incorrectly sent to the extProc server for external processing. You can temoporarily revert this change by setting the runtime guard field `envoy.reloadable_features.skip_ext_proc_on_local_reply` to `false`.
* **TLS - Envoy FIPS**: The Envoy FIPS build is updated to use the same version of BoringSSL as the Envoy standard build to meet the requirements of the revised FedRAMP policy. Default values might have changed in the encyrption settings due to the large version bump.
* **HTTP connection manager - Generate request ID**: When the `x-request-id` header is empty, `generateRequestId` now generates a request ID on the request. Previously, it only generated an ID when the header was missing.
* **Formatters**: The `%CEL%` and `%METADATA%` formatters are now considered built-in, and can be used directly in substitution format strings if the related extensions are linked.

From Envoy v1.35:
* **tracing changes**: Added `max_cache_size` to the OpenTelemetry tracer config. This limits the number of spans that can be cached before flushing. The default is 1024 spans. Previously, flushing only happens at the set interval. User should tune this setting based on their environment's expected telemetry volume.

## New features

<!-- TODO confirm 1.20 k8s and istio testing support before uncommenting these
### Kubernetes 1.33 support 

Starting in version 1.20.0, Gloo Gateway can now run on Kubernetes 1.33. For more information about supported Kubernetes, Envoy, and Istio versions, see [Supported versions]({{% versioned_link_path fromRoot="/reference/support/" %}}).

### Istio 1.26 support

Starting in version 1.20.0, Gloo Gateway can now run with Istio 1.26. For more information about supported Kubernetes, Envoy, and Istio versions, see [Supported versions]({{% versioned_link_path fromRoot="/reference/support/" %}}).
-->

### Changelogs

Check the changelogs for the type of Gloo Gateway deployment that you have. Focus especially on any **Breaking Changes** that might require a different upgrade procedure. For Gloo Gateway Enterprise, you might also review the open source changelogs because most of the proto definitions are open source.
* [Open Source changelogs]({{% versioned_link_path fromRoot="/reference/changelog/open_source/" %}})
* [Enterprise changelogs]({{% versioned_link_path fromRoot="/reference/changelog/enterprise/" %}}): Keep in mind that Gloo Gateway Enterprise pulls in Gloo Gateway Open Source as a dependency. Although the major and minor version numbers are the same for open source and enterprise, their patch versions often differ. For example, open source might use version `x.y.a` but enterprise uses version `x.y.b`. If you are unfamiliar with these versioning concepts, see [Semantic versioning](https://semver.org/). Because of the differing patch versions, you might notice different output when checking your version with `glooctl version`. For example, your API server might run Gloo Gateway Enterprise version {{< readfile file="static/content/version_gee_latest.md" markdown="true">}}, which pulls in Gloo Gateway Open Source version {{< readfile file="static/content/version_geoss_latest.md" markdown="true">}} as a dependency.
  ```bash
  ~ > glooctl version
  Client: {"version":"{{< readfile file="static/content/version_geoss_latest.md" markdown="true">}}"}
  Server: {"type":"Gateway","enterprise":true,"kubernetes":...,{"Tag":"{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}","Name":"grpcserver-ee","Registry":"quay.io/solo-io"},...,{"Tag":"{{< readfile file="static/content/version_geoss_latest.md" markdown="true">}}","Name":"discovery","Registry":"quay.io/solo-io"},...}
  ```

{{% notice tip %}}
You can use the changelogs' built-in [comparison tool]({{< versioned_link_path fromRoot="/reference/changelog/open_source/#compareversions" >}}) to compare between your current version and the version that you want to upgrade to.
{{% /notice %}}

<!--

### Feature changes {#features}

Review the following summary of important new, deprecated, or removed features.

{{% notice note %}}
The following lists consist of the changes that were initially introduced with the {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}}.0 release. These changes might be backported to earlier versions of Gloo Gateway. Additionally, there might be other changes that are introduced in later {{< readfile file="static/content/version_geoss_latest_minor.md" markdown="true">}} patch releases. For patch release changes, check the [changelogs](#changelogs).
{{% /notice %}}

**New or improved features**:

**Deprecated features**:

-->

### Helm changes {#helm}

Review the following summary of important new, deprecated, or removed Helm fields. For full details, see the [changelogs](#changelogs).

**Updated Helm fields**:

* **discovery.enabled=false**: To improve performance, discovery is now disabled by default. This change means that Upstreams are no longer automatically created for discovered Kubernetes services. If you want to continue using discovery, set `discovery.enabled=true` (or `gloo.discovery.enabled` in Gloo Gateway Enterprise) in your Helm values file before upgrading. For more information, see the [Discovery guide]({{< versioned_link_path fromRoot="/installation/advanced_configuration/fds_mode/" >}}).

<!--
### CRD changes {#crd}

New CRDs are automatically applied to your cluster when performing a `helm install` operation, but are _not_ applied when performing an `helm upgrade` operation. This is a [deliberate design choice](https://helm.sh/docs/topics/charts/#limitations-on-crds) on the part of the Helm maintainers, given the risk associated with changing CRDs. Given this limitation, you must apply new CRDs to the cluster before upgrading. 

Review the following summary of important new, deprecated, or removed CRD updates. For full details, see the [changelogs](#changelogs).

As part of the {{< readfile file="static/content/version_geoss_latest.md" markdown="true">}} release, no CRD changes were introduced.

**New and updated CRDs**:


**Deprecated CRDs**:
N/A

**Removed CRDs**:
N/A

### CLI changes {#cli}

You must upgrade `glooctl` before you upgrade Gloo Gateway. Because `glooctl` can create resources in your cluster, such as with `glooctl add route`, you might have errors in Gloo Gateway if you create resources with an older version of `glooctl`.

Review the following summary of important new, deprecated, or removed CLI options. For full details, see the [changelogs](#changelogs).

**New CLI commands or options**:

-->
## Frequently-asked questions {#faqs}

Review the following frequently-asked questions about the upgrade process. If you still aren't sure about the version upgrade impact, or if your use case doesn't quite fit the standard upgrade path, feel free to post in the `#gloo` or `#gloo-enterprise` channels of our [public Slack](https://slack.solo.io/).

### How do I upgrade Gloo Gateway in testing or sandbox environments?

If downtime is not a concern for your use case, you can follow the [Quick upgrade guide]({{< versioned_link_path fromRoot="/operations/upgrading/upgrade_steps" >}}) to update your Gloo Gateway installation.

Note that for sandbox or exploratory environments, the easiest way to upgrade is to uninstall Gloo Gateway by running `glooctl uninstall --all`. Then, re-install Gloo Gateway at the desired version by the following one of the [installation guides]({{< versioned_link_path fromRoot="/installation" >}}).
 
### How do I upgrade Gloo Gateway in a production environment, where downtime is unacceptable?

The basic `helm upgrade` process is not suitable for environments in which downtime is unacceptable. Instead, you can follow the [Canary upgrade]({{% versioned_link_path fromRoot="/operations/upgrading/canary/" %}}) guide to deploy multiple version of Gloo Gateway to your cluster, and test the upgrade version before uninstalling the existing version.

Additionally, you might need to take steps to account for other factors such as Gloo Gateway version changes, probe configurations, and external infrastructure like the load balancer that Gloo Gateway uses. Consider setting up [liveness probes and healthchecks](#downtime) in your environment.

### What happens to my Gloo Gateway CRs during an upgrade? How do I handle breaking changes?

A typical upgrade of Gloo Gateway across minor versions should not cause disruptions to the existing Gloo Gateway state. In the case of a breaking change, Solo will communicate through the upgrade guides, changelogs, or other channels if you must make a specific adjustment to perform the upgrade. Note that you can always use the `glooctl debug yaml` command to download the current Gloo Gateway state to one large YAML manifest.

### Is the upgrade procedure different if I am not a cluster administrator?

If you are not an administrator of your cluster, you might be unable to create custom resource definitions (CRDs) and other cluster-scoped resources, such as cluster roles and cluster role bindings. If you encounter an error related to these resources, you can disable their creation by including the following setting in your Helm values:
```yaml
global:
  glooRbac:
    create: false
```

Otherwise, you can try performing an installation of Gloo Gateway that is scoped to a single namespace by including the following setting in your Helm values:
```yaml
global:
  glooRbac:
    namespaced: true
```

### Why do I get an error about re-creating CRDs when upgrading using `helm install` or `helm upgrade`?

Helm v2 does not manage CRDs well, and is not supported in Gloo Gateway. Upgrade to Helm v3, delete the CRDs, and try again.

### Why do I get an error about a `gateway-certgen` job?

The upgrade creates a Kubernetes Job named `gateway-certgen` to generate a certificate for the validation webhook. The job contains the `ttlSecondsAfterFinished` value so that the cluster cleans up the job automatically, but because this setting is still in Alpha, your cluster might ignore this value. In this case, you might have an issue while upgrading in which the upgrade attempts to change the `gateway-certgen` job, but the change fails because the job is immutable. To fix this issue, you can delete the job, which already completed, and re-apply the upgrade.
