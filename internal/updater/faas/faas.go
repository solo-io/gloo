package faas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-plugins/rest"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

type FaasFunction struct {
	Name            string `json:"name"`
	Image           string `json:"image"`
	InvocationCount int64  `json:"invocationCount"`
	Replicas        int64  `json:"replicas"`
}

type FaasFunctions []FaasFunction

func GetFuncs(us *v1.Upstream) ([]*v1.Function, error) {
	fr := FassRetriever{Lister: listGatewayFunctions(httpget)}
	return fr.GetFuncs(us)
}

func IsFaas(us *v1.Upstream) bool {
	gw, err := getHost(us)
	if err != nil {
		return false
	}
	return gw != ""
}

func httpget(s string) (io.ReadCloser, error) {
	resp, err := http.Get(s)

	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func listGatewayFunctions(httpget func(string) (io.ReadCloser, error)) func(gw string) (FaasFunctions, error) {

	return func(gw string) (FaasFunctions, error) {
		u, err := url.Parse(gw)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "system", "functions")
		s := u.String()
		body, err := httpget(s)

		if err != nil {
			return nil, err
		}
		defer body.Close()
		var funcs FaasFunctions
		err = json.NewDecoder(body).Decode(&funcs)

		if err != nil {
			return nil, err
		}

		return funcs, nil
	}
}

type FassRetriever struct {
	Lister func(from string) (FaasFunctions, error)
}

func getServiceHost(us *v1.Upstream) (string, error) {
	if us.Metadata == nil {
		return "", nil
	}
	if us.Metadata.Namespace != "openfaas" || us.Name != "gateway" {
		return "", nil
	}

	spec, err := service.DecodeUpstreamSpec(us.Spec)
	if err != nil {
		return "", errors.Wrap(err, "decoding service upstream spec")
	}

	if len(spec.Hosts) == 0 {
		return "", errors.New("no hosts")
	}

	host := spec.Hosts[0]

	gw := fmt.Sprintf("http://%s:%d/", host.Addr, host.Port)
	return gw, nil
}

func getKubernetesHost(us *v1.Upstream) (string, error) {

	spec, err := kubernetes.DecodeUpstreamSpec(us.Spec)
	if err != nil {
		return "", errors.Wrap(err, "decoding service upstream spec")
	}

	if spec.ServiceNamespace != "openfaas" || spec.ServiceName != "gateway" {
		return "", nil
	}

	gw := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d/", spec.ServiceName, spec.ServiceNamespace, spec.ServicePort)
	return gw, nil
}

func getHost(us *v1.Upstream) (string, error) {

	if us.Type == service.UpstreamTypeService {
		return getServiceHost(us)
	}
	if us.Type == kubernetes.UpstreamTypeKube {
		return getKubernetesHost(us)
	}

	return "", nil
}

func (fr *FassRetriever) GetFuncs(us *v1.Upstream) ([]*v1.Function, error) {
	// decode does validation for us
	gw, err := getHost(us)
	if err != nil {
		return nil, errors.Wrap(err, "error getting faas service")
	}
	if gw == "" {
		return nil, nil
	}

	functions, err := fr.Lister(gw)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching functions")
	}

	var funcs []*v1.Function

	for _, fn := range functions {
		if fn.Name != "" {
			funcs = append(funcs, createFunction(fn))
		}
	}
	return funcs, nil
}

func createFunction(fn FaasFunction) *v1.Function {

	headersTemplate := map[string]string{":method": "POST"}

	return &v1.Function{
		Name: fn.Name,
		Spec: rest.EncodeFunctionSpec(rest.Template{
			Path:            path.Join("/function", fn.Name),
			Header:          headersTemplate,
			PassthroughBody: true,
		}),
	}
}
