---
title: "glooctl debug"
description: "Reference for the 'glooctl debug' command."
weight: 5
---
## glooctl debug

Debug Gloo (requires Gloo running on Kubernetes)

### Synopsis

Dumps state information to a local directory (or to stdout, for some subcommands). This is useful for debugging failures. If no subcommand is specified, this command dumps all the state information from the `kube`, `gloo`, and `envoy` subcommands. See the subcommand descriptions for more details.

```
glooctl debug [flags]
```

### Options

```
  -d, --directory string         directory to write debug info to (default "debug")
  -h, --help                     help for debug
  -N, --namespaces stringArray   namespaces from which to dump logs and resources (use flag multiple times to specify multiple namespaces, e.g. '-N gloo-system -N default') (default [gloo-system])
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
* [glooctl debug envoy](../glooctl_debug_envoy)	 - Dump information from the Envoy admin interface for gateway proxies to a local directory (requires Gloo running on Kubernetes)
* [glooctl debug gloo](../glooctl_debug_gloo)	 - Dump Gloo Gateway controller state to a local directory (requires Gloo running on Kubernetes)
* [glooctl debug kube](../glooctl_debug_kube)	 - Dump Kubernetes state to a local directory (requires Gloo running on Kubernetes)
* [glooctl debug logs](../glooctl_debug_logs)	 - Debug Gloo logs (requires Gloo running on Kubernetes)
* [glooctl debug yaml](../glooctl_debug_yaml)	 - Dump YAML representing the current Gloo state (requires Gloo running on Kubernetes)

