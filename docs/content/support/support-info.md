---
title: Add support information
description: Review the details to include in your support. 
weight: 930
---

Collect valuable information for Solo to review and troubleshoot your support request. 

## Environment

1. Get your Gloo Gateway version. 
    ```shell
    glooctl version -o yaml
    ```
2. Get the version of Kubernetes that you run in your cluster. 
    ```shell
    kubectl version -o yaml
    ```
3. List the infrastructure provider that hosts your environment, such as AWS or an on-premise virtual machine (VM). 

## Setup

Share the installation method that you used to install Gloo Gateway, such as Helm, `glooctl`, or Argo CD, and the configuration that you used during the installation. 

{{< tabs >}}
{{% tab name="Helm" %}}
Use the following command to extract the Helm manifests that were used for the Gloo Gateway installation and the state of each manifest. Replace `releaseName` and `namespaceName` with the values that you used for your installation.
```shell
helm get all <releaseName> --namespace <namespaceName> > gloo_edge_helm.yaml
```
  
{{% /tab %}}

{{% tab name="glooctl"%}}
Share the complete command that you used to install Gloo Gateway. The following command shows an example command that you might have used. 
```shell
glooctl install gateway enterprise --license-key <license> --values foo,bar
```
{{% /tab %}}

{{% tab name="Argo CD"%}}
Get the Argo CD applications that define the installation of Gloo Gateway by using one of the following methods. 
* `kubectl`: 
  ```shell
  kubectl get applications.argoproj.io/<applicationName> -n <namespaceName> -o yaml
  ```
* Argo CD CLI: 
  ```shell
  argocd app get <applicationName> -o yaml
  ```
{{% /tab %}}
{{< /tabs >}}

## Issue

1. Provide a detailed description of the issue. Make sure to include the following information: 
   - If reproducable, steps to reproduce the issue 
   - If applicable, sample of a client request including the payload 
   - High-level diagram of the interactions from the perspective of the client request
   - The request protocol that is handled by the application(s) in question, such as HTTP / TCP / gRPC 
   - If the issue is related to authentication or authorization, details of the auth configuration
2. Describe the impact of the issue. For example, the issue might block an update or a demo, or cause the loss of data or an entire system.
3. Export the relevant configuration files that are related to the issue.
  {{< tabs >}}
  {{% tab name="Gloo Gateway resources"%}}
  - Typically, the Gloo Gateway `Settings` object is useful to understand the configuration of Gloo Gateway.
  - For traffic management issues, include the following list of Gloo Gateway resources:
      - `Gateway`
      - `VirtualService`
      - `RouteTable`
      - `Upstream`

  - Use the following script to dump all Gloo Gateway custom resources into a file. Attach the `gloo-gateway-configuration.yaml` file to your support request. 
    ```shell
    for n in $(kubectl get crds | grep solo.io | awk '{print $1}'); do kubectl get $n --all-namespaces -o yaml >> gloo-gateway-configuration.yaml; echo "---" >> gloo-gateway-configuration.yaml; done
    ```
  {{% /tab %}}
  {{% tab name="Istio resources"%}}
  If you use Gloo Gateway as a gateway to an Istio service mesh, provide details about how Istio is configured. You can use the following command to create an Istio bug report that you can attach to your support request. 

  ```shell
  istioctl bug-report --istio-namespace <istioControlPlaneNamespace>
  ```
  {{% /tab %}}
  {{< /tabs >}}

## Product-specific details

Use the built-in Gloo Gateway tools to generate a debug report that includes all the logs and configurations for the Gloo Gateway control and data plane components. 

1. Capture the output of the `glooctl check` command.
    <br>Typically, the command output indicates any errors in the control plane components or associated resources, such as in the following example.

    ```
    Checking Deployments... 1 Errors!
    Checking Pods... 2 Errors!
    Checking Upstreams... OK
    Checking UpstreamGroups... OK
    Checking AuthConfigs... OK
    Checking RateLimitConfigs... OK
    Checking VirtualHostOptions... OK
    Checking RouteOptions... OK
    Checking Secrets... OK
    Checking VirtualServices... OK
    Checking Gateways... OK
    Checking Proxies... Skipping due to an error in checking deployments
    Skipping due to an error in checking deployments
    Error: 5 errors occurred:
    * Deployment gloo in namespace gloo-system is not available! Message: Deployment does not have minimum availability.
    * Pod gloo-8ddc4ff4c-g4mnf in namespace gloo-system is not ready! Message: containers with unready status: [gloo]
    * Not all containers in pod gloo-8ddc4ff4c-g4mnf in namespace gloo-system are ready! Message: containers with unready status: [gloo]
    * proxy check was skipped due to an error in checking deployments
    * xds metrics check was skipped due to an error in checking deployment
    ```

2. Create a Gloo Gateway debug report that collects the logs and configuration for various control plane components, such as `gloo`, `gloo-fed`, `redis`, or `observability` by using the `glooctl debug` command. The components vary depending on your Gloo Gateway setup and can be found in the `gloo-system` namespace. The command stores all the information collected in a local `debug` directory. Make sure to include these files in your support ticket. 
   ```sh
   glooctl debug -N gloo-system
   ```

3. If you experience issues in a specific namespace, make sure to also capture the logs and configurations for these resources. 
   ```sh
   glooctl debug -N <namespace>
   ```
