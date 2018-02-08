package storage

import (
	"strings"

	"github.com/pkg/errors"
)

func CreateStorageRef(namespace, name string) string {
	return namespace + "/" + name
}

// returns namespace, name from "namespace/name"
func ParseStorageRef(ref string) (string, string, error) {
	parts := strings.Split(ref, "/")
	if len(parts) != 2 {
		return "", "", errors.Errorf("invalid storage ref for kubernetes: %v", ref)
	}
	return parts[0], parts[1], nil
}
