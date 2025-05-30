---
title: Hello World
weight: 10
description: Follow this guide for a hands-on, step-by-step tutorial about creating your first Virtual Service and routing rules in Kubernetes.
---

In this guide, you deploy a sample Pet Store app to learn about some key Gloo Gateway concepts, including *Upstream* and *Virtual Service*.

* Deploy a REST service to Kubernetes as part of the Pet Store sample app
* Enable Gloo Gateway's Discovery system to find the Pet Store service and create an Upstream custom resource (CR) for it
* Create a Virtual Service and add routes sending traffic to specific paths on the Pet Store Upstream based on incoming web requests
* Verify Gloo Gateway correctly configures Envoy to route to the Upstream
* Test the routes by submitting web requests using `curl`

{{% notice note %}}
If there are no routes configured, Envoy will not be listening on the gateway port.
{{% /notice %}}

---

## Before you begin

1. Create a [Kubernetes cluster]({{< versioned_link_path fromRoot="/installation/platform_configuration/cluster_setup/" >}}).

2. [Install Gloo Gateway]({{< versioned_link_path fromRoot="/installation/gateway/kubernetes" >}}).

3. Make sure that you have the following command line tools installed.
   
   * [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
   * `glooctl`

---

## Enable service discovery {#service-discovery}

Gloo Gateway has a special custom resource called an *Upstream* that represents a service in your Kubernetes cluster. Similar to an [Envoy cluster](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/cluster.proto), an Upstream is a resource that you can use to define what traffic goes to (the destination), as well as how the traffic goes (the policy).

Gloo Gateway can automatically create Upstream resources for Kubernetes services that it detects in your cluster. Additionally, it can discover the OpenAPI (Swagger) spec of the Pet Store app and populate the Upstream with the available REST endpoints.

For performance reasons at scale, discovery is disabled by default. To try out the feature, enable discovery mode by updating the Settings. For more options, see the [Discovery guide]({{< versioned_link_path fromRoot="/installation/advanced_configuration/fds_mode/" >}}).

```bash
kubectl patch settings -n gloo-system default --type=merge --patch '{"spec":{"discovery":{"fdsMode":"BLACKLIST","enabled":true}}}'
```

--- 

## Set up the sample Pet Store app {#petstore}

In your Kubernetes cluster, deploy the sample Pet Store app. The following video provides an overview of the steps you take.

<video controls loop>
  <source src="https://solo-docs.s3.us-east-2.amazonaws.com/gloo/videos/helloworld_deploy.mp4" type="video/mp4">
</video>

### Deploy the Pet Store Application

Let's deploy the Pet Store app on Kubernetes with a YAML configuration file on GitHub. The configuration file includes two main Kubernetes resources:

* A Deployment that creates the Pet Store container
* A Service that exposes the Pet Store API

```shell
kubectl apply -f https://raw.githubusercontent.com/solo-io/gloo/v1.13.x/example/petstore/petstore.yaml
```

Example output:
```console
deployment.extensions/petstore created
service/petstore created
```

### Verify the Pet Store Application

Now, let's verify that the pod that runs the Pet Store app is running and the service is created.

1. Check the pod status. If the pod is not yet running, run the `kubectl -n default get pods -w` command and wait until it is. Then enter `Ctrl+C` or `Cmd+C` to exit the wait loop.

   ```shell
   kubectl -n default get pods
   ```
   
   Example output:

   ```console
   NAME                READY  STATUS   RESTARTS  AGE
   petstore-####-####  1/1    Running  0         30s
   ```


2. Verify that the petstore service is created.

   ```shell
   kubectl -n default get svc petstore
   ```
   
   Example output: Note that the service does not have an external IP address. It is only accessible within the Kubernetes cluster.
   
   ```console
   NAME      TYPE       CLUSTER-IP   EXTERNAL-IP  PORT(S)   AGE
   petstore  ClusterIP  10.XX.XX.XX  <none>       8080/TCP  1m
   ```

### Verify the Upstream for the Pet Store Application

The Gloo Gateway discovery services watch for new services added to the Kubernetes cluster. When the petstore service was created, Gloo Gateway automatically created an Upstream for the petstore service. If everything deployed properly, the Upstream **STATUS** should be **Accepted**. 

<video controls loop>
  <source src="https://solo-docs.s3.us-east-2.amazonaws.com/gloo/videos/helloworld_upstreams.mp4" type="video/mp4">
</video>

Let’s verify this by using the `glooctl` command line tool:

```shell
glooctl get upstreams
```
```console
+--------------------------------+------------+----------+------------------------------+
|            UPSTREAM            |    TYPE    |  STATUS  |           DETAILS            |
+--------------------------------+------------+----------+------------------------------+
| default-kubernetes-443         | Kubernetes | Pending  | svc name:      kubernetes    |
|                                |            |          | svc namespace: default       |
|                                |            |          | port:          8443          |
|                                |            |          |                              |
| default-petstore-8080          | Kubernetes | Accepted | svc name:      petstore      |
|                                |            |          | svc namespace: default       |
|                                |            |          | port:          8080          |
|                                |            |          |                              |
| gloo-system-gateway-proxy-8080 | Kubernetes | Accepted | svc name:      gateway-proxy |
|                                |            |          | svc namespace: gloo-system   |
|                                |            |          | port:          8080          |
|                                |            |          |                              |
| gloo-system-gloo-9977          | Kubernetes | Accepted | svc name:      gloo          |
|                                |            |          | svc namespace: gloo-system   |
|                                |            |          | port:          9977          |
|                                |            |          |                              |
+--------------------------------+------------+----------+------------------------------+
```

This command lists all the Upstreams Gloo Gateway has discovered, each written to an *Upstream* CR. 

The Upstream we want to see is `default-petstore-8080`. 
    
{{% notice note %}}
The Upstream was created in the `gloo-system` namespace rather than `default` because it was created by the discovery service. Upstreams and Virtual Services do not need to live in the `gloo-system` namespace to be processed by Gloo Gateway. 
{{% /notice %}}

### Investigate the YAML of the Upstream

You can view more information about the properties of a particular Upstream by specifying the output type as `kube-yaml`.

<video controls loop>
  <source src="https://solo-docs.s3.us-east-2.amazonaws.com/gloo/videos/helloworld_upstreams_2.mp4" type="video/mp4">
</video>

Let's take a closer look at the upstream that Gloo Gateway's Discovery service created:

```shell
glooctl get upstream default-petstore-8080 --output kube-yaml
```
```yaml
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  labels:
    app: petstore
    discovered_by: kubernetesplugin
  name: default-petstore-8080
  namespace: gloo-system
spec:
  discoveryMetadata: {}
  kube:
    selector:
      app: petstore
    serviceName: petstore
    serviceNamespace: default
    servicePort: 8080
status:
  statuses:
    gloo-system:
      reportedBy: gloo
      state: 1
```

By default, the created Upstream is rather simple. It represents a specific Kubernetes service. However, the Pet Store app is a Swagger service. Gloo Gateway can discover this swagger spec, but by default Gloo Gateway's function discovery features are turned off to improve performance. To enable Function Discovery Service (FDS) on our Pet Store app, we need to label the namespace.

```shell
kubectl label namespace default  discovery.solo.io/function_discovery=enabled
```

Now, FDS discovers the Swagger spec and populates the Upstream with the available REST endpoints that the Pet Store app implements.

```shell
glooctl get upstream default-petstore-8080
```

```console
+-----------------------+------------+----------+-------------------------+
|       UPSTREAM        |    TYPE    |  STATUS  |         DETAILS         |
+-----------------------+------------+----------+-------------------------+
| default-petstore-8080 | Kubernetes | Accepted | svc name:      petstore |
|                       |            |          | svc namespace: default  |
|                       |            |          | port:          8080     |
|                       |            |          | REST service:           |
|                       |            |          | functions:              |
|                       |            |          | - addPet                |
|                       |            |          | - deletePet             |
|                       |            |          | - findPetById           |
|                       |            |          | - findPets              |
|                       |            |          |                         |
+-----------------------+------------+----------+-------------------------+
```

```shell script
glooctl get upstream default-petstore-8080 --output kube-yaml
```

```yaml
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  labels:
    discovered_by: kubernetesplugin
    service: petstore
  name: default-petstore-8080
  namespace: gloo-system
spec:
  discoveryMetadata: {}
  kube:
    selector:
      app: petstore
    serviceName: petstore
    serviceNamespace: default
    servicePort: 8080
    serviceSpec:
      rest:
        swaggerInfo:
          url: http://petstore.default.svc.cluster.local:8080/swagger.json
        transformations:
          addPet:
            body:
              text: '{"id": {{ default(id, "") }},"name": "{{ default(name, "")}}","tag":
                "{{ default(tag, "")}}"}'
            headers:
              :method:
                text: POST
              :path:
                text: /api/pets
              content-type:
                text: application/json
          deletePet:
            headers:
              :method:
                text: DELETE
              :path:
                text: /api/pets/{{ default(id, "") }}
              content-type:
                text: application/json
          findPetById:
            body: {}
            headers:
              :method:
                text: GET
              :path:
                text: /api/pets/{{ default(id, "") }}
              content-length:
                text: "0"
              content-type: {}
              transfer-encoding: {}
          findPets:
            body: {}
            headers:
              :method:
                text: GET
              :path:
                text: /api/pets?tags={{default(tags, "")}}&limit={{default(limit,
                  "")}}
              content-length:
                text: "0"
              content-type: {}
              transfer-encoding: {}
status:
  statuses:
    gloo-system:
      reportedBy: gloo
      state: 1
```

The application endpoints were discovered by Gloo Gateway's Function Discovery (fds) service. This was possible because the petstore application implements OpenAPI (specifically, discovering a Swagger JSON document at `petstore/swagger.json`).

---

## Configuring Routing

We have confirmed that the Pet Store application was deployed successfully and that the Function Discovery service on Gloo Gateway automatically added an Upstream entry with all the published application endpoints of the Pet Store application. Now let's configure some routing rules on the default Virtual Service and test them to ensure we get a valid response.

<video controls loop>
  <source src="https://solo-docs.s3.us-east-2.amazonaws.com/gloo/videos/helloworld_virtualservice.mp4" type="video/mp4">
</video>

### Add a Routing Rule

Even though the Upstream has been created, Gloo Gateway will not route traffic to it until we add some routing rules on a Virtual Service. Let’s now use glooctl to create a basic route for this Upstream with the `--prefix-rewrite` flag to rewrite the path on incoming requests to match the path our petstore application expects.

```shell
glooctl add route \
  --path-exact /all-pets \
  --dest-name default-petstore-8080 \
  --prefix-rewrite /api/pets
```

{{% notice note %}}
If using Git Bash on Windows, the above will not work; Git Bash interprets the route parameters as Unix file paths and mangles them. Adding `MSYS_NO_PATHCONV=1` to the start of the above command should allow it to execute correctly.
{{% /notice %}}

We do not specify a specific Virtual Service, so the route is added to the `default` Virtual Service. If a `default` Virtual Service does not exist, `glooctl` will create one.

```console
+-----------------+--------------+---------+------+---------+-----------------+---------------------------+
| VIRTUAL SERVICE | DISPLAY NAME | DOMAINS | SSL  | STATUS  | LISTENERPLUGINS |          ROUTES           |
+-----------------+--------------+---------+------+---------+-----------------+---------------------------+
| default         |              | *       | none | Pending |                 | /all-pets -> gloo-system. |
|                 |              |         |      |         |                 | .default-petstore-8080    |
+-----------------+--------------+---------+------+---------+-----------------+---------------------------+
```

The initial **STATUS** of the petstore Virtual Service will be **Pending**. After a few seconds it should change to **Accepted**. Let’s verify that by retrieving the `default` Virtual Service with `glooctl`.

```shell
glooctl get virtualservice default
```

```console
+-----------------+--------------+---------+------+----------+-----------------+---------------------------+
| VIRTUAL SERVICE | DISPLAY NAME | DOMAINS | SSL  | STATUS   | LISTENERPLUGINS |          ROUTES           |
+-----------------+--------------+---------+------+----------+-----------------+---------------------------+
| default         |              | *       | none | Accepted |                 | /all-pets -> gloo-system. |
|                 |              |         |      |          |                 | .default-petstore-8080    |
+-----------------+--------------+---------+------+----------+-----------------+---------------------------+
```

### Verify Virtual Service Creation

Let's verify that a Virtual Service was created with that route. 

Routes are associated with Virtual Services in Gloo Gateway. When we created the route in the previous step, we didn't provide a Virtual Service, so Gloo Gateway created a Virtual Service called `default` and added the route. 

With `glooctl`, we can see that the `default` Virtual Service was created with our route:

```shell
glooctl get virtualservice default --output kube-yaml
```

```yaml
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: default
  namespace: gloo-system
  ownerReferences: []
status:
  statuses:
    gloo-system:
      reportedBy: gloo
      state: Accepted
      subresourceStatuses:
        '*v1.Proxy.gateway-proxy_gloo-system':
          reportedBy: gloo
          state: Accepted
spec:
  virtualHost:
    domains:
    - '*'
    routes:
    - matchers:
      - exact: /all-pets
      options:
        prefixRewrite: /api/pets
      routeAction:
        single:
          upstream:
            name: default-petstore-8080
            namespace: gloo-system
  ```
    
When a Virtual Service is created, Gloo Gateway immediately updates the proxy configuration. Since the status of this Virtual Service is `Accepted`, we know this route is now active. 

At this point we have a Virtual Service with a routing rule sending traffic on the path `/all-pets` to the Upstream `petstore` at a path of `/api/pets`.

### Test the Route Rule

Let’s test the route rule by retrieving the URL of Gloo Gateway, and sending a web request to the `/all-pets` path of the URL using curl. 

```shell
curl $(glooctl proxy url --name gateway-proxy)/all-pets
```

```json
[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
```

{{% notice tip %}}
If you test locally by using minikube, the load balancer service that exposes the gateway proxy is not assigned an external IP address or hostname and remains in a `<pending>` state. Because of that, the `glooctl proxy url` command returns an error similar to `Error: load balancer ingress not found on service gateway-proxy curl: (3) URL using bad/illegal format or missing URL`. To open a connection to the gateway proxy service, run `minikube tunnel`. 
{{% /notice %}}


The proxy has now been configured to route requests to the `/api/pets` REST endpoint on the Pet Store application in Kubernetes.

---

## Next Steps

Congratulations! You've successfully set up your first routing rule. That's just the tip of the iceberg though. In the next sections, we'll take a closer look at more HTTP routing capabilities, including [customizing the matching rules]({{% versioned_link_path fromRoot="/guides/traffic_management/destination_selection/" %}}), [route destination types]({{% versioned_link_path fromRoot="/guides/traffic_management/destination_types/" %}}), and [request processing features]({{% versioned_link_path fromRoot="/guides/traffic_management/request_processing/" %}}).

To learn more about the concepts behind Upstreams and Virtual Services check out the [Concepts]({{% versioned_link_path fromRoot="/introduction/architecture/concepts/" %}}) page.

If you're ready to dive deeper into routing, the next logical step is trying out different matching rules starting with [Path Matching]({{% versioned_link_path fromRoot="/guides/traffic_management/destination_selection/path_matching/" %}}).

