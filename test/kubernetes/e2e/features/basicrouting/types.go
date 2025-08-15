package basicrouting

import (
	_ "embed"
)

//go:embed testdata/vs-with-retries.yaml
var VirtualServiceWithRetriesYaml []byte
