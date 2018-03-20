
<h1 align="center">
    <img src="docs/Gloo-01.png" alt="Gloo" width="200" height="242">
  <br>
</h1>


<h3 align="center">The Function Gateway</h3>
<BR>

Gloo is a function gateway built on top of the [Envoy Proxy](https://www.Envoyproxy.io). Gloo provides a unified entry point for access to all services and serverless functions, translating from any interface spoken by a client to any interface spoken by a backend. Gloo aggregates REST APIs and events calls from clients, "glueing" together services in-cluster, out of cluster, across clusters, along with any provider of serverless functions.

What makes Gloo special is its use of function-level routing, which is made possible by the fact that Gloo intimately knows the APIs of the upstreams it routes to. This means that the client and server do not have to speak the same protocol, the same version, or the same language. Users can configure Gloo (or enable automatic discovery services) to make Gloo aware of functional back-ends (such as AWS Lambda, Google Functions, or RESTful services) and enable function-level routing. Gloo features an entirely pluggable architecture, providing the ability to extend its configuration language with plugins which add new types of upstreams and route features.

It is entirely possible to run Gloo as a traditional API gateway, without leveraging function-level capabilities. Gloo can be configured as a fully-featured API gateway, simply by using upstreams that don't support functions.
 
This repository contains the core controller of the Gloo system. The core includes the Envoy xDS server and the translator
that manages and calls [Gloo plugins](https://gloo.solo.io/introduction/architecture/).

For an in-depth breakdown of of Gloo and its features, please see our [documentation](https://gloo.solo.io).

## Getting Started

Getting started with Gloo on Kubernetes is as easy as running

```bash
kubectl apply \
  -f https://raw.githubusercontent.com/solo-io/gloo-install/master/kube/install.yaml
```

which will create the `gloo-system` namespace and deploy Envoy, Gloo, and Gloo's discovery services. To create your first 
routes with Gloo, [see the getting started page in our documentation](https://gloo.solo.io/getting_started/kubernetes/1/).

## Documentation

* [Official Documentation](https://gloo.solo.io)
* [Building hybrid app demo](https://www.youtube.com/watch?time_continue=1&v=ISR3G0CAZM0)
* [Announcement Blog](https://medium.com/solo-io/announcing-gloo-the-function-gateway-3f0860ef6600)


## Repository Guide

| Repo                                                                                  | What it does?                                                                            |
|---------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|
| [gloo](https://github.com/solo-io/gloo)                                               | The gloo control plane. Implements the ADS API for envoy                                 |
| [gloo-install](https://github.com/solo-io/gloo-install)                               | Install manifests and Helm chart.                                                                                                                                        |
| [thetool](https://github.com/solo-io/thetool)                                         | Easily build gloo+envoy with plugins enabled or disabled.                                |
| [glooctl](https://github.com/solo-io/glooctl)                                         | Command line client for gloo, for easy config manipulation.                              |
| [gloo-api](https://github.com/solo-io/gloo-api)                                       | Proto API definitions (upstreams, virtualhosts, routes...).                              |
| [gloo-function-discovery](https://github.com/solo-io/gloo-function-discovery)         | Auto discovery for functions in upstreams (i.e. lambda functions, swagger functions).    |
| [gloo-storage](https://github.com/solo-io/gloo-storage)                               | Abstracts configuration storage and change watch. kube and file are currently supported. |
| [gloo-testing](https://github.com/solo-io/gloo-testing)                               | e2e testing with kubernetes.                                                               |
| [gloo-plugins](https://github.com/solo-io/gloo-plugins)                               | Plugins that can be enabled and built into gloo using `thetool`.                         |
| [gloo-k8s-service-discovery](https://github.com/solo-io/gloo-k8s-service-discovery)   | Auto register kubernetes services as gloo upstreams.                                     |
| [gloo-ingress-controller](https://github.com/solo-io/gloo-ingress-controller)         | Kube ingress controller that generates gloo upstreams.                                   |
| [envoy-common](https://github.com/solo-io/envoy-common)                               | Common libraries that enabled functional envoy filters.                                  |
| [envoy-lambda](https://github.com/solo-io/envoy-lambda)                               | AWS lambda support for envoy.                                                            |
| [envoy-transformation](https://github.com/solo-io/envoy-transformation)               | Request and response transformation for envoy.                                           |
| [envoy-google-function](https://github.com/solo-io/envoy-google-function)             | Google Cloud Functions support for envoy.                                                |

Community
-----
Join us on our slack channel: [https://slack.solo.io/](https://slack.solo.io/)

---

### Thanks

**Gloo** would not be possible without the valuable open-source work of projects in the community. We would like to extend a special thank-you to [envoy](https://www.envoyproxy.io).
