package hubspot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	RatelimitError error = fmt.Errorf("rate limited")
)

// For more info on this data structure, see:
// https://developers.hubspot.com/docs/methods/contacts/contacts-overview
type Contact struct {
	Vid          int64   `json:"vid,omitempty"`
	CanonicalVid int64   `json:"canonical-vid,omitempty"`
	MergedVid    []int64 `json:"merged-vids,omitempty"`

	PortalId  int64 `json:"portal-id,omitempty"`
	IsContact bool  `json:"is-contact,omitempty"`

	ProfileToken string `json:"profile-token,omitempty"`
	ProfileUrl   string `json:"profile-url,omitempty"`

	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Value    string    `json:"value,omitempty"`
	Versions []Version `json:"versions,omitempty"`
}
type Version struct {
	Value       string    `json:"value,omitempty"`
	SourceType  string    `json:"source-type,omitempty"`
	SourceId    *string   `json:"source-id,omitempty"`
	SourceLabel *string   `json:"source-label,omitempty"`
	Timestamp   Timestamp `json:"timestamp,omitempty"`
}

type ContactUpdate struct {
	Properties []PropertyUpdate `json:"properties,omitempty"`
}
type PropertyUpdate struct {
	Property string `json:"property,omitempty"`
	Value    string `json:"value,omitempty"`
}

// String or null; Additional data realted to the source-type. May not be populated for all source-types.
type Timestamp time.Time

func (a *Timestamp) UnmarshalJSON(b []byte) error {
	var u int64
	if err := json.Unmarshal(b, &u); err != nil {
		return err
	}
	t := time.Unix(u/1000, (u%1000)*1000000)
	*a = Timestamp(t)

	return nil
}

func (a Timestamp) MarshalJSON() ([]byte, error) {
	t := time.Time(a)
	var u int64 = t.Unix() * 1000
	u += int64(t.Nanosecond()) / 1000000
	return json.Marshal(t)
}

const contactBaseUrl = "https://api.hubapi.com/contacts/v1/"

type Hubspot struct {
	ApiKey string
	Client *http.Client
}

func NewHubspot(token string) *Hubspot {
	return &Hubspot{
		ApiKey: token,
		Client: http.DefaultClient,
	}
}

type UpsertResponse struct {
	Vid   int64 `json:"vid,omitempty"`
	IsNew bool  `json:"isNew,omitempty"`
}

func (h *Hubspot) UpsertContact(ctx context.Context, email string, c ContactUpdate) (*UpsertResponse, error) {
	var ur UpsertResponse
	err := h.send(ctx, "contact/createOrUpdate/email/"+url.PathEscape(email), nil, &c, &ur)
	return &ur, err
}

func (h *Hubspot) GetContact(ctx context.Context, email string, properties ...string) (*Contact, error) {

	var c Contact

	query := make(url.Values)
	for _, prop := range properties {
		query.Add("property", prop)
	}

	err := h.send(ctx, "contact/email/"+url.PathEscape(email)+"/profile", query, nil, &c)

	return &c, err
}

func (h *Hubspot) send(ctx context.Context, urlpath string, query url.Values, reqbody interface{}, respBody interface{}) error {

	u, err := url.Parse(contactBaseUrl)
	if err != nil {
		return err
	}
	emailu, err := url.Parse(urlpath)
	if err != nil {
		return err
	}
	finalu := u.ResolveReference(emailu)
	if query == nil {
		query = make(url.Values)
	}
	query.Add("hapikey", h.ApiKey)
	finalu.RawQuery = query.Encode()

	var body io.Reader
	method := "GET"
	if reqbody != nil {
		method = "POST"
		var buf bytes.Buffer
		e := json.NewEncoder(&buf)
		err := e.Encode(reqbody)
		if err != nil {
			return err
		}
		body = &buf
	}

	req, err := http.NewRequest(method, finalu.String(), body)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	resp, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			return RatelimitError
		}
		return fmt.Errorf("bad response %s %d", resp.Status, resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(respBody)
	if err != nil {
		return err
	}

	return nil
}
