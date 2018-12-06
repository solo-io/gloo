package setuputils

import (
	"flag"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
)

// TODO (ilackarms): move to a flags package
type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	discoveryNamespaces arrayFlags
	setupNamespace      string
	setupName           string
	setupDir            string
)

// TODO (ilackarms): move to a flags package
func init() {
	flag.StringVar(&setupNamespace, "namespace", "gloo-system", "namespace to watch for settings crd/file")
	flag.StringVar(&setupName, "name", defaults.SettingsName, "name of settings crd/file to use")
	flag.StringVar(&setupDir, "dir", "", "directory to find bootstrap settings if not using "+
		"kubernetes crds")
	flag.Var(&discoveryNamespaces, "discover", "repeated list of namespaces to watch for configuration. "+
		"if none are specified, defaults to `default` and `gloo-system`")
}
