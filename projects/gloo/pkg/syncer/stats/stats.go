package stats

import (
	"github.com/solo-io/gloo/pkg/utils/statsutils"
	"go.opencensus.io/tag"
)

var (
	ProxyNameKey, _ = tag.NewKey("proxy_name")

	// InterruptedValidationSkips counts proxy translations discarded (no xDS update, no status write)
	// because an envoy validation fork was interrupted with a live sync context. A rising count for
	// one proxy means its xDS updates are stalled. The "validation was interrupted" warning logs
	// carry the underlying error.
	InterruptedValidationSkips = statsutils.MakeSumCounter(
		"api.gloo.solo.io/translator/validation_interruption_skips",
		"The number of proxy translations discarded because envoy config validation was interrupted",
		ProxyNameKey,
	)
)
