package openfaas

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-plugins/rest"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

type OpenFaasFunction struct {
	Name            string `json:"name"`
	Image           string `json:"image"`
	InvocationCount int64  `json:"invocationCount"`
	Replicas        int64  `json:"replicas"`
}

type OpenFaasFunctions []OpenFaasFunction

func GetFuncs(resolve resolver.Resolver, us *v1.Upstream) ([]*v1.Function, error) {
	fr := FaasRetriever{Lister: listGatewayFunctions(httpget)}
	return fr.GetFuncs(resolve, us)
}

func IsOpenFaas(us *v1.Upstream) bool {
	if us.Type == service.UpstreamTypeService {
		return isServiceHost(us)
	}
	if us.Type == kubernetes.UpstreamTypeKube {
		return isKubernetesHost(us)
	}

	return false
}

func httpget(s string) (io.ReadCloser, error) {
	resp, err := http.Get(s)

	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func listGatewayFunctions(httpget func(string) (io.ReadCloser, error)) func(gw string) (OpenFaasFunctions, error) {

	return func(gw string) (OpenFaasFunctions, error) {
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
		var funcs OpenFaasFunctions
		err = json.NewDecoder(body).Decode(&funcs)

		if err != nil {
			return nil, err
		}

		return funcs, nil
	}
}

type FaasRetriever struct {
	Lister func(from string) (OpenFaasFunctions, error)
}

func isServiceHost(us *v1.Upstream) bool {
	if us.Metadata == nil {
		return false
	}
	if us.Metadata.Namespace != "openfaas" || us.Name != "gateway" {
		return false
	}
	return true
}

func isKubernetesHost(us *v1.Upstream) bool {

	spec, err := kubernetes.DecodeUpstreamSpec(us.Spec)
	if err != nil {
		return false
	}

	if spec.ServiceNamespace != "openfaas" || spec.ServiceName != "gateway" {
		return false
	}
	return true
}

func (fr *FaasRetriever) GetFuncs(resolve resolver.Resolver, us *v1.Upstream) ([]*v1.Function, error) {

	if !IsOpenFaas(us) {
		return nil, nil
	}

	gw, err := resolve.Resolve(us)
	if err != nil {
		return nil, errors.Wrap(err, "error getting faas service")
	}

	if gw == "" {
		return nil, nil
	}

	// convert it to an http address
	httpgw := "http://" + gw

	functions, err := fr.Lister(httpgw)
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

func createFunction(fn OpenFaasFunction) *v1.Function {

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
