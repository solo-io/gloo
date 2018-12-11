package actions

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/cliutil"
	ro "github.com/solo-io/gloo/projects/gloo/cli/pkg/route/options"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins"
	grpc "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/grpc"
)

// TODO: integrate this function to the route add command
// currently unused
// TODO(mitchdraft) - get a list of available function names
// validate input and provide a selection list
func ensureGrpcParams(r *ro.Options, s *plugins.ServiceSpec_Grpc, static bool) error {
	gOpt := &r.GrpcDestinationSpec

	packages := []string{""}
	services := []string{""}
	functions := []string{""}
	if err := ensureGrpcStringParam("package", &gOpt.Package, packages, static); err != nil {
		return err
	}
	if err := ensureGrpcStringParam("service", &gOpt.Service, services, static); err != nil {
		return err
	}
	if err := ensureGrpcStringParam("function", &gOpt.Function, functions, static); err != nil {
		return err
	}

	if err := ensureCommonRestGrpcParameters(r, static); err != nil {
		return err
	}

	generateGrpcDestinationSpec(r)

	return nil
}

func generateGrpcDestinationSpec(r *ro.Options) *gloov1.DestinationSpec {
	gOpts := r.GrpcDestinationSpec
	return &gloov1.DestinationSpec{
		DestinationType: &gloov1.DestinationSpec_Grpc{
			Grpc: &grpc.DestinationSpec{
				Package:    gOpts.Package,
				Service:    gOpts.Service,
				Function:   gOpts.Function,
				Parameters: generateParametersFromInput(r),
			},
		},
	}
}

func ensureGrpcStringParam(noun string, target *string, options []string, static bool) error {
	if static {
		if *target == "" {
			return fmt.Errorf(fmt.Sprintf("Please provide a %v.", noun))
		}
		// TODO - check that target is among the options
	}
	if err := cliutil.ChooseFromList(fmt.Sprintf("Please choose a %v.", noun), target, options); err != nil {
		return err
	}

	return nil
}
