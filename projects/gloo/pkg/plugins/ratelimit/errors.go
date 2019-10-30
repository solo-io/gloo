package ratelimit

import "github.com/solo-io/go-utils/errors"

var (
	RateLimitAuthOrderingConflict = errors.New("rate limiting specified to happen before auth, but auth-based rate limits provided in the virtual host")
)
