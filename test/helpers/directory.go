package helpers

import (
	"os"
	"path/filepath"
)

func E2eDirectory() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "solo-io", "gloo", "test", "e2e")
}
