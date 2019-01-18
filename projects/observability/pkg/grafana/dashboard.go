package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

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

// GetRawDashboard loads a dashboard JSON from Grafana instance along with metadata for a dashboard.
// Contrary to GetDashboard() it not unpack loaded JSON to Board structure. Instead it
// returns it as byte slice. It guarantee that data of dashboard returned untouched by conversion
// with Board so no matter how properly fields from a current version of Grafana mapped to
// our Board fields. It useful for backuping purposes when you want a dashboard exactly with
// same data as it exported by Grafana.
func (r *Client) GetRawDashboard(uid string) ([]byte, BoardProperties, error) {
	var (
		raw    []byte
		result struct {
			Meta  BoardProperties `json:"meta"`
			Board json.RawMessage `json:"dashboard"`
		}
		code int
		err  error
	)
	if raw, code, err = r.get(fmt.Sprintf("api/dashboards/uid/%s", uid), nil); err != nil {
		return nil, BoardProperties{}, err
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

// SetRawDashboard updates existing dashboard or creates a new one.
// Contrary to SetDashboard() it accepts raw JSON instead of Board structure.
func (r *Client) SetRawDashboard(raw []byte) error {
	var (
		rawResp []byte
		resp    StatusMessage
		code    int
		err     error
	)
	if rawResp, code, err = r.post("api/dashboards/db", nil, raw); err != nil {
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

func (r *Client) DeleteDashboard(uid string) (StatusMessage, error) {
	var (
		raw   []byte
		reply StatusMessage
		err   error
	)

	if raw, _, err = r.delete(fmt.Sprintf("api/dashboards/uid/%s", uid)); err != nil {
		return StatusMessage{}, err
	}
	err = json.Unmarshal(raw, &reply)
	return reply, err
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

func (r *Client) SearchDashboards(query string, starred bool, tags ...string) ([]FoundBoard, error) {
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
	if raw, code, err = r.get("api/search", q); err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	err = json.Unmarshal(raw, &boards)
	return boards, err
}
