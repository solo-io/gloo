package basicrouting

import (
	_ "embed"
)

//go:embed testdata/vs-with-retries.yaml
var VirtualServiceWithRetriesYaml []byte

//go:embed testdata/nginx-upstream.yaml
var NginxUpstreamYaml []byte
