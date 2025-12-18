package listenerutils

import (
	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SetListenerSources(listener *v1.Listener, sources []*v1.SourceMetadata_SourceRef) {
	meta := listener.GetMetadataStatic()
	if meta == nil {
		meta = &v1.SourceMetadata{}
	}
	listener.OpaqueMetadata = &v1.Listener_MetadataStatic{
		MetadataStatic: &v1.SourceMetadata{
			Sources: sources,
		},
	}
}

// AppendSourceToListener appends a source ListenerOption to the Listener's metadata static
func AppendSourceToListener(listener *v1.Listener, source client.Object) {
	meta := listener.GetMetadataStatic()
	if meta == nil {
		meta = &v1.SourceMetadata{}
	}
	sources := meta.GetSources()
	sources = append(sources, &v1.SourceMetadata_SourceRef{
		ResourceRef: &core.ResourceRef{
			Name:      source.GetName(),
			Namespace: source.GetNamespace(),
		},
		ResourceKind:       sologatewayv1.ListenerOptionGVK.Kind,
		ObservedGeneration: source.GetGeneration(),
	})
	listener.OpaqueMetadata = &v1.Listener_MetadataStatic{
		MetadataStatic: &v1.SourceMetadata{
			Sources: sources,
		},
	}
}
