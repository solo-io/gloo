# Concepts

- [Overview](#Overview)
- [Virtual Hosts](#Virtual Hosts)
    - [Routes](#Routes)
    - [Matchers](#Matchers)
- [Upstreams](#Upstreams)
    - [Upstream Types](#Upstream Types)




<a name="Overview"></a>

### Overview

The two top-level concepts in gloo are **Virtual Hosts** and **Upstreams**.

- **Virtual Hosts** define a set of route rules that live under a domain or set of domains.
Route rules consist of a *matcher*, which specifies the kind of function calls to match (requests, events, 
and gRPC are currently supported), and the name of the destination (or destinations, for load balancing) to route them to.

- **Upstreams** define destinations for routes. Upstreams tell gloo what to route to. Upstreams may also define 
[functions](TODO) for *function-level routing*.





<a name="Virtual Hosts"></a>

### Virtual Hosts

**Virtual Hosts** define a set of route rules, an optional SNI configuration that live under a domain or set of domains.
Route rules consist of a *matcher*, which specifies the kind of function calls to match (requests, events, 
and gRPC are currently supported), and the name of the destination (or destinations, for load balancing) to route them to.




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

