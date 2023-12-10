package vault

import (
	"testing"

	"github.com/solo-io/gloo/test/gomega/assertions"
	"go.opencensus.io/stats/view"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVaultClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vault Client Suite")
}

var _ = BeforeSuite(func() {
	resetViews()
})

var _ = AfterSuite(func() {
	resetViews()
})

// resetViews resets the views used in this package
// this is useful for ensuring that we are not leaking metrics between tests within this package
// or between this package and other packages
func resetViews() {
	views := []*view.View{
		mLastLoginSuccessView,
		mLoginFailuresView,
		mLoginSuccessesView,
		mLastLoginFailureView,
	}
	view.Unregister(views...)
	_ = view.Register(views...)
	assertions.ExpectStatLastValueMatches(mLastLoginSuccess, BeZero())
	assertions.ExpectStatLastValueMatches(mLastLoginFailure, BeZero())
	assertions.ExpectStatSumMatches(mLoginSuccesses, BeZero())
	assertions.ExpectStatSumMatches(mLoginFailures, BeZero())
}
