---
title: "glooctl install"
weight: 5
---
## glooctl install

install gloo on different platforms

### Synopsis

choose which version of Gloo to install.

### Options

```
  -d, --dry-run              Dump the raw installation yaml instead of applying it to kubernetes
  -f, --file string          Install Gloo from this Helm chart archive file rather than from a release
  -h, --help                 help for install
      --license-key string   License key to activate GlooE features
  -n, --namespace string     namespace to install gloo into (default "gloo-system")
  -u, --upgrade              Upgrade an existing v1 gateway installation to use v2 CRDs. Set this when upgrading from v0.17.x or earlier versions of gloo
```

### Options inherited from parent commands

```
  -i, --interactive   use interactive mode
```

### SEE ALSO

* [glooctl](../glooctl)	 - CLI for Gloo
* [glooctl install gateway](../glooctl_install_gateway)	 - install the GlooE Gateway on kubernetes
* [glooctl install ingress](../glooctl_install_ingress)	 - install the GlooE Ingress Controller on kubernetes
* [glooctl install knative](../glooctl_install_knative)	 - install Knative with GlooE on kubernetes

