package extauth_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/solo-io/solo-projects/test/services/redis"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/solo-projects/test/services/extauth"

	"github.com/solo-io/gloo/test/services/envoy"
	glooe_envoy "github.com/solo-io/solo-projects/test/services/envoy"

	"github.com/solo-io/solo-projects/test/e2e"

	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/solo-kit/test/helpers"
)

var (
	envoyFactory       envoy.Factory
	extAuthFactory     *extauth.Factory
	redisFactory       *redis.Factory
	testContextFactory *e2e.TestContextFactory
)

var _ = BeforeSuite(func() {
	envoyFactory = glooe_envoy.NewFactory()
	extAuthFactory = extauth.NewFactory()
	redisFactory = redis.NewFactory()

	testContextFactory = &e2e.TestContextFactory{
		EnvoyFactory:   envoyFactory,
		ExtAuthFactory: extAuthFactory,
	}
})

var _ = AfterSuite(func() {
	envoyFactory.Clean()
})

func TestE2eExtAuth(t *testing.T) {
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()

	RunSpecs(t, "E2e ExtAuth Suite")
}

/*
	Helper functions used in multiple ext-auth tests
*/

func GetBasicAuthExtension() *v1.ExtAuthExtension {
	return &v1.ExtAuthExtension{
		Spec: &v1.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "basic-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getBasicAuthConfig() *v1.BasicAuth {
	return &v1.BasicAuth{
		Realm: "gloo",
		Apr: &v1.BasicAuth_Apr{
			Users: map[string]*v1.BasicAuth_Apr_SaltedHashedPassword{
				"user": {
					// Password is password
					Salt:           "0adzfifo",
					HashedPassword: "14o4fMw/Pm2L34SvyyA2r.",
				},
			},
		},
	}
}

type unsecureCookieJar struct {
	http.CookieJar
	OriginalCookies map[string]*http.Cookie
}

func (j *unsecureCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if j.OriginalCookies == nil {
		j.OriginalCookies = make(map[string]*http.Cookie)
	}
	for _, c := range cookies {
		// hack to work around go client impl that doesn't consider localhost secure.
		c.Secure = false
		// the Cookies() method from the cookie jar removes the properties of the original cookie.
		j.OriginalCookies[c.Name] = c
	}

	j.CookieJar.SetCookies(u, cookies)

}
