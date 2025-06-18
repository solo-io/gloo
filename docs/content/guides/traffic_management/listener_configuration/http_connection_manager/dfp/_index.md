---
title: Dynamic Forward Proxy
weight: 10
---

You can set up an [HTTP Dynamic Forward Proxy (DFP) filter](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/http/http_proxy) in Gloo Gateway.

In a highly dynamic environment with services coming up and down and with no service registry being able to list the available endpoints, one option is to somehow "blindly" route the client requests upstream. 

Another popular use case is to deploy a forward proxy for all egress traffic. This way, you can observe and control outbound traffic. A common security policy applied here is rate-limiting and, of course, gathering access logs.

So there are two sorts of usage for a forward proxy:
- Route all the traffic of your local network to an exit node with the forward proxy, to monitor and control all the egress traffic.
- Apply the forward proxy filter only to certain routes, which don’t have a pre-defined destination, and dynamically build the final Host value.

Before implementing a dynamic forward proxy, consider the downsides to such flexibility:
- Because no pre-defined {{< protobuf name="gloo.solo.io.Upstream" display="Upstream" >}} designates the upstream service, you cannot configure failover policies or client load-balancing.
- DNS resolution is done at runtime. Typically, when Envoy encounters a domain name the first time, Envoy pauses the request and synchronously resolves this domain to get the endpoints (IP addresses). Then, these entries are put into a local cache.

Of course, you might still decide to use a dynamic forward proxy in an API Gateway for benefits such as the following:
- You easily get metrics on the egress traffic that goes through the forward proxy.
- You can enforce authentication and authorization policies.
- You can leverage other policies available in Gloo Gateway Enterprise, like Web Application Firewall (WAF) or Data Loss Prevention (DLP).

## Enabling the Dynamic Forward Proxy

First, enable the DFP filter in your Gateway configuration.

```bash
kubectl -n gloo-system patch gw/gateway-proxy --type merge -p "
spec:
  httpGateway:
    options:
      dynamicForwardProxy: {}
"
```

Then, set the actual destination of the client request. The destination can be the `Host` header in the most basic setup. However, the destination might also be hidden in other client request headers or body parts. In this latter case, you can create the header dynamically by using a transformation template.

Review the following basic example to see how you can apply the dynamic forward option to a route.

```yaml
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: test-static
  namespace: gloo-system
spec:
  virtualHost:
    domains:
      - 'foo'
    routes:
      - matchers:
         - prefix: /
        routeAction:
          dynamicForwardProxy:
            autoHostRewriteHeader: "x-rewrite-me" # host header will be rewritten to the value of this header
```

## Set circuit breakers for dynamically discovered upstreams {#dfp-circuit-breakers}

By default, Envoy allows a maximum of 1024 connections to dynamically discovered upstreams. Starting in Gloo Gateway Enterprise version 1.19.3, you can override these settings by defining custom circuit breakers for your dynamic forward proxy. 

1. Update your gateway proxy to add custom circuit breaker settings for your dynamic forward proxy. The following example overrides the default number of connections, requests, pending requests, and retries to 40,000. For more information, see the [Circuit breaker configuration]({{% versioned_link_path fromRoot="/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/circuit_breaker/circuit_breaker.proto.sk/#circuitbreakerconfig" %}}).

   ```bash
   kubectl -n gloo-system patch gw/gateway-proxy --type merge -p "
   spec:
     httpGateway:
       options:
         dynamicForwardProxy: 
           circuitBreakers:
             maxConnections: 40000
             maxPendingRequests: 40000
             maxRequests: 40000
             maxRetries: 40000
             trackRemaining: true
   "
   ```

   | Setting | Description | 
   | -- | -- | 
   | `maxConnections` | The maximum number of connections that Envoy makes to the upstream cluster. If not specified, the default is 1024. | 
   | `maxPendingRequests` | The maximum number of pending requests that Envoy allows to the upstream cluster. If not specified, the default is 1024.| 
   | `maxRequests` | The maximum number of parallel requests that Envoy makes to the upstream cluster. If not specified, the default is 1024.| 
   | `maxRetries` | The maximum number of parallel retries that Envoy allows to the upstream cluster. If not specified, the default is 3. | 
   | `trackRemaining` | If set to `true`, then stats are published that expose the number of resources remaining until the circuit breakers open. If not specified, the default is `false`.| 

2. Verify that your circuit breaker settings were applied. 
   1. Port-forward your gateway proxy on port 19000. 
      ```sh
      kubectl port-forward deploy/gateway-proxy 19000:19000 -n gloo-system
      ```
   2. Get the config dump for your gateway proxy. 
      ```sh
      curl http://localhost:19000/config_dump | jq -r '.configs[] | select(.["@type"]=="type.googleapis.com/envoy.admin.v3.ClustersConfigDump") | .dynamic_active_clusters | map(select(.cluster.name | startswith("solo_io_generated_dfp")))'
      ```
      
      Example output: 
      ```
      ...
      "version_info": "5428280815645392517",
       "cluster": {
         "@type": "type.googleapis.com/envoy.config.cluster.v3.Cluster",
         "name": "solo_io_generated_dfp:8973398522796813204",
         "connect_timeout": "5s",
         "lb_policy": "CLUSTER_PROVIDED",
         "circuit_breakers": {
           "thresholds": [
             {
               "max_connections": 40000,
               "max_pending_requests": 40000,
               "max_requests": 40000,
               "max_retries": 40000,
               "track_remaining": true
             }
           ]
        },
      ...
      ```




