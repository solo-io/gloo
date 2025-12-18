package vhostutils

import (
	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AppendSourceToVirtualHost appends a source VirtualHostOption to the VirtualHost's metadata static
func AppendSourceToVirtualHost(vh *v1.VirtualHost, source client.Object) {
	meta := vh.GetMetadataStatic()
	if meta == nil {
		meta = &v1.SourceMetadata{}
	}
	sources := meta.GetSources()
	sources = append(sources, &v1.SourceMetadata_SourceRef{
		ResourceRef: &core.ResourceRef{
			Name:      source.GetName(),
			Namespace: source.GetNamespace(),
		},
		ResourceKind:       sologatewayv1.VirtualHostOptionGVK.Kind,
		ObservedGeneration: source.GetGeneration(),
	})
	vh.OpaqueMetadata = &v1.VirtualHost_MetadataStatic{
		MetadataStatic: &v1.SourceMetadata{
			Sources: sources,
		},
	}
}
