package setup

import (
	"time"

	check "github.com/solo-io/go-checkpoint"
	"github.com/solo-io/solo-projects/pkg/utils/setuputils"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer"
)

func Main() error {
	start := time.Now()
	check.CallCheck("gloo", version.Version, start)
	return setuputils.Main("gloo", syncer.NewSetupFunc())
}
