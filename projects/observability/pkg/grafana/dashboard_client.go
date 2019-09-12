package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

//go:generate mockgen -destination mocks/mock_dashboard_client.go -package mocks github.com/solo-io/solo-projects/projects/observability/pkg/grafana DashboardClient

type DashboardClient interface {
	GetRawDashboard(uid string) ([]byte, BoardProperties, error)
	GetDashboardVersions(dashboardId float64) ([]*Version, error)
	SetRawDashboard(raw []byte) error
	DeleteDashboard(uid string) (StatusMessage, error)
	SearchDashboards(query string, starred bool, tags ...string) ([]FoundBoard, error)
}

type dashboardClient struct {
	restClient RestClient
}

var _ DashboardClient = &dashboardClient{}

// BoardProperties keeps metadata of a dashboard.
type BoardProperties struct {
	IsStarred  bool      `json:"isStarred,omitempty"`
	IsHome     bool      `json:"isHome,omitempty"`
	IsSnapshot bool      `json:"isSnapshot,omitempty"`
	Type       string    `json:"type,omitempty"`
	CanSave    bool      `json:"canSave"`
	CanEdit    bool      `json:"canEdit"`
	CanStar    bool      `json:"canStar"`
	Expires    time.Time `json:"expires"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	UpdatedBy  string    `json:"updatedBy"`
	CreatedBy  string    `json:"createdBy"`
	Version    int       `json:"version"`
}

type Version struct {
	VersionId int
	Message   string
}

// FoundBoard keeps result of search with metadata of a dashboard.
type FoundBoard struct {
	ID        uint     `json:"id"`
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	URI       string   `json:"uri"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
}

func NewDashboardClient(restClient RestClient) DashboardClient {
	return &dashboardClient{
		restClient: restClient,
	}
}

// GetRawDashboard loads a dashboard JSON from Grafana instance along with metadata for a dashboard.
// Contrary to GetDashboard() it not unpack loaded JSON to Board structure. Instead it
// returns it as byte slice. It guarantees that data of dashboard returned untouched by conversion
// with Board so no matter how properly fields from a current version of Grafana mapped to
// our Board fields. It useful for backup purposes when you want a dashboard exactly with
// same data as it exported by Grafana.
func (d *dashboardClient) GetRawDashboard(uid string) ([]byte, BoardProperties, error) {
	var (
		raw    []byte
		result struct {
			Meta  BoardProperties `json:"meta"`
			Board json.RawMessage `json:"dashboard"`
		}
		code int
		err  error
	)
	if raw, code, err = d.restClient.Get(fmt.Sprintf("api/dashboards/uid/%s", uid), nil); err != nil {
		return nil, BoardProperties{}, err
	}
	if code == 404 {
		return nil, BoardProperties{}, DashboardNotFound(uid)
	}
	if code != 200 {
		return nil, BoardProperties{}, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&result); err != nil {
		return nil, BoardProperties{}, fmt.Errorf("unmarshal board with meta: %s\n%s", err, raw)
	}
	return []byte(result.Board), result.Meta, err
}

// get all the versions in the history for this dashboard, ordered from newest to oldest
func (d *dashboardClient) GetDashboardVersions(dashboardId float64) ([]*Version, error) {
	var (
		rawResp []byte
		err     error
	)

	// https://grafana.com/docs/http_api/dashboard_versions/#get-all-dashboard-versions
	rawResp, _, err = d.restClient.Get("/api/dashboards/id/"+fmt.Sprint(dashboardId)+"/versions", url.Values{"start": {"0"}, "limit": {""}})

	if err != nil {
		return nil, err
	}

	var allVersions []*Version
	err = json.Unmarshal(rawResp, &allVersions)
	if err != nil {
		return nil, err
	}

	return allVersions, nil
}

// SetRawDashboard updates existing dashboard or creates a new one.
// Contrary to SetDashboard() it accepts raw JSON instead of Board structure.
func (d *dashboardClient) SetRawDashboard(raw []byte) error {
	var (
		rawResp []byte
		resp    StatusMessage
		code    int
		err     error
	)
	if rawResp, code, err = d.restClient.Post("api/dashboards/db", nil, raw); err != nil {
		return err
	}
	switch code {
	case 400:
		return fmt.Errorf("%d %s", code, rawResp)
	case 401:
		return fmt.Errorf("%d %s", code, rawResp)
	case 412:
		if err = json.Unmarshal(rawResp, &resp); err != nil {
			return err
		}
		return fmt.Errorf("%d %s", code, *resp.Message)
	}
	return nil
}

func (d *dashboardClient) DeleteDashboard(uid string) (StatusMessage, error) {
	var (
		raw   []byte
		reply StatusMessage
		err   error
	)

	if raw, _, err = d.restClient.Delete(fmt.Sprintf("api/dashboards/uid/%s", uid)); err != nil {
		return StatusMessage{}, err
	}
	err = json.Unmarshal(raw, &reply)
	return reply, err
}

func (d *dashboardClient) SearchDashboards(query string, starred bool, tags ...string) ([]FoundBoard, error) {
	var (
		raw    []byte
		boards []FoundBoard
		code   int
		err    error
	)
	u := url.URL{}
	q := u.Query()
	if query != "" {
		q.Set("query", query)
	}
	if starred {
		q.Set("starred", "true")
	}
	for _, tag := range tags {
		q.Add("tag", tag)
	}
	if raw, code, err = d.restClient.Get("api/search", q); err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	err = json.Unmarshal(raw, &boards)
	return boards, err
}
