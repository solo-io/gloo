# Getting started with Kubernetes


1. What you'll need:
- [`kubectl`](TODO)
- [`glooctl`](TODO)
- Kubernetes v1.8+ deployed somewhere. [Minikube](TODO) is fine for running gloo.

1. Gloo and Envoy deployed and running on Kubernetes:
```bash
kubectl apply \
  --filename https://raw.githubusercontent.com/gloo/kube/master/install.yaml
```

 
1. Next, deploy the Pet Store app to kubernetes:
```bash
kubectl apply \
  --filename https://raw.githubusercontent.com/gloo/kube/master/install.yaml
```
