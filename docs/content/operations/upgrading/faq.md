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


### 🔥 Breaking changes

Review the breaking changes in this release. 

##### XSLT transformation removed

The XSLT transformation feature that was previously deprecated is now removed. If you use XSLT transformations, plan to use an external processing server to process this type of transformation. For more information, see [External processing]({{% versioned_link_path fromRoot="/guides/traffic_management/extproc/" %}}).  

##### Caching filter removed

The Envoy-based caching filter, which was deprecated in Gloo Gateway 1.20, is removed in 1.22. 

##### Envoy version upgrade

The Envoy dependency in Gloo Gateway 1.22 was upgraded from 1.36.x to 1.38.x. The following breaking changes are included in these releases. For detailed information, see the Envoy changelog for [1.37](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.37/v1.37) and [1.38](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.38/v1.38).

**Envoy v1.37**:

* **Container-aware CPU detection**: Envoy now sizes worker threads by using the minimum of hardware threads, CPU affinity, and cgroup CPU limits. In cgroup-limited containers, this means fewer worker threads than before, which reduces the maximum number of connections. Previously, Envoy ignored cgroup CPU limits and sized threads based only on hardware threads and CPU affinity. This change applies only when the Envoy `--concurrency` startup flag is not explicitly set. To restore the previous behavior without setting `--concurrency`, set the `ENVOY_CGROUP_CPU_DETECTION` environment variable to `false`.

**Envoy v1.38**:

