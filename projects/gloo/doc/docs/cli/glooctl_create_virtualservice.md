---
title: "glooctl create virtualservice"
weight: 5
---
## glooctl create virtualservice

Create a Virtual Service

### Synopsis

A virtual service describes the set of routes to match for a set of domains. 
Virtual services are containers for routes assigned to a domain or set of domains. 
Virtual services must not have overlapping domains, as the virtual service to match a request is selected by the Host header (in HTTP1) or :authority header (in HTTP2). When using Gloo Enterprise, virtual services can be configured with rate limiting, oauth, and apikey auth.

```
glooctl create virtualservice [flags]
```

### Options

```
      --apikey-label-selector strings              apikey label selector to identify valid apikeys for this virtual service; a comma-separated list of labels (key=value)
      --apikey-secret-name string                  name to search for in provided namespace for an individual apikey secret
      --apikey-secret-namespace string             namespace to search for an individual apikey secret
      --display-name string                        descriptive name of virtual service (defaults to resource name)
      --domains strings                            comma separated list of domains
      --enable-apikey-auth                         enable apikey auth features for this virtual service
      --enable-oidc-auth                           enable oidc auth features for this virtual service
      --enable-opa-auth                            enable opa auth features for this virtual service
      --enable-rate-limiting                       enable rate limiting features for this virtual service
  -h, --help                                       help for virtualservice
      --oidc-auth-app-url string                   the public url of your app
      --oidc-auth-callback-path string             the callback path. relative to the app url. (default "/oidc-gloo-callback")
      --oidc-auth-client-id string                 client id as registered with id provider
      --oidc-auth-client-secret-name string        name of the 'client secret' secret
      --oidc-auth-client-secret-namespace string   namespace of the 'client secret' secret
      --oidc-auth-issuer-url string                the url of the issuer
      --oidc-scope strings                         scopes to request in addition to 'openid'. optional.
      --opa-module-ref strings                     namespace.name references to a config map containing OPA modules
      --opa-query string                           The OPA query to evaluate on a request
      --rate-limit-requests uint32                 requests per unit of time (default 100)
      --rate-limit-time-unit string                unit of time over which to apply the rate limit (default "MINUTE")
```

### Options inherited from parent commands

```
      --dry-run             print kubernetes-formatted yaml rather than creating or updating a resource
  -i, --interactive         use interactive mode
      --name string         name of the resource to read or write
  -n, --namespace string    namespace for reading or writing resources (default "gloo-system")
  -o, --output OutputType   output format: (yaml, json, table, kube-yaml, wide) (default table)
```

### SEE ALSO

* [glooctl create](../glooctl_create)	 - Create a Gloo resource

