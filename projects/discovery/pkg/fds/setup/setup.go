package setup

import (
	"github.com/solo-io/solo-projects/pkg/utils/setuputils"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/syncer"
	gloosyncer "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer"
)

func Main() error {
	return setuputils.Main("fds", gloosyncer.NewSetupFuncWithRun(syncer.RunFDS))
}
