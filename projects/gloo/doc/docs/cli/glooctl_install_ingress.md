---
title: "glooctl install ingress"
weight: 5
---
## glooctl install ingress

install the GlooE Ingress Controller on kubernetes

### Synopsis

requires kubectl to be installed

```
glooctl install ingress [flags]
```

### Options

```
  -h, --help   help for ingress
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

