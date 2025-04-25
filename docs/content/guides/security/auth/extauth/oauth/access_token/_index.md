---
title: Authenticate with an Access Token
weight: 30
description: Integrating Gloo Gateway and Access Tokens
---

{{% notice note %}}
{{< readfile file="static/content/enterprise_only_feature_disclaimer" markdown="true">}}
{{% /notice %}}

You may already have an OIDC compliant authentication system in place at your organization which can issue and validate access tokens. In that case, Gloo Gateway can rely on your existing system by accepting requests with an access token and validating that token against an introspection endpoint.

In this guide we will deploy ORY Hydra, a simple OpenID Connect Provider. Hydra will serve as our existing OIDC compliant authentication system. We will generate a valid access token from the Hydra deployment and have Gloo Gateway validate that token using Hyrda's introspection endpoint.

## Step 1: Deploy a sample application

1. Create the httpbin namespace.
   ```sh
   kubectl create ns httpbin
   ```

2. Deploy the httpbin app.
   ```sh
   kubectl -n httpbin apply -f https://raw.githubusercontent.com/solo-io/gloo-mesh-use-cases/main/policy-demo/httpbin.yaml
   ```

3. Verify that the httpbin app is running.
   ```sh
   kubectl -n httpbin get pods
   ```
   
   Example output:
   ```
   NAME                      READY   STATUS    RESTARTS   AGE
   httpbin-d57c95548-nz98t   3/3     Running   0          18s
   ```

## Step 2: Create a VirtualService

