package proxy_syncer

import (
	"testing"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// vhoProxyReport builds a ProxyReport carrying a single VirtualHost error sourced by a
// VirtualHostOption, which is the shape the VirtualHostOption status plugin reads from.
func vhoProxyReport(reason string) *validation.ProxyReport {
	return &validation.ProxyReport{
		ListenerReports: []*validation.ListenerReport{{
			ListenerTypeReport: &validation.ListenerReport_AggregateListenerReport{
				AggregateListenerReport: &validation.AggregateListenerReport{
					HttpListenerReports: map[string]*validation.HttpListenerReport{
						"http": {
							VirtualHostReports: []*validation.VirtualHostReport{{
								Errors: []*validation.VirtualHostReport_Error{{
									Type:   validation.VirtualHostReport_Error_ProcessingError,
									Reason: reason,
									Metadata: &gloov1.SourceMetadata{
										Sources: []*gloov1.SourceMetadata_SourceRef{{
											ResourceKind: sologatewayv1.VirtualHostOptionGVK.Kind,
											ResourceRef: &core.ResourceRef{
												Name:      "bad-retries",
												Namespace: "default",
											},
										}},
									},
								}},
							}},
						},
					},
				},
			},
		}},
	}
}

func wrapperWith(snapVersion string, report *validation.ProxyReport) XdsSnapWrapper {
	return XdsSnapWrapper{
		snap: &xds.EnvoySnapshot{
			Clusters: cache.Resources{Version: snapVersion},
		},
		proxyKey: "gloo-system~gw-1",
		proxyWithReport: translatorutils.ProxyWithReports{
			Reports: translatorutils.TranslationReports{ProxyReport: report},
		},
	}
}

func TestXdsSnapWrapperEquals(t *testing.T) {
	tests := []struct {
		name  string
		a     XdsSnapWrapper
		b     XdsSnapWrapper
		equal bool
	}{
		{
			name:  "identical snapshot and report",
			a:     wrapperWith("v1", vhoProxyReport("invalid retry policy")),
			b:     wrapperWith("v1", vhoProxyReport("invalid retry policy")),
			equal: true,
		},
		{
			name:  "both reports empty",
			a:     wrapperWith("v1", nil),
			b:     wrapperWith("v1", nil),
			equal: true,
		},
		{
			// Rejecting an invalid option leaves the previous config in place, so only the
			// report changes. Treating this as unchanged is what left the rejected resource
			// without a status.
			name:  "identical snapshot, report gained an error",
			a:     wrapperWith("v1", nil),
			b:     wrapperWith("v1", vhoProxyReport("invalid retry policy")),
			equal: false,
		},
		{
			name:  "identical snapshot, error reason changed",
			a:     wrapperWith("v1", vhoProxyReport("invalid retry policy")),
			b:     wrapperWith("v1", vhoProxyReport("some other problem")),
			equal: false,
		},
		{
			name:  "snapshot changed",
			a:     wrapperWith("v1", nil),
			b:     wrapperWith("v2", nil),
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equals(tt.b); got != tt.equal {
				t.Errorf("Equals() = %v, want %v", got, tt.equal)
			}
			// Equality must not depend on argument order.
			if got := tt.b.Equals(tt.a); got != tt.equal {
				t.Errorf("reversed Equals() = %v, want %v", got, tt.equal)
			}
		})
	}
}
