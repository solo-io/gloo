# Concepts

- [Overview](#Overview)
- [Virtual Hosts](#Virtual Hosts)
    - [Routes](#Routes)
    - [Matchers](#Matchers)
- [Upstreams](#Upstreams)
    - [Upstream Types](#Upstream Types)




<a name="Overview"/>

### Overview

The two top-level concepts in gloo are **Virtual Hosts** and **Upstreams**.

- **Virtual Hosts** define a set of route rules, an optional SNI configuration that live under a domain or set of domains.
Route rules consist of a *matcher*, which specifies the kind of function calls to match (requests, events, 
and gRPC are currently supported), and the name of the destination (or destinations, for load balancing) to route them to.

- **Upstreams** define destinations for routes. 


<a name="Virtual Hosts"/>

### Virtual Hosts

