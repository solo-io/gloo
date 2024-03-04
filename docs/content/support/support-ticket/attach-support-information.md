---
title: Attach support information
description: Steps to gather and attach the necessary information to the support ticket.
weight: 930
---

## Details to include in your support request {#ticket-details}

### Environment

1. Gloo Edge version
    ```shell
    glooctl version -o yaml
    ```
2. Version of Kubernetes
    ```shell
    kubectl version -o yaml
    ```
3. Infrastructure provider, such as AWS or on-premise VM

### Setup

1. Installation method and configuration used for installing Gloo Edge, such as Helm or `glooctl`. This could also be GitOps tooling such as Argo CD.
   <br>Few examples of these tools are given below showing how to retrieve the configuration.

      {{< tabs >}}
      {{% tab name="Helm"%}}
  When Helm is used, the command below can extract the complete Helm state including the manifests.<br><br>
  ```shell
  helm get all <releaseName> --namespace <namespaceName> > gloo_edge_helm.yaml
  ```
  `releaseName` and `namespaceName` will need to be set to match your specific installation.
      {{% /tab %}}

      {{% tab name="Glooctl"%}}
  If `glooctl install` is used, then please share the complete command that was executed. An example of this could be,<br><br>
  ```shell
  glooctl install gateway enterprise --license-key <license> --values foo,bar
  ```
      {{% /tab %}}

      {{% tab name="Argo CD"%}}
  Argo CD applications that define the installation of Gloo Edge would be helpful using either of the following methods,<br><br>
  - With `kubectl`,
    ```shell
    kubectl get applications.argoproj.io/<applicationName> -n <namespaceName> -o yaml
    ```
  - With Argo CD CLI,
    ```shell
    argocd app get <applicationName> -o yaml
    ```
      {{% /tab %}}
      {{< /tabs >}}

### Issue
1. Description of the issue in detail. Few suggestions to include,
   - If the issue is reproducible, steps to reproduce the issue.
   - As part of the above steps, if a client request is being made please include a sample of this request including the payload where applicable.
   - Include a high level diagram of the interaction from the perspective of the client request being made.
   - What is the request protocol handled by the application(s) in question ? for  e.g. is it HTTP / TCP / gRPC ?
   - If the issue is related to authentication or authorization then details of this configuration. 
2. Impact of the issue, such as blocking an update, blocking a demo, data loss or the system is down.
3. Attach an export of relevant configuration files related to the issue.
  {{< tabs >}}
  {{% tab name="Gloo Edge Resources"%}}
  - Typically, Gloo Edge `Settings` object is useful to understand the configuration of Gloo Edge.
  - For traffic management issues the following list of Gloo Edge resources is useful.
      - `Gateway`
      - `VirtualService`
      - `RouteTable`
      - `Upstream`

    <br>Use the following script to dump all the Gloo Edge related resources,
    ```shell
    for n in $(kubectl get crds | grep solo.io | awk '{print $1}'); do kubectl get $n --all-namespaces -o yaml >> gloo-edge-configuration.yaml; echo "---" >> gloo-edge-configuration.yaml; done
    ```
    <br>You can attach the `gloo-edge-configuration.yaml` file.
  {{% /tab %}}
  {{% tab name="Istio Resources"%}}
  If you are integrating or using Gloo Edge as the gateway to Istio Service Mesh and the reported issue is related to the integration it is useful to understand how Istio is configured.

  <br><br>Please attach the output from the command below along with the Gloo Edge resources.
  ```shell
  istioctl bug-report --istio-namespace <istioControlPlaneNamespace>
  ```
  {{% /tab %}}
  {{< /tabs >}}

### Product-specific details

##### Control plane

1. Capture the output of `glooctl check` command.
    <br>Typically, the command output indicates any errors in the control plane components or associated resources.
    <br>An example output is shown below.
    ```
    Checking deployments... 1 Errors!
    Checking pods... 2 Errors!
    Checking upstreams... OK
    Checking upstream groups... OK
    Checking auth configs... OK
    Checking rate limit configs... OK
    Checking VirtualHostOptions... OK
    Checking RouteOptions... OK
    Checking secrets... OK
    Checking virtual services... OK
    Checking gateways... OK
    Checking proxies... Skipping due to an error in checking deployments
    Skipping due to an error in checking deployments
    Error: 5 errors occurred:
    * Deployment gloo in namespace gloo-system is not available! Message: Deployment does not have minimum availability.
    * Pod gloo-8ddc4ff4c-g4mnf in namespace gloo-system is not ready! Message: containers with unready status: [gloo]
    * Not all containers in pod gloo-8ddc4ff4c-g4mnf in namespace gloo-system are ready! Message: containers with unready status: [gloo]
    * proxy check was skipped due to an error in checking deployments
    * xds metrics check was skipped due to an error in checking deployment
    ```
2. Collect the logs from various control plane components, such as `gloo` by using the `debug` log level (if possible). 
    <br>To enable the `debug` log level, see [Debugging control plane]({{< versioned_link_path fromRoot="/operations/debugging_gloo#debug-control-plane" >}}).
    <br><br>Follow the steps below for `gloo` controller pod.
    1. Set the log level to `debug`.
        ```shell
        kubectl port-forward deploy/gloo -n <controlplaneNamespace> 9091:9091 > /dev/null 2>&1 &
        PID=$!
        curl -X PUT -H "Content-Type: application/json" -d '{"level": "debug"}' http://localhost:9091/logging
        ```
    2. Capture the logs when reproducing the issue.
        ```shell
        kubectl logs -f deploy/gloo -n <controlplaneNamespace> > gloo.log
        ```
    3. When log capture is done reset the log level to `info`.
        ```shell
        curl -X PUT -H "Content-Type: application/json" -d '{"level": "info"}' http://localhost:9091/logging
        kill -9 $PID
        ```
    These steps apply to all the control plane components.

##### Data plane

1. Capture the currently served xDS configuration.
   ```shell
   glooctl proxy served-config -n <controlplaneNamespace> > served-config.yaml
   ```
2. Get the configuration that is served in the `gateway-proxy` Envoy pod(s). 
   <br>For more information, see [Dumping Envoy configuration]({{< versioned_link_path fromRoot="/operations/debugging_gloo#dump-envoy-configuration" >}}) section for more details.
   ```shell
   kubectl port-forward deploy/gateway-proxy -n <proxyNamespace> 19000:19000 > /dev/null 2>&1 &
   PID=$!
   curl -s localhost:19000/config_dump\?include_eds > gateway-config.json
   kill -9 $PID
   ```
3. Get the access log(s) for failed request from the `gateway-proxy` pod(s). If Access logging is not enabled, refer to [this guide]({{< versioned_link_path fromRoot="/guides/security/access_logging" >}}) to enable it.
4. If possible, collect the logs from the `gateway-proxy` Envoy pod(s) in `debug` log level for the failed request.
   {{% notice note %}}
   Setting the log level `debug` can get very noisy. Preferred will be to set the log level on a set of log components as described in the [docs]({{< versioned_link_path fromRoot="/operations/debugging_gloo#view-envoy-logs" >}}).
   {{% /notice %}}
   For more information, see [Viewing Envoy logs]({{< versioned_link_path fromRoot="/operations/debugging_gloo#view-envoy-logs" >}}).
5. Gather the stats from the proxy pod(s).
   ```shell
   glooctl proxy stats > proxy-stats.log
   ```