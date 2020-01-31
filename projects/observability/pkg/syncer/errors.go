package syncer

import errors "github.com/rotisserie/eris"

var (
	DashboardIdNotFound = func(dashboardJson string) error {
		return errors.Errorf("Could not find dashboard id in dashboard json: %s", dashboardJson)
	}
	DashboardIdConversionError = func(rawDashboardId interface{}) error {
		return errors.Errorf("Could not convert %v to a float64", rawDashboardId)
	}
	NoGrafanaUrl = func(envVar string) error {
		return errors.Errorf("No grafana url configured in env var %s", envVar)
	}
)
