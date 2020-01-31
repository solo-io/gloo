package ratelimit

import errors "github.com/rotisserie/eris"

var (
	RateLimitAuthOrderingConflict = errors.New("rate limiting specified to happen before auth, but auth-based rate limits provided in the virtual host")
)
