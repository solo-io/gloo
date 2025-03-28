---
title: "glooctl license validate"
description: "Reference for the 'glooctl license validate' command."
weight: 5
---
## glooctl license validate

Check Gloo Gateway License Validity

### Synopsis

Checking Gloo Gateway license Validity.

Usage: `glooctl license validate [--license-key license-key]`

```
glooctl license validate [flags]
```

### Options

```
  -h, --help                 help for validate
      --license-key string   license key to validate
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

* [glooctl license](../glooctl_license)	 - subcommands for interacting with the license

