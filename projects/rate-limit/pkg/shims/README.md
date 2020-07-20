# `rate-limiter` library shims 
This package contains shims for the `rate-limiter` libraries that we need to use in `solo-projects`.

OS Gloo needs to import the `rate-limiter` Go types for code generation (so it can include `RateLimitConfigs` in the `ApiSnapshot`), 
but - since OS Gloo cannot directly depend on the private `rate-limiter` repository - we publish the rate limit protobuf 
API definitions to the `solo-api` repository and regenerate the code there.

The downside of this approach is that the `solo-api` types cannot be used directly with the `rate-limiter` libraries; 
hence the need to convert between the two. The `internal` sub-package contains the core conversion logic; it should be 
the only place in this whole repository where we directly import the following `rate-limiter` packages:

```
	"github.com/solo-io/rate-limiter/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/rate-limiter/pkg/api/ratelimit.solo.io/v1alpha1/types"
``` 