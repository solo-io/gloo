package listenerutils

import (
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

// AppendSourceToListener appends a source object to the Listener's static metadata
func AppendSourceToListener(listener *v1.Listener, source client.Object, resourceKind string) {
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
		ResourceKind:       resourceKind,
		ObservedGeneration: source.GetGeneration(),
	})
	listener.OpaqueMetadata = &v1.Listener_MetadataStatic{
		MetadataStatic: &v1.SourceMetadata{
			Sources: sources,
		},
	}
}
