---
title: "glooctl create secret apikey"
weight: 5
---
## glooctl create secret apikey

Create an ApiKey secret with the given name (Enterprise)

### Synopsis

Create an ApiKey secret with the given name. The ApiKey secret contains a single apikey. This is an enterprise-only feature.

```
glooctl create secret apikey [flags]
```

### Options

```
      --apikey string           apikey to be stored in secret
      --apikey-generate         generate an apikey
      --apikey-labels strings   comma-separated labels for the apikey secret (key=value)
  -h, --help                    help for apikey
```

### Options inherited from parent commands

```
      --dry-run             print kubernetes-formatted yaml rather than creating or updating a resource
  -i, --interactive         use interactive mode
      --name string         name of the resource to read or write
  -n, --namespace string    namespace for reading or writing resources (default "gloo-system")
  -o, --output OutputType   output format: (yaml, json, table, kube-yaml) (default table)
```

### SEE ALSO

* [glooctl create secret](../glooctl_create_secret)	 - Create a secret

