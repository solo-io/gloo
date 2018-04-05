# Quick minikube instructions:

On linux, the easiest way to start is with the kvm2 driver:
```
minikube start --vm-driver=kvm2  --extra-config=apiserver.Authorization.Mode=RBAC
kubectl create clusterrolebinding permissive-binding \
         --clusterrole=cluster-admin \
         --user=admin \
         --user=kubelet \
         --group=system:serviceaccounts
```

Then just install gloo and example:
```
kubectl apply -f install.yaml
kubectl apply -f example-gloo.yaml
```
To access:
```
export GATEWAY_ADDR=$(kubectl get po -l gloo=ingress -n gloo-system -o 'jsonpath={.items[0].status.hostIP}'):$(kubectl get svc ingress -n gloo-system -o 'jsonpath={.spec.ports[?(@.name=="http")].nodePort}')
export GATEWAY_URL=http://$GATEWAY_ADDR
```


And open in your browser:
```
xdg-open $GATEWAY_URL
```
