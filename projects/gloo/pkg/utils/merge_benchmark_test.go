package utils_test

import (
	"fmt"
	"testing"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

// These benchmarks compare the per-route cost of the two ways a RouteOption attachment can seed
// the merged options during gateway2 translation (see GetRouteOptionForRouteRule):
//
//   - Clone: the pre-#8802 behavior, deep-copying the entire options tree for every route rule on
//     every translation. For transformation-heavy RouteOptions shared by many routes this
//     dominated translation heap (~31% of a 14GB heap for one customer).
//   - ShallowCopyRouteOptions: the current behavior, allocating one top-level message per route
//     and sharing the sub-messages.
//
// Proto deep-copy cost scales with the number of messages and map/slice entries (string bytes are
// shared), so the fixture carries many header templates and dynamic metadata values like the
// customer templates in solo-io/solo-projects#8802.

func BenchmarkRouteOptionsDeepClone(b *testing.B) {
	src := largeTransformationRouteOptions()
	b.ReportAllocs()
	for b.Loop() {
		if out := src.Clone().(*v1.RouteOptions); out == nil {
			b.Fatal("expected clone")
		}
	}
}

func BenchmarkShallowCopyRouteOptions(b *testing.B) {
	src := largeTransformationRouteOptions()
	b.ReportAllocs()
	for b.Loop() {
		if out := utils.ShallowCopyRouteOptions(src); out == nil {
			b.Fatal("expected copy")
		}
	}
}

func largeTransformationRouteOptions() *v1.RouteOptions {
	headers := make(map[string]*transformation.InjaTemplate, 1000)
	for i := range 1000 {
		headers[fmt.Sprintf("x-generated-header-%d", i)] = &transformation.InjaTemplate{
			Text: fmt.Sprintf(`{{ request_header("x-input-%d") }}`, i),
		}
	}
	metadataValues := make([]*transformation.TransformationTemplate_DynamicMetadataValue, 0, 500)
	for i := range 500 {
		metadataValues = append(metadataValues, &transformation.TransformationTemplate_DynamicMetadataValue{
			MetadataNamespace: "io.solo.benchmark",
			Key:               fmt.Sprintf("key-%d", i),
			Value:             &transformation.InjaTemplate{Text: fmt.Sprintf("value-%d", i)},
		})
	}
	return &v1.RouteOptions{
		StagedTransformations: &transformation.TransformationStages{
			Regular: &transformation.RequestResponseTransformations{
				RequestTransforms: []*transformation.RequestMatch{
					{
						RequestTransformation: &transformation.Transformation{
							TransformationType: &transformation.Transformation_TransformationTemplate{
								TransformationTemplate: &transformation.TransformationTemplate{
									Headers:               headers,
									DynamicMetadataValues: metadataValues,
								},
							},
						},
					},
				},
			},
		},
	}
}
