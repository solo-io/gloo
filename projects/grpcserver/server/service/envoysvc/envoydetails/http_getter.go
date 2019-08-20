package envoydetails

import (
	"io/ioutil"
	"net/http"
)

//go:generate mockgen -destination mocks/mock_http_getter.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails HttpGetter

type HttpGetter interface {
	Get(ip, port, path string) (string, error)
}

type httpGetter struct{}

var _ HttpGetter = httpGetter{}

func NewHttpGetter() HttpGetter {
	return httpGetter{}
}

func (httpGetter) Get(ip, port, path string) (string, error) {
	response, err := http.Get("http://" + ip + ":" + port + path)
	if err != nil {
		return "", err
	}
	defer func() { response.Body.Close() }()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
