---
title: Upstream Health Checks
weight: 50
description: Automatically monitor the status of Upstreams by configuring health checks for them
---

As part of configuring an Upstream, Gloo Gateway provides the option of adding *health checks* that periodically assess the readiness of the Upstream to receive requests. For more information, see the [Envoy documentation](https://www.envoyproxy.io/docs/envoy/v1.14.1/intro/arch_overview/upstream/health_checking#arch-overview-health-checking). 

{{% notice note %}}
Upstreams with working health checks are not removed from Envoy's service directory, even due to configuration changes. To allow removal in such cases, set `ignoreHealthOnHostRemoval` in the Upstream's configuration.
{{% /notice %}}

## Configuration

Review the following sections for example configuration settings of common use cases. For descriptions of each field, refer to the [API documentation]({{% versioned_link_path fromRoot="/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/options/healthcheck/healthcheck.proto.sk/" %}}).

### Custom paths for HttpHealthChecks

To add custom paths to health check requests, review the following example.

```yaml
spec:
  healthChecks:
  - healthyThreshold: 1
    httpHealthCheck:
      path: /check/healthz
    interval: 30s
    timeout: 10s
    unhealthyThreshold: 1
```

| Setting | Description |
| --- | --- |
| `path` | The specific path to send the health check request to the upstream. Make sure that the backing destination handles health checks along this path. |
| `timeout` | The amount of time that the health check waits for a response before considering the request unsuccessful and timing out. | 
| `unhealthyThreshold` | The number of checks that can fail before the upstream is considered unhealthy. The example allows only one failed health check. |
| `healthyThreshold` | The number of checks that must succeed before an upstream is considered healthy. The example requires only one successful health check.
| `interval` | How often the health check sends requests. Set a value that will not overload your upstream service. |
