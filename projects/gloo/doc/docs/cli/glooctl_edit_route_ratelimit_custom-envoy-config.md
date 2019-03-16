---
title: "glooctl edit route ratelimit custom-envoy-config"
weight: 5
---
## glooctl edit route ratelimit custom-envoy-config

Add a custom rate limit actions (Enterprise)

### Synopsis

This allows using envoy actions to specify your rate limit descriptors.
		For available actions and more information see: https://www.envoyproxy.io/docs/envoy/v1.9.0/api-v2/api/v2/route/route.proto#route-ratelimit-action
		
		This is a Gloo Enterprise feature.

```
glooctl edit route ratelimit custom-envoy-config [flags]
```

### Options

```
  -h, --help   help for custom-envoy-config
```

### Options inherited from parent commands

```
  -x, --index uint32              edit the route with this index in the virtual service route list
  -i, --interactive               use interactive mode
      --name string               name of the resource to read or write
  -n, --namespace string          namespace for reading or writing resources (default "gloo-system")
  -o, --output string             output format: (yaml, json, table)
      --resource-version string   the resource version of the resouce we are editing. if not empty, resource will only be changed if the resource version matches
```

### SEE ALSO

* [glooctl edit route ratelimit](../glooctl_edit_route_ratelimit)	 - Configure rate limit settings (Enterprise)

