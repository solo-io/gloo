package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/hashicorp/go-multierror"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template"
	"go.uber.org/zap"
)

type grafanaState struct {
	boards    []grafana.FoundBoard
	snapshots []grafana.SnapshotListResponse
}

func (gs *grafanaState) containsBoard(upstreamUid string) bool {
	for _, v := range gs.boards {
		if v.UID == upstreamUid {
			return true
		}
	}
	return false
}

const (
	glooTag    = "gloo"
	dynamicTag = "dynamic"
)

var (
	tags = []string{glooTag, dynamicTag}

	// exponential backoff retry with an initial period of 0.1s for 7 iterations, which will mean a cumulative retry period of ~20s
	// used to avoid a race condition where observability is ready and Sync'ing before Grafana is ready to accept requests
	grafanaSyncRetryOpts = []retry.Option{retry.Delay(time.Millisecond * 100), retry.Attempts(7), retry.DelayType(retry.BackOffDelay)}
)

type GrafanaDashboardsSyncer struct {
	synced                bool
	mutex                 sync.Mutex
	dashboardClient       grafana.DashboardClient
	snapshotClient        grafana.SnapshotClient
	dashboardJsonTemplate string
}

func NewGrafanaDashboardSyncer(dashboardClient grafana.DashboardClient, snapshotClient grafana.SnapshotClient, dashboardJsonTemplate string) *GrafanaDashboardsSyncer {
	return &GrafanaDashboardsSyncer{
		dashboardClient:       dashboardClient,
		snapshotClient:        snapshotClient,
		dashboardJsonTemplate: dashboardJsonTemplate,
	}
}

func (s *GrafanaDashboardsSyncer) Synced() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.synced
}

func (s *GrafanaDashboardsSyncer) Sync(ctx context.Context, snap *v1.DashboardsSnapshot) error {
	// The observability deployment may finish and this Sync may start running before the Grafana deployment is ready
	// Retry the Sync until Grafana is ready
	// https://github.com/solo-io/solo-projects/issues/307
	ctx = contextutils.WithLogger(ctx, "observabilitySyncer")
	logger := contextutils.LoggerFrom(ctx)

	syncAttempt := 0
	return retry.Do(func() error {
		syncAttempt++

		if syncAttempt == 1 {
			logger.Info("First attempt to sync upstream state with Grafana")
		} else {
			logger.Warnf("Attempt number %d to sync with Grafana", syncAttempt)
		}
		gs, err := s.getCurrentGrafanaState()
		if err != nil {
			return err
		}
		// handle new upstreams
		networkErrs := s.createGrafanaContent(logger, snap, gs)
		if err := networkErrs.ErrorOrNil(); err != nil {
			return errors.Errorf("Encountered errors while communicating with Grafana while creating new content: %+v", err)
		}

		// handle deleted upstreams
		networkErrs = s.deleteGrafanaContent(logger, snap, gs)
		if err := networkErrs.ErrorOrNil(); err != nil {
			return errors.Errorf("Encountered errors while communicating with Grafana while deleting old content: %+v", err)
		}

		logger.Infof("finished updating all upstream dashboards")
		return nil
	}, grafanaSyncRetryOpts...)
}

// we regenerate a dashboard if:
// 1. it's been deleted by the user, ie if GetRawDashboard returns DashboardNotFound, OR
// 2. it has not been manually edited by a user
func (s *GrafanaDashboardsSyncer) shouldRegenDashboard(logger *zap.SugaredLogger, upstreamUid string) (bool, error) {
	rawDashboard, _, err := s.dashboardClient.GetRawDashboard(upstreamUid)
	if err != nil && err.Error() != grafana.DashboardNotFound(upstreamUid).Error() {
		logger.Warnf("Failed to get raw dashboard for uid %s - %s", upstreamUid, err.Error())
		return false, err
	} else if err == nil {
		isEditedByUser, err := s.isEditedByUser(rawDashboard)
		if err != nil {
			logger.Warnf("Could not verify whether dashboard has been edited- %s", err.Error())
			return false, err
		}
		if isEditedByUser {
			logger.Infof("Dashboard %s has been edited by a user - skipping", upstreamUid)
			return false, nil
		}
	}

	return true, nil
}

