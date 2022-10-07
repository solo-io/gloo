package validation

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GlooValidation interface {
	CreateDeleteRequest(ref *core.ResourceRef) *validation.GlooValidationServiceRequest
	CreateModifiedRequest(resource resources.HashableInputResource) *validation.GlooValidationServiceRequest
}

var GvkToGlooValidator = map[schema.GroupVersionKind]func() GlooValidation{
	gloov1.UpstreamGVK: func() GlooValidation { return UpstreamValidation{} },
	gloov1.SecretGVK:   func() GlooValidation { return SecretValidation{} },
}

var _ GlooValidation = SecretValidation{}
var _ GlooValidation = UpstreamValidation{}

type SecretValidation struct{}

func (sv SecretValidation) CreateDeleteRequest(ref *core.ResourceRef) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_DeletedResources{
			DeletedResources: &validation.DeletedResources{
				SecretRefs: []*core.ResourceRef{ref},
			},
		},
	}
}

func (sv SecretValidation) CreateModifiedRequest(resource resources.HashableInputResource) *validation.GlooValidationServiceRequest {
	// NOT IMPLEMENTED
	return nil
}

type UpstreamValidation struct{}

func (sv UpstreamValidation) CreateDeleteRequest(ref *core.ResourceRef) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_DeletedResources{
			DeletedResources: &validation.DeletedResources{
				UpstreamRefs: []*core.ResourceRef{ref},
			},
		},
	}
}

// TODO- fix the HashabledInputResource type.  so that it works with Secrets too.
func (sv UpstreamValidation) CreateModifiedRequest(resource resources.HashableInputResource) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_ModifiedResources{
			ModifiedResources: &validation.ModifiedResources{
				Upstreams: []*gloov1.Upstream{resource.(*gloov1.Upstream)},
			},
		},
	}
}
