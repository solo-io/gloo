# glooctl gateway-api convert 
A CLI tool to migrate from Gloo Edge APIs to the Kubernetes Gateway API.

The CLI accepts a single Kubernetes YAML file or a Gloo Gateway input snapshot. It can also scan an entire directory to find Gloo Gateway YAML files that use the Gloo Edge API. 

## Requirements
* Gloo Gateway and VirtualService objects must be provided and correctly associated for conversion. The tool matches the correct VirtualService with the Gateway based on its selectors.
* This tool defines all Listeners using `xListenerSet` which is a beta feature in Gateway API. See [ListenerSets](https://gateway-api.sigs.k8s.io/geps/gep-1713/)
* To apply the generated CRDs you must apply the latest experimental schema definition from here. [CustomResourceDefinition](https://github.com/kubernetes-sigs/gateway-api/blob/main/config/crd/experimental/gateway.networking.k8s.io_gateways.yaml)
* This must be used with Gloo Gateway version 1.19 or greater
* The generated output must be written to a new empty directory, Use `--delete-output-dir` to have the tool delete it before starting

## Use

* Read a single input yaml file and generate Gateway API Output

```shell
glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output
```

* Scan a nested directory for `.yaml` and `.yml` files and convert them to the Gateway API. 

```shell
glooctl gateway-api convert --input-dir ./gloo-configs --output-dir ./_output --retain-input-folder-structure
```

* Generate Gateway API YAML files from a Gloo Gateway input snapshot.

```shell
kubectl -n gloo-system port-forward deploy/gloo 9091
curl localhost:9091/snapshots/input > gg-input.json

glooctl gateway-api convert --input-snapshot gg-input.json --output-dir ./_output
```

## Output Formats

- **Files by namespace**: By default, a separate file is created for each generated Gateway API resource. All resources are placed into namespace-specific directories. 
- **Retain input folder structure**: When you convert files in a givenCI/CD pipeline folder structure, you might want to retain the generated configuration in the files they were converted from. To do this, add the `--retain-input-folder-structure` option.