package rlcache_test

import (
	"fmt"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	gloo_api_rl_types "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/rlcache"
)

var (
	counter = 0
)

func nextCounter() string {
	counter++
	return fmt.Sprintf("%d", counter)
}

func makeSnap() *v1snap.ApiSnapshot {
	return &v1snap.ApiSnapshot{
		Ratelimitconfigs: []*gloo_api_rl_types.RateLimitConfig{
			{
				RateLimitConfig: ratelimit.RateLimitConfig{
					ObjectMeta: v1.ObjectMeta{
						Name:      "rlc1",
						Namespace: "ns1",
						Labels: map[string]string{
							"foo": "bar1",
						},
						// make sure all objects are unique in assertions below.
						ResourceVersion: nextCounter(),
					},
				},
			},
			{
				RateLimitConfig: ratelimit.RateLimitConfig{
					ObjectMeta: v1.ObjectMeta{
						Name:      "rlc2",
						Namespace: "ns1",
						Labels: map[string]string{
							"foo": "bar2",
						},
						ResourceVersion: nextCounter(),
					},
				},
			},
			{
				RateLimitConfig: ratelimit.RateLimitConfig{
					ObjectMeta: v1.ObjectMeta{
						Name:      "rlc2",
						Namespace: "ns2",
						Labels: map[string]string{
							"foo": "bar3",
						},
						ResourceVersion: nextCounter(),
					},
				},
			},
		},
	}
}

var _ = Describe("Rlcache", func() {

	It("should collect rate limit configs to a map", func() {
		snap := makeSnap()
		ratelimits := rlcache.CollectRateLimits(snap)
		Expect(ratelimits).To(HaveLen(3))
		Expect(ratelimits[rlcache.RLKey{Name: "rlc1", Namespace: "ns1"}]).To(Equal(snap.Ratelimitconfigs[0]))
		Expect(ratelimits[rlcache.RLKey{Name: "rlc2", Namespace: "ns1"}]).To(Equal(snap.Ratelimitconfigs[1]))
		Expect(ratelimits[rlcache.RLKey{Name: "rlc2", Namespace: "ns2"}]).To(Equal(snap.Ratelimitconfigs[2]))
	})

	It("should clean up rate limits when snapshot is removed", func() {
		cache := rlcache.NewRLCache()
		// snap is dead after the function returns
		func() {
			snap := makeSnap()
			_, err := cache.FindRateLimit(snap, "ns1", "rlc1")
			Expect(err).NotTo(HaveOccurred())
			Expect(cache.Len()).To(Equal(1))
		}()
		// need one gc to run finalizer, and one gc to actually clean up.
		runtime.GC()

		runtime.GC()
		Expect(cache.Len()).To(Equal(0))

	})

	It("should find rate limits in a snapshot", func() {
		snap1 := makeSnap()
		snap2 := makeSnap()
		cache := rlcache.NewRLCache()
		rlc11, err := cache.FindRateLimit(snap1, "ns1", "rlc1")
		Expect(err).NotTo(HaveOccurred())
		Expect(rlc11).To(Equal(snap1.Ratelimitconfigs[0]))
		Expect(rlc11).NotTo(Equal(snap2.Ratelimitconfigs[0]))

		_, err = cache.FindRateLimit(snap1, "ns1", "not-here")
		Expect(err).To(MatchError("list did not find rateLimitConfig ns1.not-here"))

		rlc21, err := cache.FindRateLimit(snap2, "ns1", "rlc1")
		Expect(err).NotTo(HaveOccurred())
		Expect(rlc21).To(Equal(snap2.Ratelimitconfigs[0]))

		Expect(cache.Len()).To(Equal(2))
	})
})
