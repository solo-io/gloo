---
title: "glooctl install gateway"
weight: 5
---
## glooctl install gateway

install the GlooE Gateway on kubernetes

### Synopsis

requires kubectl to be installed

```
glooctl install gateway [flags]
```

### Options

```
  -h, --help   help for gateway
```

### Options inherited from parent commands

```
  -d, --dry-run            Dump the raw installation yaml instead of applying it to kubernetes
  -f, --file string        Install Gloo from this Helm chart archive file rather than from a release
  -i, --interactive        use interactive mode
  -n, --namespace string   namespace to install gloo into (default "gloo-system")
```

### SEE ALSO

* [glooctl install](../glooctl_install)	 - install gloo on different platforms

