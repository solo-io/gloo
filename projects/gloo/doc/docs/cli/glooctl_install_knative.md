---
title: "glooctl install knative"
weight: 5
---
## glooctl install knative

install Knative with GlooE on kubernetes

### Synopsis

requires kubectl to be installed

```
glooctl install knative [flags]
```

### Options

```
  -h, --help   help for knative
```

### Options inherited from parent commands

```
  -d, --dry-run              Dump the raw installation yaml instead of applying it to kubernetes
  -f, --file string          Install Gloo from this Helm chart archive file rather than from a release
  -i, --interactive          use interactive mode
      --license-key string   License key to activate GlooE features
  -n, --namespace string     namespace to install gloo into (default "gloo-system")
  -u, --upgrade              Upgrade an existing v1 gateway installation to use v2 CRDs. Set this when upgrading from v0.17.x or earlier versions of gloo
```

### SEE ALSO

* [glooctl install](../glooctl_install)	 - install gloo on different platforms