// returns a list of errors encountered while communicating over the network with Grafana
func (s *GrafanaDashboardsSyncer) createGrafanaContent(logger *zap.SugaredLogger, snap *v1.DashboardsSnapshot, gs *grafanaState) *multierror.Error {
	errs := &multierror.Error{}

	for _, upstream := range snap.Upstreams {
		upstreamName := upstream.Metadata.GetName()
		templateGenerator := template.NewTemplateGenerator(upstream, s.dashboardJsonTemplate)
		upstreamUid := templateGenerator.GenerateUid()

		shouldRegenDashboard, err := s.shouldRegenDashboard(logger, upstreamUid)
		if err != nil {
			err := errors.Wrapf(err, "Skipping dashboard for upstream %s", upstreamUid)
			logger.Warn(err.Error())
			errs = multierror.Append(errs, err)

			continue
		}
		if !shouldRegenDashboard {
			continue
		}

		logger.Infof("generating dashboard for upstream: %s", upstreamName)

		dash, err := templateGenerator.GenerateDashboard()

		if err != nil {
			logger.Warnf("failed to generate dashboard for upstream: %s. %s", upstreamName, err)
			continue
		}
		err = s.dashboardClient.SetRawDashboard(dash)
		if err != nil {
			err := errors.Wrapf(err, "failed to save dashboard to grafana for upstream: %s", upstreamName)
			logger.Warn(err.Error())
			errs = multierror.Append(errs, err)

			continue
		}

		missing := true
		for _, snapshot := range gs.snapshots {
			if snapshot.Name == upstreamUid {
				missing = false
			}
		}
		if missing {
			// Generate snapshot
			snap, err := templateGenerator.GenerateSnapshot()
			if err != nil {
				logger.Warnf("failed to generate snapshot for upstream: %s. %s", upstreamName, err)
				continue
			}
			_, err = s.snapshotClient.SetRawSnapshot(snap)
			if err != nil {
				err := errors.Wrapf(err, "failed to save snapshot to grafana for upstream: %s", upstreamName)
				logger.Warn(err.Error())
				errs = multierror.Append(errs, err)

				continue
			}
		}
	}

	return errs
}

func (s *GrafanaDashboardsSyncer) isEditedByUser(rawDashboard []byte) (bool, error) {
	dataMap := make(map[string]interface{})
	err := json.Unmarshal(rawDashboard, &dataMap)
	if err != nil {
		return false, err
	}

	dashboardIdInterface, ok := dataMap["id"]
	if !ok {
		return false, DashboardIdNotFound(string(rawDashboard))
	}

	dashboardId, ok := dashboardIdInterface.(float64)
	if !ok {
		return false, DashboardIdConversionError(dashboardIdInterface)
	}

	allVersions, err := s.dashboardClient.GetDashboardVersions(dashboardId)
	if err != nil {
		return false, err
	}

	// probably not possible, but just being defensive
	if len(allVersions) == 0 {
		return false, nil
	}

	mostRecentVersion := allVersions[0]

	// if the version message is not one that we automatically set, then a user must have manually updated it
	return mostRecentVersion.Message != template.DefaultCommitMessage, nil
}

// returns a list of errors encountered while communicating over the network with Grafana
func (s *GrafanaDashboardsSyncer) deleteGrafanaContent(logger *zap.SugaredLogger, snap *v1.DashboardsSnapshot, gs *grafanaState) *multierror.Error {
	errs := &multierror.Error{}

	for _, board := range gs.boards {
		missing := true
		for _, upstream := range snap.Upstreams {
			templateGenerator := template.NewTemplateGenerator(upstream, s.dashboardJsonTemplate)
			upstreamUid := templateGenerator.GenerateUid()
			if board.UID == upstreamUid {
				missing = false
			}
		}
		if missing {
			logger.Infof("deleting dashboard for missing upstream: %s", board.UID)
			_, err := s.dashboardClient.DeleteDashboard(board.UID)
			if err != nil {
				err := errors.Wrapf(err, "failed to delete dashboard for upstream: %s", board.UID)
				logger.Warn(err.Error())
				errs = multierror.Append(errs, err)

			}
		}
	}

	for _, snapshot := range gs.snapshots {
		missing := true
		for _, upstream := range snap.Upstreams {
			templateGenerator := template.NewTemplateGenerator(upstream, s.dashboardJsonTemplate)
			upstreamUid := templateGenerator.GenerateUid()
			if snapshot.Name == upstreamUid {
				missing = false
			}
		}
		if missing {
			logger.Infof("deleting snapshot for missing upstream: %s", snapshot.Name)
			err := s.snapshotClient.DeleteSnapshot(snapshot.Key)
			if err != nil {
				err := errors.Wrapf(err, "failed to delete snapshot for upstream: %s", snapshot.Name)
				logger.Warn(err.Error())
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs
}

func (s *GrafanaDashboardsSyncer) getCurrentGrafanaState() (*grafanaState, error) {
	gs := &grafanaState{}
	snapshots, err := s.snapshotClient.GetSnapshots()
	if err != nil {
		return gs, fmt.Errorf("unable to get list of current snapshots to compare against, skipping generation: %s", err)
	}
	gs.snapshots = snapshots
	boards, err := s.dashboardClient.SearchDashboards("", false, tags...)
	if err != nil {
		return gs, err
	}
	gs.boards = boards
	if err != nil {
		return gs, err
	}
	return gs, nil
}
