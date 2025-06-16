---
title: "Installing Gloo Gateway Enterprise"
menuTitle: Gloo Gateway Enterprise
description: How to install Gloo Gateway to run in Gateway Mode on Kubernetes (Default).
weight: 60
---

Review how to install Gloo Gateway Enterprise.
## Before you begin

1. Make sure that you prepared your Kubernetes cluster according to the [instructions for platform configuration]({{% versioned_link_path fromRoot="/installation/platform_configuration/cluster_setup/" %}}).
   
   {{% notice note %}}
   Pay attention to provider-specific information in the setup guide. For example, [OpenShift]({{< versioned_link_path fromRoot="/installation/platform_configuration/cluster_setup/#openshift" >}}) requires stricter multi-tenant support, so the setup guide includes an example Helm chart `values.yaml` file that you must supply while installing Gloo Gateway Enterprise.
   {{% /notice %}}

2. Get your Gloo Gateway Enterprise license key. {{< readfile file="static/content/license-key" markdown="true">}}

3. Check whether `glooctl`, the Gloo Gateway command line tool (CLI), is installed.
   ```bash
   glooctl version
   ```
   * If `glooctl` is not installed, [install it](#install-glooctl).
   * If `glooctl` is installed, [update it to the latest version](#update-glooctl).

{{< readfile file="installation/glooctl_setup.md" markdown="true" >}}

## Installing Gloo Gateway Enterprise on Kubernetes {#install-steps}

Review the following steps to install Gloo Gateway Enterprise with `glooctl` or with Helm.

### Installing on Kubernetes with `glooctl`

Once your Kubernetes cluster is up and running, run the following command to deploy the Gloo Gateway to the `gloo-system` namespace:

```bash
glooctl install gateway enterprise --license-key YOUR_LICENSE_KEY
```

{{% notice note %}}
For OpenShift clusters, make sure to include the `--values values.yaml` option to point to the [Helm chart custom values file]({{< versioned_link_path fromRoot="/installation/platform_configuration/cluster_setup/#openshift" >}}) that you created.
{{% /notice %}}

<details>
<summary>Special Instructions to Install Gloo Gateway Enterprise on Kind</summary>
If you followed the cluster setup instructions for Kind <a href="{{< versioned_link_path fromRoot="/installation/platform_configuration/cluster_setup/#kind" >}}">here</a>, then you should have exposed custom ports 31500 (for http) and 32500 (https) from your cluster's Docker container to its host machine. The purpose of this is to make it easier to access your service endpoints from your host workstation.  Use the following custom installation for Gloo Gateway to publish those same ports from the proxy as well.

```bash
glooctl install gateway enterprise --license-key YOUR_LICENSE_KEY --values - << EOF
gloo:
  gatewayProxies:
    gatewayProxy:
      service:
        type: NodePort
        httpPort: 31500
        httpsPort: 32500
        httpNodePort: 31500
        httpsNodePort: 32500
EOF
```

```
Creating namespace gloo-system... Done.
Starting Gloo Gateway Enterprise installation...

Gloo Gateway Enterprise was successfully installed!
```

Note also that the url to invoke services published via Gloo Gateway will be slightly different with Kind-hosted clusters.  Much of the Gloo Gateway documentation instructs you to use `$(glooctl proxy url)` as the header for your service url.  This will not work with kind.  For example, instead of using curl commands like this:

```bash
curl $(glooctl proxy url)/all-pets
```

You will instead route your request to the custom port that you configured above for your docker container to publish. For example:

```bash
curl http://localhost:31500/all-pets
```
</details>

Once you've installed Gloo Gateway, please be sure [to verify your installation](#verify-your-installation).


{{% notice note %}}
You can run the command with the flag `--dry-run` to output 
the Kubernetes manifests (as `yaml`) that `glooctl` will 
apply to the cluster instead of installing them.
{{% /notice %}}

### Installing on Kubernetes with Helm

This is the recommended method for installing Gloo Gateway Enterprise to your production environment as it offers rich customization to
the Gloo Gateway control plane and the proxies Gloo Gateway manages.

As a first step, you have to add the Gloo Gateway repository to the list of known chart repositories:

```shell
helm repo add glooe https://storage.googleapis.com/gloo-ee-helm
```

Finally, install Gloo Gateway using the following command:

```shell
helm install gloo glooe/gloo-ee --namespace gloo-system \
  --create-namespace --set-string license_key=YOUR_LICENSE_KEY
```

{{% notice note %}}
For OpenShift clusters, make sure to include the `--values values.yaml` option to point to the [Helm chart custom values file]({{< versioned_link_path fromRoot="/installation/platform_configuration/cluster_setup/#openshift" >}}) that you created.
{{% /notice %}}

{{% notice warning %}}
Using Helm 2 is not supported in Gloo Gateway.
{{% /notice %}}

Once you've installed Gloo Gateway, please be sure [to verify your installation](#verify-your-installation).

### Argo CD installation

[Argo Continuous Delivery (Argo CD)](https://argo-cd.readthedocs.io/en/stable/) is a declarative, Kubernetes-native continuous deployment tool that can read and pull code from Git repositories and deploy it to your cluster. Because of that, you can integrate Argo CD into your GitOps pipeline to automate the deployment and synchronization of your apps. 

Planning to use Argo CD version 7.8 or later? Review the [known issue]({{< versioned_link_path fromRoot="/installation/gateway/argo/#settings-issue" >}}) first.

**Set up Argo CD**

1. Install Argo CD in your cluster. 
   ```sh
   kubectl create namespace argocd
   until kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.12.3/manifests/install.yaml > /dev/null 2>&1; do sleep 2; done
   # wait for deployment to complete
   kubectl -n argocd rollout status deploy/argocd-applicationset-controller
   kubectl -n argocd rollout status deploy/argocd-dex-server
   kubectl -n argocd rollout status deploy/argocd-notifications-controller
   kubectl -n argocd rollout status deploy/argocd-redis
   kubectl -n argocd rollout status deploy/argocd-repo-server
   kubectl -n argocd rollout status deploy/argocd-server
   ```

2. Update the default Argo CD password for the admin user to solo.io.
   ```sh
   # bcrypt(password)=$2a$10$79yaoOg9dL5MO8pn8hGqtO4xQDejSEVNWAGQR268JHLdrCw6UCYmy
   # password: solo.io
   kubectl -n argocd patch secret argocd-secret \
     -p '{"stringData": {
       "admin.password": "$2a$10$79yaoOg9dL5MO8pn8hGqtO4xQDejSEVNWAGQR268JHLdrCw6UCYmy",
       "admin.passwordMtime": "'$(date +%FT%T%Z)'"
     }}'
   ```
   
3. Port-forward the Argo CD server on port 9999.
   ```sh
   kubectl port-forward svc/argocd-server -n argocd 9999:443
   ```

4. Open the [Argo CD UI](https://localhost:9999/).

5. Log in with the `admin` username and `solo.io` password.

**Install Gloo Gateway**
   
1. Use the following YAML file to create an Argo CD application and deploy the Gloo Gateway Enterprise Helm chart. Make sure to enter your license key in the `license_key` field. You can add custom settings to the `spec.source.helm.values` section. 
         
   {{% notice note %}}
   Argo CD does not have a concept of installs or upgrades. All updates are executed by using syncs. Because of that, the value of `gateway.certGenJob.runOnUpdate` (if set) is ignored. Instead, the job runs on every sync. 
   {{% /notice %}}
   
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: argoproj.io/v1alpha1
   kind: Application
   metadata:
     name: gloo-gateway-ee-helm
     namespace: argocd
   spec:
     destination:
       namespace: gloo-system
       server: https://kubernetes.default.svc
     project: default
     source:
       chart: gloo-ee
       helm:
         skipCrds: false
         values: |
           gloo:
             kubeGateway:
               enabled: false
           observability:
             enabled: false
           prometheus:
             enabled: false
           license_key: <enterprise-license-key>
       repoURL: https://storage.googleapis.com/gloo-ee-helm
       targetRevision: {{< readfile file="static/content/version_gee_latest.md" markdown="true">}}
     syncPolicy:
       automated:
         # Prune resources during auto-syncing (default is false)
         prune: true 
         # Sync the app in part when resources are changed only in the target Kubernetes cluster
         # but not in the git source (default is false).
         selfHeal: true 
       syncOptions:
       - CreateNamespace=true
   EOF
   ```

2. Verify that the `gloo` control plane components are up and running.
   ```sh
   kubectl get pods -n gloo-system 
   ```
   
   Example output: 
   ```
   NAME                                  READY   STATUS      RESTARTS   AGE
   extauth-7449fc4b67-wcgn5              1/1     Running     0          3m6s
   gloo-b9ff69d5d-c85wx                  1/1     Running     0          3m6s
   gloo-resource-migration-c86hp         0/1     Completed   0          3m33s
   gloo-resource-rollout-9vnjj           0/1     Completed   0          3m6s
   gloo-resource-rollout-check-m8rhk     0/1     Completed   0          2m50s
   gloo-resource-rollout-cleanup-z6ktl   0/1     Completed   0          3m22s
   rate-limit-84656bddb7-dpmql           1/1     Running     0          3m6s
   redis-54757c7964-dnh67                1/1     Running     0          3m6s
   ```

3. Open the Argo CD UI and verify that you see the Argo CD application with a `Healthy` and `Synced` status.

**Optional: Cleanup** </br>

If you no longer need this quick-start Gloo Gateway environment, you can uninstall your setup by following these steps: 

{{< tabs >}}
{{% tab name="Argo CD UI" %}}
1. Port-forward the Argo CD server on port 9999.
   ```sh
   kubectl port-forward svc/argocd-server -n argocd 9999:443
   ```

2. Open the [Argo CD UI](https://localhost:9999/applications).

3. Log in with the `admin` username and `solo.io` password.
4. Find the application that you want to delete and click **x**. 
5. Select **Foreground** and click **Ok**. 
6. Verify that the pods were removed from the `gloo-system` namespace. 
   ```sh
   kubectl get pods -n gloo-system
   ```
   
   Example output: 
   ```  
   No resources found in gloo-system namespace.
   ```

{{% /tab %}}
{{% tab name="Argo CD CLI" %}}
1. Port-forward the Argo CD server on port 9999.
   ```sh
   kubectl port-forward svc/argocd-server -n argocd 9999:443
   ```
   
2. Log in to the Argo CD UI. 
   ```sh
   argocd login localhost:9999 --username admin --password solo.io --insecure
   ```
   
3. Delete the application.
   ```sh
   argocd app delete gloo-gateway-ee-helm --cascade --server localhost:9999 --insecure
   ```
   
   Example output: 
   ```
   Are you sure you want to delete 'gloo-gateway-ee-helm' and all its resources? [y/n] y
   application 'gloo-gateway-ee-helm' deleted   
   ```

4. Verify that the pods were removed from the `gloo-system` namespace. 
   ```sh
   kubectl get pods -n gloo-system
   ```
   
   Example output: 
   ```  
   No resources found in gloo-system namespace.
   ```
{{% /tab %}}
{{< /tabs >}}

### Airgap installation

You can install Gloo Gateway Enterprise in an air-gapped environment, such as an on-premises datacenter, clusters that run on an intranet or private network only, or other disconnected environments.

Before you begin, make sure that you have the following setup:
* A connected device that can pull the required images from the internet.
* An air-gapped or disconnected device that you want to install Gloo Gateway Enterprise in.
* A private image registry such as Sonatype Nexus Repository or JFrog Artifactory that both the connected and disconnected devices can connect to.

To install Gloo Gateway Enterprise in an air-gapped environment:

1. Set the Gloo Gateway Enterprise version that you want to use as an environment variable, such as the latest version in the following example.
   ```shell
   export GLOO_EE_VERSION={{< readfile file="static/content/version_gee_latest.md" markdown="true">}}
   ```
2. On the connected device, download the Gloo Gateway Enterprise images.
   ```shell
   helm template glooe/gloo-ee --version $GLOO_EE_VERSION --set-string license_key=$GLOO_LICENSE_KEY | yq e '. | .. | select(has("image"))' - | grep image: | sed 's/image: //'
   ```
   
   The example output includes the list of images.
   ```
   quay.io/solo-io/gloo-fed-apiserver:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:07cfb09d3f57e7ff13f0b77a27f85541c323fbbb64611a6d103d4869cd9e53b2
   quay.io/solo-io/gloo-federation-console:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:a114aac5fb45a82f9496e080abc6491998bd582a2d108ba7e283e8ae74bb9f31
   quay.io/solo-io/gloo-fed-apiserver-envoy:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:5d81b0a5ea7366536870a26348722d8ca54a814fef5b6bdb5c925b6a5e866980
   quay.io/solo-io/gloo-fed:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:1dae4f26058cc7e189b3d67832732d0126765f3e9336e437634791ea43081580
   quay.io/solo-io/gloo-ee:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:5a189ff231d4a6e7d3d17b4c59a8173f9f93e3f0e8e899ba570e9d2fa3006751
   quay.io/solo-io/discovery-ee:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:12a9408f36955aa9b632259556a0055fc3fa55a40d4fccf64c5ba1c9b3593fbf
   quay.io/solo-io/gloo-ee-envoy-wrapper:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:92e8389cfa6589d2b5d76ebc25071566117fb98cc58916013f342e1a4bc366d7
   "grafana/grafana:8.2.1"
   "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.6.0"
   "jimmidyson/configmap-reload:v0.5.0"
   "quay.io/prometheus/prometheus:v2.39.1"
   docker.io/busybox:1.28
   docker.io/redis:7.0.11
   quay.io/solo-io/rate-limit-ee:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:18be159bbabce83d085359761682204e79b3dd1e4e7e250cb504658c4828dbc8
   quay.io/solo-io/extauth-ee:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:2ad07cd3cd734aa48f332c8c8d453d501e6ad8f0ab70bb7e9c0bd518e084d41d
   quay.io/solo-io/observability-ee:{{< readfile file="static/content/version_gee_latest.md" markdown="true">}}@sha256:e3b3e61e0d9b4801f4b4cf0b7c2dd80b43181017ef3d24fad778c60e6c771b82
   quay.io/solo-io/kubectl:1.15.11@sha256:a32099a598211312c12f5dc5162d3dbfb3046e59b8099ad2dead2068d4dd453d
   quay.io/solo-io/kubectl:1.15.11@sha256:a32099a598211312c12f5dc5162d3dbfb3046e59b8099ad2dead2068d4dd453d
   quay.io/solo-io/kubectl:1.15.11@sha256:a32099a598211312c12f5dc5162d3dbfb3046e59b8099ad2dead2068d4dd453d
   quay.io/solo-io/kubectl:1.15.11@sha256:a32099a598211312c12f5dc5162d3dbfb3046e59b8099ad2dead2068d4dd453d
   quay.io/solo-io/kubectl:1.15.11@sha256:a32099a598211312c12f5dc5162d3dbfb3046e59b8099ad2dead2068d4dd453d
   quay.io/solo-io/certgen:1.15.11@sha256:4a30d309158ddf575878508c38f38745c483853da61cb69f669b3fee8da35541
   quay.io/solo-io/kubectl:1.15.11@sha256:a32099a598211312c12f5dc5162d3dbfb3046e59b8099ad2dead2068d4dd453d
   ```

3. Push the images from the connected device to a private registry that the disconnected device can pull from. For instructions and any credentials you must set up to complete this step, consult your registry provider, such as [Nexus Repository Manager](https://help.sonatype.com/repomanager3/formats/docker-registry/pushing-images) or [JFrog Artifactory](https://www.jfrog.com/confluence/display/JFROG/Getting+Started+with+Artifactory+as+a+Docker+Registry).
4. Optional: You might want to set up your private registry so that you can also pull the Helm charts. For instructions, consult your registry provider, such as [Nexus Repository Manager](https://help.sonatype.com/repomanager3/formats/helm-repositories) or [JFrog Artifactory](https://www.jfrog.com/confluence/display/JFROG/Kubernetes+Helm+Chart+Repositories).
5. When you [install Gloo Gateway Enterprise with a custom Helm chart values file](#customizing-your-installation-with-helm), make sure to use the specific images that you downloaded and stored in your private registry in the previous steps.

## Customizing your installation with Helm

You can customize the Gloo Gateway installation by providing your own Helm chart values file.

For example, you can create a file named `value-overrides.yaml` with the following content.

```yaml
global:
  glooRbac:
    # do not create kubernetes rbac resources
    create: false
settings:
  # configure gloo to write generated custom resources to a custom namespace
  writeNamespace: my-custom-namespace
  watchNamespaces:
  - default
  - my-custom-namespace
gloo:
  gateway:
    # For multiple gateways: read Gateway config in all 'watchNamespaces',
    # not just the namespace that the gateway controller is deployed to
    readGatewaysFromAllNamespaces: true 
```

Then, refer to the file during installation to override default values in the Gloo Gateway Helm chart.

```shell
helm install gloo glooe/gloo-ee --namespace gloo-system \
  -f value-overrides.yaml --create-namespace --set-string license_key=YOUR_LICENSE_KEY
```

{{% notice warning %}}
Using Helm 2 is not supported in Gloo Gateway.
{{% /notice %}}

### List of Gloo Gateway Helm chart values

The following table describes the most important enterprise-only values that you can override in your custom values file.

For more information, see the following resources:
* [Gloo Gateway Open Source overrides]({{< versioned_link_path fromRoot="/reference/helm_chart_values/" >}}) (also available in Enterprise). 
* [Advanced customization guide]({{% versioned_link_path fromRoot="/installation/gateway/kubernetes/helm_advanced/" %}}).
* [Enterprise Helm chart reference document]({{% versioned_link_path fromRoot="/reference/helm_chart_values/enterprise_helm_chart_values/" %}}).

{{% notice note %}}
Gloo Gateway Open Source Helm values in Enterprise must be prefixed with `gloo`, unless they are the Gloo Gateway settings, such as `settings.<rest of helm value>`.
{{% /notice %}}

| Option | Type | Description |
| --- | --- | --- |
| global.extensions.caching.enabled                         | bool     | Deploy the caching server in the `gloo-system` namespace. Default is `false`. |
| global.extensions.extAuth.enabled                         | bool     | Deploy the ext-auth server in the `gloo-system` namespace. Default is `true`. |
| global.extensions.extAuth.envoySidecar                    | bool     | Deploy ext-auth in the `gateway-proxy` pod as a sidecar to Envoy. Communicates over a Unix domain socket instead of TCP. Default is `false`. |
| gloo.gatewayProxies.NAME.tcpKeepaliveTimeSeconds | unit32 | The amount of time in seconds for connections to be idle before sending keep-alive probes. Defaults to 60s. You might use this to prevent sync issues due to network connectivity glitches. For more information, see [the Knowledge Base help article](https://support.solo.io/hc/en-us/articles/12066701909524).|
| gloo.gloo.disableLeaderElection | bool | Leave this field set to the default value of `false` when you have multiple replicas of the `gloo` deployment. This way, Gloo Gateway elects a leader from the replicas, with the other replicas ready to become leader if needed in case the elected leader pod fails or restarts. If you want to run only one replica of `gloo`, you can set this value to `true`.|
| grafana.defaultInstallationEnabled                        | bool     | Deploy Grafana in the `gloo-system` namespace. Default is `true`. |
| observability.enabled                                     | bool     | Deploy Grafana in the `gloo-system` namespace. Default is `true`. |
| observability.customGrafana.enabled                       | bool     | Use your own Grafana instance instead of the default Gloo Gateway Grafana instance. Default is `false`. |
| observability.customGrafana.username                      | string   | Authenticate to your custom Grafana instance using this username for basic auth. |
| observability.customGrafana.password                      | string   | Authenticate to your custom Grafana instance using this password basic auth. |
| observability.customGrafana.apiKey                        | string   | Authenticate to your custom Grafana instance using this API key. |
| observability.customGrafana.url                           | string   | The URL for your custom Grafana instance. |
| prometheus.enabled                                        | bool     | Deploy Prometheus in the `gloo-system` namespace. Default is `true`. |
| rateLimit.enabled                                         | bool     | Deploy the rate-limiting server in the `gloo-system` namespace. Default is `true`. |
---

## Enterprise UI

Gloo Gateway Enterprise comes with a built-in UI that you can use to view information about your cluster and the Gloo Gateway instance that you installed. You can enable the Gloo Gateway Enterprise UI by using the `gloo-fed.glooFedApiserver.enable=true` setting during the installation. 

{{< tabs >}}
{{% tab name="glooctl install" %}}
```shell script
echo "gloo-fed:
  glooFedApiserver:
    enable: true" > values.yaml
glooctl install gateway enterprise --values values.yaml --license-key=<LICENSE_KEY>
```
{{% /tab %}}
{{% tab name="helm install" %}}
```shell script
helm install gloo glooe/gloo-ee --namespace gloo-system --set gloo-fed.glooFedApiserver.enable=true --set license_key=<LICENSE_KEY>
```
{{% /tab %}}
{{< /tabs >}}

Note that when you also enable Gloo Federation by using the `gloo-fed.enabled=true` setting, the UI does not show any federation data until you [register one or more clusters]({{< versioned_link_path fromRoot="/guides/gloo_federation/cluster_registration/" >}}). 


## Verify your Installation

Check that the Gloo Gateway pods and services have been created. Depending on your install option, you may see some differences
from the following example. And if you choose to install Gloo Gateway into a different namespace than the default `gloo-system`,
then you will need to query your chosen namespace instead.

```shell
kubectl --namespace gloo-system get all
```

```noop
NAME                                                          READY   STATUS    RESTARTS   AGE
pod/discovery-6dbb5fd8bc-gk2th                                1/1     Running   0          2m5s
pod/extauth-68bb4745fc-2rs7b                                  1/1     Running   0          2m5s
pod/gateway-proxy-7c49898fdf-blxps                            1/1     Running   0          2m5s
pod/gloo-7748b94989-dj85p                                     1/1     Running   0          2m5s
pod/gloo-fed-76c85d689b-q62k4                                 1/1     Running   0          2m5s
pod/gloo-fed-console-dd5f877bd-jgg8n                          3/3     Running   0          2m5s
pod/glooe-grafana-6f95948945-pvbcg                            1/1     Running   0          2m4s
pod/glooe-prometheus-kube-state-metrics-v2-6c79cc9554-hlhns   1/1     Running   0          2m5s
pod/glooe-prometheus-server-757dc7d8f7-x489q                  2/2     Running   0          2m5s
pod/observability-78cb7bddf7-kcrbm                            1/1     Running   0          2m5s
pod/rate-limit-5ddd4b69d-84d6b                                1/1     Running   0          2m5s
pod/redis-888f4d9b5-p76wk                                     1/1     Running   0          2m4s

NAME                                             TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                                                AGE
service/extauth                                  ClusterIP      10.xxx.xx.xx    <none>          8083/TCP                                               2m6s
service/gateway-proxy                            LoadBalancer   10.xxx.xx.xx    34.xx.xxx.xxx   80:30437/TCP,443:31651/TCP                             2m6s
service/gloo                                     ClusterIP      10.xxx.xx.xx    <none>          9977/TCP,9976/TCP,9988/TCP,9966/TCP,9979/TCP,443/TCP   2m7s
service/gloo-fed-console                         ClusterIP      10.xxx.xx.xx    <none>          10101/TCP,8090/TCP,8081/TCP                            2m6s
service/glooe-grafana                            ClusterIP      10.xxx.xx.xxx   <none>          80/TCP                                                 2m6s
service/glooe-prometheus-kube-state-metrics-v2   ClusterIP      10.xxx.xx.xxx   <none>          8080/TCP                                               2m6s
service/glooe-prometheus-server                  ClusterIP      10.xxx.xx.xx    <none>          80/TCP                                                 2m7s
service/rate-limit                               ClusterIP      10.xxx.xx.xxx   <none>          18081/TCP                                              2m7s
service/redis                                    ClusterIP      10.xxx.xx.xx    <none>          6379/TCP                                               2m6s

NAME                                                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/discovery                                1/1     1            1           2m7s
deployment.apps/extauth                                  1/1     1            1           2m7s
deployment.apps/gateway-proxy                            1/1     1            1           2m7s
deployment.apps/gloo                                     1/1     1            1           2m7s
deployment.apps/gloo-fed                                 1/1     1            1           2m7s
deployment.apps/gloo-fed-console                         1/1     1            1           2m7s
deployment.apps/glooe-grafana                            1/1     1            1           2m7s
deployment.apps/glooe-prometheus-kube-state-metrics-v2   1/1     1            1           2m7s
deployment.apps/glooe-prometheus-server                  1/1     1            1           2m7s
deployment.apps/observability                            1/1     1            1           2m7s
deployment.apps/rate-limit                               1/1     1            1           2m7s
deployment.apps/redis                                    1/1     1            1           2m7s

NAME                                                                DESIRED   CURRENT   READY   AGE
replicaset.apps/discovery-6dbb5fd8bc                                1         1         1       2m6s
replicaset.apps/extauth-68bb4745fc                                  1         1         1       2m7s
replicaset.apps/gateway-proxy-7c49898fdf                            1         1         1       2m6s
replicaset.apps/gloo-7748b94989                                     1         1         1       2m7s
replicaset.apps/gloo-fed-76c85d689b                                 1         1         1       2m7s
replicaset.apps/gloo-fed-console-dd5f877bd                          1         1         1       2m6s
replicaset.apps/glooe-grafana-6f95948945                            1         1         1       2m6s
replicaset.apps/glooe-prometheus-kube-state-metrics-v2-6c79cc9554   1         1         1       2m6s
replicaset.apps/glooe-prometheus-server-757dc7d8f7                  1         1         1       2m6s
replicaset.apps/observability-78cb7bddf7                            1         1         1       2m7s
replicaset.apps/rate-limit-5ddd4b69d                                1         1         1       2m7s
replicaset.apps/redis-888f4d9b5                                     1         1         1       2m6s
```

#### Looking for opened ports?
You will NOT have any open ports listening on a default install. For Envoy to open the ports and actually listen, you need to have a Route defined in one of the VirtualServices that will be associated with that particular Gateway/Listener. Please see the [Hello World tutorial to get started]({{% versioned_link_path fromRoot="/guides/traffic_management/hello_world/" %}}). 

{{% notice note %}}
NOT opening the listener ports when there are no listeners (routes) is by design with the intention of not over-exposing your cluster by accident (for security). If you feel this behavior is not justified, please let us know.
{{% /notice %}}

## Uninstall {#uninstall}

To uninstall Gloo Gateway, you can use the `glooctl` CLI. If you installed Gloo Gateway to a different namespace, include the `-n` option.

```shell
glooctl uninstall -n my-namespace
```

{{% notice warning %}}
Make sure that your cluster has no other instances of Gloo Gateway running, such as by running `kubectl get pods --all-namespaces`. If you remove the CRDs while Gloo Gateway is still installed, you will experience errors.
{{% /notice %}}

```shell
glooctl uninstall --all
```

## Next Steps

After you install Gloo Gateway, check out the [User Guides]({{< versioned_link_path fromRoot="/guides/" >}}).

{{< readfile file="static/content/upgrade-note.md" markdown="true">}}
