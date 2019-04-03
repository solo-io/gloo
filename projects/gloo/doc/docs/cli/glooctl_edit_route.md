---
title: "glooctl edit route"
weight: 5
---
## glooctl edit route



### Synopsis



### Options

```
  -h, --help           help for route
  -x, --index uint32   edit the route with this index in the virtual service route list
```

### Options inherited from parent commands

```
  -i, --interactive               use interactive mode
      --name string               name of the resource to read or write
  -n, --namespace string          namespace for reading or writing resources (default "gloo-system")
  -o, --output string             output format: (yaml, json, table)
      --resource-version string   the resource version of the resource we are editing. if not empty, resource will only be changed if the resource version matches
```

### SEE ALSO

* [glooctl edit](../glooctl_edit)	 - Edit a Gloo resource
* [glooctl edit route externalauth](../glooctl_edit_route_externalauth)	 - Configure disable external auth on a route (Enterprise)
* [glooctl edit route ratelimit](../glooctl_edit_route_ratelimit)	 - Configure rate limit settings (Enterprise)

