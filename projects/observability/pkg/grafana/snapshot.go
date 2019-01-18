package grafana

import (
	"encoding/json"
	"fmt"
	"net/url"
)

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

func (r *Client) SetRawSnapshot(raw []byte) (*SnapshotResponse, error) {
	var (
		rawResp []byte
		resp    StatusMessage
		code    int
		err     error
		sresp   = &SnapshotResponse{}
	)
	if rawResp, code, err = r.post("api/snapshots", nil, raw); err != nil {
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

func (r *Client) GetSnapshots() ([]SnapshotListResponse, error) {
	var (
		rawResp   []byte
		err       error
		code      int
		slistresp []SnapshotListResponse
	)

	if rawResp, code, err = r.get("/api/dashboard/snapshots", url.Values{}); err != nil {
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

func (r *Client) DeleteSnapshot(key string) error {
	var (
		err     error
		rawResp []byte
		code    int
	)

	if rawResp, code, err = r.delete(fmt.Sprintf("api/snapshots/%s", key)); err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("%d error: \n %s", code, string(rawResp))
	}
	return nil
}
