# Getting Started on Kubernetes


### What you'll need
- [`kubectl`](TODO)
- [`glooctl`](TODO)
- Kubernetes v1.8+ deployed somewhere. [Minikube](TODO) is a great way to get a cluster up quickly.


### Steps

1. Gloo and Envoy deployed and running on Kubernetes:

        kubectl apply \
          -f https://raw.githubusercontent.com/gloo/kube/master/install.yaml

 
1. Next, deploy the Pet Store app to kubernetes:

        kubectl apply \
          -f https://raw.githubusercontent.com/solo-io/gloo-install/master/kube/example-swagger.yaml

1. The discovery services should have already created an Upstream for the petstore service.
Let's verify this:

        kubectl get upstreams -n gloo-system
        
        NAME                        AGE
        default-petstore-8080       1h
        gloo-system-gloo-8081       1h
        gloo-system-ingress-8080    1h
        gloo-system-ingress-8443    1h

    The upstream we want to see is `default-petstore-8080`. Digging a little deeper,
    we can verify that Gloo's function discovery populated our upstream with 
    the available rest endpoints it implements. Note: the upstream was created in 
    the `gloo-system` namespace rather than `default` because it was created by a
    discovery service. Upstreams and virtualhosts do not need to live in the `gloo-system`
    namespace to be processed by Gloo.
    
1. Let's take a closer look at the functions that are available on this upstream (edited here to reduce verbosity):
    
        kubectl get upstream -n gloo-system default-petstore-8080 -o yaml
        
        apiVersion: gloo.solo.io/v1
        kind: Upstream
        metadata:
          annotations:
            generated_by: kubernetes-upstream-discovery
            gloo.solo.io/service-type: swagger
            gloo.solo.io/swagger_url: http://petstore.default.svc.cluster.local:8080/swagger.json
          name: default-petstore-8080
          namespace: gloo-system
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
            service_port: 8080
          type: kubernetes
        status:
          state: 1
    
1. Let's now use `glooctl` to create a route for this upstream.

        glooctl route create \
          --path-exact /api/pets \
          --upstream default-petstore-8080

    With `glooctl`, we can see that a virtual host was created with our route:

        glooctl virtualhost get -o yaml
        
        metadata:
          namespace: gloo-system
          resource_version: "3052"
        name: default
        routes:
        - request_matcher:
            path_exact: /api/pets
          single_destination:
            upstream:
              name: default-petstore-8080
        status:
          state: Accepted

1. Let's test the route `/api/pets` using `curl`:

        export GATEWAY_ADDR=$(kubectl get po -l gloo=ingress -n gloo-system -o 'jsonpath={.items[0].status.hostIP}'):$(kubectl get svc ingress -n gloo-system -o 'jsonpath={.spec.ports[?(@.name=="http")].nodePort}')
        export GATEWAY_URL=http://$GATEWAY_ADDR
            
        curl ${GATEWAY_URL}/api/pets
        
        [{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
        
        
Great! our gateawy is up and running. Let's make things a bit more sophisticated in the next section with [Function Routing](TODO).

<!-- 1. Let's create route to `addPet`.

    Notice the functions that have been discovered:
    
        functions:
          - name: addPet
            spec:
              body: '{"id": {{pet.id}},"name": {{pet.name}},"tag": {{pet.tag}}}'
              headers:
                :method: POST
              path: /api/pets
          ...
    
    This function will apply request transformation on incoming requests that route to this function.
    In order to take advantage of this transformation, we have to make sure out route specifies sources for the 
-->