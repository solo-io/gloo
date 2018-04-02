package azure

import (
	"github.com/gogo/protobuf/types"
	multierror "github.com/hashicorp/go-multierror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Function app name", func() {
	Describe("validation", func() {
		Context("of an empty name", func() {
			It("should fail", func() {
				Expect(isValidFunctionAppName("")).To(Equal(false))
			})
		})
		Context("of a valid name", func() {
			It("should succeed", func() {
				Expect(isValidFunctionAppName("azure-function-app-1")).To(Equal(true))
			})
		})
		Context("of a name starting with a digit", func() {
			It("should succeed", func() {
				Expect(isValidFunctionAppName("1-azure-function-app")).To(Equal(true))
			})
		})
		Context("of a name starting with a dash", func() {
			It("should succeed", func() {
				Expect(isValidFunctionAppName("-azure-function-app-1")).To(Equal(true))
			})
		})
		Context("of a name containing an underscore", func() {
			It("should fail", func() {
				Expect(isValidFunctionAppName("azure-function_app-1")).To(Equal(false))
			})
		})
		Context("of a name containing an invalid character", func() {
			It("should fail", func() {
				Expect(isValidFunctionAppName("azure-function-1!")).To(Equal(false))
			})
		})
	})
})

var _ = Describe("Function name", func() {
	Describe("validation", func() {
		Context("of an empty name", func() {
			It("should fail", func() {
				Expect(isValidFunctionName("")).To(Equal(false))
			})
		})
		Context("of a valid name", func() {
			It("should succeed", func() {
				Expect(isValidFunctionName("azure-function-1")).To(Equal(true))
			})
		})
		Context("of a name starting with a digit", func() {
			It("should fail", func() {
				Expect(isValidFunctionName("1-azure-function")).To(Equal(false))
			})
		})
		Context("of a name starting with a dash", func() {
			It("should fail", func() {
				Expect(isValidFunctionName("-azure-function-1")).To(Equal(false))
			})
		})
		Context("of a name starting with an underscore", func() {
			It("should fail", func() {
				Expect(isValidFunctionName("_azure-function-1")).To(Equal(false))
			})
		})
		Context("of a name containing an invalid character", func() {
			It("should fail", func() {
				Expect(isValidFunctionName("azure-function-1!")).To(Equal(false))
			})
		})
	})
})

var _ = Describe("Authentication level", func() {
	Describe("validation", func() {
		Context("of an empty string", func() {
			It("should fail", func() {
				Expect(isValidAuthLevel("")).To(Equal(false))
			})
		})
		Context("of \"anonymous\"", func() {
			It("should succeed", func() {
				Expect(isValidAuthLevel("anonymous")).To(Equal(true))
			})
		})
		Context("of \"function\"", func() {
			It("should succeed", func() {
				Expect(isValidAuthLevel("function")).To(Equal(true))
			})
		})
		Context("of \"admin\"", func() {
			It("should succeed", func() {
				Expect(isValidAuthLevel("admin")).To(Equal(true))
			})
		})
		Context("of an invalid value", func() {
			It("should fail", func() {
				Expect(isValidAuthLevel("invalid")).To(Equal(false))
			})
		})
	})
})

func decodeUpstreamSpec(functionAppName string, secretRef string) (*UpstreamSpec, error) {
	m := &types.Struct{
		Fields: map[string]*types.Value{
			"function_app_name": {Kind: &types.Value_StringValue{StringValue: functionAppName}},
			"secret_ref":        {Kind: &types.Value_StringValue{StringValue: secretRef}},
		},
	}
	return DecodeUpstreamSpec(m)
}

var _ = Describe("UpstreamSpec", func() {
	Describe("decoding from a map", func() {
		Context("with valid parameters", func() {
			It("should succeed", func() {
				spec, err := decodeUpstreamSpec("azure-function-app-1", "my-azure-sec")
				Expect(err).NotTo(HaveOccurred())
				Expect(spec.FunctionAppName).To(Equal("azure-function-app-1"))
				Expect(spec.SecretRef).To(Equal("my-azure-sec"))
			})
		})
		Context("with an invalid function app name", func() {
			It("should error", func() {
				_, err := decodeUpstreamSpec("_azure-function-app-1", "my-azure-sec")
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("retrieving hostname", func() {
		It("should succeed ", func() {
			spec, err := decodeUpstreamSpec("azure-function-app-1", "my-azure-sec")
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.GetHostname()).To(Equal("azure-function-app-1.azurewebsites.net"))
		})
	})
})

func decodeFunctionSpec(functionName string, authLevel string) (*FunctionSpec, error) {
	m := &types.Struct{
		Fields: map[string]*types.Value{
			"function_name": {Kind: &types.Value_StringValue{StringValue: functionName}},
			"auth_level":    {Kind: &types.Value_StringValue{StringValue: authLevel}},
		},
	}
	return DecodeFunctionSpec(m)
}

var _ = Describe("FunctionSpec", func() {
	Describe("decoding from a map", func() {
		Context("with valid parameters", func() {
			It("should succeed", func() {
				spec, err := decodeFunctionSpec("azure-function-1", "anonymous")
				Expect(err).NotTo(HaveOccurred())
				Expect(spec.FunctionName).To(Equal("azure-function-1"))
				Expect(spec.AuthLevel).To(Equal("anonymous"))
			})
		})
		Context("with an invalid function name", func() {
			It("should result in one error", func() {
				_, err := decodeFunctionSpec("_azure-function-1", "anonymous")
				Expect(err).To(HaveOccurred())
				merr, ok := err.(*multierror.Error)
				Expect(ok).To(Equal(true))
				Expect(len(merr.Errors)).To(Equal(1))
			})
		})
		Context("with an invalid authorization level", func() {
			It("should result in one error", func() {
				_, err := decodeFunctionSpec("azure-function-1", "invalid")
				Expect(err).To(HaveOccurred())
				merr, ok := err.(*multierror.Error)
				Expect(ok).To(Equal(true))
				Expect(len(merr.Errors)).To(Equal(1))
			})
		})
		Context("with an invalid function name and an invalid authorization level", func() {
			It("should result in two errors", func() {
				_, err := decodeFunctionSpec("_azure-function-1", "invalid")
				Expect(err).To(HaveOccurred())
				merr, ok := err.(*multierror.Error)
				Expect(ok).To(Equal(true))
				Expect(len(merr.Errors)).To(Equal(2))
			})
		})
	})
})
