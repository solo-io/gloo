---
title: Installing Gloo Gateway to Multiple Namespaces
weight: 20
description: Multi-tenant Gloo Gateway installations by installing to multiple namespaces
---

In the default deployment scenario, a single deployment of the Gloo Gateway control plane and Envoy proxy are installed for the entire cluster. However, in some cases, it may be desirable to deploy multiple instances of the Gloo Gateway control plane and proxies in a single cluster.

This is useful when multiple tenants or applications want control over their own instance of Gloo Gateway. Some deployment scenarios may involve a Gloo Gateway per-application architecture. Additionally, different Gloo Gateway instances living in their own namespace may be given different levels of RBAC permissions.

In this document, we will review how to deploy multiple instances of Gloo Gateway to their own namespaces within a single Kubernetes cluster. 

---

## Scope Gloo Gateway to specific namespaces

By default, Gloo Gateway watches all namespaces in a cluster for Kubernetes services and Gloo Gateway custom resources. However, if you have many namespaces, translation time might slow down due to the amount of resources that Gloo Gateway must attempt to process.

In this case, you might want to scope Gloo Gateway to watch only specific namespaces. Gloo Gateway can detect Kubernetes services and Gloo Gateway resources only in namespaces that you list, and cannot detect services and custom resources in any other namespaces. This setting can improve translation time by reducing the number of resources that Gloo Gateway attempts to process across all namespaces in your cluster.

Additionally, Gloo Gateway reads configuration for the Gateway custom resource only in the namespace that the gateway controller is deployed to by default. For Gateway configuration in other namespaces, such as to support multiple gateways, you can enable the `gateway.readGatewaysFromAllNamespaces` setting.

### Specify namespaces to watch for Kuberenetes services and Gloo Gateway CRs

To configure namespaces for Gloo Gateway to watch for Kubernetes services and Gloo Gateway CRs, you can use one of the following Helm settings:
* `settings.watchNamespaces`: Create a static list of namespaces for Gloo Gateway to watch. This setting is recommended if you have a set of configuration namespaces that is unlikely to change often.
* `settings.watchNamespaceSelectors.matchLabels` and `settings.watchNamespaceSelectors.matchExpressions`: Configure Gloo Gateway to watch namespaces based on a namespace selector, such as a label or an expression. This setting helps you dynamically determine the list of namespaces to watch, instead of using a static list of namespaces that you must update if you need to add or remove a namespace.

For example, to dynamically watch namespaces that have the `gloo-translate: enabled` label, the `env: prod` label, or the `env: dev` label, your Helm settings might look like the following.
```yaml
...
settings:
  watchNamespaceSelectors:
  - matchLabels:
      gloo-translate: enabled
  - matchExpressions:
    - key: env
      operator: In
      values:
        - prod
        - dev
```

### Specify namespaces to watch for Gateway configuration

To configure namespaces for Gloo Gateway to watch for the Gateway custom resource configuration, you can enable the `{{< protobuf name="gloo.solo.io.Settings" display="gateway.readGatewaysFromAllNamespaces">}}` setting. This setting allows Gloo Gateway to read Gateway configuration from any namespaces that are watched, such as the namespaces you specify in the `watchNamespaces` or `watchNamespaceSelectors` settings.

```yaml
...
gloo:
  gateway:
    readGatewaysFromAllNamespaces: true
```

### Observability

To help you monitor the namespaces that Gloo Gateway watches, you can use the following related logs and metrics.

**Logs**: These debug-level logs in the `gloo` deployment indicate when the set of watched namespaces has changed.
* `"received updated list of namespaces", zap.Any("namespaces", newSnapshot.Kubenamespaces)`
* `"list of namespaces to watch", zap.Any("oldNamespacesToWatch", oldNamespacesToWatch), zap.Any("newNamespacesToWatch", newNamespacesToWatch), zap.Any("namespacesChanged", namespacesChanged)`

**Metrics**: This metric indicates how many namespaces Gloo Gateway currently watches. You can use this metric in conjunction with the `gloo.solo.io/setups_run` metric, which indicates how often the control plane configuration reloads due to Settings configuration changes. These metrics can help you determine whether namespace cycling causes more control plane configuration reloads to run than necessary.
* `gloo.solo.io/namespaces_watched`

---

## Install namespace-scoped Gloo Gateway with Helm

By leveraging namespace scoping options, you can install Gloo Gateway to as many namespaces as needed without overlap.

{{% notice note %}}
`watchNamespaces` can be shared between Gloo Gateway instances, as long as any Virtual Services are not written to a shared namespace. When this happens, both Gloo Gateway instances will attempt to apply the same routing config, which can cause domain conflicts.
{{% /notice %}}

In this section we'll deploy Gloo Gateway twice, each instance to a different namespace, with two different Helm value files. 

For Gloo Gateway Enterprise users, you often use Gloo with the enterprise observability tools, Grafana and Prometheus. However, you cannot use the same observability instance for multiple Gloo instances. You can disable the observability tool for additional Gloo instances, or create separate observability tool instances by using name and RBAC overrides, as shown in the following examples.

For Gloo Gateway Open Source users, remove the [Grafana]({{% versioned_link_path fromRoot="/guides/observability/grafana/" %}}) and [Prometheus]({{% versioned_link_path fromRoot="/guides/observability/prometheus/" %}}) settings from the examples. Grafana and Prometheus are enterprise-only features.

