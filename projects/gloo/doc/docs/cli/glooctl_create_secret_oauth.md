---
title: "glooctl create secret oauth"
weight: 5
---
## glooctl create secret oauth

Create an OAuth secret with the given name (Enterprise)

### Synopsis

Create an OAuth secret with the given name. The OAuth secrets contains the client_secret as defined in RFC 6749. This is an enterprise-only feature.

```
glooctl create secret oauth [flags]
```

### Options

```
      --client-secret string   oauth client secret
  -h, --help                   help for oauth
```

### Options inherited from parent commands

```
      --dry-run            print kubernetes-formatted yaml rather than creating or updating a resource
  -i, --interactive        use interactive mode
      --name string        name of the resource to read or write
  -n, --namespace string   namespace for reading or writing resources (default "gloo-system")
  -o, --output string      output format: (yaml, json, table)
      --yaml               print basic (non-kubernetes) yaml rather than creating or updating a resource
```

### SEE ALSO

* [glooctl create secret](../glooctl_create_secret)	 - Create a secret

