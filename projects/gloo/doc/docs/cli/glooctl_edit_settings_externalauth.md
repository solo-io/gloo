---
title: "glooctl edit settings externalauth"
weight: 5
---
## glooctl edit settings externalauth

Configure external auth settings (Enterprise)

### Synopsis

Let gloo know the location of the ext auth server. This is a Gloo Enterprise feature.

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
  -o, --output OutputType         output format: (yaml, json, table, kube-yaml) (default kube-yaml)
      --resource-version string   the resource version of the resource we are editing. if not empty, resource will only be changed if the resource version matches
```

### SEE ALSO

* [glooctl edit settings](../glooctl_edit_settings)	 - root command for settings

