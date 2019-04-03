---
title: "glooctl edit route ratelimit"
weight: 5
---
## glooctl edit route ratelimit

Configure rate limit settings (Enterprise)

### Synopsis

Let gloo know the location of the rate limit server. This is a Gloo Enterprise feature.

### Options

```
  -h, --help   help for ratelimit
```

### Options inherited from parent commands

```
  -x, --index uint32              edit the route with this index in the virtual service route list
  -i, --interactive               use interactive mode
      --name string               name of the resource to read or write
  -n, --namespace string          namespace for reading or writing resources (default "gloo-system")
  -o, --output string             output format: (yaml, json, table)
      --resource-version string   the resource version of the resource we are editing. if not empty, resource will only be changed if the resource version matches
```

### SEE ALSO

* [glooctl edit route](../glooctl_edit_route)	 - 
* [glooctl edit route ratelimit custom-envoy-config](../glooctl_edit_route_ratelimit_custom-envoy-config)	 - Add a custom rate limit actions (Enterprise)

