---
title: Set up the Gloo UI
weight: 10
description: Install the Gloo UI and start gaining insights into the configuration and health of your Gloo Gateway setup and the workloads in your cluster. 
---
Install the Gloo UI to get an at-a-glance view of the configuration, health, and compliance status of your Gloo Gateway setup and the workloads in your cluster. 

To learn more about the features of the Gloo UI, see [About the Gloo UI]({{< versioned_link_path fromRoot="/guides/observability/ui/explore#about-the-gloo-ui" >}}).

## Before you begin

Follow the steps to [install Gloo Gateway with the Enterprise Edition]({{< versioned_link_path fromRoot="/installation/enterprise/" >}}). 

## Set up the Gloo UI

Use these instructions to install the Gloo UI in the same cluster as Gloo Gateway. The Gloo UI analyzes your Gloo Gateway setup and provides metrics and insights to you. 

1. Set the name of your cluster and your Gloo Gateway license key as an environment variable.
   ```sh
   export CLUSTER_NAME=<cluster-name>
   export GLOO_GATEWAY_LICENSE_KEY=<license-key>
   ```

2. Add the Helm repo for the Gloo UI. 
   ```sh
   helm repo add gloo-platform https://storage.googleapis.com/gloo-platform/helm-charts
   helm repo update  
   ```
   
3. Install the custom resources for the Gloo UI. 
   ```sh
   helm upgrade -i gloo-platform-crds gloo-platform/gloo-platform-crds \
    --namespace=gloo-system \
    --version={{< readfile file="/static/content/version-platform.md" markdown="true">}} \
    --set installEnterpriseCrds=false
   ```

4. Install the Gloo UI and configure it for Gloo Gateway.
   ```yaml 
   helm upgrade -i gloo-platform gloo-platform/gloo-platform \
   --namespace gloo-system \
   --version={{< readfile file="static/content/version-platform.md" markdown="true">}} \
   -f - <<EOF
   common:
     adminNamespace: "gloo-system"
     cluster: $CLUSTER_NAME
   featureGates:
     insightsConfiguration: true
   glooInsightsEngine:
     enabled: true
   glooUi:
     enabled: true
   licensing:
     glooGatewayLicenseKey: $GLOO_GATEWAY_LICENSE_KEY
   prometheus:
     enabled: true
   telemetryCollector:
     enabled: true
     mode: deployment
     replicaCount: 1
   EOF
   ```

5. Verify that the Gloo UI components are successfully installed. 
   ```sh
   kubectl get pods -n gloo-system
   ```
   
   Example output: 
   {{< highlight yaml "hl_lines=4-6" >}}
   NAME                                        READY   STATUS    RESTARTS   AGE
   extauth-f7695bf7f-f6dkt                     1/1     Running   0          10m
   gloo-587b79d556-tpvfj                       1/1     Running   0          10m
   gloo-mesh-ui-66db8d9584-kgjld               3/3     Running   0          72m
   gloo-telemetry-collector-68b8cf6f49-zhx87   1/1     Running   0          57m
   prometheus-server-7484d8bfd-tx5s4           2/2     Running   0          72m
   rate-limit-557dcb857f-9zq2t                 1/1     Running   0          10m
   redis-5d6c6bcd4-cnmbm                       1/1     Running   0          10m
   {{< /highlight >}}
   
   
## Next

Continue with [exploring the features of the Gloo UI]({{< versioned_link_path fromRoot="/guides/observability/ui/explore" >}}). 