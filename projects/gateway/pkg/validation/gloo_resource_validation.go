package validation

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

/*
	Currently having issues with the following resources as they are not Input Resources
	// Artifacts          gloo_solo_io.ArtifactList
	// Endpoints          gloo_solo_io.EndpointList
	// Secrets            gloo_solo_io.SecretList
	// Ratelimitconfigs   github_com_solo_io_gloo_projects_gloo_pkg_api_external_solo_ratelimit.RateLimitConfigList

	// will have to check out RateLimitConfigs, not sure what causes it to not be a Hashable Input Resource
	it is a custom resource, but not sure what is causing this...
	projects/gloo/api/external/solo/ratelimit/solo-kit.json
*/

// DeleteGlooResourceValidator allows validation of deletion of gloo.solo.io resources
type DeleteGlooResourceValidator interface {
	CreateDeleteRequest(ref *core.ResourceRef) *validation.GlooValidationServiceRequest
}

// DeleteGlooResourceValidator allows validation of gloo.solo.io resources
type GlooResourceValidator interface {
	CreateModifiedRequest(resource resources.HashableInputResource) *validation.GlooValidationServiceRequest
}

var GvkToGlooValidator = map[schema.GroupVersionKind]GlooResourceValidator{
	gloov1.UpstreamGVK: UpstreamValidator{},
}

var GvkToDeleteGlooValidator = map[schema.GroupVersionKind]DeleteGlooResourceValidator{
	gloov1.UpstreamGVK: UpstreamValidator{},
	gloov1.SecretGVK:   SecretValidator{},
}

var _ GlooResourceValidator = UpstreamValidator{}
var _ DeleteGlooResourceValidator = UpstreamValidator{}
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

type UpstreamValidator struct{}

func (sv UpstreamValidator) CreateDeleteRequest(ref *core.ResourceRef) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_DeletedResources{
			DeletedResources: &validation.DeletedResources{
				UpstreamRefs: []*core.ResourceRef{ref},
			},
		},
	}
}

// TODO- fix the HashabledInputResource type.  so that it works with Secrets too.
func (sv UpstreamValidator) CreateModifiedRequest(resource resources.HashableInputResource) *validation.GlooValidationServiceRequest {
	return &validation.GlooValidationServiceRequest{
		Resources: &validation.GlooValidationServiceRequest_ModifiedResources{
			ModifiedResources: &validation.ModifiedResources{
				Upstreams: []*gloov1.Upstream{resource.(*gloov1.Upstream)},
			},
		},
	}
}
