package failover

import (
	"context"

	"github.com/solo-io/go-utils/stringutils"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	v1sets "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewFailoverDependencyCalculator(
	failoverSchemeClient fedv1.FailoverSchemeClient,
	glooInstanceClient fedv1.GlooInstanceClient,
) FailoverDependencyCalculator {
	return &failoverDependencyCalculatorImpl{
		failoverSchemeClient: failoverSchemeClient,
		glooInstanceClient:   glooInstanceClient,
	}
}

type failoverDependencyCalculatorImpl struct {
	failoverSchemeClient fedv1.FailoverSchemeClient
	glooInstanceClient   fedv1.GlooInstanceClient
}

func (f *failoverDependencyCalculatorImpl) ForUpstream(
	ctx context.Context,
	upstream *skv2v1.ClusterObjectRef,
) ([]*fedv1.FailoverScheme, error) {
	failoverSchemeList, err := f.failoverSchemeClient.ListFailoverScheme(ctx)
	if err != nil {
		return nil, err
	}
	var result []*fedv1.FailoverScheme
	for idx := range failoverSchemeList.Items {
		scheme := &failoverSchemeList.Items[idx]
		for _, usRef := range f.getUpstreams(scheme) {
			if usRef.Equal(upstream) {
				result = append(result, scheme)
				break
			}
		}
	}
	return result, nil
}

func (f *failoverDependencyCalculatorImpl) ForGlooInstance(
	ctx context.Context,
	glooInstanceRef *skv2v1.ObjectRef,
) ([]*fedv1.FailoverScheme, error) {
	failoverSchemeList, err := f.failoverSchemeClient.ListFailoverScheme(ctx)
	if err != nil {
		return nil, err
	}
	glooInstance, err := f.glooInstanceClient.GetGlooInstance(ctx, client.ObjectKey{
		Namespace: glooInstanceRef.GetNamespace(),
		Name:      glooInstanceRef.GetName(),
	})
	if err != nil {
		// if the gloo instance cannot be found than return all failover schemes as this means the resource has been
		// deleted, and we can't know which ones need to be updated.
		return v1sets.NewFailoverSchemeSetFromList(failoverSchemeList).List(), client.IgnoreNotFound(err)
	}
	var result []*fedv1.FailoverScheme
	for idx := range failoverSchemeList.Items {
		scheme := &failoverSchemeList.Items[idx]
		for _, usRef := range f.getUpstreams(scheme) {
			// Only add a gloo instance to the list, if one of the upstreams "belongs" to it
			if IsGlooInstanceUpstream(glooInstance, usRef) {
				result = append(result, scheme)
				break
			}
		}
	}
	return result, nil
}

func (f *failoverDependencyCalculatorImpl) getUpstreams(
	failoverScheme *fedv1.FailoverScheme,
) []*skv2v1.ClusterObjectRef {
	// Always add primary upstream
	result := []*skv2v1.ClusterObjectRef{failoverScheme.Spec.GetPrimary()}
	// Iterate through all failover targets, and priority groups, to determine all possible upstreams
	for _, priorityGroup := range failoverScheme.Spec.GetFailoverGroups() {
		for _, groupMember := range priorityGroup.GetPriorityGroup() {
			for _, usRef := range groupMember.GetUpstreams() {
				result = append(result, &skv2v1.ClusterObjectRef{
					Name:        usRef.GetName(),
					Namespace:   usRef.GetNamespace(),
					ClusterName: groupMember.GetCluster(),
				})
			}
		}
	}
	return result
}

// IsGlooInstanceUpstream determines whether a particular upstream belongs to the passed in gloo instance
func IsGlooInstanceUpstream(instance *fedv1.GlooInstance, usRef *skv2v1.ClusterObjectRef) bool {
	if instance.Spec.GetCluster() == usRef.GetClusterName() {
		switch {

		case len(instance.Spec.GetControlPlane().GetWatchedNamespaces()) == 0:
			return true
		case stringutils.ContainsString(usRef.GetNamespace(), instance.Spec.GetControlPlane().GetWatchedNamespaces()):
			return true
		}
	}
	return false
}
