---
title: "glooctl proxy logs"
description: "Reference for the 'glooctl proxy logs' command."
weight: 5
---
## glooctl proxy logs

dump Envoy logs from one of the proxy instancesNote: this will enable verbose logging on Envoy

```
glooctl proxy logs [flags]
```

### Options

```
  -d, --debug    enable debug logging on the proxy as part of this command (default true)
  -f, --follow   enable debug logging on the proxy as part of this command
  -h, --help     help for logs
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
      --name string                the name of the proxy pod/deployment to use
  -n, --namespace string           namespace for reading or writing resources (default "gloo-system")
      --port string                the name of the service port to connect to (default "http")
      --use-consul                 use Consul Key-Value storage as the backend for reading and writing config (VirtualServices, Upstreams, and Proxies)
```

### SEE ALSO

* [glooctl proxy](../glooctl_proxy)	 - interact with proxy instances managed by Gloo

