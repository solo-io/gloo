package functiontypes

import (
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/internal/updater/swagger"
	"github.com/solo-io/gloo-plugins/aws"
	"github.com/solo-io/gloo-plugins/google"
)

type FunctionType string

const (
	FunctionTypeLambda  FunctionType = "functionTypeLambda"
	FunctionTypeGfuncs  FunctionType = "functionTypeGfuncs"
	FunctionTypeSwagger FunctionType = "functionTypeSwagger"
	NonFunctional       FunctionType = "nonFunctional"
)

func GetFunctionType(us *v1.Upstream) FunctionType {
	switch {
	case us.Type == aws.UpstreamTypeAws:
		return FunctionTypeLambda
	case us.Type == gfunc.UpstreamTypeGoogle:
		return FunctionTypeGfuncs
	case swagger.IsSwagger(us):
		return FunctionTypeSwagger
	}
	return NonFunctional
}
