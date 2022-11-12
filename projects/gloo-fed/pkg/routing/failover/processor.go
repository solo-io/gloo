package failover

import (
	"context"
	"fmt"

	"github.com/rotisserie/eris"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/ezkube"
	gloo_api_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	v1sets "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/sets"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/fields"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	EmptyPrimaryTargetError   = eris.New("Primary target cannot be nil")
	EmptyFailoverTargetsError = eris.New("Failover groups must contain at least one entry")

	PrimaryTargetAlreadyInUseError = func(us ezkube.ClusterResourceId, fs *fedv1.FailoverScheme) error {
		return eris.Errorf(
			"Upstream %s.%s on cluster %s is already the primary target on FailoverScheme %s.%s",
			us.GetName(),
			us.GetNamespace(),
			us.GetClusterName(),
			fs.GetName(),
			fs.GetNamespace(),
		)
	}
	GlooInstanceError = func(glooInstances int, usRef ezkube.ClusterResourceId) error {
		return eris.Errorf(
			"Expected 1 Gloo Instance found %d which can be associated with the upstream %s.%s on cluster %s",
			glooInstances,
			usRef.GetName(),
			usRef.GetNamespace(),
			usRef.GetClusterName(),
		)
	}

	NoEndpointsError = func(glooInstance *fedv1.GlooInstance, proxy *fed_types.GlooInstanceSpec_Proxy) error {
		return eris.Errorf(
			"No available endpoints on Gloo Instance located on cluster: %s, in namespace: %s for proxy %s.%s",
			glooInstance.Spec.GetCluster(),
			glooInstance.Spec.GetControlPlane().GetNamespace(),
			proxy.GetName(),
			proxy.GetNamespace(),
		)
	}

	NoFailoverPortError = func(glooInstance *fedv1.GlooInstance, proxy *fed_types.GlooInstanceSpec_Proxy) error {
		return eris.Errorf(
			"No failover port (%d) found on Gloo Instance located on cluster: %s, in namespace: %s for proxy %s.%s",
			PortNumber,
			glooInstance.Spec.GetCluster(),
			glooInstance.Spec.GetControlPlane().GetNamespace(),
			proxy.GetName(),
			proxy.GetNamespace(),
		)
	}
)

func NewFailoverProcessor(
	glooClientset gloov1.MulticlusterClientset,
	glooInstanceClient fedv1.GlooInstanceClient,
	failoverSchemeClient fedv1.FailoverSchemeClient,
	statusManager *StatusManager,
) FailoverProcessor {
	return &failoverProcessorImpl{
		glooClientset:        glooClientset,
		glooInstanceClient:   glooInstanceClient,
		failoverSchemeClient: failoverSchemeClient,
		statusManager:        statusManager,
	}
}

type failoverProcessorImpl struct {
	glooClientset        gloov1.MulticlusterClientset
	glooInstanceClient   fedv1.GlooInstanceClient
	failoverSchemeClient fedv1.FailoverSchemeClient
	statusManager        *StatusManager
}

func (f failoverProcessorImpl) ProcessFailoverUpdate(
	ctx context.Context,
	obj *fedv1.FailoverScheme,
) (*gloov1.Upstream, StatusBuilder) {
	statusBuilder := f.statusManager.NewStatusBuilder(obj)
	if obj.Spec.GetPrimary() == nil {
		return nil, statusBuilder.Invalidate(EmptyPrimaryTargetError)
	}
	primaryClusterClient, err := f.glooClientset.Cluster(obj.Spec.GetPrimary().GetClusterName())
	if err != nil {
		return nil, statusBuilder.Fail(err)
	}

	primaryUpstream, err := primaryClusterClient.Upstreams().GetUpstream(ctx, client.ObjectKey{
		Namespace: obj.Spec.GetPrimary().GetNamespace(),
		Name:      obj.Spec.GetPrimary().GetName(),
	})
	if err != nil {
		// If the primary upstream cannot be found, the object is invalid and cannot be retried.
		if errors.IsNotFound(err) {
			return nil, statusBuilder.Invalidate(err)
		}
		return nil, statusBuilder.Fail(err)
	}

	failoverSchemeList, err := f.failoverSchemeClient.ListFailoverScheme(ctx)
	if err != nil {
		return nil, statusBuilder.Fail(err)
	}
	// For each existing failover scheme, check if the primary upstream is already a primary target.
	// Make sure to not check the scheme which is currently being reconciled.
	for idx := range failoverSchemeList.Items {
		failoverScheme := failoverSchemeList.Items[idx]
		if failoverScheme.Spec.GetPrimary().Equal(obj.Spec.GetPrimary()) &&
			sets.Key(&failoverScheme) != sets.Key(obj) {
			return nil, statusBuilder.
				Invalidate(PrimaryTargetAlreadyInUseError(obj.Spec.GetPrimary(), &failoverScheme))
		}

	}

	if len(obj.Spec.GetFailoverGroups()) == 0 {
		return nil, statusBuilder.Invalidate(EmptyFailoverTargetsError)
	}

	failoverCfg, updatedStatusBuilder := f.buildFailoverConfig(ctx, obj, statusBuilder)
	if updatedStatusBuilder != nil {
		return nil, updatedStatusBuilder
	}

	primaryUpstream.Spec.Failover = failoverCfg
	return primaryUpstream, nil
}

