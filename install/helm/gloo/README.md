# gloo-chart

Helm chart for Gloo

This Helm chart can be used to manually deploy Gloo to Kubernetes.
It is also used by `thetool` to deploy Gloo. When deployed using
`thetool`, it generates `gloo-chart.yaml` which replaces the 
`values.yaml`

To use monitoring and open tracing, please ensure you have at least
4GB RAM assigned to the minikube VM. Using Prometheus needs minikube
to be started in RBAC mode.

    minikube start --extra-config=apiserver.Authorization.Mode=RBAC --memory 4096
    kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default


