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
            retryOn: "retriable-status-codes"
            retriable_status_codes:
               - 502
               - 503
               - 504
               - 429
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
| `retries.retryOn` | Retries a request if the upstream service returns one of the retriable HTTP response codes that are defined in `retriable_status_codes`.  | 
| `retries.retryBackoff` |Defines a general retry backoff strategy. The `baseInterval` is the time that Envoy waits before a retry is performed. An optional jitter is applied to the `baseInterval`. The `maxInterval` setting is the maximum time that Envoy can wait before performing a retry.  | 
| `retries.rateLimitedRetryBackOff` | A separate retry backoff strategy for requests that were rate limited. The backoff strategy is applied only when specific headers are present in the response. For example, if a 429 HTTP response code is returned, which is typical for rate limited requests, but the headers that are defined in the rate limited backoff strategy are not present, the strategy is not applied. Instead, the general backoff strategy that is defined in `retries.retryBackoff` is applied. </br> </br> For Envoy to apply the rate limited retry backoff strategy, Envoy first determines whether a request must be retried by checking the defined HTTP response codes in the `retryOn` setting. If a retriable HTTP response code is found, Envoy then checks for the backoff reset headers that are defined in the `rateLimitedRetryBackOff` setting. This example defines the  `X-RateLimit-Reset` and `Retry-After` headers. If one of these headers is present in a response, the general backoff strategy is ignored and the rate limited retry backoff strategy is applied. To determine the time to wait before a retry, Envoy uses the value in the defined reset header. {{% notice note %}} If a reset header contains an interval that is larger than the interval that is defined in `rateLimitedRetryBackOff.maxInterval`, the header is discarded and the next header in the list is tried. If the interval in all headers is larger than the `maxInterval`, the `maxInterval` is applied. {{% /notice %}}| 

For more information, see the [Envoy documentation](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#envoy-v3-api-field-config-route-v3-retrypolicy-rate-limited-retry-back-off). 
