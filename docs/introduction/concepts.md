# Concepts

- [Overview](#Overview)
- [Virtual Hosts](#Virtual Hosts)
    - [Routes](#Routes)
    - [Matchers](#Matchers)
    - [Destinations](#Destinations)
    - [SNI Config](#SNI Config)
- [Upstreams](#Upstreams)
    - [Upstream Types](#Upstream Types)




<a name="Overview"></a>

### Overview

The two top-level concepts in gloo are **Virtual Hosts** and **Upstreams**.

- **Virtual Hosts** define a set of route rules that live under a domain or set of domains.
Route rules consist of a *matcher*, which specifies the kind of function calls to match (requests, events, 
and gRPC are currently supported), and the name of the destination (or destinations) to route them to.

- **Upstreams** define destinations for routes. Upstreams tell gloo what to route to. Upstreams may also define 
[functions](TODO) for *function-level routing*.





<a name="Virtual Hosts"></a>

### Virtual Hosts

**Virtual Hosts** define a set of route rules, an optional SNI configuration for a given domain or set of domains.

gloo will select the appropriate virtual host (set of routes) based on the domain specified in a request's `Host` header 
(in HTTP 1.1) or `:authority` header (HTTP 2.0). 

Virtual Hosts support wildcard domains (starting with `*`).

gloo will create a `default` virtual host for the user if the user does not provide one. The `default` virtual host
matches the `*` domain, which will serve routes for any request that does not include a `Host`/`:authority` header,
or a request that requests a domain that does not match another virtual host.

The each domain specified for a virtualhost must be unique across the set of all virtual hosts provided to gloo.

For many use cases, it may be sufficient to let all routes live on a single virtual host. In thise scenario,
gloo will use the same set of route rules to for requests, regardless of their `Host` or `:authority` header.

Route rules consist of a *matcher*, which specifies the kind of function calls to match (requests, events, 
and gRPC are currently supported), and the name of the destination (or destinations, for load balancing) to route them to.



<a name="Routes"></a>

#### Routes

**Routes** are the primary building block of the virtual host. A route contains a single **matcher** and one of: a 
**single destination**, or a **list of weighted destinations**.

In short, a route is essentially a rule which tells gloo: *if* the request matches this matcher, *then* route it to this 
destination.

Because multiple matchers can match a single request, the order of routes in the virtual host matters. gloo
will select the first route which matches the request when making routing decisions. It is therefore important to place
fallback routes (e.g. matching any request for path `/` with a custom 404 page) towards the bottom of the route list.



<a name="Matchers"></a>

#### Matchers

Matchers currently have two types supported types of requests to match on:

* **Request Matchers** match on properties of HTTP requests. This includes the request path (`:path` header in HTTP 2.0),
method (`:method` in HTTP 2.0) headers (their keys and optionally their values), and query parameters.

* **Event Matchers** match properties of HTTP events, as per the [CloudEvents specification](https://github.com/cloudevents/spec/blob/master/spec.md).
*Note: the CloudEvents spec is in version 0.1 and likely to be changed in the future*. The only property **Event Matcher**
currently matches on is the *event-type* of an event (specified by the `x-event-type` request header). 




<a name="Destinations"></a>

#### Destinations

Destinations specify where to route a request once a matching route has been selected. A route can point to a single destination,
or it can split traffic for that route among a series of weighted destinations.

A destination can be either an *upstream destination* or a *function destination*.

**Upstream Destinations** are analogous to [envoy clusters](TODO). Requests routed to upstream destinations will be routed
to a server which is expected to handle the request once it has been admitted (and possibly transformed) by gloo.

**Function Destinations** allow requests to be routed directly to *functions* that live on various upstreams. A function
can be a serverless function call (e.g. Lambda, Google Cloud Function, OpenFaaS function, etc.), an API call on a service
(e.g. a REST API call, OpenAPI operation, gRPC invocation, etc.), or publishing to a message queue (e.g. NATS, AMQP, etc.).
Function-level routing is enabled in Envoy by gloo's [functional filters](TODO). gloo supports the addition of new upstream
types as well as new function types through our [plugin interface](TODO).



<a name="Upstreams"></a>

### Upstreams

**Upstreams** define destinations for routes. Upstreams tell gloo what to route to and how to route to them. gloo determines
how to handle routing for the upstream based on its `type` field. Upstreams have a `type`-specific `spec` field which must 
be used to provide routing information to gloo based on the type of upstream.

Let's walk through an example of a kubernetes upstream in order to understand how this works.

gloo reads in a configuration that looks like the following: 

```yaml
name: my-k8s-service
type: kubernetes
spec:
  service_name: my-k8s-service
  service_namespace: default
  openapi: OPENAPI_URL #TODO!@!! 
  - TODO
```

- `name` tells gloo what the identifier for this upstream will be (for routes that point to it).
- `type: kubernetes` tells gloo that the [kubernetes plugin](TODO) knows how to route to this upstream
- `spec: ...` tells the kubernetes plugin the service name and namespace, which 
- `swagger` tells the [swagger plugin](TODO) how 


The most basic upstream type is the [`service` upstream type](TODO), which simply tells gloo
of which hosts an upstream consists. More sophisticated upstream types include the [kubernetes upstream type](TODO), the 
[NATS upstream type](TODO), and the [AWS Lambda upstream type](TODO).

