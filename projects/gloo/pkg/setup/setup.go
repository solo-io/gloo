package setup

import (
	"context"

	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils/statuses"

	"github.com/solo-io/gloo/pkg/version"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
)

func Main(
	customCtx context.Context,
	statusReporter statuses.StatusReportChan,
) error {
	return startSetupLoop(customCtx, statusReporter)
}

func StartGlooInTest(customCtx context.Context) error {
	return startSetupLoop(customCtx, nil)
}

func startSetupLoop(
	ctx context.Context,
	statusReporter statuses.StatusReportChan,
) error {
	return setuputils.Main(
		setuputils.SetupOpts{
			LoggerName: "gloo",
			Version:    version.Version,
			SetupFunc: setup.NewSetupFuncWithExtensions(
				setup.Extensions{
					StatusReporter: statusReporter,
				},
			),
			ExitOnError: true,
			CustomCtx:   ctx,
		},
	)
}
