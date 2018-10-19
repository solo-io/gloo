package graphql

import (
	"os"

	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/models"
)

func getOAuthEndpoint() (models.OAuthEndpoint, error) {
	oauthUrl := os.Getenv("OAUTH_SERVER") // ip:port of openshift server
	if oauthUrl == "" {
		return models.OAuthEndpoint{}, errors.Errorf("apiserver configured improperly, OAUTH_SERVER environment variable is not set")
	}
	oauthClient := os.Getenv("OAUTH_CLIENT") // ip:port of openshift server
	if oauthClient == "" {
		return models.OAuthEndpoint{}, errors.Errorf("apiserver configured improperly, OAUTH_CLIENT environment variable is not set")
	}
	return models.OAuthEndpoint{
		URL:        oauthUrl,
		ClientName: oauthClient,
	}, nil
}
