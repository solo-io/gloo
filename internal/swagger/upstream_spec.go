package swagger

import (
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/common/annotations"
)

const ServiceTypeSwagger = "swagger"

const (
	AnnotationKeySwaggerURL = "gloo.solo.io/swagger_url"
	AnnotationKeySwaggerDoc = "gloo.solo.io/swagger_doc"
)

// TODO: create service spec on upstreams themselves
// this is needed for NATS, etc. various service types
// that can be a subclass of an upstream type
type Annotations struct {
	SwaggerURL       string
	InlineSwaggerDoc string
}

// TODO: discover & set this annotation key on upstreams by checking for user-provided & common swagger urls
func GetSwaggerAnnotations(us *v1.Upstream) (*Annotations, error) {
	swaggerUrl, urlOk := us.Metadata.Annotations[AnnotationKeySwaggerURL]
	swaggerDoc, docOk := us.Metadata.Annotations[AnnotationKeySwaggerDoc]
	if !urlOk && !docOk {
		return nil, errors.Errorf("one of %v or %v must be set in the annotation for a swagger upstream", AnnotationKeySwaggerURL, AnnotationKeySwaggerDoc)
	}
	return &Annotations{
		SwaggerURL:       swaggerUrl,
		InlineSwaggerDoc: swaggerDoc,
	}, nil
}

func IsSwagger(us *v1.Upstream) bool {
	return us.Metadata.Annotations[annotations.ServiceType] == ServiceTypeSwagger
}
