package helpers

import (
	"io"
	"os"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/printers"
	"github.com/solo-io/solo-kit/pkg/utils/cliutils"
)

func PrintVirtualServices(virtualServices gatewayv1.VirtualServiceList, outputType string) {
	cliutils.PrintList(outputType, "", virtualServices,
		func(data interface{}, w io.Writer) error {
			printers.VirtualServiceTable(data.(gatewayv1.VirtualServiceList), w)
			return nil
		}, os.Stdout)
}
