
<h1 align="center">
    <img src="docs/Gloo-01.png" alt="Gloo" width="200" height="242">
  <br>
</h1>


<h3 align="center">The Function Gateway</h3>
<BR>

Gloo is an ingress controller and function gateway built on the [Envoy Proxy](https://www.Envoyproxy.io), written in Go.
 
Gloo expands upon the model of traditional API Gateways as (currently) the only gateway to support 
[routing on the function level](https://gloo.solo.io/introduction/introduction/). This allows users to compose 
gateway APIs from more granular components, i.e. the functions that compose services, rather than the services themselves.

This repository contains the core controller of the Gloo system. The core includes the Envoy xDS server and the translator
that manages and calls [Gloo plugins](https://gloo.solo.io/introduction/architecture/).

For an in-depth breakdown of of Gloo and its features, please see our [documentation](https://gloo.solo.io).

## Documentation

* [Official Documentation](https://gloo.solo.io)
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

---

### Thanks

**Gloo** would not be possible without the valuable open-source work of projects in the community. We would like to extend a special thank-you to [envoy](https://www.envoyproxy.io).
