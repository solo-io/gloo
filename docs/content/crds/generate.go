package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/cmdutils"
	"github.com/solo-io/skv2/codegen/util"
	"os"
	"path/filepath"
)

var (
	crdocBinary = filepath.Join(util.GetModuleRoot(), "_output", ".bin", "crdoc")
)

type CrdReference struct {
	SourceCrdFile               string
	DestinationCrdReferenceFile string
}

func main() {
	err := generateCrdReferenceDocs(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateCrdReferenceDocs(ctx context.Context) error {
	crds := []CrdReference{
		{
			SourceCrdFile:               "install/helm/gloo/crds/gateway.gloo.solo.io_gatewayparameters.yaml",
			DestinationCrdReferenceFile: "docs/content/reference/api/github.com/solo-io/gloo/projects/gateway2/api/v1alpha1/gateway_parameters.md",
		},
		{
			SourceCrdFile:               "install/helm/gloo/crds/gateway.gloo.solo.io_directresponses.yaml",
			DestinationCrdReferenceFile: "docs/content/reference/api/github.com/solo-io/gloo/projects/gateway2/api/v1alpha1/direct_response_action.md",
		},
	}

	outErr := &multierror.Error{}
	for _, crd := range crds {
		err := generateCrdReferenceMarkdown(ctx, crd)
		outErr = multierror.Append(outErr, err)
	}

	return outErr.ErrorOrNil()
}

func generateCrdReferenceMarkdown(ctx context.Context, crd CrdReference) error {
	cmd := cmdutils.Command(ctx, crdocBinary,
		"--resources",
		crd.SourceCrdFile,
		"--output",
		crd.DestinationCrdReferenceFile,
		"--template",
		filepath.Join("docs", "content", "crds", "templates", "markdown.tmpl"),
	)
	runErr := cmd.Run()

	if runErr.Cause() != nil {
		return eris.Wrapf(runErr.Cause(), "%s produced error %s", runErr.PrettyCommand(), runErr.Error())
	}
	return nil
}
