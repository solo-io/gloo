---
title: "glooctl create secret oauth"
weight: 5
---
## glooctl create secret oauth

Create an OAuth secret with the given name

### Synopsis

Create an OAuth secret with the given name. The OAuth secrets contains the client_secret as defined in RFC 6749.

```
glooctl create secret oauth [flags]
```

### Options

```
      --client-secret string   oauth client secret
  -h, --help                   help for oauth
      --name string            name of the resource to read or write
  -n, --namespace string       namespace for reading or writing resources (default "gloo-system")
```

### Options inherited from parent commands

```
  -i, --interactive     use interactive mode
  -o, --output string   output format: (yaml, json, table)
```

### SEE ALSO

* [glooctl create secret](../glooctl_create_secret)	 - Create a secret

