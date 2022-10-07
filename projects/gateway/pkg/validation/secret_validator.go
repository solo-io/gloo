package validation

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ DeleteGlooResourceValidator = SecretValidator{}

type SecretValidator struct{}

func (sv SecretValidator) CreateDeleteRequest(ref *core.ResourceRef) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_DeletedResources{
			DeletedResources: &validation.DeletedResources{
				SecretRefs: []*core.ResourceRef{ref},
			},
		},
	}
}
