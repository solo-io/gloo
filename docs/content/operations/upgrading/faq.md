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

**Envoy version upgrade**

The Envoy dependency in Gloo Gateway 1.19 was upgraded from 1.31.x to 1.33.x. This change includes the following upstream breaking changes. For more information about these changes, see the [Envoy changelog documentation](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.33/v1.33).

* **Trust internal addresses**: By default Envoy is configured to trust internal addresses. However, in an upcoming release of Envoy, this behavior is changed and internal addresses must be added to the `internal_address_config` for Envoy to trust them. For example, if you have tooling, such as probes on your private network, make sure to include these IP addresses or CIDR ranges in the `internal_address_config` field. To try out or enable this upcoming behavior change, you can set the `envoy.reloadable_features.explicit_internal_address_config` runtime guard to `true`. For more information, see the related [pull request](https://github.com/envoyproxy/envoy/pull/36221/files) in Envoy.
* **Access log handlers**: Access log handlers that are added by filters are now evaluated before access log handlers that are configured in the `access_log` configuration. To disable this behavior, you can set the `envoy.reloadable_features.filter_access_loggers_first` runtime guard flag to `false`.
* **Cluster name change in Kuberetes Gateway API**: When using the Kubernetes Gateway API alongside Gloo Edge APIs and you route to Kubernetes Services or Upstreams, the Envoy cluster name format is changed to extract more details about the service. The new format uses underscores to list service details, such as `upstreamName_upstreamNs_svcNs_svcName_svcPort`. If you enabled the Kubernetes Gateway API integration with `kubeGateway.enabled=true`, both the Gloo Edge and Kubernetes Gateway API proxies use the same format for these cluster names.
* **Extproc tracing changes**: In previous releases, tracing spans that were generated by the extProc filter were sampled by default. Now, these traces are not automatically sampled. Instead, the tracing decision is inherited from the parent span.
* **Opencensus tracing extension removed**: Support for the Opencensus tracing extension is removed in Envoy 1.33. 
* **Opentracing extension removed**: Support for the Opentracing extension is removed in Envoy 1.32. For more information, see the related [Envoy issue](https://github.com/envoyproxy/envoy/issues/27401). 


## New features

### Set authority header for gRPC OpenTelemetry collectors

When referencing a gRPC OpenTelemetry collector in your Gateway, Gloo Gateway automatically generates an Envoy configuration that sets the cluster name as the `:authority` pseudo-header. If your collector expects a different `:authority` header, you can specify that by setting the `spec.httpGateway.options.httpConnectionManagerSettings.tracing.openTelemetryConfig.grpcService.authority` value on your Gateway as shown in the following example. 

{{< highlight yaml "hl_lines=23-24" >}}
apiVersion: gateway.solo.io/v1
kind: Gateway
metadata:
  labels:
    app: gloo
    app.kubernetes.io/name: gateway-proxy-tracing-authority
  name: gateway-proxy-tracing-authority
spec:
  bindAddress: '::'
  bindPort: 18082
  proxyNames:
    - gateway-proxy
  httpGateway:
    virtualServiceSelector:
      gateway-type: tracing
    options:
      httpConnectionManagerSettings:
        tracing:
          openTelemetryConfig:
            collectorUpstreamRef:
              name: opentelemetry-collector
              namespace: default
            grpcService:
              authority: my-authority
{{< /highlight >}}
        
### Log filter state in gRPC access logs

You can enable logging of the filter state when performing gRPC access logging. The filter state logger calls the `FilterState::Object::serializeAsProto` to serialize the filter state object. 

You can enable the filter in your Helm values file or a Gateway resource directly. 

{{< tabs >}}
{{% tab name="Helm" %}}
The following example adds the modsecurity object from a WAF policy to the filter state object in the access log. 
```yaml
gloo:
  accessLogger:
    enabled: true
    image:
      registry: quay.io/solo-io
      tag: 1.0.0-ci1
  gatewayProxies:
    gatewayProxy:
      gatewaySettings:
        customHttpGateway:
          options:
            waf:
              auditLogging:
                action: ALWAYS
                location: FILTER_STATE
              customInterventionMessage: 'ModSecurity intervention! Custom message details here..'
              ruleSets:
              - ruleStr: |
                  # Turn rule engine on
                  SecRuleEngine On
                  SecRule REQUEST_HEADERS:User-Agent "scammer" "deny,status:403,id:107,phase:1,msg:'blocked scammer'"
        accessLoggingService:
          accessLog:
          - grpcService:
              logName: example
              staticClusterName: access_log_cluster
              filterStateObjectsToLog:
              - io.solo.modsecurity.audit_log
```
{{% /tab %}}
{{% tab name="Gateway resource" %}}
The following example adds the modsecurity object from a WAF policy to the filter state object in the access log. 
```yaml  
apiVersion: gateway.solo.io/v1
kind: Gateway
metadata:
  name: gateway-proxy
  namespace: gloo-system
spec:
  bindAddress: '::'
  bindPort: 8080
  proxyNames: 
    - gateway-proxy
  httpGateway: 
    options:
      waf:
        customInterventionMessage: 'ModSecurity intervention! Custom message details here..'
        ruleSets:
        - ruleStr: |
            # Turn rule engine on
            SecRuleEngine On
            SecRule REQUEST_HEADERS:User-Agent "scammer" "deny,status:403,id:107,phase:1,msg:'blocked scammer'" 
  useProxyProto: false
  options:
    accessLoggingService:
      accessLog:
      - grpcService:
          logName: example
          staticClusterName: access_log_cluster
          filterStateObjectsToLog:
          - io.solo.modsecurity.audit_log
```

{{% /tab %}}
{{< /tabs >}}

Example access log output: 
```
"filter_state_objects":{"io.solo.modsecurity.audit_log":{"type_url":"type.googleapis.com/google.protobuf.StringValue","value":"Ct0DLS0tQS0tClsxMC9NYXIvMjAyNToxNDoyNTowMSArMDAwMF0gMTc0MTYxNjcwMTQ1LjkwMTMzMCAxMjcuMC4wLjEgMzg5NTQgIDAKLS0tQi0tCkdFVCAvYWxsLXBldHMgSFRUUC8xLjEKOm1ldGhvZDogR0VUCmhvc3Q6IGxvY2FsaG9zdDo4MDgwCjpwYXRoOiAvYWxsLXBldHMKOmF1dGhvcml0eTogbG9jYWxob3N0OjgwODAKeC1mb3J3YXJkZWQtcHJvdG86IGh0dHAKOnNjaGVtZTogaHR0cAp1c2VyLWFnZW50OiBjdXJsLzguMTEuMAphY2NlcHQ6ICovKgp4LXJlcXVlc3QtaWQ6IDUwZTBjZTVmLTk0NjktNGQ3NS1hMzUxLWU4MmI0ODFlYTFmYgoKLS0tRi0tCkhUVFAvMS4xIDIwMAo6c3RhdHVzOiAyMDAKY29udGVudC10eXBlOiBhcHBsaWNhdGlvbi94bWwKZGF0ZTogTW9uLCAxMCBNYXIgMjAyNSAxNDoyNTowMSBHTVQKY29udGVudC1sZW5ndGg6IDg2CngtZW52b3ktdXBzdHJlYW0tc2VydmljZS10aW1lOiAyCgotLS1ILS0KCi0tLVotLQoK"}}}
```

For more information, see the [Access logging API]({{% versioned_link_path fromRoot="/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/options/als/als.proto.sk/" %}}). 


### Match conditions on validation webhook

You can now specify match conditions on the Gloo Gateway or Kubernetes validating admission webhook level to filter the resources that you want to include or exclude from validation. Match conditions are written in [CEL (Common Expression Language)](https://cel.dev). 

Examples: 

To exclude secrets or other resources with the `foo` label on the Kubernetes validating admission webhook, add the following values to your Helm values file. 
```yaml
gateway:
  validation:
    kubeCoreFailurePolicy: Fail # For "strict" validation mode, fail the validation if webhook server is not available
    kubeCoreMatchConditions:
    - name: 'not-a-secret-or-secret-with-foo-label-key'
       expression: 'request.resource.resource != "secrets" || ("labels" in oldObject.metadata && "foo" in oldObject.metadata.labels)'
```

To exclude all Upstream resources on the Gloo Gateway validating admission webhook, add the following values to your Helm values file.  
```yaml
gateway:
  validation:
    failurePolicy: Fail # For "strict" validation mode, fail the validation if webhook server is not available
    matchConditions:
      - name: skip-upstreams
        expression: '!(request.resource.group == "gloo.solo.io" && request.resource.resource == "upstreams")' # Match non-upstream resources.
    webhook:
      skipDeleteValidationResources: []
```

For more information, see the Helm reference for [OSS]({{< versioned_link_path fromRoot="/reference/helm_chart_values/open_source_helm_chart_values/" >}}) and [Enterprise]({{< versioned_link_path fromRoot="/reference/helm_chart_values/enterprise_helm_chart_values/" >}}).

### Kubernetes 1.32 support 

Starting in version 1.19.0, Gloo Gateway can now run on Kubernetes 1.32. For more information about supported Kubernetes, Envoy, and Istio versions, see [Supported versions]({{% versioned_link_path fromRoot="/reference/support/" %}}).

### Istio 1.25 support

Starting in version 1.19.0, Gloo Gateway can now run with Istio 1.25. For more information about supported Kubernetes, Envoy, and Istio versions, see [Supported versions]({{% versioned_link_path fromRoot="/reference/support/" %}}).

<!-- ggv2-related changes:
ggv2 - Disable Istio Envoy proxy from running by default and only rely on proxyless Istio agent mtls integration. Note: Although this is a change to the default behavior of the istio integration, this should not have any impact on most users as the sidecar proxy was unused in the data path. (https://github.com/solo-io/solo-projects/issues/5711)

ggv2 - glooctl get proxy will not work if you have persisted Proxy CRs in etcD and you are querying and older server version (1.16 and below). In general, we recommend that you keep your client and server versions in sync. You can verify the client/server versions you are currently running by calling glooctl version. (https://github.com/solo-io/gloo/pull/9226)
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


### Helm changes {#helm}

Review the following summary of important new, deprecated, or removed Helm fields. For full details, see the [changelogs](#changelogs).

**New Helm fields**:

**Updated Helm fields**:

**Deprecated Helm fields**:

<!--GGv2 changes:

ggv2 - Introduced new fields to kubeGateway top-level field which configure the deployed Gateway proxies generated from a Gateway. Also introduced a new default GatewayParameters to be rendered when kubeGateway.enabled=true. This contains defaults for Istio/SDS, as well as things like envoy image, deployment replicas, and extra labels in the pod template. (https://github.com/solo-io/solo-projects/issues/6107)

ggv2 - Add k8s Gateway Istio integration values to the Gloo Gateway Helm chart under kubeGateway.gatewayParameters.glooGateway. (https://github.com/solo-io/solo-projects/issues/5743)

ggv2 - Rename the kube gateway envoy container image helm value from kubeGateway.gatewayParameters.glooGateway.image to kubeGateway.gatewayParameters.glooGateway.envoyContainer.image. (https://github.com/solo-io/solo-projects/issues/6107)

ggv2 - Introduce gateway.validation.webhook.enablePolicyApi which controls whether or not RouteOptions and VirtualHostOptions CRs are subject to validation. By default, this value is true. The validation of these Policy APIs only runs if the Kubernetes Gateway integration is enabled (kubeGateway.enabled). (https://github.com/solo-io/solo-projects/issues/6352)


### CRD changes {#crd}

New CRDs are automatically applied to your cluster when performing a `helm install` operation, but are _not_ applied when performing an `helm upgrade` operation. This is a [deliberate design choice](https://helm.sh/docs/topics/charts/#limitations-on-crds) on the part of the Helm maintainers, given the risk associated with changing CRDs. Given this limitation, you must apply new CRDs to the cluster before upgrading. 

Review the following summary of important new, deprecated, or removed CRD updates. For full details, see the [changelogs](#changelogs).

As part of the {{< readfile file="static/content/version_geoss_latest.md" markdown="true">}} release, no CRD changes were introduced.

**New and updated CRDs**:


**Deprecated CRDs**:
N/A

**Removed CRDs**:
N/A
-->

### CLI changes {#cli}

You must upgrade `glooctl` before you upgrade Gloo Gateway. Because `glooctl` can create resources in your cluster, such as with `glooctl add route`, you might have errors in Gloo Gateway if you create resources with an older version of `glooctl`.


Review the following summary of important new, deprecated, or removed CLI options. For full details, see the [changelogs](#changelogs).

**New CLI commands or options**:

* [`glooctl debug`]({{< versioned_link_path fromRoot="/reference/cli/glooctl_debug/" >}}) and [`glooctl debug yaml`]({{< versioned_link_path fromRoot="/reference/cli/glooctl_debug_yaml/" >}}): Collect Kubernetes, Gloo Gateway controller, and Envoy information from your environment, such as logs, YAML manifests, metrics, and snapshots. This information can be used to debug issues in your environment or to provide this information to the Solo.io support team.
* [`glooctl gateway api convert`]({{< versioned_link_path fromRoot="/reference/cli/glooctl_gateway-api_convert/" >}}): Use this command to convert Gloo Edge APIs to Kubernetes Gateway API YAML files so that you can preview and run the migration to Gloo Gateway with the Kubernetes Gateway API. 


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
