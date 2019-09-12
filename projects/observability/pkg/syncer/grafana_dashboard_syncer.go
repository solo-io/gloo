package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go.uber.org/zap"

	v1 "github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"
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
	GLOO_TAG    = "gloo"
	DYNAMIC_TAG = "dynamic"

	SERVICE_LINK = "http://glooe-grafana.gloo-system.svc.cluster.local"
	SERVICE_PORT = "80"
)

var TAGS = []string{GLOO_TAG, DYNAMIC_TAG}

type GrafanaDashboardsSyncer struct {
	synced          bool
	mutex           sync.Mutex
	dashboardClient grafana.DashboardClient
	snapshotClient  grafana.SnapshotClient
}

func NewGrafanaDashboardSyncer(dashboardClient grafana.DashboardClient, snapshotClient grafana.SnapshotClient) *GrafanaDashboardsSyncer {
	return &GrafanaDashboardsSyncer{
		dashboardClient: dashboardClient,
		snapshotClient:  snapshotClient,
	}
}

func (s *GrafanaDashboardsSyncer) Synced() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.synced
}

func (s *GrafanaDashboardsSyncer) Sync(ctx context.Context, snap *v1.DashboardsSnapshot) error {
	ctx = contextutils.WithLogger(ctx, "observabilitySyncer")
	logger := contextutils.LoggerFrom(ctx)

	gs, err := s.getCurrentGrafanaState()
	if err != nil {
		return err
	}
	// handle new upstreams
	s.createGrafanaContent(logger, snap, gs)

	// handle deleted upstreams
	s.deleteGrafanaContent(logger, snap, gs)

	logger.Infof("finished updating all upstream dashboards")
	return nil
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

func (s *GrafanaDashboardsSyncer) createGrafanaContent(logger *zap.SugaredLogger, snap *v1.DashboardsSnapshot, gs *grafanaState) {
	for _, upstream := range snap.Upstreams {
		upstreamName := upstream.Metadata.GetName()
		templateGenerator := grafana.NewTemplateGenerator(upstream)
		upstreamUid := templateGenerator.GenerateUid()

		shouldRegenDashboard, err := s.shouldRegenDashboard(logger, upstreamUid)
		if err != nil {
			logger.Errorf("Skipping dashboard for upstream %s because of error - %s", upstreamUid, err.Error())
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
			logger.Warnf("failed to save dashboard to grafana for upstream: %s. %s", upstreamName, err)
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
				logger.Warnf("failed to save snapshot to grafana for upstream: %s. %s", upstreamName, err)
				continue
			}
		}
	}
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
	return mostRecentVersion.Message != grafana.DefaultCommitMessage, nil
}

func (s *GrafanaDashboardsSyncer) deleteGrafanaContent(logger *zap.SugaredLogger, snap *v1.DashboardsSnapshot, gs *grafanaState) {
	for _, board := range gs.boards {
		missing := true
		for _, upstream := range snap.Upstreams {
			templateGenerator := grafana.NewTemplateGenerator(upstream)
			upstreamUid := templateGenerator.GenerateUid()
			if board.UID == upstreamUid {
				missing = false
			}
		}
		if missing {
			logger.Infof("deleting dashboard for missing upstream: %s", board.UID)
			_, err := s.dashboardClient.DeleteDashboard(board.UID)
			if err != nil {
				logger.Warnf("failed to delete dashboard for upstream: %s. %s", board.UID, err)

			}
		}
	}

	for _, snapshot := range gs.snapshots {
		missing := true
		for _, upstream := range snap.Upstreams {
			templateGenerator := grafana.NewTemplateGenerator(upstream)
			upstreamUid := templateGenerator.GenerateUid()
			if snapshot.Name == upstreamUid {
				missing = false
			}
		}
		if missing {
			logger.Infof("deleting snapshot for missing upstream: %s", snapshot.Name)
			err := s.snapshotClient.DeleteSnapshot(snapshot.Key)
			if err != nil {
				logger.Warnf("failed to delete snapshot for upstream: %s. %s", snapshot.Name, err)
			}
		}
	}

}

func (s *GrafanaDashboardsSyncer) getCurrentGrafanaState() (*grafanaState, error) {
	gs := &grafanaState{}
	snapshots, err := s.snapshotClient.GetSnapshots()
	if err != nil {
		return gs, fmt.Errorf("unable to get list of current snapshots to compare against, skipping generation: %s", err)
	}
	gs.snapshots = snapshots
	boards, err := s.dashboardClient.SearchDashboards("", false, TAGS...)
	if err != nil {
		return gs, err
	}
	gs.boards = boards
	if err != nil {
		return gs, err
	}
	return gs, nil
}
