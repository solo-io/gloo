package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
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
	// this needs to remain in sync with what's listed in the documentation of
	// the default id variable located in settings.grafanaConfiguration.defaultDashboardFolderId
	upstreamFolderIdAnnotationKey = "observability.solo.io/dashboard_folder_id"
	// this is the default folder ID used by Grafana, which they refer to as the general folder.
	// It is always a valid folder to create dashboards in. We defer to this the moment it we determine
	// that we can't determine if a non-general folder id is valid or not, since behavior for
	// creating/moving dashboards towards invalid folder ids is inconsistent and riddled with silent
	// errors that Grafana doesn't propagate.
	generalFolderId = uint(0)
)

var (
	tags = []string{glooTag, dynamicTag}

	// exponential backoff retry with an initial period of 0.1s for 7 iterations, which will mean a cumulative retry period of ~20s
	// used to avoid a race condition where observability is ready and Sync'ing before Grafana is ready to accept requests
	grafanaSyncRetryOpts = []retry.Option{retry.Delay(time.Millisecond * 100), retry.Attempts(7), retry.DelayType(retry.BackOffDelay)}
)

type GrafanaDashboardsSyncer struct {
	synced                        bool
	mutex                         sync.Mutex
	dashboardClient               grafana.DashboardClient
	snapshotClient                grafana.SnapshotClient
	upstreamDashboardJsonTemplate string
	defaultDashboardUids          map[string]struct{} // expected dashboards that do not correspond with upstreams
	defaultDashboardFolderId      uint
}

