package v1

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpstreamToCRD(meta metav1.ObjectMeta, us v1.Upstream) Upstream {
	return Upstream{
		ObjectMeta: meta,
		Spec:       DeepCopyUpstream(us),
	}
}

func UpstreamFromCRD(crd *Upstream) v1.Upstream {
	return v1.Upstream(crd.Spec)
}

func VirtualHostToCRD(meta metav1.ObjectMeta, us v1.VirtualHost) VirtualHost {
	return VirtualHost{
		ObjectMeta: meta,
		Spec:       DeepCopyVirtualHost(us),
	}
}

func VirtualHostFromCRD(crd *VirtualHost) v1.VirtualHost {
	return v1.VirtualHost(crd.Spec)
}
