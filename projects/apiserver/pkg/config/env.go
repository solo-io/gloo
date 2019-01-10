package config

import (
	"os"
)

// ValidateEnvVars must be called during initialization
func ValidateEnvVars() {
	oauthURL := getOauthURL()
	if oauthURL == "" {
		panic("apiserver configured improperly, OAUTH_SERVER environment variable is not set")
	}
	oauthClient := getOauthClient()
	if oauthClient == "" {
		panic("apiserver configured improperly, OAUTH_CLIENT environment variable is not set")
	}
}

// GetOAuthEndpointValues returns the URL and Client needed for OAuth
// ValidateEnvVars must be called beforehand (on initialization)
func GetOAuthEndpointValues() (string, string) {
	return getOauthURL(), getOauthClient()
}

func getOauthURL() string {
	return os.Getenv("OAUTH_SERVER") // ip:port of openshift server
}

func getOauthClient() string {
	return os.Getenv("OAUTH_CLIENT")
}
