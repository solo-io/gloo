# Bootstrap Options

When deployed using [Helm Charts](TODO), [TheTool](TODO), or a plain [Kubernetes YAML File](TODO),
Gloo's bootstrap options will be set automatically.

Here are the bootstrap options (set via CLI flags) available for configuring Gloo:

| flag         | purpose                                                                             | possible values | notes                                                                                                                                                 |   |
|--------------|-------------------------------------------------------------------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|---|
| `storage.type` | indicates the type of storage backend Gloo should monitor for configuration objects | "kube", "file"  | if using "kube", Gloo must be either run in-cluster, or provided a valid kubeconfig and master url. "file" requires the `--file.config.dir` to be set |   |
| `storage.refreshrate` | the polling interval when monitoring for new / updated config objects. if using kubernetes, this value will instead set the ressyncperoid for the config object controller. read more about kubernetes controllers [here](http://borismattijssen.github.io/articles/kubernetes-informers-controllers-reflectors-stores) | a valid duration (e.g. 5s, 10m) |                                           |   |
| `secrets.type` | indicates the type of secret storage backend Gloo should monitor for secrets | "kube", "vault", "file" | if using "kube", Gloo must be either run in-cluster, or provided a valid kubeconfig and master url.  "file" requires `--file.secret.dir` to be set "vault" requires `--vault.addr` and `--vault.token` to be set |   |
| `secrets.refreshrate` | the polling interval when monitoring for new / updated secrets. if using kubernetes, this value will instead set the ressyncperoid for the secrets controller. read more about kubernetes controllers [here](http://borismattijssen.github.io/articles/kubernetes-informers-controllers-reflectors-stores)              | a valid duration (e.g. 5s, 10m) |                                           |   |
| `kube.namespace` | set the kubernetes namespace to watch for config objects. if left empty, this will default to `gloo-system`.   | any valid namespace                                  | required if using `--storage.type=kube`                          |   |
| `kubeconfig`     | path to kubeconfig to use if using kubernetes features (secrets, storage, or the [kubernetes plugin](TODO)).   | path to kubeconfig. defaults to ${HOME}/.kube/config | required if using kubernetes features and running out-of-cluster |   |
| `master`         | url of a kubernetes master, if using kubernetes features (secrets, storage, or the [kubernetes plugin](TODO)). | a valid url for a kubernetes master                  | required if using kubernetes features and running out-of-cluster |   |
| `vault.addr`          | address of a vault server to monitor for secret storage                                                                                                                                                                                                                                                                 | a valid url                     | required if using vault as a secret store |   |
| `vault.token`   | the token to use for authenticating to vault. Gloo doesn't currently support Vault authentication methods other than token auth | a valid auth token  | required if using "vault" secret type                                                                                            |   |
| `vault.retries` | the number of times the vault poller should retry API requests to vault                                                         | uint > 0            | default to 3                                                                                                                     |   |
| `xds.port`      | port on which to serve Envoy v2 gRPC API requests                                                                               | a valid port number | defaults to 8081. if you edit this option, be sure to change the [bootstrap config for Envoy](TODO) to point at the new xDS port |   |


