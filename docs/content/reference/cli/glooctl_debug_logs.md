---
title: "glooctl debug logs"
description: "Reference for the 'glooctl debug logs' command."
weight: 5
---
## glooctl debug logs

Print Gloo logs from a Kubernetes cluster

```
glooctl debug logs [flags]
```

### Options

```
      --errors-only        filter for error logs only
  -f, --file string        file to be read or written to
  -h, --help               help for logs
  -n, --namespace string   namespace for reading or writing resources (default "gloo-system")
      --zip                save logs to a tar file (specify location with -f)
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
  -d, --directory string           directory to write debug info to (default "debug")
  -i, --interactive                use interactive mode
      --kube-context string        kube context to use when interacting with kubernetes
      --kubeconfig string          kubeconfig to use, if not standard one
  -N, --namespaces stringArray     namespaces from which to dump logs and resources (use flag multiple times to specify multiple namespaces, e.g. '-N gloo-system -N default') (default [gloo-system])
      --use-consul                 use Consul Key-Value storage as the backend for reading and writing config (VirtualServices, Upstreams, and Proxies)
```

### SEE ALSO

* [glooctl debug](../glooctl_debug)	 - Debug Gloo Gateway (requires Gloo running on Kubernetes)

