package helpers

import (
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

func NewTestSecrets() secretwatcher.SecretMap {
	return map[string]map[string]string{
		"user": {"username": "me@example.com", "password": "foobar"},
	}
}