* **RSA key usage enforcement**: Envoy 1.38 sets [`enforce_rsa_key_usage`](https://www.envoyproxy.io/docs/envoy/v1.38.0/api-v3/extensions/transport_sockets/tls/v3/tls.proto.html#envoy-v3-api-field-extensions-transport-sockets-tls-v3-upstreamtlscontext-enforce-rsa-key-usage) to `true` by default for upstream TLS connections. If the `keyUsage` extension is present in the upstream certificate and is incompatible with the TLS usage, the TLS handshake fails. In a future version of Envoy, this option will be removed and enforcing behavior will always apply. This setting is specific to upstream TLS connections (not downstream client connections). The `keyUsage` extension tells consumers what the certificate's public key is allowed to be used for. If the extension is present but does not match the TLS role, the upstream handshake fails. Note that kgateway does not expose this setting in `BackendConfigPolicy`, so it cannot be set back to `false`. Common RSA key usage values compatible with TLS are:
  * `digitalSignature`
  * `keyEncipherment`
  * `keyCertSign` (CA certs only)
  * `cRLSign` (CA certs that sign revocation lists only)

  Verify your upstream certificates include compatible `keyUsage` values before upgrading.

* **Circuit breaker metrics**: Added a new `upstream_rq_active_overflow` counter that is incremented when a request is rejected because the `max_requests` circuit breaker is exhausted. Previously, this condition incorrectly incremented the `upstream_rq_pending_overflow` metric, making it impossible to distinguish between pending queue saturation and active request saturation. After the upgrade, only the `upstream_rq_active_overflow` is incremented for this case, so you might see a drop in `upstream_rq_pending_overflow` counts. If you have existing dashboards or alerts that rely on the `upstream_rq_pending_overflow` metric to detect `max_requests` circuit breaker trips, set the Envoy runtime flag `envoy.reloadable_features.skip_pending_overflow_count_on_active_rq` to `false` to increment both counters while you migrate your monitoring to the `upstream_rq_active_overflow` metric.
* **Memory management**: Replaced the custom timer-based tcmalloc memory release with tcmalloc's native `ProcessBackgroundActions` and `SetBackgroundReleaseRate` APIs. This provides more comprehensive background memory management, including per-CPU cache reclamation, cache shuffling, and size class resizing, in addition to memory release. The `tcmalloc.released_by_timer` stat is removed.
* **RBAC header matching**: Fixed the RBAC header matcher to validate each header value individually instead of concatenating multiple header values into a single string. This prevents potential policy bypasses when requests contain multiple values for the same header. The new behavior is enabled by default and controlled by the runtime guard `envoy.reloadable_features.rbac_match_headers_individually`.

##### Upstream admission validation enforced

Previously, an Upstream resource with a corrupt or invalid `protoDescriptorBin` in `grpcJsonTranscoder` was accepted by the admission webhook. Envoy later rejected the entire `RouteConfiguration` and entered a continuous retry loop that caused silent routing failures. Setting `fullEnvoyValidation: true` had no effect in this mode.

Upstream admission validation is now enforced. Invalid `protoDescriptorBin` values are now rejected at admission time.

**Action required**: Before upgrading, verify that any Upstream resources that set the `grpcJsonTranscoder` field have a valid `protoDescriptorBin`. Upstreams with corrupt descriptors are rejected on update or re-apply after the upgrade.


## ⚒️ Installation changes 

In addition to comparing differences across versions in the changelog, review the following installation changes. 

### Default Redis image upgraded to v8

The default Redis image that is bundled with Gloo Gateway was upgraded from `7.2.11-alpine` to `8.6.4-alpine` to address high and critical CVEs present in Redis 7.x. Review the [Redis 8 release notes](https://raw.githubusercontent.com/redis/redis/8.0/00-RELEASENOTES) for any breaking changes that are relevant to your usage. If you need to stay on a different version, you can override the image tag by setting `gloo.redis.deployment.image.tag` in your Helm values.

```yaml
gloo:
  redis:
    deployment:
      image:
        tag: <version>
```

## 🌟 New features

Review the following new features are introduced in this version.

### Forward headers to client on ext-auth success

The `PassThroughHttp.Response` configuration in the AuthConfig resource now supports an `allowedClientHeadersOnSuccess` field. When set, headers that are returned by the HTTP passthrough auth server are forwarded to the downstream client on a successful auth check. Headers that are already present on the client response are overwritten. Previously, you could only forward headers to the client on denial (`allowedClientHeadersOnDenied`) or to the upstream on success (`allowedUpstreamHeaders`). For more information, see the [HTTP passthrough auth guide]({{% versioned_link_path fromRoot="/guides/security/auth/extauth/passthrough_auth/http/" %}}).

### Configure connection pool settings for HTTP passthrough auth

The `PassThroughHttp` configuration in the AuthConfig resource now supports a `connectionPool` block and a `responseHeaderTimeout` field to tune the HTTP client that the ext-auth service uses when calling the passthrough auth server. Use `connectionPool.maxConns` to cap the number of concurrent connections per host, `connectionPool.idleTimeout` to control how long idle connections are kept open, and `responseHeaderTimeout` to set the maximum time to wait for the auth server to return response headers before the request is considered timed out. For more information, see the [HTTP passthrough auth guide]({{% versioned_link_path fromRoot="/guides/security/auth/extauth/passthrough_auth/http/" %}}).

### IPv4-only listener bind addresses

By default, Gloo Gateway binds listeners to the IPv6 wildcard address `::`, which listens for traffic on all interfaces. In IPv4-only Kubernetes environments, this causes gateways to fail to bind and proxies to receive an invalid xDS fallback configuration. You can now set `gatewayProxies.<NAME>.gatewaySettings.ipv4Only: true` in your Helm values to switch the Gateway spec's bind address to `0.0.0.0` instead. For more information, see [IPv4-only environments]({{% versioned_link_path fromRoot="/guides/traffic_management/listener_configuration/#ipv4-only" %}}).

### Set the gateway proxy workload type

You can now configure the type of workload that you want to use for your gateway proxies by using the `gloo.gatewayProxies.<NAME>.kind.workloadType` Helm value. Set this field to `"Deployment"` or `"DaemonSet"` to control whether the gateway proxy is deployed as a Deployment or a DaemonSet. For more information, see the [Helm values reference]({{% versioned_link_path fromRoot="/reference/helm_chart_values/enterprise_helm_chart_values/" %}}).

### Tune GOMEMLIMIT settings

You can now set `goMemLimitPercent` on Gloo Gateway deployments to compute `GOMEMLIMIT` as a percentage of `resources.limits.memory` at Helm render time. Previously, `GOMEMLIMIT` was set to 100% of the cgroup limit, leaving no headroom for the Go garbage collector and causing OOM kills during HPA prescale events. The recommended value is `80`–`90`. Setting the field to `0` or leaving it unset preserves the existing behavior. For all available Helm values, see the [Enterprise Helm values reference]({{% versioned_link_path fromRoot="/reference/helm_chart_values/enterprise_helm_chart_values/" %}}).

| Deployment | Helm field | Notes |
|---|---|---|
| `gloo` | `gloo.deployment.goMemLimitPercent` | |
| `extAuth` | `global.extensions.extAuth.deployment.goMemLimitPercent` | |
| `rateLimit` | `global.extensions.rateLimit.deployment.goMemLimitPercent` | |
| `discovery` | `discovery.deployment.goMemLimitPercent` | Not deployed by default. Requires opt-in. |
| `caching` | `global.extensions.caching.deployment.goMemLimitPercent` | Not deployed by default. Requires opt-in. |
| `observability` | `observability.deployment.goMemLimitPercent` | Not deployed by default. Requires opt-in. |

### Control the Host header for shadow traffic

The `RouteShadowing` configuration in the RouteOption resource now supports two new fields for controlling the `Host`/`:authority` header of mirrored requests.

- `disableShadowHostSuffixAppend`: By default, Envoy appends `-shadow` to the `Host`/`:authority` header of mirrored requests. Set this field to `true` to send the original `Host` header unchanged. This is useful when the shadow destination has strict host-based routing rules that reject the modified header. For more information, see [Disable the `-shadow` host suffix]({{% versioned_link_path fromRoot="/guides/traffic_management/request_processing/shadowing/#disable-shadow-suffix" %}}).
- `hostRewriteLiteral`: Replaces the `Host`/`:authority` header of mirrored requests with the specified value. Include a port if the shadow destination needs one, as the port from the original request is not carried over. For more information, see [Rewrite the Host header for shadow traffic]({{% versioned_link_path fromRoot="/guides/traffic_management/request_processing/shadowing/#host-rewrite-literal" %}}).


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