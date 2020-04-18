## GRPC web

### Example usage

[example react app](https://github.com/improbable-eng/grpc-web/tree/master/client/grpc-web-react-example)

## Development against the backend

The UI can talk to the backend through...

### Running against an API server locally

### Running against an API server in Kubernetes

If the API server and proxy are already deployed to Kubernetes, then you can use Kubernetes port-forwarding
to be able to talk to them.

```
...
```

### Accessing private npm packages
You might run into this error:
> error An unexpected error occurred: "Failed to replace env in config: ${NPM_TOKEN}".

Since we are importing the private `@solo-io/dev-portal-grpc` module, you need to authenticate with npm.
First, let's get an npm token (you need to be a member of the solo-io npm org for this to work):

```
# login to npm
npm login

# request a read-only token
npm token create --read-only

```

Then, you have to add the resulting token to your environment as explained 
[here](https://docs.npmjs.com/using-private-packages-in-a-ci-cd-workflow):