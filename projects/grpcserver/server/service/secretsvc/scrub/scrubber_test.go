package scrub_test

import (
	"context"

	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc/scrub"
)

var (
	scrubber     scrub.Scrubber
	testMetadata = core.Metadata{
		Name:      "doo",
		Namespace: "bar",
	}
)

var _ = Describe("Scrubber Test", func() {
	Describe("Scrubber", func() {

		BeforeEach(func() {
			scrubber = scrub.NewScrubber()
		})

		It("works", func() {
			testCases := []struct {
				desc            string
				input, expected *gloov1.Secret
			}{
				{
					desc: "aws",
					input: &gloov1.Secret{
						Kind: &gloov1.Secret_Aws{
							Aws: &gloov1.AwsSecret{
								AccessKey: "hello",
								SecretKey: "world",
							}},
						Metadata: testMetadata,
					},
					expected: &gloov1.Secret{
						Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
						Metadata: testMetadata,
					},
				},
				{
					desc: "azure",
					input: &gloov1.Secret{
						Kind: &gloov1.Secret_Azure{
							Azure: &gloov1.AzureSecret{
								ApiKeys: map[string]string{"hello": "world"},
							}},
						Metadata: testMetadata,
					},
					expected: &gloov1.Secret{
						Kind:     &gloov1.Secret_Azure{Azure: &gloov1.AzureSecret{}},
						Metadata: testMetadata,
					},
				},
				{
					desc: "tls",
					input: &gloov1.Secret{
						Kind: &gloov1.Secret_Tls{
							Tls: &gloov1.TlsSecret{
								PrivateKey: "asdf",
							}},
						Metadata: testMetadata,
					},
					expected: &gloov1.Secret{
						Kind:     &gloov1.Secret_Tls{Tls: &gloov1.TlsSecret{}},
						Metadata: testMetadata,
					},
				},
			}

			for _, tc := range testCases {
				scrubber.Secret(context.Background(), tc.input)
				Expect(tc.input).To(Equal(tc.expected), tc.desc)
			}
		})
	})

})
