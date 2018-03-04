# Getting started with Kubernetes

- [Simple Installation](#Simple Installation)
- [Advanced  Installation](#Advanced Installation)



<a name="Simple Installation"></a>
### Simple Installation

#### What you'll need

1. Kubernetes v1.7 or higher deployed. We recommend using [minikube](TODO) to get a demo cluster up quickly.
1. [`kubectl`](TODO) installed on your local machine.
1. [`glooctl`](TODO) installed on your local machine.

Once your Kubernetes cluster is up and running, run the following command to deploy Gloo and Envoy to the `gloo-system` namespace:

```bash
kubectl apply \
  --filename https://raw.githubusercontent.com/gloo-install/kube/master/install.yaml
```

Check that the Gloo pods and services have been created: 

```bash
kubectl get all -n gloo-system

NAME                           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/envoy                   1         1         1            0           4s
deploy/function-discovery      1         1         1            1           4s
deploy/gloo                    1         1         1            0           4s
deploy/k8s-service-discovery   1         1         1            1           4s

NAME                                 DESIRED   CURRENT   READY     AGE
rs/envoy-687ff7867d                  1         1         0         4s
rs/function-discovery-5876db67df     1         1         1         4s
rs/gloo-6f68b9f7d6                   1         1         0         4s
rs/k8s-service-discovery-c548ccd57   1         1         1         4s

NAME                           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/envoy                   1         1         1            0           4s
deploy/function-discovery      1         1         1            1           4s
deploy/gloo                    1         1         1            0           4s
deploy/k8s-service-discovery   1         1         1            1           4s

NAME                                 DESIRED   CURRENT   READY     AGE
rs/envoy-687ff7867d                  1         1         0         4s
rs/function-discovery-5876db67df     1         1         1         4s
rs/gloo-6f68b9f7d6                   1         1         0         4s
rs/k8s-service-discovery-c548ccd57   1         1         1         4s

NAME                                       READY     STATUS              RESTARTS   AGE
po/envoy-687ff7867d-8lc8f                  0/1       ContainerCreating   0          5s
po/function-discovery-5876db67df-nxc45     1/1       Running             0          5s
po/gloo-6f68b9f7d6-fm47k                   0/1       ContainerCreating   0          5s
po/k8s-service-discovery-c548ccd57-snm88   1/1       Running             0          5s

NAME        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
svc/envoy   NodePort    10.101.33.224   <none>        8080:31014/TCP,8443:31391/TCP   5s
svc/gloo    ClusterIP   10.104.86.18    <none>        8081/TCP                        5s

```

See [Getting Started on Kubernetes](../getting_started/kubernetes.md) to get started creating routes with Gloo