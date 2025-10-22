---
weight: 99
title: Gloo Gateway
---

# Gloo Gateway (Gloo Edge API)

Welcome to the Gloo Gateway (Gloo Edge API) documentation.

## What is Gloo Gateway?

Gloo Gateway is a feature-rich, Envoy-powered, Kubernetes-native ingress controller, and next-generation API gateway. Gloo Gateway is exceptional in its function-level routing; its support for legacy apps, microservices and serverless; its discovery capabilities; its numerous features; and its tight integration with leading open-source projects. Gloo Gateway is uniquely designed to support hybrid applications, in which multiple technologies, architectures, protocols, and clouds can coexist.

{{% notice tip %}}Want to use Gloo Gateway with the Kubernetes Gateway API? Check out the [separate product documentation set](https://docs.solo.io/gateway/latest/). This product documentation set is to use Gloo Gateway with Solo's `Gateway` API custom resource.{{% /notice %}}

![Gloo Gateway Architecture]({{% versioned_link_path fromRoot="/img/gloo-architecture-envoys.png" %}})

## Next Generation API Gateway

Although the idea of the API Gateway has been around for a bit, the role of the API Gateway is going through an identity crisis as we adopt more automated, self-service, platforms like Kubernetes, Cloud Foundry, and public-cloud. Your existing API Management solutions weren't built for highly dynamic environments like Kubernetes and require a lot of additional infrastructure to keep up, be highly-available, and production ready. Additionally, a lot of organizations have deployed these pieces of infrastructure in a highly centralized way that plays to the assumptions of the solution and not the desires of the the organization.

When we say Gloo Gateway is a "next-generation" gateway, we mean that it was purpose-built for a highly dynamic, ephemeral environment like Kubernetes (or other workload orchestration platforms) and is built with the assumption of decentralized ownership. Gloo Gateway can provide powerful API Gateway functionality for both existing, on-premises investments (like VM deployments or physical hardware), as well as Kubernetes, and even including forward-leaning compute options like Function as a Service. Legacy API Management vendors would have to completely re-write their solutions to play nicely in this new cloud-native world.

---

## Getting to know Gloo Gateway

* [**Getting Started**]({{% versioned_link_path fromRoot="/getting_started/" %}})
* [**Installation**]({{% versioned_link_path fromRoot="/installation" %}})
* [**Developers**]({{% versioned_link_path fromRoot="/guides/dev/" %}})

## Blogs & Demos

* [Blog Articles](https://www.solo.io/blog/announcing-gloo-the-function-gateway/)
* [Watch Gloo Gateway Videos and Demos](https://www.youtube.com/watch?v=HEN8IHCOqSo&list=PLBOtlFtGznBgrbLwyrPNIsdXLuq-mqa6U)

## Community

* Join us on our Slack channel: [https://slack.solo.io/](https://slack.solo.io/)
* Follow us on Twitter: [https://twitter.com/soloio_inc](https://twitter.com/soloio_inc)

---

## Thanks

Gloo Gateway would not be possible without the valuable open-source work of projects in the community. We would like to extend a special thank-you to [Envoy](https://www.envoyproxy.io).
