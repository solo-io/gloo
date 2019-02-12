

# Rate Limit

Rate limits are defined on the virtual service specification as `spec.virtualHost.virtualHostPlugins.extensions.configs.rate-limit`.


```yaml
rate-limit:
    anonymous_limits:
        requests_per_unit: 1000
        unit: HOUR
    authorized_limits:
        requests_per_unit: 200
        unit: MINUTE
```

- Rate limits can be set for anonymous requests, authorized requests, both, or neither.
- `authorized_requests` represent the rate limits imposed on requests that are associated with a known user id
- `anonymous_requests` represent the rate limits imposed on requests that are not associated with a known user id. In this case, the limit is applied to the request's remote address.

- `requests_per_unit` takes an integer value
- `unit` must be one of these strings: `SECOND`, `MINUTE`, `HOUR`, `DAY`

## Create a new virtual service with rate limits enabled

The minimum required configuration in order to create a new virtual service with anonymous and authorized rate limits enabled is shown below.

In this example, we restrict authorized users to 200 requests per minute and anonymous users to 1000 requests per hour.

```yaml
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: myvs
  namespace: gloo-system
spec:
  virtualHost:
    domains:
    - example.com
    name: gloo-system.myvs
    virtualHostPlugins:
      extensions:
        configs:
          rate-limit:
            anonymous_limits:
              requests_per_unit: 1000
              unit: HOUR
            authorized_limits:
              requests_per_unit: 200
              unit: MINUTE
```

- run `kubectl apply -f <filename>` to create this virtualservice

## Edit the rate limit config on an existing virtual service

Print your virtual service specification with:

`kubectl get virtualservice -n <namespace> <virtual_service_name> -o yaml -f <filename>`

in our example above, we expect to see something like:
```yaml
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"gateway.solo.io/v1","kind":"VirtualService","metadata":{"annotations":{},"name":"myvs","namespace":"gloo-system"},"spec":{"virtualHost":{"domains":["example.com"],"name":"gloo-system.myvs","virtualHostPlugins":{"extensions":{"configs":{"rate-limit":{"anonymous_limits":{"requests_per_unit":1000,"unit":"HOUR"},"authorized_limits":{"requests_per_unit":200,"unit":"MINUTE"}}}}}}}}
  creationTimestamp: 2019-02-11T19:41:36Z
  generation: 1
  name: myvs
  namespace: gloo-system
  resourceVersion: "14919"
  selfLink: /apis/gateway.solo.io/v1/namespaces/gloo-system/virtualservices/myvs
  uid: 0b01d943-2e35-11e9-b196-0800271c7f63
spec:
  virtualHost:
    domains:
    - example.com
    name: gloo-system.myvs
    virtualHostPlugins:
      extensions:
        configs:
          rate-limit:
            anonymous_limits:
              requests_per_unit: 1000
              unit: HOUR
            authorized_limits:
              requests_per_unit: 200
              unit: MINUTE
status:
  reported_by: gateway
  state: 1
  subresource_statuses:
    '*v1.Proxy gloo-system gateway-proxy':
      reported_by: gloo
      state: 1
```


### Disable all rate limiting

Edit the spec, removing the rate limit block. Your updated spec should look as follows:

```yaml
spec:
  virtualHost:
    domains:
    - example.com
    name: gloo-system.myvs
    virtualHostPlugins:
      extensions:
        configs:
```
run `kubectl apply -f <filename>` to update the virtualservice

### Disable authorized rate limiting

Edit the spec, removing the authorized rate limit block. Your updated spec should look as follows:

```yaml
spec:
  virtualHost:
    domains:
    - example.com
    name: gloo-system.myvs
    virtualHostPlugins:
      extensions:
        configs:
            anonymous_limits:
              requests_per_unit: 1000
              unit: HOUR
```
run `kubectl apply -f <filename>` to update the virtualservice

### Disable anonymous rate limiting

Edit the spec, removing the anonymous rate limit block. Your updated spec should look as follows:

```yaml
spec:
  virtualHost:
    domains:
    - example.com
    name: gloo-system.myvs
    virtualHostPlugins:
      extensions:
        configs:
            authorize_limits:
              requests_per_unit: 1000
              unit: HOUR
```
run `kubectl apply -f <filename>` to update the virtualservice
