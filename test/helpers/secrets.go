package helpers

import "github.com/solo-io/glue/secrets/watcher"

func NewTestSecrets() watcher.SecretMap {
	return map[string]map[string]string{
		"user": {"username": "me@example.com", "password": "foobar"},
	}
}