func (f failoverProcessorImpl) buildFailoverConfig(
	ctx context.Context,
	obj *fedv1.FailoverScheme,
	statusBuilder StatusBuilder,
) (*gloo_api_v1.Failover, StatusBuilder) {
	failoverCfg := &gloo_api_v1.Failover{}
	for _, priority := range obj.Spec.GetFailoverGroups() {
		glooPriority := &gloo_api_v1.Failover_PrioritizedLocality{}
		for _, groupMember := range priority.GetPriorityGroup() {
			glooLocality := &gloo_api_v1.LocalityLbEndpoints{
				LoadBalancingWeight: groupMember.GetLocalityWeight(),
				Locality:            &gloo_api_v1.Locality{},
			}
			instances, err := f.glooInstanceClient.ListGlooInstance(
				ctx,
				fields.BuildClusterFieldMatcher(groupMember.GetCluster()),
			)
			if err != nil {
				return nil, statusBuilder.Fail(err)
			}
			instanceSet := v1sets.NewGlooInstanceSetFromList(instances)

			glooClusterClient, err := f.glooClientset.Cluster(groupMember.GetCluster())
			if err != nil {
				return nil, statusBuilder.Fail(err)
			}

			for _, usRef := range groupMember.GetUpstreams() {
				us, err := glooClusterClient.Upstreams().GetUpstream(ctx, client.ObjectKey{
					Namespace: usRef.GetNamespace(),
					Name:      usRef.GetName(),
				})
				if err != nil {
					if errors.IsNotFound(err) {
						return nil, statusBuilder.Invalidate(err)
					}
					return nil, statusBuilder.Fail(err)
				}
				usLocality, err := f.computeEndpoints(&skv2v1.ClusterObjectRef{
					Name:        us.GetName(),
					Namespace:   us.GetNamespace(),
					ClusterName: groupMember.GetCluster(),
				}, instanceSet)
				if err != nil {
					return nil, statusBuilder.Invalidate(err)
				}
				// If usLocality has zone (and/or) region replace them.
				if usLocality.GetLocality().GetZone() != "" {
					glooLocality.GetLocality().Zone = usLocality.GetLocality().GetZone()
				}
				if usLocality.GetLocality().GetRegion() != "" {
					glooLocality.GetLocality().Region = usLocality.GetLocality().GetRegion()
				}
				// append endoints
				glooLocality.LbEndpoints = append(glooLocality.LbEndpoints, usLocality.LbEndpoints...)

			}

			glooPriority.LocalityEndpoints = append(glooPriority.LocalityEndpoints, glooLocality)
		}
		failoverCfg.PrioritizedLocalities = append(failoverCfg.PrioritizedLocalities, glooPriority)
	}
	return failoverCfg, nil
}

