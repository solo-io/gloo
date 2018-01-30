package helpers

import "github.com/solo-io/glue/secrets"

func NewTestSecrets() secrets.SecretMap {
	return map[string]map[string][]byte{
		"user": {"username": []byte("me@example.com"), "password": []byte("foobar")},
	}
}
