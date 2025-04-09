---
title: "glooctl install"
description: "Reference for the 'glooctl install' command."
weight: 5
---
## glooctl install

install gloo on different platforms

### Synopsis

choose which version of Gloo to install.

### Options

```
      --create-namespace      Create the namespace to install gloo into (default true)
  -d, --dry-run               Dump the raw installation yaml instead of applying it to kubernetes
  -f, --file string           Install Gloo from this Helm chart archive file rather than from a release
  -h, --help                  help for install
  -n, --namespace string      namespace to install gloo into (default "gloo-system")
      --release-name string   helm release name (default "gloo")
      --values strings        List of files with value overrides for the Gloo Helm chart, (e.g. --values file1,file2 or --values file1 --values file2)
  -v, --verbose               If true, output from kubectl commands will print to stdout/stderr
      --version string        version to install (e.g. 1.4.0, defaults to latest)
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
* [glooctl install gateway](../glooctl_install_gateway)	 - install the Gloo Gateway on Kubernetes
* [glooctl install ingress](../glooctl_install_ingress)	 - install the Gloo Ingress Controller on Kubernetes