func NewGrafanaDashboardSyncer(dashboardClient grafana.DashboardClient, snapshotClient grafana.SnapshotClient, upstreamDashboardJsonTemplate string, defaultDashboardFolderId uint, defaultDashboardUids map[string]struct{}) *GrafanaDashboardsSyncer {
	return &GrafanaDashboardsSyncer{
		dashboardClient:               dashboardClient,
		snapshotClient:                snapshotClient,
		upstreamDashboardJsonTemplate: upstreamDashboardJsonTemplate,
		defaultDashboardFolderId:      defaultDashboardFolderId,
		defaultDashboardUids:          defaultDashboardUids,
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
func (s *GrafanaDashboardsSyncer) shouldRegenDashboard(logger *zap.SugaredLogger, uid string) (bool, error) {
	rawDashboard, _, err := s.dashboardClient.GetRawDashboard(uid)
	if err != nil && err.Error() != grafana.DashboardNotFound(uid).Error() {
		logger.Warnf("Failed to get raw dashboard for uid %s - %s", uid, err.Error())
		return false, err
	} else if err == nil {
		isEditedByUser, err := s.isEditedByUser(rawDashboard)
		if err != nil {
			logger.Warnf("Could not verify whether dashboard has been edited- %s", err.Error())
			return false, err
		}
		if isEditedByUser {
			logger.Infof("Dashboard %s has been edited by a user - skipping", uid)
			return false, nil
		}
	}

	return true, nil
}

// returns a list of errors encountered while communicating over the network with Grafana
func (s *GrafanaDashboardsSyncer) createGrafanaContent(logger *zap.SugaredLogger, snap *v1.DashboardsSnapshot, gs *grafanaState) *multierror.Error {
	errs := &multierror.Error{}
	// This map shouldn't be persisted across syncs, but it should also only be
	// populated once per sync to minimize API calls, and only if needed to validate custom folder ids.
	folderIds := make(map[uint]bool)

	// validate the default id, and set it to the always-safe generalFolderId if something goes wrong
	isValid, err := s.validateFolderId(logger, folderIds, s.defaultDashboardFolderId)
	if !isValid {
		// whatever the reason, we can't use the configured folder id, so use the safe default, then log a reason.
		s.defaultDashboardFolderId = generalFolderId
		err = errors.Errorf("Default dashboard folder ID of \"%d\" is invalid. All upstreams without "+
			"their own annotated dashboard folder IDs will have their dashboards sent to the general folder (id: %d) "+
			"The list of valid folderIds returned by Grafana includes the following values: [%s]",
			s.defaultDashboardFolderId,
			generalFolderId,
			folderIdsToString(folderIds))
		if err != nil { // replace more-common invalid id error with data retrieval error the first time it occurs.
			err = errors.Wrapf(err, "Due to an error in data retrieval, the configured default dashboard id %d "+
				"can't be validated. All upstreams dashboards will be generated in the general folder (id: %d).",
				s.defaultDashboardFolderId,
				generalFolderId)
		}
		// log, but don't propagate/return folder id-related errors, as this could cause endless re-sync loops.
		logger.Warn(err.Error())
	}

	for _, upstream := range snap.Upstreams {
		folderIdToUse := s.defaultDashboardFolderId
		upstreamName := upstream.Metadata.GetName()

		// If this upstream is annotated with a custom Grafana folder id, we need to validate it in the same
		// manner that we checked the default id.
		if annotatedFolderId, annotationFound := upstream.Metadata.Annotations[upstreamFolderIdAnnotationKey]; annotationFound {
			// start by making sure the annotation value is actually an int
			intFolderId, err := strconv.Atoi(annotatedFolderId)

			if err != nil {
				err = errors.Wrapf(err, "Annotated folder id %s for the upstream %s could not be converted to an "+
					"integer. This upstream's dashboard will be created in the general folder (id %d) instead.",
					annotatedFolderId,
					upstreamName,
					generalFolderId)
				logger.Warn(err.Error())
			} else { // it's an int, make sure that it's valid next.
				uintFolderId := uint(intFolderId)
				isValid, err := s.validateFolderId(logger, folderIds, uintFolderId)

				if !isValid {
					err = errors.Errorf("Annotated dashboard folder ID \"%d\" for upstream %s is invalid."+
						"This upstream's dashboard will be created in the general folder (id %d) instead. "+
						"The list of valid folderIds returned by Grafana includes the following values: [%s]",
						uintFolderId,
						upstreamName,
						generalFolderId,
						folderIdsToString(folderIds))
					if err != nil { // replace more-common invalid id error with data retrieval error the first time it occurs.
						err = errors.Wrapf(err, "Due to an error in data retrieval, we cannot determine if the given dashboard "+
							"folder id %d for the upstream %s is valid. This upstream's dashboard will be created in the general folder (id %d) instead.",
							uintFolderId,
							upstreamName,
							generalFolderId)
					}

					logger.Warn(err.Error())
				} else { // annotated folder id has been validated, use it instead of s.defaultDashboardFolderId
					folderIdToUse = uintFolderId
				}
			}
		}

		templateGenerator := template.NewUpstreamTemplateGenerator(upstream, s.upstreamDashboardJsonTemplate)
		uid := templateGenerator.GenerateUid()

		shouldRegenDashboard, err := s.shouldRegenDashboard(logger, uid)
		if err != nil {
			err := errors.Wrapf(err, "Skipping dashboard for upstream %s", uid)
			logger.Warn(err.Error())
			errs = multierror.Append(errs, err)
			continue
		}
		if !shouldRegenDashboard {
			continue
		}

		logger.Infof("generating dashboard for upstream: %s", upstreamName)

		dashPost, err := templateGenerator.GenerateDashboardPost(folderIdToUse)

		if err != nil {
			logger.Warnf("failed to generate dashboard for upstream: %s. %s", upstreamName, err)
			continue
		}
		err = s.dashboardClient.PostDashboard(dashPost)
		if err != nil {
			err := errors.Wrapf(err, "failed to save dashboard to grafana for upstream: %s", upstreamName)
			logger.Warn(err.Error())
			errs = multierror.Append(errs, err)

			continue
		}

		missing := true
		for _, snapshot := range gs.snapshots {
			if snapshot.Name == uid {
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
		if _, ok := s.defaultDashboardUids[board.UID]; ok {
			continue // default dashboard should not be deleted
		}

		missing := true
		for _, upstream := range snap.Upstreams {
			templateGenerator := template.NewUpstreamTemplateGenerator(upstream, s.upstreamDashboardJsonTemplate)
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
			templateGenerator := template.NewUpstreamTemplateGenerator(upstream, s.upstreamDashboardJsonTemplate)
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

func folderIdsToString(folderIds map[uint]bool) string {
	res := ""
	for key, _ := range folderIds {
		res += fmt.Sprintf("%d, ", key)
	}
	return res
}

// call the grafana API to get a list of all folder IDs
func (s *GrafanaDashboardsSyncer) populateValidFolderIds(folderIds map[uint]bool) error {
	folderIds[generalFolderId] = true // Default of 0 is always valid, but isn't returned by the folder API.
	folders, err := s.dashboardClient.GetAllFolderIds()
	if err != nil {
		return errors.Wrapf(err, "folder info retrieval failed, cannot generate dashboards with "+
			"non-zero folderIds")
	} else {
		for _, folder := range folders {
			folderIds[folder.ID] = true
		}
		return nil
	}
}

// Grafana silently ignores invalid folderId values on dashboard API calls, and either sends the
// resulting dashboard to the general folder (on dashboard creation), or does nothing (on modification).
// Rather than let that slide, we cross check any inputted folder ID against
// a list of all known folder IDs.
// Returns true if the inputted ID is valid, and propagates an error if grafana folder id retrieval failed.
// As a side effect, this populates the folderIds map if it is empty and folderId != generalFolderId.
func (s *GrafanaDashboardsSyncer) validateFolderId(logger *zap.SugaredLogger, folderIds map[uint]bool, folderId uint) (bool, error) {
	if folderId != generalFolderId {
		// This is only true if the map hasn't been populated, since even if the API call to gradana fails,
		// we still populate this map with [generalFolderId] = true at a minimum.
		if len(folderIds) == 0 {
			logger.Infof("Detected first non-nil folder id (%d), "+
				"retrieving folderId list for verification.", folderId)
			err := s.populateValidFolderIds(folderIds)
			if err != nil { // something went wrong
				return false, err
			}
		}
		// now that we actually have folder ids, validate the inputted id
		if _, ok := folderIds[folderId]; !ok {
			// If the current upstream's folderId is not in the folderIds set, then it's invalid and can't be
			// used.
			return false, nil
		}
	}
	return true, nil // folder id either equals generalFolderId or is valid.
}
