package defaults

import (
	"os"
	"path/filepath"
)

var (
	WasmConfigDir       = home() + "/.gloo/wasm"
	WasmImageDir        = filepath.Join(WasmConfigDir, "store")
	WasmCredentialsFile = filepath.Join(WasmConfigDir, "credentials.json")
)

func home() string {
	dir, _ := os.UserHomeDir()
	return dir
}
