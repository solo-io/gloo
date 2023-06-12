package transformer

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	v32 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	envoytransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	xslt "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformers/xslt"
	osTransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/transformation"
)

var _ = Describe("Plugin", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		p      *plugin
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
	})

	Context("translate transformation", func() {
		BeforeEach(func() {
			p = NewPlugin()
			p.Init(plugins.InitParams{Ctx: ctx})
		})

		It("translates header body transform", func() {
			headerBodyTransform := &envoytransformation.HeaderBodyTransform{}

			input := &osTransformation.Transformation{
				TransformationType: &osTransformation.Transformation_HeaderBodyTransform{
					HeaderBodyTransform: headerBodyTransform,
				},
			}

			expectedOutput := &envoytransformation.Transformation{
				TransformationType: &envoytransformation.Transformation_HeaderBodyTransform{
					HeaderBodyTransform: headerBodyTransform,
				},
			}
			output, err := transformation.TranslateTransformation(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(expectedOutput))
		})

		Context("LogRequestResponseInfo", func() {

			var (
				expectedOutput      *envoytransformation.Transformation
				inputTransformation *osTransformation.Transformation
			)

			BeforeEach(func() {
				inputTransformation = &osTransformation.Transformation{
					TransformationType: &osTransformation.Transformation_XsltTransformation{
						XsltTransformation: &xslt.XsltTransformation{
							Xslt: "test",
						},
					},
				}

				expectedOutput = &envoytransformation.Transformation{
					TransformationType: &envoytransformation.Transformation_TransformerConfig{
						TransformerConfig: &v32.TypedExtensionConfig{
							// Arbitrary name for TypedExtension, will error if left empty
							Name: XsltTransformerFactoryName,
							TypedConfig: &anypb.Any{
								TypeUrl: "type.googleapis.com/envoy.config.transformer.xslt.v2.XsltTransformation",
								// Value:   "CgR0ZXN0",
								Value: []byte{10, 4, 116, 101, 115, 116},
							},
						},
					},
				}
			})

			It("can set log_request_response_info on transformation level", func() {
				inputTransformation.LogRequestResponseInfo = true
				expectedOutput.LogRequestResponseInfo = &wrapperspb.BoolValue{Value: true}
				output, err := p.plugin.TranslateTransformation(inputTransformation)

				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(Equal(expectedOutput))
			})
		})
	})
})
