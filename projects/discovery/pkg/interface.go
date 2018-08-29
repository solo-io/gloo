package pkg

import (
	"context"
	"net/url"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
)

type UpstreamMutator func(*v1.Upstream) error

/*
upstreams can be obviously funcitonal like AWS λ, fission,...  or an upstream that was already detected and marked as such.
or potentially like static upstreams.
*/
// detectors detect a specific type of functional service
// if they detect the service, they return service info and
// annotations (optional) for the service
// we want to bake sure that detect upstream for aws doesn't do anything
// perhaps we can do just that
type FuncitonDiscovery interface {
	IsUpstreamFunctional(u *v1.Upstream) bool //
	// if this returns true we can skip DetectUpstreamType and go straight to DetectFunctions
	// if this returns false we should call detect upstream type.
	// if detect upstream type retrurns true, we have the type!
	// if it returns false and nil error, it means it was detected to not be of this type -
	// ideally this means that this detector will no longer be used with this upstream. in practice this can be logged\ignored.
	// if it returns false and some error, try again later with back-off \ timeout.

	// YES = USE ME TO UpdateUpstreamWithDetectedFunctions
	// NO = DON'T USE ME TO DetectUpstreamType
	// MAYBE = USE ME TO DetectUpstreamType
	// in practice a maybe type will always return maybe. perhaps better to call detectupstreamtype and have it return error
	// for non relevant upstream types

	// try to resolve a url from the upstream
	// if we have a url continue, if no then no

	// error if not us, nil if yes.
	DetectUpstreamType(ctx context.Context, url *url.URL) (*plugins.ServiceSpec, error)

	// potentially waits for change.
	// if it detects the upstream is a known functional type, give us the
	// service info and annotations to mark it with
	// DetectFunctionalService(ctx context.Context, us *v1.Upstream, addr string) (*FuncitonDetectorResult, error)
	DetectFunctions(ctx context.Context, secrets func() v1.SecretList, in *v1.Upstream, out func(UpstreamMutator) error) error
}

type Resolver interface {
	/*
		tcp if not known
		http \ https if known or perhaps nats?
	*/
	Resolve(u *v1.Upstream) (*url.URL, error)
}

// STEP ONE, for generic upstream, detect
// NEW -> DETECTING -> TYPED()

// flow:
// upstream type: aws
// detector type: swagger (can only be used with upstreams that have a url that's resolavable)
