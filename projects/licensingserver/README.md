# Licensing server

## Explanation
The licensing server runs in the kubernetes cluster along side gloo, or potentially any other solo-io product. 
The server itself is an executable which when run spins up a GRPC server, as well as an http server along side it.
The http server along side simply serves as an alternate healthcheck endpoint for the GRPC server itself. The GRPC server 
has only 2 services, each with a single endpoint. A healthcheck endpoint, and a validate endpoint which takes a key,
and returns whether or not it is valid.

The licensing server has 2 main parts, the server itself, and the client to access it. The server does not have an opinion 
about the licensing backend it interacts with, upon creation a `LicensingClient` must be passed in which performs the actual validation.
As of right now, instantiating the server creates a keygen client which validates the keys, but this may change in the future.

The client is located in the root of the package and provider a simple API for connecting to, and calling the server itself. 

## Running
As of right now the kube.yaml script for running gloo will create a deployment to run the licensing server. It requires 
the licensing secret to be present and will error out without it.
Temporarily this secret can be added with the following command:
```bash
kubectl create secret generic licensing --from-literal=accountid=solo-io --from-literal=auth    token=admi-0cb3d220cb8cd27d732e6c4f54ae16826704cc55022d7a15eb473c92d9660b6d26b1b0ce35b3ccd0816cb72dc3f1abe475b08    1288873c22972762f69db65a9v2
```
As I say later on the [TODO](#todo) section I would like to improve this, I added a secret.tmpl file in the hack folder for
this purpose and it will act as a template file for deploying this secret but I have yet to flesh out that functionality.

## Testing
There are unit and e2e tests contained within the licensing server package.
The unit tests can run without any setup, or credentials. However the e2e tests require an authentication method for keygen.

The credentials object is defined in `pkg/clients/keygen_impl.go KeygenAuthConfig`. 

A couple options exist to pass these credentials in.
1) create a json file and pass the path to that file in the env var `KEYGEN_CREDENTIALS_FILE`
2) pass the creds directly in through env vars `KEYGEN_ACCOUNT_ID, KEYGEN_AUTH_TOKEN`

To run the tests: 
```bash
make licensing-server-test
```

## Proto generation
In order to re-generate the proto files for the licensing server run: 
```bash
make licensing-server-generate
```


## TODO

1. Decide on a better long-term solution for creating and deploying the auth secret for keygen
2. Perform full e2e testing within a kube cluster
3. Create a CLI to make common keygen processes easier for devs