1. Create the VirtualService to expose the httpbin app on the gateway proxy. The gateway forwards traffic on the `/` prefix path to the httpbin app. 
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: gateway.solo.io/v1
   kind: VirtualService
   metadata:
     name: httpbin
     namespace: gloo-system
   spec:
     virtualHost:
       domains:
       - '*'
       routes:
       - matchers:
         - prefix: /
         routeAction:
           single:
             kube:
               ref:
                 name: httpbin
                 namespace: httpbin
               port: 8000
   EOF
   ```

2. Port-forward your gateway proxy. 
   ```sh
   kubectl -n gloo-system port-forward svc/gateway-proxy 8080:80
   ```

3. Send a request to the httpbin app. Verify that you get back a 200 HTTP response code. 
   ```sh
   curl -I http://localhost:8080/headers 
   ```
   
   Example output: 
   ```
   HTTP/1.1 200 OK
   access-control-allow-credentials: true
   access-control-allow-origin: *
   content-type: application/json; encoding=utf-8
   x-envoy-upstream-service-time: 1
   server: envoy
   transfer-encoding: chunked   
   ```
   
## Step 3: Secure the VirtualService

In the next step, you learn how to secure access to the httpbin app with an access token that you get back from an OpenID Connect povider. 

### Install Ory Hydra

To implement the authentication flow, an OpenID Connect provider must be available to Gloo Gateway. For demonstration purposes, you set up the [Ory Hydra](https://www.ory.sh/hydra/docs) provider in the same cluster as Gloo Gateway.

1. Add the Ory Helm repository. 
   ```bash
   helm repo add ory https://k8s.ory.sh/helm/charts
   helm repo update
   ```

2. Install Hydra. In this example setup, an in-memory database is used for Hydra. The development mode is enabled so that the HTTP protocol can be used to interact with Hydra. Note that this setup is for demonstration purposes only. In production setups, use a secure Hydra setup. 
   ```bash
   helm install \
       --set 'hydra.config.secrets.system={$(LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | base64 | head -c 32)}' \
       --set 'hydra.config.dsn=memory' \
       --set 'hydra.config.urls.self.issuer=http://public.hydra.localhost/' \
       --set 'hydra.config.urls.login=http://example-idp.localhost/login' \
       --set 'hydra.config.urls.consent=http://example-idp.localhost/consent' \
       --set 'hydra.config.urls.logout=http://example-idp.localhost/logout' \
       --set 'ingress.public.enabled=true' \
       --set 'ingress.admin.enabled=true' \
       --set 'hydra.dangerousForceHttp=true' \
       --set 'hydra.dev=true' \
       hydra-example \
       ory/hydra --version 0.53.0
   ```

3. Verify that Hydra is up and running. 
   ```sh
   kubectl get pods
   ```
   
   Example output: 
   ```
   NAME                                     READY   STATUS    RESTARTS   AGE
   hydra-example-58cd5bf699-9jgz5           1/1     Running   0          10m
   hydra-example-maester-75c985dd5b-s4b27   1/1     Running   0          10m
   petclinic-0                              1/1     Running   1          25h
   petclinic-db-0                           1/1     Running   0          25h
   ```

Hydra automatically exposes its admin endpoint on port 4445 and a public endpoint on port 4444. You use the admin endpoint to create a client ID and password and to validate the access token. The public endpoint is used to generate an access token. 

### Create the client ID and access token

1. Port-forward the Hydra admin endpoint. 
   ```bash
   kubectl port-forward deploy/hydra-example 4445 &
   portForwardPid1=$! # Store the port-forward pid so we can kill the process later
   ```
   
   Example output: 
   ```
   [1] 1417
   ~ >Forwarding from 127.0.0.1:4445 -> 4445
   ```

2. Send a request to the `/admin/clients` endpoint to create a client ID and secret. 
   ```sh
   curl -X POST http://127.0.0.1:4445/admin/clients \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d '{"client_id": "my-client", "client_secret": "secret", "grant_types": ["client_credentials"]}'
   ```
   
   Example output: 
   ```json
   {"client_id":"my-client","client_name":"","client_secret":"secret","redirect_uris":null,"grant_types":["client_credentials"],"response_types":null,"scope":"offline_access offline openid","audience":[],"owner":"","policy_uri":"","allowed_cors_origins":[],"tos_uri":"","client_uri":"","logo_uri":"","contacts":null,"client_secret_expires_at":0,"subject_type":"public","jwks":{},"token_endpoint_auth_method":"client_secret_basic","userinfo_signed_response_alg":"none","created_at":"2025-04-24T17:51:41Z","updated_at":"2025-04-24T17:51:41.286258Z","metadata":{},"registration_access_token":"ory_at_dhTxXG1O05HsnLvQ1u9XzWBksf9SNSWiiFOv0bGc1gM.rG0leFJUD6-p2Q0B--ueHTGVqYPQGlwtEgg6R7MJmAY","registration_client_uri":"http://public.hydra.localhost/oauth2/register/my-client","skip_consent":false,"skip_logout_consent":null,"authorization_code_grant_access_token_lifespan":null,"authorization_code_grant_id_token_lifespan":null,"authorization_code_grant_refresh_token_lifespan":null,"client_credentials_grant_access_token_lifespan":null,"implicit_grant_access_token_lifespan":null,"implicit_grant_id_token_lifespan":null,"jwt_bearer_grant_access_token_lifespan":null,"refresh_token_grant_id_token_lifespan":null,"refresh_token_grant_access_token_lifespan":null,"refresh_token_grant_refresh_token_lifespan":null}
   ```

3. Port-forward the public Hydra endpoint on port 4444. You use this endpoint to generate an acecss token. 
   ```bash
   kubectl port-forward deploy/hydra-example 4444 &
   portForwardPid2=$! # Store the port-forward pid so we can kill the process later
   ```
   
   Example output: 
   ```
   [2] 1431
   ~ >Forwarding from 127.0.0.1:4444 -> 4444
   ```

4. Send a request to the public endpoint to generate an access token for the client ID and secret that you previously registered. 
   ```sh
   curl -X POST http://127.0.0.1:4444/oauth2/token \
     -u my-client:secret \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=client_credentials"
   ```
   
   Example output: 
   ```
   {"access_token":"ory_at_e3LV3kgUrwV5SBD5qU6aoXMEaob1HHdxLHBrRElfg2E.HC6NtLqOW77njmV_biZLc4sqUjDYik4W0iOKalgHs14","expires_in":3599,"scope":"","token_type":"bearer"}
   ```

5. Store the access token in an environment variable. 
   ```sh
   ACCESS_TOKEN=ory_at_e3LV3kgUrwV5SBD5qU6aoXMEaob1HHdxLHBrRElfg2E.HC6NtLqOW77njmV_biZLc4sqUjDYik4W0iOKalgHs14
   ```

6. Use the introspection path on the admin endpoint to validate the access token. The same introspection path is used in Gloo Gateway later to validate your access tokens. 
   ```bash
   curl -X POST http://127.0.0.1:4445/admin/oauth2/introspect \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     -H 'Accept: application/json' \
     -d "token=$ACCESS_TOKEN"
   ```

   Example output:
   ```
   {"active":true,"client_id":"my-client","sub":"my-client","exp":1745523386,"iat":1745519786,"nbf":1745519786,"aud":[],"iss":"http://public.hydra.localhost/","token_type":"Bearer","token_use":"access_token"}
   ```
  

### Create an AuthConfig

Now that all the necessary resources are in place, add the introspection path to an `AuthConfig` to secure your VirtualService with Gloo Gateway.

1. Create a Kubernetes secret to store the client secret that was returned when you registered the client with the Hydra admin endpoint. If you followed this guide, the secret is `secret`. 
   ```sh
   kubectl create secret generic client-secret \
     --from-literal=client-secret=secret \
     --namespace=gloo-system \
     --type=extauth.solo.io/oauth
   ```

2. Create an AuthConfig that references the introspection path. This configuration instructs Gloo Gateway to use the `introspectionUrl` to validate access tokens that are provided in an `authorization` header with the request. If the token is missing or invalid, Gloo Gateway denies the request. You use the internal hostname of the Hydra administrative service as the request comes from Gloo Gateway's exauth pod, which has access to Kubernetes DNS.
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: enterprise.gloo.solo.io/v1
   kind: AuthConfig
   metadata:
     name: oidc-hydra
     namespace: gloo-system
   spec:
     configs:
     - oauth2:
         accessTokenValidation:
           introspection: 
              introspectionUrl: http://hydra-example-admin.default:4445/admin/oauth2/introspect
              clientId: my-client
              clientSecretRef: 
                name: client-secret
                namespace: gloo-system
   EOF
   ```

