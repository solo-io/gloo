package pkg_test

import (
	"context"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-kit/projects/discovery/pkg"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
)

type testDiscovery struct {
	isUpstreamFunctionalResult bool
	serviceSpec                *plugins.ServiceSpec
	detectUpstreamTypeError    error
	detectFunctionsError       error
	mutate                     UpstreamMutator

	isUpstreamFunctional bool
	detectUpstreamType   bool
	detectFunctions      bool
}

func (t *testDiscovery) IsUpstreamFunctional(u *v1.Upstream) bool {
	t.isUpstreamFunctional = true
	return t.isUpstreamFunctionalResult
}
func (t *testDiscovery) DetectUpstreamType(ctx context.Context, url *url.URL) (*plugins.ServiceSpec, error) {
	t.detectUpstreamType = true
	return t.serviceSpec, t.detectUpstreamTypeError
}

func (t *testDiscovery) DetectFunctions(ctx context.Context, secrets func() v1.SecretList, in *v1.Upstream, out func(UpstreamMutator) error) error {
	t.detectFunctions = true
	if t.mutate != nil {
		out(t.mutate)
	}

	return t.detectFunctionsError
}

type fakeResolver struct {
	resolveUrl   *url.URL
	resolveError error
}

func (t *fakeResolver) Resolve(u *v1.Upstream) (*url.URL, error) {
	return t.resolveUrl, t.resolveError
}

var _ = Describe("Updater", func() {

	var (
		ctx      context.Context
		cancel   context.CancelFunc
		resolver *fakeResolver
		testDisc *testDiscovery
		updater  *Updater
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		resolver = &fakeResolver{}
		testDisc = &testDiscovery{}
		updater = NewUpdater(ctx, resolver, 0, []FunctionDiscovery{testDisc})
	})
	AfterEach(func() {
		cancel()
	})

	It("should detect functions when upstream type is known", func() {
		testDisc.isUpstreamFunctionalResult = true
		updater.UpstreamAdded(nil)
		time.Sleep(time.Second / 10)
		Expect(testDisc.isUpstreamFunctional).To(BeTrue())
		Expect(testDisc.detectUpstreamType).To(BeFalse())
		Expect(testDisc.detectFunctions).To(BeTrue())
	})

})
