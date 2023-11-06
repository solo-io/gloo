# External Authn/z

## Authn/z in Envoy
Read the [Envoy documentation about External Authorization](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter#arch-overview-ext-authz)

There is an External Authorization HTTP Filter which calls out to an external gRPC or HTTP service to determine whether to authorize the request or not.

There are 3 aspects to Extauth:
- The [HttpFilter configuration](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_authz/v3/ext_authz.proto.html#extensions-filters-http-ext-authz-v3-extauthz)
- The [Route level configuration](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_authz/v3/ext_authz.proto.html#extensions-filters-http-ext-authz-v3-extauthzperroute)
- The [service that performs authorization](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter#service-definition)

## Open Source API

Gloo [Extauth Plugin which converts Gloo resources into Envoy resources](https://github.com/solo-io/gloo/blob/main/projects/gloo/pkg/plugins/extauth/plugin.go)

### Http Filter (Gateway CR)
On a Gateway resource, you can define [Http specific configuration](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gateway/api/v1/gateway.proto#L68). There are [HttpListenerOptions](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gateway/api/v1/http_gateway.proto#L45), and within these options you can define [ExtAuth configuration](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/options.proto#L114).

This [configuration](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto#L96) is converted into the HttpFilter that lives on the Envoy Listener.

### Http Filter (Settings CR)
You can define this configuration per Gateway (which becomes a listener). Alternatively, if you want to define Global configuration, you can place it on the [Settings object](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/settings.proto#L378)

### Route Level Configuration
The [route level configuration](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto#L84) can be defined on either:
- A [virtual host](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/options.proto#L262)
- A [route](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/options.proto#L418)
- A [weighted destination](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/options.proto#L508)

In the Route level API, only the `customAuth` and `disable` config are actually used in the open source implementation

## Open Source Limitations

1. No server. Open source users must build their own authn/z service, and reference it in their filter config
2. 1 filter. Only a single filter can be configured for open source user. This means if you wanted slightly different filters (ie different timeouts and failure modes) for different routes, you could not do that

## Enterprise API

### Multiple ExtAuth Filters (Settings CR)
On a Settings resource, you can define [multiple sets of extauth configuration](https://github.com/solo-io/gloo/blob/dc55e357ca32c0639147d26dbc5742e982b871d1/projects/gloo/api/v1/settings.proto#L453). Each option represents the configuration for a single ExtAuth filter. Therefore, having the ability to define multiple filters can be powerful if you want to have certain requests be processed by a different auth service.

### ExtAuth Server API
The open source product requires that you ship your own ExtAuth service. The Enterprise product [ships a service](https://github.com/solo-io/ext-auth-service) for you. We need some way of enabling users to define the configuration which the server will use.

Our [Server API](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto#L26) is defined by users as a standalone resource.

We have a user facing API and an internal API. It is confusing, but both of these are defined in the same file:
- [External API](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto#L26): Users configure this. Note the `status` and `metadata` fields which we indicated are signals that an object is a Custom Resource.
- [Internal API](https://github.com/solo-io/gloo/blob/364a05b040fd416bda3a934eb20aef4e091198fc/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto#L1169): Gloo Edge converts the external API into this object before sending to the ExtAuth service via xDS

### ExtAuth Filter Configuration Translation
Gloo Enterprise [Extauth Plugin which converts Gloo resources into Envoy resources](https://github.com/solo-io/solo-projects/blob/main/projects/gloo/pkg/plugins/extauth/plugin.go)

### ExtAuth Server Configuration Translation
Gloo Edge does not require Kubernetes to run. Therefore, instead of depending on a Kubernetes controller to process AuthConfig's, the Gloo component translates AuthConfig's (external API) into ExtAuthConfig (internal API) and we use xDs to send config to the ext-auth-service

**A common mistake is to assume that the [translation package](https://github.com/solo-io/ext-auth-service/tree/master/pkg/controller/translation) in the ext-auth-service is responsible for translating the AuthConfig. This is not the case. The translation logic all lives in this repository, either within [projects/gloo](/projects/gloo) or [projects/extauth](/projects/extauth) package.**

Below are the rough steps that a user defined AuthConfig API object takes in the Gloo component (Pod). It builds on the concepts outlined above, and attempts to provide a high level guide to the order of events, not a low level evaluation of how those steps take place.

For this example, we will use an AuthConfig that defines behavior for [OIDC](https://github.com/solo-io/gloo/blob/4f6982adb277fc561ba06019df2cfe4b840eebfb/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto#L300)
1. User has Gloo Edge running. When starting up Gloo, you can define in the [Settings API object where resource configuration will be persisted](https://github.com/solo-io/gloo/blob/4f6982adb277fc561ba06019df2cfe4b840eebfb/projects/gloo/api/v1/settings.proto#L47). For simplicity, we will assume Kubernetes, but examine the API to see the different persistence layers we support for configuration.
2. User applies an AuthConfig CR
3. [Kubernetes API server](https://kubernetes.io/docs/concepts/overview/components/#kube-apiserver) receives the request, and [persists the configuration](https://kubernetes.io/docs/concepts/overview/components/#etcd)
4. The Gloo [API Snapshot Emitter identifies that the set of available AuthConfig objects has changed](https://github.com/solo-io/gloo/blob/4f6982adb277fc561ba06019df2cfe4b840eebfb/projects/gloo/pkg/api/v1/gloosnapshot/api_snapshot_emitter.sk.go#L1057), and pushes a new Snapshot 
5. The Snapshot is received by the [Gloo API Event Loop and passed to the Syncer to process](https://github.com/solo-io/gloo/blob/4f6982adb277fc561ba06019df2cfe4b840eebfb/projects/gloo/pkg/api/v1/gloosnapshot/api_event_loop.sk.go#L107)
6. The [Gloo Syncer processes the Snapshot](https://github.com/solo-io/gloo/blob/4f6982adb277fc561ba06019df2cfe4b840eebfb/projects/gloo/pkg/syncer/translator_syncer.go#L101), first building the Envoy resources from the Snapshot, and then processing the [Syncer Extensions](https://github.com/solo-io/gloo/blob/4f6982adb277fc561ba06019df2cfe4b840eebfb/projects/gloo/pkg/syncer/translator_syncer.go#L121). The [Extauth Translator Syncer](https://github.com/solo-io/solo-projects/blob/f35e3be72b44deb4199747e1c8c5132be8b25e7e/projects/gloo/pkg/syncer/extauth/extauth_translator_syncer.go#L74) is one of those extensions. 
7. The Extauth Translator Syncer [iterates over the set of AuthConfigs in the Snapshot](https://github.com/solo-io/solo-projects/blob/f35e3be72b44deb4199747e1c8c5132be8b25e7e/projects/gloo/pkg/syncer/extauth/xds_snapshot_producer.go#L263)
8. The OIDC AuthConfig configured by the user, is [converted from the External (user-facing) API into our Internal (auth-service) API](https://github.com/solo-io/solo-projects/blob/3a2b4e3928c35e2193a7f139128fbe5b870468f7/projects/gloo/pkg/syncer/extauth/translate.go#L66). To do this, we first expand the reference to the ClientSecret by looking up in our Snapshot for the Secret with that <name,namespace>. In Gloo Edge, Secrets can be maintained in a variety of persistence layers. Again, for simplicity let’s assume that Secrets are persisted in Kubernetes. 
9. The full set of translated AuthConfigs are [bundled into an xDS Snapshot, versioned by the hash of the configuration, and persisted in our in-memory xDS Cache](https://github.com/solo-io/solo-projects/blob/f35e3be72b44deb4199747e1c8c5132be8b25e7e/projects/gloo/pkg/syncer/extauth/extauth_translator_syncer.go#L130).

