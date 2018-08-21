package azure

import (
	"fmt"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
)

func GetHostname(s *azure.UpstreamSpec) string {
	return fmt.Sprintf("%s.azurewebsites.net", s.FunctionAppName)
}
