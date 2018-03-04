package swagger

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-plugins/common/annotations"
	kubeplugin "github.com/solo-io/gloo-plugins/kubernetes"
	serviceplugin "github.com/solo-io/gloo/pkg/coreplugins/service"
	"github.com/solo-io/gloo/pkg/log"
)

// local cache to avoid retrying the same upstream
// may occasionally cause races where an update fails but we have
// already marked the upstream as "tried"
// TODO: create a cleaner way to cache upstreams that we haave marked as swagger,
// but for which updating failed (e.g. because of out of date resource version)
var triedUpstreams = make(map[string]bool)

var commonSwaggerURIs = []string{
	"/swagger.json",
	"/swagger/docs/v1",
	"/swagger/docs/v2",
	"/v1/swagger",
}

// adds swagger annotations to upstreams it discovers
func DiscoverSwaggerUpstreams(resolver *resolver.Resolver, swaggerUrisToTry []string, upstreams []*v1.Upstream) {
	for _, us := range upstreams {
		if !shouldTryDiscovery(us) {
			continue
		}
		log.Debugf("initiating swagger detection for %v", us.Name)
		triedUpstreams[us.Name] = true
		if err := discoverSwaggerUpstream(resolver, swaggerUrisToTry, us); err != nil {
			log.Warnf("unable to discover whether upstream %v implements swagger or not.", us.Name)
		}
	}
}

func shouldTryDiscovery(us *v1.Upstream) bool {
	switch {
	case us.Type != kubeplugin.UpstreamTypeKube && us.Type != serviceplugin.UpstreamTypeService:
		fallthrough
	case IsSwagger(us): //already discovered
		fallthrough
	case triedUpstreams[us.Name]:
		return false
	}
	return true
}

func discoverSwaggerUpstream(resolver *resolver.Resolver, swaggerUrisToTry []string, us *v1.Upstream) error {
	// only discover for kube or service
	// TODO: add more types here
	switch us.Type {
	default:
		return nil
	case kubeplugin.UpstreamTypeKube:
	case serviceplugin.UpstreamTypeService:
	}

	addr, err := resolver.Resolve(us)
	if err != nil {
		return err
	}
	if addr == "" {
		return nil
	}
	for _, uri := range append(swaggerUrisToTry, commonSwaggerURIs...) {
		url := "http://" + addr + uri
		log.Debugf("querying swagger url %v", url)
		res, err := http.Get(url)
		if err != nil {
			return errors.Wrapf(err, "could not perform HTTP GET on resolved addr: %v", addr)
		}
		// found a swagger service
		if res.StatusCode == 200 {
			setSwaggerAnnotations(url, us)
			return nil
		}
	}
	// not a swagger upstream
	return nil
}

func setSwaggerAnnotations(url string, us *v1.Upstream) {
	log.Debugf("swagger service detected: %v", url)
	us.Metadata.Annotations[annotations.ServiceType] = ServiceTypeSwagger
	us.Metadata.Annotations[AnnotationKeySwaggerURL] = url
	us.Metadata.Annotations[AnnotationKeySwaggerURL] = url
}
