package utils

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rltypes "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/retries"
)

var _ = Describe("Merge", func() {
	It("merges top-level route options fields", func() {
		dst := &v1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{Value: "preserve-me"},
			Retries: &retries.RetryPolicy{
				RetryOn:    "5XX",
				NumRetries: 0, // should not overwrite this field
			},
		}
		d := prototime.DurationToProto(time.Minute)
		src := &v1.RouteOptions{
			Timeout: d,
			Retries: &retries.RetryPolicy{
				RetryOn:    "5XX",
				NumRetries: 3, // do not overwrite 0 value above
				// do not merge PerTryTimeout because it doesn't exist in dst
				PerTryTimeout: d,
			},
			RatelimitBasic: &ratelimit.IngressRateLimit{
				AuthorizedLimits: &rltypes.RateLimit{
					Unit:            1,
					RequestsPerUnit: 2,
				},
			},
		}
		expected := &v1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{Value: "preserve-me"},
			Timeout:       d,
			Retries: &retries.RetryPolicy{
				RetryOn:    "5XX",
				NumRetries: 0,
			},
			RatelimitBasic: &ratelimit.IngressRateLimit{
				AuthorizedLimits: &rltypes.RateLimit{
					Unit:            1,
					RequestsPerUnit: 2,
				},
			},
		}

		actual, overwrote := ShallowMergeRouteOptions(dst, src)
		Expect(actual).To(Equal(expected))
		Expect(overwrote).To(BeTrue())
	})

	Describe("ShallowCopyRouteOptions", func() {
		It("returns nil for a nil source", func() {
			Expect(ShallowCopyRouteOptions(nil)).To(BeNil())
		})

		It("copies top-level fields by value without deep-copying sub-messages", func() {
			src := &v1.RouteOptions{
				PrefixRewrite: &wrappers.StringValue{Value: "rewrite-me"},
				Retries: &retries.RetryPolicy{
					RetryOn:    "5XX",
					NumRetries: 3,
				},
			}

			out := ShallowCopyRouteOptions(src)

			// The copy is a distinct top-level message that compares equal by value.
			Expect(out).NotTo(BeIdenticalTo(src))
			Expect(out).To(Equal(src))

			// Sub-messages are shared by pointer rather than deep-cloned: this is the
			// allocation saving that keeps translation heap bounded when many routes
			// reference the same RouteOption.
			Expect(out.GetRetries()).To(BeIdenticalTo(src.GetRetries()))
			Expect(out.GetPrefixRewrite()).To(BeIdenticalTo(src.GetPrefixRewrite()))
		})

		It("isolates top-level field reassignment on the copy from the source", func() {
			src := &v1.RouteOptions{
				PrefixRewrite: &wrappers.StringValue{Value: "original"},
			}

			out := ShallowCopyRouteOptions(src)
			// Route plugins reassign top-level fields on the merged options; this must not
			// leak back into the shared source RouteOption.
			out.PrefixRewrite = &wrappers.StringValue{Value: "changed"}

			Expect(src.GetPrefixRewrite().GetValue()).To(Equal("original"))
		})
	})
})
