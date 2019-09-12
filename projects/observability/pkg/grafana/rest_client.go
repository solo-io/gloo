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

//go:generate mockgen -destination mocks/mock_grafana_client.go -package mocks github.com/solo-io/solo-projects/projects/observability/pkg/grafana RestClient

// DefaultHTTPClient initialized Grafana with appropriate conditions.
// It allows you globally redefine HTTP client.
var DefaultHTTPClient = http.DefaultClient

type RestClient interface {
	Get(query string, params url.Values) ([]byte, int, error)
	Patch(query string, params url.Values, body []byte) ([]byte, int, error)
	Put(query string, params url.Values, body []byte) ([]byte, int, error)
	Post(query string, params url.Values, body []byte) ([]byte, int, error)
	Delete(query string) ([]byte, int, error)
}

// RestClient uses Grafana REST API for interacting with Grafana server.
type restClient struct {
	baseURL   string
	key       string
	basicAuth bool
	client    *http.Client
}

var _ RestClient = &restClient{}

// StatusMessage reflects status message as it returned by Grafana REST API.
type StatusMessage struct {
	ID      *uint   `json:"id"`
	OrgID   *uint   `json:"orgId"`
	Message *string `json:"message"`
	Version *int    `json:"version"`
	Status  *string `json:"resp"`
}

// initialize a RestClient for interacting with an instance of Grafana's api server;
// apiKeyOrBasicAuth accepts either 'username:password' basic authentication credentials,
// or a Grafana API key
func NewRestClient(apiURL, apiKeyOrBasicAuth string, httpClient *http.Client) RestClient {
	key := ""
	basicAuth := strings.Contains(apiKeyOrBasicAuth, ":")
	baseURL, _ := url.Parse(apiURL)
	if !basicAuth {
		key = fmt.Sprintf("Bearer %s", apiKeyOrBasicAuth)
	} else {
		parts := strings.Split(apiKeyOrBasicAuth, ":")
		baseURL.User = url.UserPassword(parts[0], parts[1])
	}
	return &restClient{baseURL: baseURL.String(), basicAuth: basicAuth, key: key, client: httpClient}
}

func (r *restClient) Get(query string, params url.Values) ([]byte, int, error) {
	return r.doRequest("GET", query, params, nil)
}

func (r *restClient) Patch(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("PATCH", query, params, bytes.NewBuffer(body))
}

func (r *restClient) Put(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("PUT", query, params, bytes.NewBuffer(body))
}

func (r *restClient) Post(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("POST", query, params, bytes.NewBuffer(body))
}

func (r *restClient) Delete(query string) ([]byte, int, error) {
	return r.doRequest("DELETE", query, nil, nil)
}

func (r *restClient) doRequest(method, query string, params url.Values, buf io.Reader) ([]byte, int, error) {
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