Create a file named `gloo1-overrides.yaml` and paste the following inside:

```yaml
settings:
  create: true
  writeNamespace: gloo1
  watchNamespaces:
  - default
  - gloo1
gloo:
  gateway:
    readGatewaysFromAllNamespaces: true # Read Gateway config in all 'watchNamespaces'
grafana: # Remove the grafana settings for Gloo Gateway OSS
  rbac:
    namespaced: true
prometheus: # Remove the prometheus settings for Gloo Gateway OSS
  kube-state-metrics:
    fullnameOverride: glooe-prometheus-kube-state-metrics-1
  server:
    fullnameOverride: glooe-prometheus-server-1
```

Now, let's install Gloo Gateway. Review our [Kubernetes installation guide]({{% versioned_link_path fromRoot="/installation/gateway/kubernetes/" %}}) if you need a refresher.

First create the namespace for our first Gloo Gateway deployment:

```shell script
kubectl create ns gloo1
```

Then install Gloo Gateway using one of the following methods:

{{< tabs >}}
{{< tab name="glooctl" codelang="shell" >}}
glooctl install gateway -n gloo1 --values gloo1-overrides.yaml
{{< /tab >}}
{{< tab name="Helm 3" codelang="shell">}}
helm install gloo gloo/gloo --namespace gloo1 -f gloo1-overrides.yaml
{{< /tab >}}
{{< /tabs >}}

{{% notice warning %}}
Using Helm 2 is not supported in Gloo Gateway.
{{% /notice %}}

Check that gloo pods are running: 

```bash
kubectl get pod -n gloo1
```

```bash
NAME                             READY   STATUS    RESTARTS   AGE
discovery-798cdd5499-z7rrt       1/1     Running   0          37s
gateway-5fc999b847-jf4xp         1/1     Running   0          32s
gateway-proxy-67f4c7dfb6-hc5kg   1/1     Running   0          27s
gloo-dd5bcdc8f-bvtjh             1/1     Running   0          39s
```

And we should see that Gloo Gateway is only creating Upstreams from services in `default` and `gloo1`:

```bash
kubectl get us -n gloo1                                              
```

```bash
NAME                      AGE
default-kubernetes-443    1h
gloo1-gateway-proxy-443   1h
gloo1-gateway-proxy-80    1h
gloo1-gloo-9977           1h
```

Let's repeat the above process, substituting `gloo2` for `gloo1`:

Create a file named `gloo2-overrides.yaml` and paste the following inside:

```yaml
settings:
  create: true
  writeNamespace: gloo2
  watchNamespaces:
  - default
  - gloo2
gloo:
  gateway:
    readGatewaysFromAllNamespaces: true # Read Gateway config in all 'watchNamespaces'
grafana: # Remove the grafana settings for Gloo Gateway OSS
  rbac:
    namespaced: true
prometheus: # Remove the prometheus settings for Gloo Gateway OSS
  kube-state-metrics:
    fullnameOverride: glooe-prometheus-kube-state-metrics-2
  server:
    fullnameOverride: glooe-prometheus-server-2
```

Now, let's install Gloo Gateway for the second time. First create the second namespace:

```shell script
# create the namespace for our second gloo deployment
kubectl create ns gloo2
```

Then perform the second installation using one of the following methods:

{{< tabs >}}
{{< tab name="glooctl" codelang="shell" >}}
glooctl install gateway -n gloo2 --values gloo2-overrides.yaml
{{< /tab >}}
{{< tab name="Helm 3" codelang="shell">}}
helm install gloo gloo/gloo --namespace gloo2 -f gloo2-overrides.yaml
{{< /tab >}}
{{< /tabs >}}

Check that gloo pods are running: 

```bash
kubectl get pod -n gloo2
```

```bash
NAME                             READY   STATUS    RESTARTS   AGE
discovery-798cdd5499-kzmkc       1/1     Running   0          8s
gateway-5fc999b847-pn2tk         1/1     Running   0          8s
gateway-proxy-67f4c7dfb6-284wv   1/1     Running   0          8s
gloo-dd5bcdc8f-krp5p             1/1     Running   0          9s
```

And we should see that the second installation of Gloo Gateway is only creating Upstreams from services in `default` and `gloo2`:

```bash
kubectl get us -n gloo2
```

```bash
NAME                      AGE
default-kubernetes-443    53s
gloo2-gateway-proxy-443   53s
gloo2-gateway-proxy-80    53s
gloo2-gloo-9977           53s
```

And that's it! We can now create routes for Gloo Gateway #1 by creating our Virtual Services in the `gloo1` namespace, and routes for Gloo Gateway #2 by creating Virtual Services in the `gloo2` namespace. We can add `watchNamespaces` to our liking; the only catch is that a Virtual Service which lives in a shared namespace will be applied to both gateways (which can lead to undesired behavior if this was not the intended effect).

{{% notice warning %}}
When uninstalling a single instance of Gloo Gateway when multiple instances are installed, you should only delete the namespace into which that instance is installed. Running `glooctl uninstall` can cause cluster-wide resources to be deleted, which will break any remaining Gloo Gateway installation in your cluster
{{% /notice %}}