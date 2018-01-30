package helpers

import "github.com/solo-io/glue/secrets/watcher"

func NewTestSecrets() watcher.SecretMap {
	return map[string]map[string][]byte{
		"user": {"username": []byte("me@example.com"), "password": []byte("foobar")},
	}
}
