package kubernetes

import (
	"context"
	"fmt"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/test/setup/helpers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/solo-io/gloo/test/setup/types"
)

func DeployApplication(ctx context.Context, clusterName string, versions []string, app *types.App, cluster *Cluster) error {
	deployInfo := app.Deployment
	if deployInfo == nil {
		return nil
	}

	var (
		logger     = cluster.GetLogger()
		controller = cluster.GetController()
	)

	logger.Infof("creating deployment: %s", deployInfo.Name)

	stopTimer := helpers.TimerFunc(fmt.Sprintf("[%s] %s creation", clusterName, deployInfo.Name))
	defer stopTimer()

	for _, cm := range app.ConfigMaps {
		if err := CreateOrUpdate[*corev1.ConfigMap](ctx, cm, controller); err != nil {
			return errors.Wrap(err, "failed creating config map")
		}
	}

	if len(versions) == 0 {
		// Single Version
		if app.ServiceAccount != nil {
			if err := CreateOrUpdate[*corev1.ServiceAccount](ctx, app.ServiceAccount, controller); err != nil {
				return errors.Wrap(err, "failed creating service account")
			}
		}
		if app.Deployment != nil {
			if err := CreateOrUpdate[*appsv1.Deployment](ctx, deployInfo, controller); err != nil {
				return errors.Wrap(err, "failed creating deployment")
			}
		}
		if app.Service != nil {
			if err := CreateOrUpdate[*corev1.Service](ctx, app.Service, controller); err != nil {
				return errors.Wrap(err, "failed creating service")
			}
		}
		return nil
	}

	for _, version := range versions {
		if app.ServiceAccount != nil {
			if err := CreateOrUpdate[*corev1.ServiceAccount](ctx, versionServiceAccount(version, app.ServiceAccount), controller); err != nil {
				return errors.Wrap(err, "failed creating service account")
			}
		}
		if app.Service != nil {
			if err := CreateOrUpdate[*corev1.Service](ctx, versionService(version, app.Service), controller); err != nil {
				return errors.Wrap(err, "failed creating service")
			}
		}
		if app.Deployment != nil {
			if err := CreateOrUpdate[*appsv1.Deployment](ctx, versionDeployment(version, deployInfo), controller); err != nil {
				return errors.Wrap(err, "failed creating deployment")
			}
		}
	}

	return nil
}
