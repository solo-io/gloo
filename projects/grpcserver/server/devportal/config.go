package devportal

import "os"

const (
	EnabledEnvName = "DEV_PORTAL_ENABLED"
	Enabled        = "true"
)

// Returns true if the developer portal is enabled
func IsEnabled() bool {
	return os.Getenv(EnabledEnvName) == Enabled
}
