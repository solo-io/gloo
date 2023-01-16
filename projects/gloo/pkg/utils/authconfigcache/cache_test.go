package authconfigcache_test

import (
	"fmt"
	"runtime"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/authconfigcache"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
)

var (
	counter = 0
)

func nextCounter() string {
	counter++
	return fmt.Sprintf("%d", counter)
}

var _ = Describe("AuthConfigCache", func() {

	It("should collect auth  configs to a map", func() {
		snap := makeSnap()
		authConfigs := authconfigcache.CollectAuthConfigs(snap.AuthConfigs)
		Expect(authConfigs).To(HaveLen(3))
		Expect(authConfigs[authconfigcache.ACKey{Name: "ac-1", Namespace: defaults.GlooSystem}]).To(Equal(snap.AuthConfigs[0]))
		Expect(authConfigs[authconfigcache.ACKey{Name: "ac-2", Namespace: defaults.GlooSystem}]).To(Equal(snap.AuthConfigs[1]))
		Expect(authConfigs[authconfigcache.ACKey{Name: "ac-3", Namespace: defaults.GlooSystem}]).To(Equal(snap.AuthConfigs[2]))
	})

	It("should clean up auth configs when snapshot is removed", func() {
		cache := authconfigcache.NewAuthConfigCache(10)
		// snap is dead after the function returns
		func() {
			snap := makeSnap()
			_, err := cache.FindAuthConfig(snap, defaults.GlooSystem, "ac-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(cache.Len()).To(Equal(1))
		}()
		// need one gc to run finalizer, and one gc to actually clean up.
		runtime.GC()

		runtime.GC()
		Expect(cache.Len()).To(Equal(0))

	})

	It("should find auth configs in a snapshot", func() {
		snap1 := makeSnap()
		snap2 := makeSnap()
		cache := authconfigcache.NewAuthConfigCache(10)
		authConfig11, err := cache.FindAuthConfig(snap1, defaults.GlooSystem, "ac-1")
		Expect(err).NotTo(HaveOccurred())
		Expect(authConfig11).To(Equal(snap1.AuthConfigs[0]))
		Expect(authConfig11).NotTo(Equal(snap2.AuthConfigs[0]))

		_, err = cache.FindAuthConfig(snap1, defaults.GlooSystem, "not-here")
		Expect(err).To(MatchError("list did not find authConfig gloo-system.not-here"))

		authConfig21, err := cache.FindAuthConfig(snap2, defaults.GlooSystem, "ac-1")
		Expect(err).NotTo(HaveOccurred())
		Expect(authConfig21).To(Equal(snap2.AuthConfigs[0]))

		Expect(cache.Len()).To(Equal(2))
	})
})

func makeSnap() *v1snap.ApiSnapshot {
	return &v1snap.ApiSnapshot{
		AuthConfigs: []*extauth.AuthConfig{
			{
				Metadata: &core.Metadata{
					Name:      "ac-1",
					Namespace: defaults.GlooSystem,
					// make sure all objects are unique in assertions
					ResourceVersion: nextCounter(),
				},
				Configs: []*extauth.AuthConfig_Config{},
			},
			{
				Metadata: &core.Metadata{
					Name:            "ac-2",
					Namespace:       defaults.GlooSystem,
					ResourceVersion: nextCounter(),
				},
				Configs: []*extauth.AuthConfig_Config{},
			},
			{
				Metadata: &core.Metadata{
					Name:            "ac-3",
					Namespace:       defaults.GlooSystem,
					ResourceVersion: nextCounter(),
				},
				Configs: []*extauth.AuthConfig_Config{},
			},
		},
	}
}
