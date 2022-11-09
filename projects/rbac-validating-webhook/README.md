# gloo-fed Multicluster RBAC
Gloo Federation RBAC Validation Webhook

This is a custom Kubernetes Validating Webhook: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/

For more information, see:
https://docs.solo.io/gloo-edge/latest/guides/gloo_federation/multicluster_rbac/

To install with multicluster RBAC enabled, install Gloo Federation with
the helm value `enableMultiClusterRbac=true`.

Once you see the `rbac-validating-webhook` pod up and running,
the Kubernetes Validating Webhook is registered. Without any permissions applied
by default, running:
```shell script
kubectl apply -f projects/gloo-fed/example/resources/fed-us.yaml
```
should fail.

To give ourselves permissions in a kind cluster, run:
```shell script
kubectl apply -f projects/rbac-validating-webhook/example/resources/multicluster-rbac.yaml
```

This will create the roles and role-bindings:
```
multiclusterrolebinding.multicluster.solo.io/kind-admin created
multiclusterrole.multicluster.solo.io/kind-admin created
```

Now, you can apply the federated upstream: 
```shell script
kubectl apply -f projects/gloo-fed/example/resources/fed-us.yaml
```
and it will now work.