---
title: "glooctl edit virtualservice ratelimit"
weight: 5
---
## glooctl edit virtualservice ratelimit

Configure rate limit settings (Enterprise)

### Synopsis

Let gloo know the location of the rate limit server. This is a Gloo Enterprise feature.

### Options

```
  -h, --help   help for ratelimit
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

* [glooctl edit virtualservice](../glooctl_edit_virtualservice)	 - edit a virtualservice in a namespace
* [glooctl edit virtualservice ratelimit custom-envoy-config](../glooctl_edit_virtualservice_ratelimit_custom-envoy-config)	 - Add a custom rate limit actions (Enterprise)