func (f *failoverProcessorImpl) computeEndpoints(
	usRef *skv2v1.ClusterObjectRef,
	instanceSet v1sets.GlooInstanceSet,
) (*gloo_api_v1.LocalityLbEndpoints, error) {
	glooLocality := &gloo_api_v1.LocalityLbEndpoints{}
	// continuously merge with the gateway instance set to build a list of possibe gateway instances
	usInstanceSet := GetGlooInstanceForUpstream(usRef, instanceSet)
	if usInstanceSet.Length() != 1 {
		return nil, GlooInstanceError(usInstanceSet.Length(), usRef)
	}
	instance := usInstanceSet.List()[0]

	proxy := discovery.GetAdminProxyForInstance(instance)
	writeNamespace := instance.Spec.GetAdmin().GetWriteNamespace()
	// It is fine to set this every time. Every instance on a cluster must have identical regions
	if instance.Spec.GetRegion() != "" {
		glooLocality.Locality = &gloo_api_v1.Locality{
			Region: instance.Spec.GetRegion(),
		}
	}

	// For now just default to first zone in the list
	if len(proxy.GetZones()) > 0 {
		glooLocality.Locality.Zone = proxy.GetZones()[0]
	}

	if len(proxy.GetIngressEndpoints()) == 0 {
		return nil, NoEndpointsError(instance, proxy)
	}

	// Search all discovered endpoints
	endpoint, port := f.getEndpointAndPort(proxy)
	if endpoint == nil {
		return nil, NoFailoverPortError(instance, proxy)
	}
	// Append the LbEndpoint for the Upstream
	glooLocality.LbEndpoints = append(glooLocality.LbEndpoints, &gloo_api_v1.LbEndpoint{
		Address: endpoint.GetAddress(),
		Port:    port,
		UpstreamSslConfig: &gloo_api_v1.UpstreamSslConfig{
			SslSecrets: &gloo_api_v1.UpstreamSslConfig_SecretRef{
				// TODO: Allow configuration of Upstream/Downstream ssl secrets
				SecretRef: &core.ResourceRef{
					Name:      UpstreamSecretName,
					Namespace: writeNamespace,
				},
			},
			// Set the SNI name to the gloo translated envoy Cluster name of the Upstream. This allows the request to
			// get forwarded properly by the Upstream Gloo using the forward_sni_cluster_name feature.
			// https://github.com/solo-io/gloo/blob/a956a85f339259e561ec76abaeed6036820de209/projects/gloo/api/v1/proxy.proto#L137
			Sni: UpstreamToClusterName(&skv2v1.ObjectRef{
				Name:      usRef.GetName(),
				Namespace: usRef.GetNamespace(),
			}),
			// TODO: Configure SAN verification on Upstream/Downstream
			VerifySubjectAltName: nil,
		},
	})
	return glooLocality, nil
}

func (f *failoverProcessorImpl) getEndpointAndPort(
	proxy *fed_types.GlooInstanceSpec_Proxy,
) (*fed_types.GlooInstanceSpec_Proxy_IngressEndpoint, uint32) {
	for _, ep := range proxy.GetIngressEndpoints() {
		for _, epPort := range ep.GetPorts() {
			// If a port is found with the designated name, failover, assign that endpoint and port for later use
			if epPort.GetName() == PortName {
				return ep, epPort.GetPort()
			}
		}
	}
	return nil, 0
}

func (f *failoverProcessorImpl) ProcessFailoverDelete(
	ctx context.Context,
	obj *fedv1.FailoverScheme,
) (*gloov1.Upstream, error) {
	// If the primary is nil during a deletion, it means that the object would already be in an error state, and there
	// is nothing to be done about it here, so we can simply return
	if obj.Spec.GetPrimary() == nil {
		return nil, nil
	}
	primaryClusterClient, err := f.glooClientset.Cluster(obj.Spec.GetPrimary().GetClusterName())
	if err != nil {
		return nil, err
	}

	primaryUpstream, err := primaryClusterClient.Upstreams().GetUpstream(ctx, client.ObjectKey{
		Namespace: obj.Spec.GetPrimary().GetNamespace(),
		Name:      obj.Spec.GetPrimary().GetName(),
	})
	// If the upstream has already been removed, no need to keep trying here
	if err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	primaryUpstream.Spec.Failover = nil
	return primaryUpstream, nil
}

// GetGlooInstanceForUpstream returns all gloo instances that an Upstream belongs to, this should only be 1
func GetGlooInstanceForUpstream(us ezkube.ClusterResourceId, instances v1sets.GlooInstanceSet) v1sets.GlooInstanceSet {
	result := v1sets.NewGlooInstanceSet()
	for _, instance := range instances.List() {
		if IsGlooInstanceUpstream(instance, us) {
			result.Insert(instance)
		}
	}
	return result
}

// UpstreamToClusterName returns the name of the cluster created for a given upstream
// Copy paste from https://github.com/solo-io/gloo/blob/a956a85f339259e561ec76abaeed6036820de209/projects/gloo/pkg/translator/utils.go#L19
func UpstreamToClusterName(upstream *skv2v1.ObjectRef) string {

	// For non-namespaced resources, return only name
	if upstream.GetNamespace() == "" {
		return upstream.GetName()
	}

	// Don't use dots in the name as it messes up prometheus stats
	return fmt.Sprintf("%s_%s", upstream.GetName(), upstream.GetNamespace())
}
