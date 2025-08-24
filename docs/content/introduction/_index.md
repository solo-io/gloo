---
title: What is Gloo Gateway?
menuTitle: Concepts
weight: 10
---

Gloo Gateway is a feature-rich, Kubernetes-native ingress controller, and next-generation API gateway. Gloo Gateway is exceptional in its function-level routing; its support for legacy apps, microservices and serverless; its discovery capabilities; its numerous features; and its tight integration with leading open-source projects. Gloo Gateway is uniquely designed to support hybrid applications, in which multiple technologies, architectures, protocols, and clouds can coexist.

![Gloo Gateway Architecture]({{% versioned_link_path fromRoot="/img/gloo-architecture-envoys.png" %}})

* **Kubernetes ingress controller**: Gloo Gateway can function as a feature-rich ingress controller when deployed on Kubernetes and especially simplifies routing capabilities when deployed into public clouds like AWS EKS.

* **Hybrid apps**: Gloo Gateway creates applications that route to backends implemented as microservices, serverless functions, and legacy apps. This feature can help users to gradually migrate from their legacy code to microservices and serverless; can let users add new functionalities using cloud-native technologies while maintaining their legacy codebase; can be used in cases where different teams in an organization choose different architectures; and more. See [here](https://www.solo.io/blog/building-hybrid-apps-with-gloo) for more on the Hybrid App paradigm.

* **Service Mesh Ingress**: Service mesh technologies solve problems with service-to-service communications across cloud networks. Problems such as service identity, consistent L7 network telemetry gathering, service resilience, traffic routing between services, as well as policy enforcement (like quotas, rate limiting, etc) can be solved with a service mesh. For a service mesh to operate correctly, it needs a way to get traffic into the mesh. The problems with getting traffic from the edge into the cluster are a bit different from service-to-service problems. Things like edge caching, first-hop security and traffic management, Oauth and end-user authentication/authorization, per-user rate limiting, web-application firewalling, etc are all things an Ingress gateway can and should help with. Gloo Gateway solves these problems and complements any service mesh including Istio, Linkerd, Consul Connect, and AWS App Mesh.

---

## What makes Gloo Gateway unique?

* **Function-level routing allows integration of legacy applications, microservices and serverless**: Gloo Gateway can route requests directly to _functions_, which can be: a serverless function call (e.g. Lambda, Google Cloud Function, OpenFaaS function, etc.); an API call on a microservice or a legacy service (e.g. a REST API call, OpenAPI operation, XML/SOAP request etc.); or publishing to a message queue. This unique ability is what makes Gloo Gateway the only API gateway that supports hybrid apps, as well as the only one that does not tie the user to a specific paradigm.

* **Gloo Gateway incorporates vetted open-source projects to provide broad functionality**: Gloo Gateway supports high-quality features by integrating with top open-source projects, including gRPC, OpenTracing, and more. Gloo Gateway's architecture allows rapid integration of future popular open-source projects as they emerge.

* **Full automated discovery lets users move fast**: Upon launch, Gloo Gateway creates a catalog of all available destinations, and continuously maintains it up to date. This takes the responsibility for 'bookkeeping' away from the developers, and guarantees that new features become available as soon as they are ready. Gloo Gateway discovers across IaaS, PaaS and FaaS providers, as well as Swagger and gRPC.

* **Gloo Gateway integrates intimately with the user's environment**: with Gloo Gateway, users are free to choose their favorite tools for scheduling (such as K8s, Nomad, OpenShift, etc), persistence (K8s, Consul, etcd, etc) and security (K8s, Vault).

---

## Routing Features

* **Dynamic Load Balancing**: Load balance traffic across multiple upstream services.

* **Health Checks**: Active and passive monitoring of your upstream services.

* **OpenTracing**: Monitor requests using the well-supported OpenTracing standard.

* **Monitoring**: Export HTTP metrics to Prometheus or Statsd.

* **SSL**: Highly customizable options for adding SSL encryption to upstream services with full support for SNI.

* **Transformations**: Add, remove, or manipulate HTTP requests and responses.

---

## Why Gloo Gateway?

Gloo Gateway makes it easy to solve your challenges of managing ingress traffic into your application architectures (not just Kubernetes) regardless of where they run. Backend services can be discovered when running or registered in Kubernetes, AWS Lambda, VMs, EC2, Consul, et. al. Please see the [Gloo Gateway announcement](https://medium.com/solo-io/announcing-gloo-the-function-gateway-3f0860ef6600) for more on its origin. 

* **Solve difficult cloud-native and hybrid challenges**: Microservices make understanding an application's API difficult. Gloo Gateway implements the [API Gateway pattern](https://microservices.io/patterns/apigateway.html) to add shape and structure to your architecture.

* **Build on Envoy proxy the right way**: Gloo Gateway is the decoupled control plane for Envoy Proxy enabling developers and operators to dynamically update Envoy using the xDS gRPC APIs in a declarative format. Please see our blogs on [building a control plane for Envoy](https://medium.com/solo-io/guidance-for-building-a-control-plane-to-manage-envoy-proxy-at-the-edge-as-a-gateway-or-in-a-mesh-badb6c36a2af) and [control plane deployment strategies.](https://medium.com/solo-io/guidance-for-building-a-control-plane-for-envoy-part-5-deployment-tradeoffs-a6ef55c06327)

* **Stepping stone to Service Mesh**: Gloo Gateway adds service-mesh capabilities to your cluster ingress without being a service mesh itself. Gloo Gateway allows you to iteratively take small steps towards advanced features and ties in with systems like Flagger for [canary automation](https://docs.flagger.app/usage/gloo-progressive-delivery), and plugs in natively to [service-mesh implementations]({{% versioned_link_path fromRoot="/guides/integrations/service_mesh/" %}}) like Istio, Linkerd or Consul.

* **Integration of legacy applications**: Gloo Gateway can route requests directly to _functions_, an API call on a microservice or a legacy service, or publishing to a message queue. This unique ability makes Gloo Gateway the only API gateway supporting hybrid apps without tying the user to a specific paradigm.

* **Incorporate vetted open-source projects for broad functionality**: Gloo Gateway support high-quality features by integrating with top open-source projects, including gRPC, OpenTracing, and more. Gloo Gateway's architecture allows rapid integration of future popular open-source projects as they emerge.

* **Fully automated discovery lets users move fast**: Upon launch, Gloo Gateway creates a catalog of all available destinations, and continuously maintains it up to date. Gloo Gateway discovers across IaaS, PaaS and FaaS providers, as well as Swagger and gRPC.

* **Integration with existing tools**: with Gloo Gateway, users are free to choose their favorite tools for scheduling (such as K8s, Nomad, OpenShift, etc), persistence (K8s, Consul, etcd, etc) and security (K8s, Vault).

---

## Next Steps

* [Getting Started]({{% versioned_link_path fromRoot="/getting_started/" %}}): Get started with your own Gloo Gateway deployment
* [Architecture]({{% versioned_link_path fromRoot="/introduction/architecture/" %}}): Learn about the high-level architecture behind Gloo Gateway
* [Deployment Options]({{% versioned_link_path fromRoot="/introduction/architecture/deployment_options/" %}}): Learn about the various deployment options for Gloo Gateway
* [Concepts]({{% versioned_link_path fromRoot="/introduction/architecture/concepts/" %}}): Learn about the core concepts behind Gloo Gateway

