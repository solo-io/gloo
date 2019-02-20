---
title: "glooctl install kube"
weight: 5
---
## glooctl install kube

install Gloo on kubernetes to the gloo-system namespace

### Synopsis

requires kubectl to be installed

```
glooctl install kube [flags]
```

### Options

```
  -h, --help   help for kube
```

### Options inherited from parent commands

```
      --docker-email string      Email for docker registry. Use for pulling private images.
      --docker-password string   Password for docker registry authentication. Use for pulling private images.
      --docker-server string     Docker server to use for pulling images (default "https://index.docker.io/v1/")
      --docker-username string   Username for Docker registry authentication. Use for pulling private images.
  -d, --dry-run                  Dump the raw installation yaml instead of applying it to kubernetes
  -f, --file string              Install Gloo from this kubernetes manifest yaml file rather than from a release
  -i, --interactive              use interactive mode
      --release string           install using this release version. defaults to the latest github release
```

### SEE ALSO

* [glooctl install](../glooctl_install)	 - install gloo on different platforms

