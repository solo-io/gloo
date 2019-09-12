package grafana

import (
	"encoding/json"
	"fmt"
	"net/url"
)

//go:generate mockgen -destination mocks/mock_snapshot_client.go -package mocks github.com/solo-io/solo-projects/projects/observability/pkg/grafana SnapshotClient

type SnapshotClient interface {
	SetRawSnapshot(raw []byte) (*SnapshotResponse, error)
	GetSnapshots() ([]SnapshotListResponse, error)
	DeleteSnapshot(key string) error
}

type SnapshotResponse struct {
	DeleteKey string `json:"deleteKey"`
	DeleteUrl string `json:"deleteUrl"`
	Key       string `json:"key"`
	Url       string `json:"url"`
}

type SnapshotListResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Key      string `json:"key"`
	OrgId    int    `json:"orgId"`
	UserId   int    `json:"userId"`
	External bool   `json:"external"`
	Expires  string `json:"expires"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
}

type snapshotClient struct {
	restClient RestClient
}

var _ SnapshotClient = &snapshotClient{}

func NewSnapshotClient(restClient RestClient) SnapshotClient {
	return &snapshotClient{
		restClient: restClient,
	}
}

func (s *snapshotClient) SetRawSnapshot(raw []byte) (*SnapshotResponse, error) {
	var (
		rawResp []byte
		resp    StatusMessage
		code    int
		err     error
		sresp   = &SnapshotResponse{}
	)
	if rawResp, code, err = s.restClient.Post("api/snapshots", nil, raw); err != nil {
		return sresp, err
	}
	switch code {
	case 400:
		return sresp, fmt.Errorf("%d %s", code, rawResp)
	case 401:
		return sresp, fmt.Errorf("%d %s", code, rawResp)
	case 412:
		if err = json.Unmarshal(rawResp, &resp); err != nil {
			return sresp, err
		}
		return sresp, fmt.Errorf("%d %s", code, *resp.Message)
	default:
		err := json.Unmarshal(rawResp, sresp)
		if err != nil {
			return sresp, err
		}
		return sresp, nil
	}
}

func (s *snapshotClient) GetSnapshots() ([]SnapshotListResponse, error) {
	var (
		rawResp   []byte
		err       error
		code      int
		slistresp []SnapshotListResponse
	)

	if rawResp, code, err = s.restClient.Get("/api/dashboard/snapshots", url.Values{}); err != nil {
		return slistresp, err
	}
	switch code {
	case 200:
		err := json.Unmarshal(rawResp, &slistresp)
		if err != nil {
			return slistresp, err
		}
		return slistresp, nil
	default:

		return slistresp, fmt.Errorf("error while trying to get snapshots (code: %d),\n response:\n %s", code, string(rawResp))
	}
}

func (s *snapshotClient) DeleteSnapshot(key string) error {
	var (
		err     error
		rawResp []byte
		code    int
	)

	if rawResp, code, err = s.restClient.Delete(fmt.Sprintf("api/snapshots/%s", key)); err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("%d error: \n %s", code, string(rawResp))
	}
	return nil
}
