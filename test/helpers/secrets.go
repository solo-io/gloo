package helpers

import (
	"github.com/solo-io/glue/pkg/secretwatcher"
)

func NewTestSecrets() secretwatcher.SecretMap {
	return map[string]map[string][]byte{
		"user": {"username": []byte("me@example.com"), "password": []byte("foobar")},
	}
}
