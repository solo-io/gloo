package vault

import (
	"github.com/solo-io/gloo/pkg/utils"
	"go.opencensus.io/stats/view"
)

func init() {
	_ = view.Register(
		mLastLoginSuccessView,
		mLastLoginFailureView,
		mLoginSuccessesView,
		mLoginFailuresView,
	)
}

var (
	mLastLoginSuccess     = utils.Int64Measure("gloo.solo.io/vault/last_login_success", "Timestamp of last successful authentication of vault")
	mLastLoginSuccessView = utils.ViewForCounter(mLastLoginSuccess, view.LastValue())

	mLastLoginFailure     = utils.Int64Measure("gloo.solo.io/vault/last_login_failure", "Timestamp of last failed authentication of vault")
	mLastLoginFailureView = utils.ViewForCounter(mLastLoginFailure, view.LastValue())

	mLoginSuccesses     = utils.Int64Measure("gloo.solo.io/vault/login_successes", "Number of successful authentications of vault")
	mLoginSuccessesView = utils.ViewForCounter(mLoginSuccesses, view.Sum())

	mLoginFailures     = utils.Int64Measure("gloo.solo.io/vault/login_failures", "Number of failed authentications of vault")
	mLoginFailuresView = utils.ViewForCounter(mLoginFailures, view.Sum())
)
