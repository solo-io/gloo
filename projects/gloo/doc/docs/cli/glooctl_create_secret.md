---
title: "glooctl create secret"
weight: 5
---
## glooctl create secret

Create a secret

### Synopsis

Create a secret

```
glooctl create secret [flags]
```

### Options

```
  -h, --help   help for secret
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

* [glooctl create](../glooctl_create)	 - Create a Gloo resource
* [glooctl create secret apikey](../glooctl_create_secret_apikey)	 - Create an ApiKey secret with the given name (Enterprise)
* [glooctl create secret aws](../glooctl_create_secret_aws)	 - Create an AWS secret with the given name
* [glooctl create secret azure](../glooctl_create_secret_azure)	 - Create an Azure secret with the given name
* [glooctl create secret oauth](../glooctl_create_secret_oauth)	 - Create an OAuth secret with the given name (Enterprise)
* [glooctl create secret tls](../glooctl_create_secret_tls)	 - Create a secret with the given name

