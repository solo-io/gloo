# Kubernetes Service Plugin for Gloo


#### Description

The [Kubernetes Service Plugin for Gloo](https://github.com/solo-io/gloo-plugins/tree/master/kubernetes) leverages Kubernetes
Services to discover endpoints for Envoy to route to. The **Kubernetes Upstream Type** supports routing to pods by label.

#### Upstream Spec Configuration

The **Upstream Type** for Kubernetes upstreams is `kubernetes`. 

The [upstream spec](../v1/upstream.md#v1.Upstream) for Kubernetes Upstreams has the following structure:

```yaml
service_name: string
service_namespace: string
service_port: int
labels: map<string, string>
```

| Field | Type |  Description |
| ----- | ---- |  ----------- |
| service_name | string |  The name of the [kubernetes service](https://kubernetes.io/docs/concepts/services-networking/service/) this upstream routes to. **required** |
| service_namespace | string |  The namespace that contains the [kubernetes service](https://kubernetes.io/docs/concepts/services-networking/service/) this upstream routes to. **required** |
| service_port | int | The service port the upstream should route to. **required if the service has more than one port defined** *Note: create an upstream for each service port to which you want to route* 
| labels | map<string, string> | If specified, this upstream will route only to pods [selected for these labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/). **optional** *Note: create an upstream for each set of labels to which you want to route*   



#### Function Spec Configuration
The Kubernetes Plugin does not support functions. However, other function plugins (such as the [transformation plugin](request_transformation.md)) 
can process functions on kubernetes upstreams.  


#### Example Kubernetes Upstream

The following is an example of a valid Kubernetes [Upstream](../introduction/concepts.md#Upstreams):

```yaml
name: petstore
spec:
  service_name: "petstore"
  service_namespace: "default"
type: kubernetes
```

With labels:

```yaml
name: petstore-v1
spec:
  service_name: "petstore"
  service_namespace: "default"
  labels:
    version: v1
type: kubernetes
```


#### Discovery

The Gloo Kubernetes Service Discovery Service<!--(TODO)--> will automatically discover upstreams from Kubernetes Services if it is running.
Upstreams discovered in this way will not contain labels.
