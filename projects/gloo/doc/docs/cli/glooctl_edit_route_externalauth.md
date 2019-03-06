---
title: "glooctl edit route externalauth"
weight: 5
---
## glooctl edit route externalauth

Configure disable external auth on a route (Enterprise)

### Synopsis

Allows disabling external auth on specific routes. External auth is a gloo enterprise feature.

```
glooctl edit route externalauth [flags]
```

### Options

```
  -d, --disable        set to true to disable authentication on this route
  -h, --help           help for externalauth
  -x, --index uint32   edit the route with this index in the virtual service route list
```

### Options inherited from parent commands

```
  -i, --interactive               use interactive mode
      --name string               name of the resource to read or write
  -n, --namespace string          namespace for reading or writing resources (default "gloo-system")
  -o, --output string             output format: (yaml, json, table)
      --resource-version string   the resource version of the resouce we are editing. if not empty, resource will only be changed if the resource version matches
```

### SEE ALSO

* [glooctl edit route](../glooctl_edit_route)	 - 

