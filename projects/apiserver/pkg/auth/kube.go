package auth

import (
	"strings"

	"github.com/solo-io/go-utils/errors"
	authv1 "k8s.io/api/authentication/v1"
	"k8s.io/client-go/kubernetes/typed/authentication/v1"
)

// Attempts to authenticate a token to a known user.
// Note: TokenReview requests may be cached by the webhook token authenticator plugin in the kube-apiserver.
func GetUsername(auth v1.AuthenticationV1Interface, token string) (username string, err error) {

	token = strings.TrimPrefix(token, "Bearer ")

	tokenReview, err := auth.TokenReviews().Create(&authv1.TokenReview{Spec: authv1.TokenReviewSpec{Token: token}})
	if err != nil {
		return "", err
	}

	if !tokenReview.Status.Authenticated {
		return "", errors.Errorf(tokenReview.Status.Error)
	}

	return tokenReview.Status.User.Username, nil
}
