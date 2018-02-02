package main

import (
	_ "github.com/solo-io/glue/internal/configwatcher/file"
	_ "github.com/solo-io/glue/internal/configwatcher/kube"
	_ "github.com/solo-io/glue/internal/endpointdiscovery/kube"
	_ "github.com/solo-io/glue/internal/ingressconverter/kube"
	_ "github.com/solo-io/glue/internal/secretwatcher/file"
	_ "github.com/solo-io/glue/internal/secretwatcher/kube"
	_ "github.com/solo-io/glue/internal/xds"
	"github.com/solo-io/glue/pkg/log"
)

func main() {
	log.Printf("Hi Ashish")
}
