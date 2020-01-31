package grafana

import (
	errors "github.com/rotisserie/eris"
)

var (
	DashboardNotFound = func(upstreamUid string) error {
		return errors.Errorf("could not find dashboard for upstream %s", upstreamUid)
	}
	IncompleteGrafanaCredentials = errors.New("Incomplete grafana credentials provided")
)
