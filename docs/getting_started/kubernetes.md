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
      --filename https://raw.githubusercontent.com/solo-io/gloo-install/master/kube/example-swagger.yaml
    ```

1. The discovery services should have already created an Upstream for the petstore service.
Let's verify this:
    ```bash
    kubectl get upstreams -n gloo-system
    NAME                        AGE
    default-petstore-8344       1h
    gloo-system-gloo-8081       1h
    gloo-system-ingress-8080    1h
    gloo-system-ingress-8443    1h
    ```
    
    The upstream we want to see is `default-petstore-8344`. Digging a little deeper,
    we can verify that Gloo's function discovery populated our upstream with 
    the available rest endpoints it implements. Note: the upstream was created in 
    the `gloo-system` namespace rather than `default` because it was created by a
    discovery service. Upstreams and virtualhosts do not need to live in the `gloo-system`
    namespace to be processed by Gloo.
    
1. Let's take a closer look at the functions that are available on this upstream:
    ```bash
    kubectl get upstream -n gloo-system default-petstore-8344 -o yaml
    apiVersion: gloo.solo.io/v1
    kind: Upstream
    metadata:
      annotations:
        generated_by: kubernetes-upstream-discovery
        gloo.solo.io/service-type: swagger
        gloo.solo.io/swagger_url: http://petstore.default.svc.cluster.local:8344/swagger.json
        kubectl.kubernetes.io/last-applied-configuration: |
          {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"sevice":"petstore"},"name":"petstore","namespace":"default"},"spec":{"ports":[{"port":8344,"protocol":"TCP"}],"selector":{"app":"petstore"}}}
      clusterName: ""
      creationTimestamp: 2018-03-04T19:25:07Z
      generation: 0
      name: default-petstore-8344
      namespace: gloo-system
      resourceVersion: "9843"
      selfLink: /apis/gloo.solo.io/v1/namespaces/gloo-system/upstreams/default-petstore-8344
      uid: bf664e7d-1fe1-11e8-a7b4-08002757c4e6
    spec:
      functions:
      - name: addPet
        spec:
          body: '{"id": {{pet.id}},"name": {{pet.name}},"tag": {{pet.tag}}}'
          headers:
            :method: POST
          path: /api/pets
      - name: deletePet
        spec:
          body: ""
          headers:
            :method: DELETE
          path: /api/pets/{{id}}
      - name: findPetById
        spec:
          body: ""
          headers:
            :method: GET
          path: /api/pets/{{id}}
      - name: findPets
        spec:
          body: ""
          headers:
            :method: GET
          path: /api/pets?tags={{tags}}&limit={{limit}}
      spec:
        labels: null
        service_name: petstore
        service_namespace: default
        service_port: 8344
      type: kubernetes
    status:
      state: 1
    ```
    
1. Let's now use `glooctl` to create a route for this upstream.
