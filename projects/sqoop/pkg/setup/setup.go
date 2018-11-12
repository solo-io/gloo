package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-projects/pkg/utils/setuputils"
	"github.com/solo-io/solo-projects/projects/sqoop/pkg/syncer"
)

func Main(settingsDir string) error {
	return setuputils.Main("sqoopp", syncer.NewSetupSyncer(memory.NewInMemoryResourceCache(), kube.NewKubeCache()), settingsDir)
}
