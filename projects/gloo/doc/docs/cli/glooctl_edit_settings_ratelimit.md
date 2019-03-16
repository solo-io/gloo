---
title: "glooctl edit settings ratelimit"
weight: 5
---
## glooctl edit settings ratelimit

Configure rate limit settings (Enterprise)

### Synopsis

Let gloo know the location of the rate limit server. This is a Gloo Enterprise feature.

```
glooctl edit settings ratelimit [flags]
```

### Options

```
  -h, --help                                help for ratelimit
      --ratelimit-server-name string        name of the ext rate limit upstream
      --ratelimit-server-namespace string   namespace of the ext rate limit upstream
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

* [glooctl edit settings](../glooctl_edit_settings)	 - root command for settings
* [glooctl edit settings ratelimit custom-server-config](../glooctl_edit_settings_ratelimit_custom-server-config)	 - Add a custom rate limit settings (Enterprise)

