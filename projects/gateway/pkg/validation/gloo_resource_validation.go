package validation

import (
	"context"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type GlooValidation interface {
	CreateDeleteRequest(ctx context.Context, ref *core.ResourceRef) *validation.GlooValidationServiceRequest
	CreateModifiedRequest(ctx context.Context, resource resources.HashableInputResource) *validation.GlooValidationServiceRequest
}

var _ GlooValidation = SecretValidation{}
var _ GlooValidation = UpstreamValidation{}

type SecretValidation struct{}

func (sv SecretValidation) CreateDeleteRequest(ctx context.Context, ref *core.ResourceRef) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_DeletedResources{
			DeletedResources: &validation.DeletedResources{
				SecretRefs: []*core.ResourceRef{ref},
			},
		},
	}
}

func (sv SecretValidation) CreateModifiedRequest(ctx context.Context, resource resources.HashableInputResource) *validation.GlooValidationServiceRequest {
	// NOT IMPLEMENTED
	return nil
}

type UpstreamValidation struct{}

func (sv UpstreamValidation) CreateDeleteRequest(ctx context.Context, ref *core.ResourceRef) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_DeletedResources{
			DeletedResources: &validation.DeletedResources{
				UpstreamRefs: []*core.ResourceRef{ref},
			},
		},
	}
}

// TODO- fix the HashabledInputResource type.  so that it works with Secrets too.
func (sv UpstreamValidation) CreateModifiedRequest(ctx context.Context, resource resources.HashableInputResource) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_ModifiedResources{
			ModifiedResources: &validation.ModifiedResources{
				Upstreams: []*gloov1.Upstream{resource.(*gloov1.Upstream)},
			},
		},
	}
}
