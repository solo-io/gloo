---
title: "glooctl edit settings"
weight: 5
---
## glooctl edit settings

root command for settings

### Synopsis

root command for settings

### Options

```
  -h, --help   help for settings
```

### Options inherited from parent commands

```
  -i, --interactive               use interactive mode
      --name string               name of the resource to read or write
  -n, --namespace string          namespace for reading or writing resources (default "gloo-system")
  -o, --output OutputType         output format: (yaml, json, table, kube-yaml) (default kube-yaml)
      --resource-version string   the resource version of the resource we are editing. if not empty, resource will only be changed if the resource version matches
```

### SEE ALSO

* [glooctl edit](../glooctl_edit)	 - Edit a Gloo resource
* [glooctl edit settings externalauth](../glooctl_edit_settings_externalauth)	 - Configure external auth settings (Enterprise)
* [glooctl edit settings ratelimit](../glooctl_edit_settings_ratelimit)	 - Configure rate limit settings (Enterprise)

