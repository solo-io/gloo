---
title: "glooctl edit settings ratelimit custom-server-config"
weight: 5
---
## glooctl edit settings ratelimit custom-server-config

Add a custom rate limit settings (Enterprise)

### Synopsis

This allows using lyft rate limit server configuration language to configure the rate limit server.
		For more information see: https://github.com/lyft/ratelimit
		Note: do not add the 'domain' configuration key.
		This is a Gloo Enterprise feature.

```
glooctl edit settings ratelimit custom-server-config [flags]
```

### Options

```
  -h, --help   help for custom-server-config
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

* [glooctl edit settings ratelimit](../glooctl_edit_settings_ratelimit)	 - Configure rate limit settings (Enterprise)

