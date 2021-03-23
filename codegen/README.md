# gloo-fed Code Generation
Gloo Fed skv2 code generation with custom templates.

## SKV2 vs SKV1 (solo-kit) Code Generation
Note that Gloo Fed codegen uses skv2, and uses protos that are in a different format:

```
message UpstreamSpec {}
message UpstreamStatus {}
```
(e.g. in github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream.proto)

instead of the skv1-compatible proto:
```
message Upstream{}
```
(e.g. in github.com/solo-io/gloo/projects/gloo/api/v1/upstream.proto)

To prevent conflicts, skv2 code generation is done in a separate make target
that clears away the previous vendor_any directory.

Additionally, The generated go code from these protos must be kept in separate packages.
The package names are registered to a global registry. If they are both pulled into the same binary,
an error like this will appear:
```
2021/02/11 00:25:12 WARNING: proto: file "github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto" has a name conflict over enterprise.gloo.solo.io.ExtAuthConfig
	previously from: "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	currently from:  "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
```

A consequence of this is that no solo-kit types may be used in these binaries:
gloo-fed, apiserver, rbac-validating-webhook, glooctl-fed.

No skv2 types may be used in these binaries:
grpcserver, gloo, extauth, rate-limit, observability.

## Code Generation make targets

To generate code, run:
```
make clean install-go-tools generate-gloo-fed-code
```

# Custom Templates
A `model.Group` object in the groups.go file defines a logical
'grouping' of CRDs. In gloo fed, we have the following groups:

    FedGroup,
    GatewayGroup,
    FedGatewayGroup,
    GlooGroup,
    FedGlooGroup,
    EnterpriseGlooGroup,
    FedEnterpriseGroup,
    EnterpriseGlooGroup,
    FedEnterpriseGroup,
    RateLimitGroup,
    FedRateLimitGroup

Each group except for the `FedGroup` is based on an existing Gloo package type:

    gateway.solo.io
    gloo.solo.io
    enterprise.gloo.solo.io
    ratelimit.api.solo.io

Each group has a field called CustomTemplates:

    // data for providing custom templates to generate custom code for groups
	CustomTemplates []CustomTemplates
	
This is where all the custom templates defined in this
directory will go.
