package routeutils

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

func AppendRouteSources(route *v1.Route, newSources []*gloov1.SourceMetadata_SourceRef) {
	meta := route.GetMetadataStatic()
	if meta == nil {
		meta = &gloov1.SourceMetadata{}
	}
	sources := meta.GetSources()
	sources = append(sources, newSources...)
	route.OpaqueMetadata = &gloov1.Route_MetadataStatic{
		MetadataStatic: &gloov1.SourceMetadata{
			Sources: sources,
		},
	}
}

func SetRouteSources(route *v1.Route, sources []*gloov1.SourceMetadata_SourceRef) {
	meta := route.GetMetadataStatic()
	if meta == nil {
		meta = &gloov1.SourceMetadata{}
	}
	route.OpaqueMetadata = &gloov1.Route_MetadataStatic{
		MetadataStatic: &gloov1.SourceMetadata{
			Sources: sources,
		},
	}
}
