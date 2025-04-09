---
title: "glooctl create"
description: "Reference for the 'glooctl create' command."
weight: 5
---
## glooctl create

Create a Gloo resource

### Synopsis

Gloo resources be created from files (including stdin)

```
glooctl create [flags]
```

### Options

```
      --dry-run             print kubernetes-formatted yaml rather than creating or updating a resource
  -f, --file string         file to be read or written to
  -h, --help                help for create
      --name string         name of the resource to read or write
  -n, --namespace string    namespace for reading or writing resources (default "gloo-system")
  -o, --output OutputType   output format: (yaml, json, table, kube-yaml, wide) (default table)
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
* [glooctl create authconfig](../glooctl_create_authconfig)	 - Create an Auth Config
* [glooctl create secret](../glooctl_create_secret)	 - Create a secret
* [glooctl create upstream](../glooctl_create_upstream)	 - Create an Upstream
* [glooctl create upstreamgroup](../glooctl_create_upstreamgroup)	 - Create an Upstream Group
* [glooctl create virtualservice](../glooctl_create_virtualservice)	 - Create a Virtual Service

