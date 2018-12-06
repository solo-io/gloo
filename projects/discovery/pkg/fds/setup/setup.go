package setup

import (
	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds/syncer"
	gloosyncer "github.com/solo-io/gloo/projects/gloo/pkg/syncer"
)

func Main() error {
	return setuputils.Main("fds", gloosyncer.NewSetupFuncWithRun(syncer.RunFDS))
}
