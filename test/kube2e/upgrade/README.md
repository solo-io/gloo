# Test Overview 
- All upgrade tests run twice. 
    - current minor last patch -> current code (ex: 1.12.3 -> 1.12.4)
    - previous minor latest patch -> current code (ex 1.11.latest -> 1.12.4)
- Each upgrade test installs gloo, upgrades CRDs and then upgrades gloo
- Each test runs a data setup and validation on installation of gloo and on upgrade of gloo

# Assets package 
- Assets are YAML definitions of resources should have filename in the format `resourceName_resourceType.yaml`
- There should be one resource per YAML file 
- The assets used in the caching tests reference a sandbox caching image that we build based on an image provided for envoy testing - more can be found here [/test/kube2e/assets/caching/README.md](/test/kube2e/assets/caching/README.md)

# Running Tests
- You will need to generate the `_test` and `_output` folders in order to run these tests locally
  - running `./ci/kind/setup-kind` is an easy way to set up your environment to run the tests
- If you are debugging and cancel the test early there may be leftover pods and namespace resources, these commands can help clean things up for you
```
# remove gloo-system namespace resources and crds
kubectl delete ns gloo-system && kubectl delete crd --all

# delete old roles and bindings
kubectl delete clusterrolebinding,clusterrole -l app=gloo
kubectl delete clusterrolebinding,clusterrole -l app.kubernetes.io/instance=gloo
kubectl delete clusterrolebinding,clusterrole -l app=glooe-prometheus 
kubectl delete clusterrolebinding,clusterrole -l app.kubernetes.io/instance=gloo-ee
```
