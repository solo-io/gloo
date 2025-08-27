---
title: Retries
weight: 100
description: Implement upstream retries when experiencing transient network errors
---

Specifies the retry policy for the route where you can say for a specific error condition how many times to retry and
for how long to try.

* `retryOn` : specifies the condition under which to retry the forward request to the upstream. Same as [Envoy x-envoy-retry-on](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#x-envoy-retry-on).
* `numRetries` : (default: 1) optional attribute that specifies the allowed number of retries.
* `perTryTimeout` : optional attribute that specifies the timeout per retry attempt. Is of type [Google.Protobuf.WellKnownTypes.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration).

{{< highlight yaml "hl_lines=20-23" >}}
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: 'default'
  namespace: 'gloo-system'
spec:
  virtualHost:
    domains:
    - '*'
    routes:
    - matchers:
       - prefix: '/petstore'
      routeAction:
        single:
          upstream:
            name: 'default-petstore-8080'
            namespace: 'gloo-system'
      options:
        retries:
          retryOn: 'connect-failure'
          numRetries: 3
          perTryTimeout: '5s'
{{< /highlight >}}


## Retries for rate limited requests

You can configure a separate retry backoff strategy for requests that are rate limited. This way, you can prevent retries on already rate limited services. 

To determine if a request was rate limited, a specific reset header is used, such as the `X-RateLimit-Reset` or `Retry-After` headers. If the header is present in the response, the backoff strategy for rate limited requests is applied and the general retry backoff strategy is ignored.

The following VirtualService configures a general retry and retry backoff strategy, and a separate retry backoff strategy for rate limited services. 

```yaml
kubectl apply -f- <<EOF
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: vs-with-retries
spec:
  virtualHost:
    domains:
      - nginx.example.com
    routes:
      - matchers:
          - prefix: /retries
        routeAction:
          single:
            upstream:
              name: nginx-upstream-basicroute
        options:
          retries:
            retryOn: "5xx"
            retryBackOff:
              baseInterval: "1s"
              maxInterval: "1s"
            rateLimitedRetryBackOff:
              maxInterval: "1s"
              resetHeaders:
                - name: "X-RateLimit-Reset"
                  format: UNIX_TIMESTAMP
                - name: "Retry-After"
                  format: SECONDS
EOF
```

| Setting | Description | 
| -- | -- | 
| `retries.retryOn` | Retries a request if the upstream service returned a 5XX HTTP response code.  | 
| `retries.retryBackoff` |Defines a general retry backoff strategy for requests that return a 5XX HTTP response code. The `baseInterval` is the time that Envoy waits before a retry is performed. An optional jitter is applied to the `baseInterval`. The `maxInterval` setting is the maximum time that Envoy can wait before performing a retry.  | 
| `retries.rateLimitedRetryBackOff` | A separate backoff strategy for requests that were rate limited. The backoff strategy defines the reset headers that must be returned to trigger the backoff strategy. In this example, the  `X-RateLimit-Reset` and `Retry-After` headers are defined. Note that the Envoy proxy matches these headers in order and case insensitive. If one of these headers is present in a response, the general backoff strategy is ignored and the rate limited backoff strategy is applied. To determine the time to wait before a retry, Envoy uses the value in the defined reset header. {{% notice note %}} If a reset header contains an interval that is larger than the interval that is defined in `rateLimitedRetryBackOff.maxInterval`, the header is discarded and the next header in the list is tried. If the interval in all headers is larger than the `maxInterval`, the `maxInterval` is applied. {{% /notice %}}| 

For more information, see the [Envoy documentation](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#envoy-v3-api-field-config-route-v3-retrypolicy-rate-limited-retry-back-off). 

