# Installing on Kubernetes

### What you'll need

1. Kubernetes v1.8+ or higher deployed. We recommend using [minikube](https://kubernetes.io/docs/getting-started-guides/minikube/) to get a demo cluster up quickly.
1. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed on your local machine.


### Install

#### 1. Install Glooctl

Deploy `glooctl` binary onto your `PATH`. If you don't have the enterprise Gloo CLI, please contact Solo 
at https://www.solo.io/enterprise.


#### 2. Install Gloo to your Kubernetes Cluster using Glooctl

Once your Kubernetes cluster is up and running, run the following command to deploy Gloo and Envoy to the `gloo-system` namespace:

```bash

glooctl install kube \
   --docker-username=your-authorized-docker-username \
   --docker-password=your-authorized-docker-password \
   --docker-email==your-authorized-docker-email

```

Since you are installing Gloo Enterprise, you'll need to install Gloo using 
Docker credentials that have been authenticated to the Gloo Enterprise 
docker registry. 

Check that the Gloo pods and services have been created:

```bash
kubectl get all -n gloo-system

NAME                                 READY     STATUS    RESTARTS   AGE
pod/discovery-77467d765f-rzrvg       1/1       Running   1          34m
pod/gateway-676d756695-752xr         1/1       Running   0          34m
pod/gateway-proxy-596c4bd9f7-d5q56   1/1       Running   0          34m
pod/gloo-665d768998-s2c2c            1/1       Running   0          34m
pod/rate-limit-748d974968-zktxl      1/1       Running   1          34m
pod/redis-66db7fdf56-bc58j           1/1       Running   0          34m

NAME                    TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/gateway-proxy   LoadBalancer   10.106.186.210   <pending>     8080:30316/TCP   34m
service/gloo            ClusterIP      10.108.221.193   <none>        9977/TCP         34m
service/redis           ClusterIP      10.103.233.99    <none>        6379/TCP         34m

NAME                            DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/discovery       1         1         1            1           34m
deployment.apps/gateway         1         1         1            1           34m
deployment.apps/gateway-proxy   1         1         1            1           34m
deployment.apps/gloo            1         1         1            1           34m
deployment.apps/rate-limit      1         1         1            1           34m
deployment.apps/redis           1         1         1            1           34m

NAME                                       DESIRED   CURRENT   READY     AGE
replicaset.apps/discovery-77467d765f       1         1         1         34m
replicaset.apps/gateway-676d756695         1         1         1         34m
replicaset.apps/gateway-proxy-596c4bd9f7   1         1         1         34m
replicaset.apps/gloo-665d768998            1         1         1         34m
replicaset.apps/rate-limit-748d974968      1         1         1         34m
replicaset.apps/redis-66db7fdf56           1         1         1         34m
```

Everything should be up and running. If this process does not work, please [open an issue](https://github.com/solo-io/gloo/issues/new). We are happy to answer
questions on our [diligently staffed Slack channel](https://slack.solo.io/).

See [Getting Started on Kubernetes](../getting_started/kubernetes/1.md) to get started creating routes with Gloo.


### Uninstall 

To tear down the installation at any point, you can simply run

```bash
kubectl delete namespace gloo-system
```

<!-- end -->