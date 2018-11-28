package setup

import (
	"github.com/solo-io/solo-projects/pkg/utils/setuputils"
	"github.com/solo-io/solo-projects/projects/sqoop/pkg/syncer"
)

func Main() error {
	return setuputils.Main("sqoop", syncer.Setup)
}
