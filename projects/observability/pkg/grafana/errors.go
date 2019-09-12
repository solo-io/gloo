package grafana

import (
	"github.com/solo-io/go-utils/errors"
)

var (
	DashboardNotFound = func(upstreamUid string) error {
		return errors.Errorf("could not find dashboard for upstream %s", upstreamUid)
	}
)
