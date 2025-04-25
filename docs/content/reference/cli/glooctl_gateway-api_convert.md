---
title: "glooctl gateway-api convert"
description: "Reference for the 'glooctl gateway-api convert' command."
weight: 5
---
## glooctl gateway-api convert

Convert Gloo Edge APIs to Gateway API

### Synopsis

Convert Gloo Edge APIs to Gateway APIs by either providing Kubernetes YAML files or a Gloo Gateway input snapshot.

```
glooctl gateway-api convert [flags]
```

### Examples

```
# This command converts Gloo Edge APIs to Kubernetes Gateway API YAML files and places them in the '--output-dir' directory, grouped by namespace.
# To generate Gateway API YAML files from a Gloo Gateway snapshot that is retrieved from a running 'gloo' pod. The 'output-dir' must not exist.
  glooctl gateway-api convert --gloo-control-plane deploy/gloo --output-dir ./_output

# To generate Gateway API YAML files from a single Kubernetes YAML file. The 'output-dir' must not exist.
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output

# To delete and recreate the content in the 'output-dir', add the 'delete-output-dir'' option.
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --delete-output-dir

# To generate Gateway API YAML files from a single Kubernetes YAML file, but place all the output configurations in to the same file. 
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --retain-input-folder-structure

# To load a bunch of '*.yaml' or '*.yml' files in nested directories. You can also use the '--retain-input-folder-structure' option to keep the original file structure, which can be helpful in CI/CD pipelines.
  glooctl gateway-api convert --input-dir ./gloo-configs --output-dir ./_output --retain-input-folder-structure

To download a Gloo Gateway snapshot from a running 'gloo' pod (verison 1.17+) and generate Gateway API YAML files from that snapshot. 
  kubectl -n gloo-system port-forward deploy/gloo 9091
  curl localhost:9091/snapshots/input > gg-input.json
  
  glooctl gateway-api convert --input-snapshot gg-input.json --output-dir ./_output

# To get the stats for each migration, such as the number of configuration files that were generated, add the '--print-stats' option. 
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --print-stats

# To retain non-Gloo Gateway API YAML files, add  the '--include-unknown' option. 
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --include-unknown
```

### Options

```
      --combine-route-options                 Combine routeOptions that are exactly the same and share them among the HTTPRoutes
      --create-namespaces                     Create namespaces for the objects in a file
      --delete-output-dir                     Delete the output directory if it already exists
      --gloo-control-plane string             Name of the Gloo control plane pod
  -n, --gloo-control-plane-namespace string   Namespace of the Gloo control plane pod (default "gloo-system")
  -h, --help                                  help for convert
      --include-unknown                       Copy unknown resources to output files (if files contain resources that are not Gloo APIs)
      --input-dir string                      InputDir to read yaml/yml files recursively
      --input-file string                     Convert single file to Gateway API
      --input-snapshot string                 Gloo input snapshot file location
      --output-dir string                     Output directory to write Gateway API configurations, it must not exist before or can be deleted/recreated with --recreate-output-dir (default "./_output")
      --print-stats                           Print stats about the conversion
      --retain-input-folder-structure         When writing the output write the Gateway API configurations in the same folder structure they were read from (input-dir only)
```

### Options inherited from parent commands

```
  -c, --config string              set the path to the glooctl config file (default "<home_directory>/.gloo/glooctl-config.yaml")
      --consul-address string      address of the Consul server. Use with --use-consul (default "127.0.0.1:8500")
      --consul-allow-stale-reads   Allows reading using Consul's stale consistency mode.
      --consul-datacenter string   Datacenter to use. If not provided, the default agent datacenter is used. Use with --use-consul
      --consul-root-key string     key prefix for the Consul key-value storage. (default "gloo")
      --consul-scheme string       URI scheme for the Consul server. Use with --use-consul (default "http")
      --consul-token string        Token is used to provide a per-request ACL token which overrides the agent's default token. Use with --use-consul
  -i, --interactive                use interactive mode
      --kube-context string        kube context to use when interacting with kubernetes
      --kubeconfig string          kubeconfig to use, if not standard one
      --use-consul                 use Consul Key-Value storage as the backend for reading and writing config (VirtualServices, Upstreams, and Proxies)
```

### SEE ALSO

* [glooctl gateway-api](../glooctl_gateway-api)	 - Gateway API specific commands

