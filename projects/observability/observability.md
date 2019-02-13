# Observability

Gloo's core routing capabilities are built on top of [Envoy](https://www.envoyproxy.io/), and therefore can take
advantage of Envoy's powerful observability features. Envoy exposes a multitude of statistics as described 
[here](https://www.envoyproxy.io/docs/envoy/latest/configuration/http_conn_man/stats), which gloo in turn exposes to the user
along with more gloo-specific metrics.

All of the statistics gloo exposes are in the [prometheus](https://prometheus.io/) format. By default gloo adds the necessary
annotations to be discovered by prometheus. More information about prometheus service discovery can be found 
[here](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#kubernetes_sd_config).


### Statistics via Prometheus

Gloo uses all of the standard prometheus annotations to be discovered, therefore any prometheus deployment can be used. 
By default GlooE ships with a prometheus deployment, based off of the official prometheus helm chart, located 
[here](https://github.com/helm/charts/tree/master/stable/prometheus).

The default prometheus deployment can be reached by running the following:
```bash
kubectl port-forward -n gloo-system deployment/glooe-prometheus-server 9090
```
After running the above command navigate to [localhost:9090](localhost:9090) to view the statistics as well as admin 
information for prometheus.

### Visualization via Grafana

[Grafana](https://grafana.com/) is one of the most ubiquitous pieces of open source visualization software currently available.
Grafana works with many different data sources, including prometheus. Similar to prometheus, GlooE ships with a simple deployment 
of Grafana based off of the official Grafana helm chart, located [here](https://github.com/helm/charts/tree/master/stable/grafana).

In addition to including Grafana in the default installation. The GlooE UI integrates closely with Grafana. The grafana server 
can be reached by navigating to `https:<gloo-ui-location>/grafana`. However, the UI also exposes some metrics automatically.
These can be reached by navigating to the stats tab of the GlooE UI. These include graphs detailing the traffic being routed
to the various services Gloo has discovered, as well as some basic information about the health of the cluster.

## Admin pages

As a convenience Gloo ships with optional admin pages, similar to Envoy's admin page. By default these pages are turned off
for the sake of efficiency, but they can be turned on very easily. In order to enable the stats page for a given pod set the env 
variable `START_STATS_SERVER=true` in the deployment template spec.
 
By default this pod will not be scraped by the prometheus deployment, but this functionality can be added by adding the following
annotations to the deployment template spec.
```yaml

annotations:
  prometheus.io/path: /metrics
  prometheus.io/port: "8081"
  prometheus.io/scrape: "true"
```
 
Am example deployment might look something like below: 

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: <pod-name>
  namespace: gloo-system
spec:
  replicas: 1
  template:
    metadata: 
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: "8081"
        prometheus.io/scrape: "true"
    spec:
      containers:
      - image: "<image-name>"
        imagePullPolicy: Always
        name: <continer-name>
        env:
        - name: "START_STATS_SERVER"
          value: true
```

Another option for exposing these stats to prometheus is using our helm chart. Instruction on how to install Gloo via our helm chart
can be found [here](https://github.com/solo-io/gloo/blob/master/docs/installation/install_with_helm.md). The same rules outlined in 
that doc apply to GlooE, as well as more customization options for our enterprise features.
