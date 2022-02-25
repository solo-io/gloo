package syncer

import (
	"context"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template"
)

func loadDefaultDashboard(ctx context.Context, templateGenerator template.TemplateGenerator, folderId uint, dashboardClient grafana.DashboardClient) {
	logger := contextutils.LoggerFrom(ctx)

	uid := templateGenerator.GenerateUid()
	_, _, err := dashboardClient.GetRawDashboard(uid)
	if err == nil || err.Error() != grafana.DashboardNotFound(uid).Error() {
		logger.Infof("default dashboard already exists: %s", uid)
		return // do not create dashboard if it already exists
	}

	logger.Infof("generating default dashboard: %s", uid)

	dashPost, err := templateGenerator.GenerateDashboardPost(folderId)
	if err != nil {
		logger.Warnf("failed to generate default dashboard: %s. %s", uid, err)
		return
	}

	err = dashboardClient.PostDashboard(dashPost)
	if err != nil {
		err := errors.Wrapf(err, "failed to save default dashboard to grafana: %s", uid)
		logger.Warn(err.Error())

		return
	}
	logger.Infof("saved default dashboard: %s", uid)
}
