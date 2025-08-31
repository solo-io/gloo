package serviceconverter

import (
	"context"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	corev1 "k8s.io/api/core/v1"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

const GlooDnsIpFamilyAnnotation = "gloo.solo.io/dns_lookup_ip_family"

var ipFamilies = map[string]v1.DnsIpFamily{
	"v4":      v1.DnsIpFamily_V4_ONLY,
	"v4_pref": v1.DnsIpFamily_V4_PREFERRED,
	"v6":      v1.DnsIpFamily_V6_ONLY,
	"v6_pref": v1.DnsIpFamily_V6_PREFERRED,
	"dual":    v1.DnsIpFamily_DUAL_IP_FAMILY,
}

// DnsIpFamilyConverter sets the upstream ip family for DNS lookups:
// (1) the service has the "dns_ip_family" annotation; or
// (2) the global "dns_ip_family" annotation defined in Settings.UpstreamOptions; or
type DnsIpFamilyConverter struct{}

func (u *DnsIpFamilyConverter) ConvertService(ctx context.Context, svc *corev1.Service, port corev1.ServicePort, us *v1.Upstream) error {
	us.DnsLookupIpFamily = processIpFamily(ctx, svc)
	return nil
}

func processIpFamily(ctx context.Context, svc *corev1.Service) v1.DnsIpFamily {
	if svc.Annotations != nil {
		return lookupDnsIpFamily(svc.Annotations[GlooDnsIpFamilyAnnotation])
	}
	if globalAnnotations := settingsutil.MaybeFromContext(ctx).GetUpstreamOptions().GetGlobalAnnotations(); globalAnnotations != nil {
		return lookupDnsIpFamily(globalAnnotations[GlooDnsIpFamilyAnnotation])
	}
	return v1.DnsIpFamily_DEFAULT
}

func lookupDnsIpFamily(family string) v1.DnsIpFamily {
	if value, ok := ipFamilies[family]; ok {
		return value
	}
	return v1.DnsIpFamily_DEFAULT
}
