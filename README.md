
<h1 align="center">
    <img src="docs/Gloo-01.png" alt="Gloo" width="200" height="242">
  <br>
</h1>


<h3 align="center">The Function Gateway</h3>
<BR>

Gloo is a function gateway built on top of the [Envoy Proxy](https://www.Envoyproxy.io). Gloo provides a unified entry point
for access to all services and serverless functions, translating from any interface spoken by a client to any interface
spoken by a backend. Gloo aggregates REST APIs and events calls from clients, "glueing" together services in-cluster, 
out of cluster, across clusters, along with any provider of serverless functions.

This Repo
-----
This repository contains the components that compose the core Gloo storage watcher, Envoy xDS server, and config translator.
For a better understanding of Gloo and its features, please see our [documentation](https://gloo.solo.io).

Blog
-----

To learn more about the motivation behind creating Gloo read our [blog](https://medium.com/solo-io/announcing-gloo-the-function-gateway-3f0860ef6600)

Documentation
-----

Get started by reading our docs here: [https://gloo.solo.io/](https://gloo.solo.io/)

Quick Repository Guide:
-----
| Repo                                                                                  | What it does?                                                                            |
|---------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|
| [gloo](https://github.com/solo-io/gloo)                                               | The gloo control plane. Implements the ADS API for envoy                                 |
| [gloo-install](https://github.com/solo-io/-install)                                   | Install manifests.                                                                       |
| [gloo-chart](https://github.com/solo-io/gloo-chart)                                   | Helm charts for gloo.                                                                    |
| [thetool](https://github.com/solo-io/thetool)                                         | Easily build gloo+envoy with plugins enabled or disabled.                                |
| [glooctl](https://github.com/solo-io/glooctl)                                         | Command line client for gloo, for easy config manipulation.                              |
| [gloo-api](https://github.com/solo-io/gloo-api)                                       | Proto API definitions (upstreams, virtualhosts, routes...).                              |
| [gloo-function-discovery](https://github.com/solo-io/gloo-function-discovery)         | Auto discovery for functions in upstreams (i.e. lambda functions, swagger functions).    |
| [gloo-storage](https://github.com/solo-io/gloo-storage)                               | Abstracts configuration storage and change watch. kube and file are currently supported. |
| [gloo-testing](https://github.com/solo-io/gloo-testing)                               | e2e testing with minikube.                                                               |
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
