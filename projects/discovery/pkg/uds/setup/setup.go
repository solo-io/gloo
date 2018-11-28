package setup

import (
	"github.com/solo-io/solo-projects/pkg/utils/setuputils"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/uds/syncer"
	gloosyncer "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer"
)

func Main() error {
	return setuputils.Main("uds", gloosyncer.NewSetupFuncWithRun(syncer.RunUDS))
}
