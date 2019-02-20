---
title: "glooctl edit settings externalauth"
weight: 5
---
## glooctl edit settings externalauth

Configure external auth settings

### Synopsis

Let gloo know the location of the ext auth server

```
glooctl edit settings externalauth [flags]
```

### Options

```
      --extauth-server-name string        name of the ext auth server upstream
      --extauth-server-namespace string   namespace of the ext auth server upstream
  -h, --help                              help for externalauth
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

