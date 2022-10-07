package validation

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

/*
	Resources that are implemented in Validation must implement the DeleteGlooResourceValidator and/or GlooResourceValidator interfaces.
	These are used to handle the specific logic used in validation of gateway resources.
	Have to add  to GVK maps once supported GvkToGlooValidator and GvkToDeleteGlooValidator.
*/

const GlooGroup = "gloo.solo.io"

// DeleteGlooResourceValidator allows validation of deletion of gloo.solo.io resources
type DeleteGlooResourceValidator interface {
	CreateDeleteRequest(ref *core.ResourceRef) *validation.GlooValidationServiceRequest
}

// DeleteGlooResourceValidator allows validation of gloo.solo.io resources
type GlooResourceValidator interface {
	CreateModifiedRequest(resource resources.Resource) *validation.GlooValidationServiceRequest
}

var GvkToGlooValidator = map[schema.GroupVersionKind]GlooResourceValidator{
	gloov1.UpstreamGVK: UpstreamValidator{},
}

var GvkToDeleteGlooValidator = map[schema.GroupVersionKind]DeleteGlooResourceValidator{
	gloov1.UpstreamGVK: UpstreamValidator{},
	gloov1.SecretGVK:   SecretValidator{},
}
