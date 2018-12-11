package constants

import (
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
)

var DefaultDomains = []string{"*"}

// TODO(mitchdraft) get this from a function call
var WatchNamespaces = []string{defaults.GlooSystem}
