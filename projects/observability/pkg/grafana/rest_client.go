package grafana

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

//go:generate mockgen -destination mocks/mock_rest_client.go -package mocks github.com/solo-io/solo-projects/projects/observability/pkg/grafana RestClient

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
	client        *http.Client
	credentials   *Credentials
	grafanaApiUrl string
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

type BasicAuth struct {
	Username string
	Password string
}

// must have either the username and password set or the api key set
type Credentials struct {
	BasicAuth *BasicAuth
	ApiKey    string
}

// initialize a RestClient for interacting with an instance of Grafana's api server;
// apiKeyOrBasicAuth accepts either 'username:password' basic authentication credentials,
// or a Grafana API key
func NewRestClient(grafanaApiUrl string, httpClient *http.Client, credentials *Credentials) RestClient {
	return &restClient{
		client:        httpClient,
		credentials:   credentials,
		grafanaApiUrl: grafanaApiUrl,
	}
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
	u, _ := url.Parse(r.grafanaApiUrl)

	u.Path = path.Join(u.Path, query)
	if params != nil {
		u.RawQuery = params.Encode()
	}

	var apiKeyHeader string
	if r.credentials.BasicAuth != nil {
		u.User = url.UserPassword(r.credentials.BasicAuth.Username, r.credentials.BasicAuth.Password)
	} else if r.credentials.ApiKey != "" {
		apiKeyHeader = fmt.Sprintf("Bearer %s", r.credentials.ApiKey)
	} else {
		return nil, 0, IncompleteGrafanaCredentials
	}

	req, err := http.NewRequest(method, u.String(), buf)

	if apiKeyHeader != "" {
		req.Header.Set("Authorization", apiKeyHeader)
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
