package helpers

import (
	"path/filepath"
	"os"
)

func E2eDirectory() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "solo-io", "glue", "test", "e2e")
}