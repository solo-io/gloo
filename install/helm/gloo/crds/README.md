The CRDs in this directory are partially generated by solo-kit.
That implementation is explained [here](https://github.com/solo-io/solo-kit/tree/main/pkg/code-generator/schemagen#implementation).

It is worth noting that solo-kit only generates the schemas for these CRDs.
Other spec fields such as the `categories` can be updated manually.

Gloo Gateway CRD `categories`:
- All Gloo Gateway CRDs should include the "gloo-gateway" category.
- Any Gloo Gateway CRDs which are only used by enterprise customers should additionally include the "solo-io" category.
  - Currently, these are the AuthConfig, RateLimitConfig, and GraphQLApi CRDs.