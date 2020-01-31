package configsvc

import (
	"time"

	"github.com/pkg/errors"
)

var (
	LicenseIsInvalidError = func(err error) error {
		return errors.Wrap(err, "License is invalid")
	}

	FailedToReadSettingsError = func(err error) error {
		return errors.Wrap(err, "Failed to read settings")
	}

	FailedToUpdateSettingsError = func(err error) error {
		return errors.Wrap(err, "Failed to update settings")
	}

	FailedToParseSettingsFromYamlError = func(err error) error {
		return errors.Wrap(err, "Failed to parse settings from YAML")
	}

	InvalidRefreshRateError = func(d time.Duration) error {
		return errors.Errorf("Refresh rate must be at least one second. %v seconds provided", d.Seconds())
	}

	FailedToListNamespacesError = func(err error) error {
		return errors.Wrap(err, "Failed to list namespaces")
	}
)
