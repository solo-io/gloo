package graphql

import (
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/config"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

func getOAuthEndpoint() (models.OAuthEndpoint, error) {
	oauthURL, oauthClient := config.GetOAuthEndpointValues()
	return models.OAuthEndpoint{
		URL:        oauthURL,
		ClientName: oauthClient,
	}, nil
}

func getAPIVersion() string {
	return config.APIVersion
}
