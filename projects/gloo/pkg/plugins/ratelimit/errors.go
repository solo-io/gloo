package ratelimit

import errors "github.com/rotisserie/eris"

var (
	RouteTypeMismatchErr = errors.Errorf("internal error: input route has route action but output route has not")
	ConfigNotFoundErr    = func(ns, name string) error {
		return errors.Errorf("could not find RateLimitConfig resource with name [%s] in namespace [%s]", name, ns)
	}
	ReferencedConfigErr = func(err error, ns, name string) error {
		return errors.Wrapf(err, "failed to process RateLimitConfig resource with name [%s] in namespace [%s]", name, ns)
	}
	MissingNameErr     = errors.Errorf("Cannot configure basic rate limit for resource without name.")
	DuplicateNameError = func(name string) error {
		return errors.Errorf("Basic rate limit already configured for resource with name [%s]; routes and virtual hosts must have distinct names if configured with basic ratelimits.", name)
	}
	AuthOrderingConflict = errors.New("rate limiting specified to happen before auth, but auth-based rate limits provided in the virtual host")
)
