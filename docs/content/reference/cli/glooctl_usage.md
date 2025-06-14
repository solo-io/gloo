---
title: "glooctl usage"
description: "Reference for the 'glooctl usage' command."
weight: 5
---
## glooctl usage

Scan Gloo for feature usage

### Synopsis

glooctl usage will evaluate Gloo Gateway snapshots and collect usage stats. It also has the ability to scan for Gloo Gateway proxies and grab their current throughput stats.

```
glooctl usage [flags]
```

### Examples

```
# This command scans Gloo Gateway for feature usage.
# To get usage stats from a running Gloo Gateway control plane.
  glooctl usage

# To get usage stats from a Gloo Gateway snapshot file.
  glooctl usage --input-snapshot ./gg-input.json

# To get usage stats from a Gloo Gateway snapshot file in json format.
  glooctl usage --input-snapshot ./gg-input.json --output-format json

# To get throughput stats from a Gloo Gateway proxy pods.
  glooctl usage --scan-proxies deploy/gateway-proxy

# To get throughput stats from a Gloo Gateway proxy running in a different namespace than the control plane
  glooctl usage --scan-proxies deploy/gateway-proxy --proxy-namespaces gloo-system
	
# To print all the backend endpoint stats per Gloo Gateway proxy (requires --scan-proxies)
  glooctl usage --scan-proxies deploy/gateway-proxy --include-endpoint-stats
```

### Options

```
      --gloo-control-plane string             Name of the Gloo control plane pod (default "deploy/gloo")
  -n, --gloo-control-plane-namespace string   Namespace of the Gloo control plane pod (default "gloo-system")
  -h, --help                                  help for usage
      --include-endpoint-stats                Include endpoint stats in the output (default true)
      --input-snapshot string                 Gloo input snapshot file location
      --output-format string                  Output format (text, json, yaml) (default "yaml")
      --proxy-namespaces strings              Namespaces that contain gloo proxies (default gloo-system or gloo-control-plane-namespace)
      --scan-proxies strings                  Scan for Gloo proxies and grab their routing information
```

### Options inherited from parent commands

```
  -c, --config string              set the path to the glooctl config file (default "<home_directory>/.gloo/glooctl-config.yaml")
      --consul-address string      address of the Consul server. Use with --use-consul (default "127.0.0.1:8500")
      --consul-allow-stale-reads   Allows reading using Consul's stale consistency mode.
      --consul-datacenter string   Datacenter to use. If not provided, the default agent datacenter is used. Use with --use-consul
      --consul-root-key string     key prefix for the Consul key-value storage. (default "gloo")
      --consul-scheme string       URI scheme for the Consul server. Use with --use-consul (default "http")
      --consul-token string        Token is used to provide a per-request ACL token which overrides the agent's default token. Use with --use-consul
  -i, --interactive                use interactive mode
      --kube-context string        kube context to use when interacting with kubernetes
      --kubeconfig string          kubeconfig to use, if not standard one
      --use-consul                 use Consul Key-Value storage as the backend for reading and writing config (VirtualServices, Upstreams, and Proxies)
```

### SEE ALSO

* [glooctl](../glooctl)	 - CLI for Gloo

