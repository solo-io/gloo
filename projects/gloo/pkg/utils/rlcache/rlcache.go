package rlcache

import (
	"runtime"
	"unsafe"

	errors "github.com/rotisserie/eris"
	gloo_api_rl_types "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"k8s.io/utils/lru"
)

type RLCache struct {
	rateLimitsLruCache *lru.Cache
}

func NewRLCache() *RLCache {
	// we only ever expect to need to look at the most recent snapshot.
	// we use a cache of size 10 to have some wiggle room in case we need more for a short amount of time
	// a runtime finalizer is used to make sure we release memory as soon as it is not needed.
	return &RLCache{rateLimitsLruCache: lru.New(10)}
}

type RLKey struct{ Name, Namespace string }
type RLMap map[RLKey]*gloo_api_rl_types.RateLimitConfig

func (c *RLCache) FindRateLimit(snap *v1snap.ApiSnapshot, namespace, name string) (*gloo_api_rl_types.RateLimitConfig, error) {
	var ratelimits RLMap
	key := uintptr(unsafe.Pointer(snap))
	if cached, ok := c.rateLimitsLruCache.Get(key); ok {
		ratelimits = cached.(RLMap)
	} else {
		ratelimits = CollectRateLimits(snap)
		c.rateLimitsLruCache.Add(key, ratelimits)
		runtime.SetFinalizer(snap, c.finalize)
	}
	rlc, ok := ratelimits[RLKey{Name: name, Namespace: namespace}]
	if !ok {
		return nil, errors.Errorf("list did not find rateLimitConfig %v.%v", namespace, name)
	}
	return rlc, nil
}

func (c *RLCache) Len() int {
	return c.rateLimitsLruCache.Len()
}

func (c *RLCache) finalize(snap *v1snap.ApiSnapshot) {
	key := uintptr(unsafe.Pointer(snap))
	c.rateLimitsLruCache.Remove(key)
}

func CollectRateLimits(snap *v1snap.ApiSnapshot) RLMap {

	ratelimits := RLMap{}
	for _, rlc := range snap.Ratelimitconfigs {
		ratelimits[RLKey{Name: rlc.Name, Namespace: rlc.Namespace}] = rlc
	}
	return ratelimits
}
