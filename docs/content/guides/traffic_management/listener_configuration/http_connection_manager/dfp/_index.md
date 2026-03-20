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
kubectl apply -f- <<EOF
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
EOF
```

## HTTPS tunneling with Dynamic Forward Proxy

To route traffic to HTTPS targets via the dynamic forward proxy, set the `connectTerminate` upgrade type on the VirtualService. This setting instructs Envoy to terminate the HTTP `CONNECT` request, resolve the target host using DFP DNS, and forward the raw TCP payload upstream. This configuration allows the client to complete a TLS handshake directly with the destination without Envoy decrypting or re-encrypting the traffic.

To enable HTTPS tunneling, you need two settings:
1. **Gateway**: Enable the `CONNECT` upgrade type on the **httpConnectionManagerSettings**. By default, Envoy rejects HTTP `CONNECT` requests unless the `CONNECT` upgrade is explicitly allowed at the listener level. This Gateway-level setting tells the HttpConnectionManager to accept `CONNECT` as a valid upgrade protocol, which is a prerequisite for any route-level handling.
2. **VirtualService**: Enable `connectTerminate` in the route's upgrade options. This setting instructs Envoy to terminate the `CONNECT` handshake and forward the raw TCP stream to the upstream, rather than proxying the `CONNECT` request as a regular HTTP request.

{{% notice warning %}}
`connectTerminate` creates a raw TCP tunnel between the client and upstream. Because Envoy forwards the payload as opaque bytes without inspecting HTTP headers or applying HTTP-level policies, any HTTP request filters (such as header manipulation, WAF rules, or authorization checks) configured on the route do not apply to the tunneled traffic. An attacker could use the tunnel to bypass those controls and reach the upstream directly. Ensure that network-level controls or separate mTLS policies are in place to restrict what clients can tunnel to before enabling this feature in production.
{{% /notice %}}

### Configure CONNECT tunneling

1. Apply the Gateway resource to enable the CONNECT upgrade for the dynamic forward proxy. This example updates the default `gateway-proxy` Gateway resource. Note that `bindAddress` and `bindPort` must be unique across all Gateways on the same proxy.

   ```shell
   kubectl apply -f- <<EOF
   apiVersion: gateway.solo.io/v1
   kind: Gateway
   metadata:
     name: gateway-proxy
     namespace: gloo-system
   spec:
     bindAddress: "::"
     bindPort: 8080
     httpGateway:
       options:
         dynamicForwardProxy:
           dnsCacheConfig:
             dnsLookupFamily: V4_ONLY
             hostTtl: 86400s
         httpConnectionManagerSettings:
           upgrades:
             - connect:
                 enabled: true
       virtualServiceSelector:
         app: my-forward-proxy
     ssl: false
     useProxyProto: false
   EOF
   ```

   {{% notice note %}}
   To get the current values for `bindAddress`, `bindPort`, `ssl`, `useProxyProto`, and `proxyNames` on your existing Gateway resource before applying, run `kubectl get gw gateway-proxy -n gloo-system -o yaml`.
   {{% /notice %}}

2. Create a VirtualService that matches `CONNECT` requests and terminates them by using the `connectTerminate` setting. 
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: gateway.solo.io/v1
   kind: VirtualService
   metadata:
     name: forward-proxy-vs
     namespace: gloo-system
     labels: 
       app: my-forward-proxy
   spec:
     virtualHost:
       domains:
         - '*'
       routes:
         - matchers:
             - connectMatcher: {}   # matches HTTP CONNECT requests
           routeAction:
             dynamicForwardProxy: {}   # host is taken from the CONNECT request authority
           options:
             upgrades:
               - connectTerminate:
                   enabled: true   # terminates CONNECT and forwards payload as raw TCP
   EOF
   ```

3. Send an HTTPS request to the `httpbin.org` site through the gateway proxy. When you target an HTTPS URL through an HTTP dynamic forward proxy, `curl` automatically sends a `CONNECT` request to establish the tunnel before performing the TLS handshake.
   ```sh
   export INGRESS_GW_ADDRESS=$(kubectl get svc -n gloo-system gateway-proxy \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
   curl -v -x http://$INGRESS_GW_ADDRESS:80 https://httpbin.org/get 
   ```

   In the CLI output, verify that you see that the `CONNECT` tunnel was established and that you successfully connected to the httpbin.org site via HTTPS.

   {{< highlight shell "hl_lines=2 4 5 10 13 14" >}}
   * Connected to 172.18.0.4 (172.18.0.4) port 80
   * CONNECT tunnel: HTTP/1.1 negotiated
   * allocate connect buffer
   * Establish HTTP proxy tunnel to httpbin.org:443
   > CONNECT httpbin.org:443 HTTP/1.1
   > Host: httpbin.org:443
   > User-Agent: curl/8.7.1
   > Proxy-Connection: Keep-Alive
   >
   < HTTP/1.1 200 OK
   < server: envoy
   <
   * CONNECT phase completed
   * CONNECT tunnel established, response 200
   * ALPN: curl offers h2,http/1.1
   * (304) (OUT), TLS handshake, Client hello (1):
   *  CAfile: /etc/ssl/cert.pem
   *  CApath: none
   * (304) (IN), TLS handshake, Server hello (2):
   * TLSv1.2 (IN), TLS handshake, Certificate (11):
   ...
   * using HTTP/2
   * [HTTP/2] [1] OPENED stream for https://httpbin.org/get
   * [HTTP/2] [1] [:method: GET]
   * [HTTP/2] [1] [:scheme: https]
   * [HTTP/2] [1] [:authority: httpbin.org]
   * [HTTP/2] [1] [:path: /get]
   * [HTTP/2] [1] [user-agent: curl/8.7.1]
   * [HTTP/2] [1] [accept: */*]
   > GET /get HTTP/2
   > Host: httpbin.org
   > User-Agent: curl/8.7.1
   > Accept: */*
   >
   * Request completely sent off
   < HTTP/2 200
   < content-type: application/json
   < content-length: 252
   < server: gunicorn/19.9.0
   < access-control-allow-origin: *
   < access-control-allow-credentials: true
   <
   {
     "args": {},
     "headers": {
       "Accept": "*/*",
       "Host": "httpbin.org",
       "User-Agent": "curl/8.7.1",
       "X-Amzn-Trace-Id": "Root=1-69baf890-425cbdb422823e25314bd8bf"
     },
     "origin": "69.X.XX.XXX",
     "url": "https://httpbin.org/get"
     }
   * Connection #0 to host 172.18.0.4 left intact
   {{< /highlight >}}


## Set circuit breakers for dynamically discovered upstreams {#dfp-circuit-breakers}

By default, Envoy allows a maximum of 1024 connections to dynamically discovered upstreams. You can override these settings by defining custom circuit breakers for your dynamic forward proxy. 

1. Update your gateway proxy to add custom circuit breaker settings for your dynamic forward proxy. The following example overrides the default number of connections, requests, pending requests, and retries to 40,000. For more information, see the [Circuit breaker configuration]({{% versioned_link_path fromRoot="/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/circuit_breaker/circuit_breaker.proto.sk/#circuitbreakerconfig" %}})

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



