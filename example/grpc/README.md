gRPC Demo
==========
This is a brief guide to transcode JSON/HTTP requests to gRPC using Gloo

# Prerequisites
In this demo, we will use the following command line tools:
- `minikube` to create a kubernetes test environment.
- `kubectl` to interact with kubernetes.
- `glooctl` to interact with gloo.

# Setup the environment

# And now for a brief demo of Gloo providing a JSON-to-gRPC Bridge

## Install Kubernetes

```shell
minikube start
```

## Install Gloo
```shell
kubectl apply \
          -f https://raw.githubusercontent.com/solo-io/gloo/master/install/kube/install.yaml
```

Wait \ Verify that all the pods are in Running status:
```
kubectl get pods -n gloo-system
```

## Deploy the GRPC Bookstore

```
kubectl apply \
    -f https://raw.githubusercontent.com/solo-io/gloo/master/example/grpc/deploy.yaml
```

## Get the URL of the gRPC service 

```
GRPC_ADDR=$(kubectl get po -l app=grpc-bookstore -n default -o 'jsonpath={.items[0].status.hostIP}'):$(kubectl get svc grpc-bookstore -n default -o 'jsonpath={.spec.ports[0].nodePort}')
export GRPC_URL=http://$GRPC_ADDR
```


## Try the gRPC Service directly with cURL

```
curl $GRPC_URL
curl $GRPC_URL/bookstore.Bookstore/ListShelves
```

# Use Gloo to translate cURL to gRPC

## Create a route for a function

```
glooctl route create \
    --path-exact /shelves \
    --http-method GET \
    --upstream default-grpc-bookstore-8080 \
    --function ListShelves 
```

## Get the url of the ingress
If you installed kubernetes using minikube as mentioned above, you can use this command:
```shell
GATEWAY_ADDR=$(kubectl get po -l gloo=ingress -n gloo-system -o 'jsonpath={.items[0].status.hostIP}'):$(kubectl get svc ingress -n gloo-system -o 'jsonpath={.spec.ports[?(@.name=="http")].nodePort}')
export GATEWAY_URL=http://$GATEWAY_ADDR
```

## Try with curl again

```
curl $GATEWAY_URL/shelves
```

## Add a route to create shelves
```
glooctl route create \
    --path-exact /shelves \
    --http-method POST \
    --upstream default-grpc-bookstore-8080 \
    --function CreateShelf 
```

## Create a Shelf using JSON

```
curl $GATEWAY_URL/shelves \
    -d '{"shelf": {"id": 1, "theme": "cloud computing"}}'
```

## See that it was created
```
curl $GATEWAY_URL/shelves
```
