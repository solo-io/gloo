# Architecture

- [Overview](#Overview)
- [Component Architecture](#Component Architecture)
- [Discovery Architecture](#Discovery Architecture)


<a name="Overview"></a>

### Overview

Gloo aggregates back end services and provides function-to-function translation for clients, allowing decoupling from back end APIs

![Overview](high_level_architecture.png "High Level Architecture")

Clients issue requests or [emit events](TODO link to our events sdk) to routes defined on Gloo. These routes are mapped
to functions on upstream services by Gloo's configuration (provided by clients of Gloo's API). 

Gloo performs the necessary transformation between the routes defined by clients and the back end functions. Gloo is able 
to support various upstream functions through its extendable [function plugin interface](TODO).

Gloo offers first-class API management features on all functions:

- Timeouts
- Metrics & Tracing
- Health Checks
- Retries
- Advanecd load balancing
- TLS Termination with SNI Support
- HTTP Header modification
<!-- TODO: -Authentication -->
<!-- TODO: -JWT/Oauth2 -->




<a name="Component Architecture"></a>

### Component Architecture

In the most basic sense, Gloo is a translation engine and [envoy xDS server](TODO) providing advanced configuration for Envoy (including [Gloo's
custom Envoy filters](TODO)). Gloo follows an event-based architecture, watching various sources of configuration for
updates and responding immediately with v2 gRPC updates to Envoy. 

core

secret watcher <- secret sources     <- users

config watcher <- config sources

address discovery <- service registries

|
V

translator ( plugins ) --> reporter --> users   
   |
   V
   xds
      \
        \
         V
clients ->  Envoy ->  services






<a name="Discovery Architecture"></a>

### Discovery Architecture

Gloo is supported by it is a suite of optional [discovery services](TODO) that automatically discover and configure 
gloo with upstreams and functions to simplify routing for users and self-service.  



discovery sources 
                          discovery services ->      config sources - > gloo