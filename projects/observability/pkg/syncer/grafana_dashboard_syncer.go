package syncer

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"

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
)

var TAGS = []string{GLOO_TAG, DYNAMIC_TAG}

type GrafanaDashboardsSyncer struct {
	synced bool
	mutex  sync.Mutex
	client *grafana.Client
}

func NewGrafanaDashboardSyncer(client *http.Client) *GrafanaDashboardsSyncer {
	grafCli := grafana.NewClient("http://grafana.gloo-system.svc.cluster.local:3000", "admin:admin", client)
	return &GrafanaDashboardsSyncer{
		client: grafCli,
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

func (s *GrafanaDashboardsSyncer) createGrafanaContent(logger *zap.SugaredLogger, snap *v1.DashboardsSnapshot, gs *grafanaState) {
	for _, upstream := range snap.Upstreams.List() {

		upstreamName := upstream.Metadata.GetName()
		upstreamUid := grafana.NameToUid(upstream.Metadata.GetName())
		// check grafana boards for presence of current upstream
		if gs.containsBoard(upstreamUid) {
			continue
		}
		logger.Infof("generating dashboard for upstream: %s", upstreamName)
		dash, err := grafana.GenerateDashboard(upstream)
		if err != nil {
			logger.Warnf("failed to generate dashboard for upstream: %s. %s", upstreamName, err)
			continue
		}
		err = s.client.SetRawDashboard(dash)
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
			snap, err := grafana.GenerateSnapshot(upstream)
			if err != nil {
				logger.Warnf("failed to generate snapshot for upstream: %s. %s", upstreamName, err)
				continue
			}
			_, err = s.client.SetRawSnapshot(snap)
			if err != nil {
				logger.Warnf("failed to save snapshot to grafana for upstream: %s. %s", upstreamName, err)
				continue
			}
		}
	}
}

func (s *GrafanaDashboardsSyncer) deleteGrafanaContent(logger *zap.SugaredLogger, snap *v1.DashboardsSnapshot, gs *grafanaState) {
	for _, board := range gs.boards {
		missing := true
		for _, upstream := range snap.Upstreams.List() {
			upstreamUid := grafana.NameToUid(upstream.Metadata.GetName())
			if board.UID == upstreamUid {
				missing = false
			}
		}
		if missing {
			logger.Infof("deleting dashboard for missing upstream: %s", board.UID)
			_, err := s.client.DeleteDashboard(board.UID)
			if err != nil {
				logger.Warnf("failed to delete dashboard for upstream: %s. %s", board.UID, err)

			}
		}
	}

	for _, snapshot := range gs.snapshots {
		missing := true
		for _, upstream := range snap.Upstreams.List() {
			upstreamUid := grafana.NameToUid(upstream.Metadata.GetName())
			if snapshot.Name == upstreamUid {
				missing = false
			}
		}
		if missing {
			logger.Infof("deleting snapshot for missing upstream: %s", snapshot.Name)
			err := s.client.DeleteSnapshot(snapshot.Key)
			if err != nil {
				logger.Warnf("failed to delete snapshot for upstream: %s. %s", snapshot.Name, err)
			}
		}
	}

}

func (s *GrafanaDashboardsSyncer) getCurrentGrafanaState() (*grafanaState, error) {
	gs := &grafanaState{}
	snapshots, err := s.client.GetSnapshots()
	if err != nil {
		return gs, fmt.Errorf("unable to get list of current snapshots to compare against, skipping generation: %s", err)
	}
	gs.snapshots = snapshots
	boards, err := s.client.SearchDashboards("", false, TAGS...)
	if err != nil {
		return gs, err
	}
	gs.boards = boards
	if err != nil {
		return gs, err
	}
	return gs, nil
}
