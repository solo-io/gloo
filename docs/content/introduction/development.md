---
title: Development
weight: 60
---

Developers can work on Gloo Gateway in a number of different ways. You can contribute to Gloo Gateway open-source, help with documentation, or extend Gloo Gateway's functionality via the addition of plugins.

---

## Contributing

You can contribute to the Gloo Gateway open-source project by logging issues, generating PRs, or helping with documentation. More information about contributing to Gloo Gateway can be found in our [Contributing]({{% versioned_link_path fromRoot="/contributing/" %}}) section.

---

## Plugins

Gloo Gateway invites developers to extend Gloo Gateway's functionality and adapt to new use cases via the addition of plugins. 

Gloo Gateway's plugin based architecture makes it easy to extend functionality in a variety of areas:

- [Gloo Gateway's API](https://github.com/solo-io/gloo/tree/main/projects/gloo/api/v1): extensible through the use of [Protocol Buffers](https://developers.google.com/protocol-buffers/) along with [Solo-Kit](https://github.com/solo-io/solo-kit)
- [Service Discovery Plugins](https://github.com/solo-io/gloo/blob/main/projects/gloo/pkg/discovery/discovery.go#L21): automatically discover service endpoints from catalogs such as [Kubernetes](https://github.com/solo-io/gloo/tree/main/projects/gloo/pkg/plugins/kubernetes) and [Consul](https://github.com/solo-io/gloo/tree/main/projects/gloo/pkg/plugins/consul)
- [Function Discovery Plugins](https://github.com/solo-io/gloo/blob/main/projects/discovery/pkg/fds/interface.go#L31): annotate services with information discovered by polling services directly (such as OpenAPI endpoints and gRPC methods).
- [Routing Plugins](https://github.com/solo-io/gloo/blob/main/projects/gloo/pkg/plugins/plugin_interface.go#L53): customize what happens to requests when they match a route or virtual host
- [Upstream Plugins](https://github.com/solo-io/gloo/tree/main/projects/gloo/pkg/plugins): customize what happens to requests when they are routed to a service

To get started with development around Gloo Gateway, we recommend taking a look at our [Developer Guides]({{% versioned_link_path fromRoot="/guides/dev/" %}}).

