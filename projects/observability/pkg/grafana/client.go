package grafana

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// DefaultHTTPClient initialized Grafana with appropriate conditions.
// It allows you globally redefine HTTP client.
var DefaultHTTPClient = http.DefaultClient

// Client uses Grafana REST API for interacting with Grafana server.
type Client struct {
	baseURL   string
	key       string
	basicAuth bool
	client    *http.Client
}

// StatusMessage reflects status message as it returned by Grafana REST API.
type StatusMessage struct {
	ID      *uint   `json:"id"`
	OrgID   *uint   `json:"orgId"`
	Message *string `json:"message"`
	Version *int    `json:"version"`
	Status  *string `json:"resp"`
}

// NewClient initializes client for interacting with an instance of Grafana server;
// apiKeyOrBasicAuth accepts either 'username:password' basic authentication credentials,
// or a Grafana API key
func NewClient(apiURL, apiKeyOrBasicAuth string, client *http.Client) *Client {
	key := ""
	basicAuth := strings.Contains(apiKeyOrBasicAuth, ":")
	baseURL, _ := url.Parse(apiURL)
	if !basicAuth {
		key = fmt.Sprintf("Bearer %s", apiKeyOrBasicAuth)
	} else {
		parts := strings.Split(apiKeyOrBasicAuth, ":")
		baseURL.User = url.UserPassword(parts[0], parts[1])
	}
	return &Client{baseURL: baseURL.String(), basicAuth: basicAuth, key: key, client: client}
}

func (r *Client) get(query string, params url.Values) ([]byte, int, error) {
	return r.doRequest("GET", query, params, nil)
}

func (r *Client) patch(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("PATCH", query, params, bytes.NewBuffer(body))
}

func (r *Client) put(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("PUT", query, params, bytes.NewBuffer(body))
}

func (r *Client) post(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("POST", query, params, bytes.NewBuffer(body))
}

func (r *Client) delete(query string) ([]byte, int, error) {
	return r.doRequest("DELETE", query, nil, nil)
}

func (r *Client) doRequest(method, query string, params url.Values, buf io.Reader) ([]byte, int, error) {
	u, _ := url.Parse(r.baseURL)
	u.Path = path.Join(u.Path, query)
	if params != nil {
		u.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if !r.basicAuth {
		req.Header.Set("Authorization", r.key)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, resp.StatusCode, err
}