3. Update the VirtualService to add the AuthConfig. 
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: gateway.solo.io/v1
   kind: VirtualService
   metadata:
     name: httpbin
     namespace: gloo-system
   spec:
     virtualHost:
       domains:
       - '*'
       routes:
       - matchers:
         - prefix: /
         routeAction:
           single:
             kube:
               ref:
                 name: httpbin
                 namespace: httpbin
               port: 8000
       options:
         extauth:
           configRef:
             name: oidc-hydra
             namespace: gloo-system
   EOF
   ```

### Test the AuthConfig

Test access to the httpbin app with and without an access token. 

1. Port-forward the Gloo Gateway proxy service so that it is reachable from your machine at `localhost:8080`:
   ```sh
   kubectl -n gloo-system port-forward svc/gateway-proxy 8080:80 &
   portForwardPid3=$! # Store the port-forward pid so we can kill the process later
   ```

2. Send a request to the httpbin app without a token. Verify that the request is denied with a 403 Forbidden HTTP response code. 
   ```bash
   curl http://localhost:8080/headers -v
   ```
   
   Example output: 
   ```
   ...
   < HTTP/1.1 403 Forbidden
   ```

3. Send another request to the app. This time, use an invalid access token. Verify that the request is also denied with a 403 HTTP response code. 
   ```sh
   curl http://localhost:8080/headers \
     -H "Authorization: Bearer qwertyuio23456789" -v
   ```
   
   Example output: 
   ```
   ...
   < HTTP/1.1 403 Forbidden
   ```
   
4. Send another request to the app and include the access token that you got from Hydra. Verify that the request is successfully authenticated. Verify that you get back 200 HTTP response code. 
   ```sh
   curl http://localhost:8080/headers \
     -H "Authorization: Bearer $ACCESS_TOKEN" -v
   ```
   
   Example output: 
   ```
   < HTTP/1.1 200 OK
   ...
   {
     "headers": {
       "Accept": [
         "*/*"
       ],
       "Authorization": [
         "Bearer ory_at_e3LV3kgUrwV5SBD5qU6aoXMEaob1HHdxLHBrRElfg2E.HC6NtLqOW77njmV_biZLc4sqUjDYik4W0iOKalgHs14"
       ],
       "Host": [
         "localhost:8080"
       ],
       "User-Agent": [
         "curl/8.7.1"
       ],
       "X-Envoy-Expected-Rq-Timeout-Ms": [
         "15000"
       ],
       "X-Forwarded-Proto": [
         "http"
       ],
       "X-Request-Id": [
         "2c5c68dd-4367-4c68-ad0d-3dd99b0a9eb9"
       ]
   }
   ```

## Step 4: Extract claims from the introspection URL

You can extract the claims that are returned by the introspection URL and add them as dynamic metadata to your authorization response. Then, you can map these claims to request headers by using a transformation policy. 

1. Review the claims that were previously returned from the introspection URL and choose the claims that you want to extract.
   ```
   {"active":true,"client_id":"my-client","sub":"my-client","exp":1745523386,"iat":1745519786,"nbf":1745519786,"aud":[],"iss":"http://public.hydra.localhost/","token_type":"Bearer","token_use":"access_token"}
   ```

2. Update your AuthConfig to specify the claims that you want to add as dynamic metadata. Dynamic metadata enrich the information that is passed to the upstream service and can be used for further processing or decision making. In this example, you instruct Gloo Gateway to extract the value of the `sub` claim and to map it to the `sub` dynamic metadata key. 

   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: enterprise.gloo.solo.io/v1
   kind: AuthConfig
   metadata:
     name: oidc-hydra
     namespace: gloo-system
   spec:
     configs:
     - oauth2:
         accessTokenValidation:
           introspection:
             introspectionUrl: http://hydra-example-admin.default:4445/admin/oauth2/introspect
             clientId: my-client
             clientSecretRef: 
               name: client-secret
               namespace: gloo-system
           dynamicMetadataFromClaims:
             sub: sub
   EOF
   ```
   
   {{% notice note %}}
   When you extract the `sub` claim from the introspection URL and add it as dynamic metadata, the claim might be added as a nested claim. For example, if you follow this guide, `sub` key is added as a nested key under `config_0` as shown in the following example.
   
   ```
   {"level":"debug","ts":"2025-04-25T03:32:33Z","logger":"ext-auth.ext-auth-service",
   "msg":"dynamic metadata on response","version":"undefined",
   "x-request-id":"af1a5bb9-6356-4651-95b5-335db571a4ba","dynamic metadata":"fields:
   {key:\"config_0\"  value:{struct_value:{fields:{key:\"sub\"  value:{string_value:\"my-client\"}}
   }}}"}
   ```
   
   It is important to know the syntax of the authorization response, because you later use a transformation policy to extract the dynamic metadata key-value pair and map it to a header in the request. To view the authorization response, enable debug logging for the extauth pod and review the pod's logs.
   
   {{% /notice %}}

3. Update your VirtualService to add a transformation policy that extracts the dynamic metadata value from the `sub` key and maps it to the `x-customer` header in your request. Because the `sub` key is returned as a nested key, you must include all the keys that are part of the chain and separate them with `:`. In this example, `config_0:sub` is used as the `sub` key is nested under the `config_0` key in the authorization response.
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: gateway.solo.io/v1
   kind: VirtualService
   metadata:
     name: httpbin
     namespace: gloo-system
   spec:
     virtualHost:
       domains:
       - '*'
       routes:
       - matchers:
         - prefix: /
         routeAction:
           single:
             kube:
               ref:
                 name: httpbin
                 namespace: httpbin
               port: 8000
       options:
         extauth:
           configRef:
             name: oidc-hydra
             namespace: gloo-system
         stagedTransformations:
           regular:
             requestTransforms:
             - requestTransformation:
                 transformationTemplate:
                   headers:
                     x-customer:
                       text: '{{ dynamic_metadata("config_0:sub", "envoy.filters.http.ext_authz")}}'
   EOF
   ```
   
4. Send a request to the app and include the access token that you got from Hydra. Verify that the request is successfully authenticated with a 200 HTTP response code and that you see the `x-customer: my-client` header in your response. 
   ```sh
   curl http://localhost:8080/headers \
     -H "Authorization: Bearer $ACCESS_TOKEN" -v
   ```
   
   Example output: 
   {{< highlight yaml "hl_lines=17-18">}}
   < HTTP/1.1 200 OK
   ...
   {
     "headers": {
       "Accept": [
         "*/*"
       ],
       "Authorization": [
         "Bearer ory_at_e3LV3kgUrwV5SBD5qU6aoXMEaob1HHdxLHBrRElfg2E.HC6NtLqOW77njmV_biZLc4sqUjDYik4W0iOKalgHs14"
       ],
       "Host": [
         "localhost:8080"
       ],
       "User-Agent": [
        "curl/8.7.1"
       ],
       "X-Customer": [
         "my-client"
       ],
       "X-Envoy-Expected-Rq-Timeout-Ms": [
         "15000"
       ],
       "X-Forwarded-Proto": [
         "http"
       ],
      "X-Request-Id": [
         "6a110bb1-13ec-44e8-ad97-81cc6c68ad3c"
       ]
     }
   }
   {{< /highlight >}}


## Cleanup
You can clean up the resources created in this guide by running:

```
kill $portForwardPid1
kill $portForwardPid2
kill $portForwardPid3
kubectl delete virtualservice -n gloo-system httpbin
kubectl delete secret client-secret -n gloo-system
kubectl delete authconfig -n gloo-system oidc-hydra
helm uninstall hydra-example
```

## Summary and Next Steps

In this guide you saw how Gloo Gateway could be used with an existing OIDC system to validate access tokens and grant access to a VirtualService. You may want to also check out the authentication guides that use [Dex]({{< versioned_link_path fromRoot="/guides/security/auth/extauth/oauth/dex/" >}}) and [Google]({{< versioned_link_path fromRoot="/guides/security/auth/extauth/oauth/dex/" >}}) for more alternatives when it comes to OAuth-based authentication.