package validation

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
)

// GlooValidationRequest is sent to the validation form the validator.  This is replacing the
// use of what used to be a request sent to the Gloo service, when Gloo and Gateway where
// both different pods.  This contains a modified and deleted set of resources contained in the api snapshot.
type GlooValidationRequest struct {
	// ModifiedSnapshot contains the list of resources that are undergoing modifications
	ModifiedSnapshot gloov1snap.ApiSnapshot
	// DeletedSnapshot contains the list of resources that need to be deleted
	DeletedResources gloov1snap.ApiSnapshot
	// Proxy is the current proxy.  If nil it represents a single resource
	Proxy *gloov1.Proxy
}
