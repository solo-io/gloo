package vault

import (
    "github.com/solo-io/gloo/pkg/utils"
    "go.opencensus.io/stats"
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
    mLastLoginSuccess     = stats.Int64("gloo.solo.io/vault/last_login_success", "Timestamp of last successful authentication of vault", stats.UnitDimensionless)
    mLastLoginSuccessView = utils.ViewForCounter(mLastLoginSuccess, view.LastValue())

    mLastLoginFailure     = stats.Int64("gloo.solo.io/vault/last_login_failure", "Timestamp of last failed authentication of vault", stats.UnitDimensionless)
    mLastLoginFailureView = utils.ViewForCounter(mLastLoginFailure, view.LastValue())

    mLoginSuccesses     = stats.Int64("gloo.solo.io/vault/login_successes", "Number of successful authentications of vault", stats.UnitDimensionless)
    mLoginSuccessesView = utils.ViewForCounter(mLoginSuccesses, view.Sum())

    mLoginFailures     = stats.Int64("gloo.solo.io/vault/login_failures", "Number of failed authentications of vault", stats.UnitDimensionless)
    mLoginFailuresView = utils.ViewForCounter(mLoginFailures, view.Sum())
)