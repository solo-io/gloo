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
      --docker-email string      Email for docker registry. Use for pulling private images.
      --docker-password string   Password for docker registry authentication. Use for pulling private images.
      --docker-server string     Docker server to use for pulling images (default "https://index.docker.io/v1/")
      --docker-username string   Username for Docker registry authentication. Use for pulling private images.
  -d, --dry-run                  Dump the raw installation yaml instead of applying it to kubernetes
  -f, --file string              Install Gloo from this kubernetes manifest yaml file rather than from a release
  -h, --help                     help for install
      --release string           install using this release version. defaults to the latest github release
```

### Options inherited from parent commands

```
  -i, --interactive   use interactive mode
```

### SEE ALSO

* [glooctl](../glooctl)	 - CLI for Gloo
* [glooctl install kube](../glooctl_install_kube)	 - install Gloo on kubernetes to the gloo-system namespace

