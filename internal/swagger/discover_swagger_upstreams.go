package swagger

import (
	"net/http"

	"github.com/pkg/errors"

	"sync"

	"github.com/cenkalti/backoff"
	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-plugins/common/annotations"
	kubeplugin "github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-plugins/transformation"
	serviceplugin "github.com/solo-io/gloo/pkg/coreplugins/service"
	"github.com/solo-io/gloo/pkg/log"
)

// local cache to avoid retrying the same upstream
var triedUpstreams = make(map[string]bool)
var mapLock sync.RWMutex

var commonSwaggerURIs = []string{
	"/swagger.json",
	"/swagger/docs/v1",
	"/swagger/docs/v2",
	"/v1/swagger",
}

// adds swagger annotations to upstreams it discovers
func DiscoverSwaggerUpstream(resolver *resolver.Resolver, swaggerUrisToTry []string, us *v1.Upstream) {
	log.Debugf("should try swagger ? %v: %v", us.Name, shouldTryDiscovery(us))
	if !shouldTryDiscovery(us) {
		return
	}
	log.Debugf("initiating swagger detection for %v", us.Name)
	err := backoff.Retry(func() error {
		return discoverSwaggerUpstream(resolver, swaggerUrisToTry, us)
	}, backoff.NewExponentialBackOff())
	if err != nil {
		log.Warnf("unable to discover whether upstream %v implements swagger or not.\n%v", us.Name, err.Error())
	}
	mapLock.Lock()
	triedUpstreams[us.Name] = true
	mapLock.Unlock()
}

func shouldTryDiscovery(us *v1.Upstream) bool {
	mapLock.RLock()
	defer mapLock.RUnlock()

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
	var errs error
	for _, uri := range append(swaggerUrisToTry, commonSwaggerURIs...) {
		url := "http://" + addr + uri
		log.Debugf("querying swagger url %v", url)
		res, err := http.Get(url)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "could not perform HTTP GET on resolved addr: %v", addr))
			continue
		}
		// found a swagger service
		if res.StatusCode == 200 {
			setSwaggerAnnotations(url, us)
			return nil
		}
	}
	// not a swagger upstream
	return errs
}

func setSwaggerAnnotations(url string, us *v1.Upstream) {
	log.Debugf("swagger service detected: %v", url)
	if us.Metadata == nil {
		us.Metadata = &v1.Metadata{}
	}
	if us.Metadata.Annotations == nil {
		us.Metadata.Annotations = make(map[string]string)
	}
	us.Metadata.Annotations[annotations.ServiceType] = transformation.ServiceTypeTransformation
	us.Metadata.Annotations[AnnotationKeySwaggerURL] = url
}
