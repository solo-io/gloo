package openApi

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-openapi/swag"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	graphql "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi/graphqlschematranslation"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/go-multierror"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	rest_plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/rest"
	"github.com/solo-io/go-utils/contextutils"
)

var commonOpenApiURIs = []string{
	"/openapi.json",
	"/swagger.json",
	"/swagger/docs/v1",
	"/swagger/docs/v2",
	"/v1/openApi",
	"/v2/openApi",
	"/api-docs",
}

func NewFunctionDiscoveryFactory() fds.FunctionDiscoveryFactory {
	return &OpenApiFunctionDiscoveryFactory{
		DetectionTimeout: time.Second * 75,
		FunctionPollTime: time.Second,
	}
}

type OpenApiFunctionDiscoveryFactory struct {
	DetectionTimeout time.Duration
	FunctionPollTime time.Duration
}

var _ fds.FunctionDiscoveryFactory = new(OpenApiFunctionDiscoveryFactory)

func (f *OpenApiFunctionDiscoveryFactory) FunctionDiscoveryFactoryName() string {
	return "OpenApiFunctionDiscoveryFactory"
}

func (f *OpenApiFunctionDiscoveryFactory) NewFunctionDiscovery(u *v1.Upstream, clients fds.AdditionalClients) fds.UpstreamFunctionDiscovery {
	return &OpenApiFunctionDiscovery{
		detectionTimeout: f.DetectionTimeout,
		functionPollTime: f.FunctionPollTime,
		openApiUrisToTry: commonOpenApiURIs,
		upstream:         u,
		graphqlClient:    clients.GraphqlClient,
	}
}

type OpenApiFunctionDiscovery struct {
	detectionTimeout time.Duration
	functionPollTime time.Duration
	upstream         *v1.Upstream
	openApiUrisToTry []string

	graphqlClient v1alpha1.GraphQLSchemaClient
}

func getswagspec(u *v1.Upstream) *rest_plugins.ServiceSpec_SwaggerInfo {
	spec, ok := u.GetUpstreamType().(v1.ServiceSpecGetter)
	if !ok {
		return nil
	}
	serviceSpec := spec.GetServiceSpec()
	if serviceSpec == nil {
		return nil
	}
	restwrapper, ok := serviceSpec.GetPluginType().(*plugins.ServiceSpec_Rest)
	if !ok {
		return nil
	}
	rest := restwrapper.Rest
	return rest.GetSwaggerInfo()
}

func (f *OpenApiFunctionDiscovery) IsFunctional() bool {
	a := getswagspec(f.upstream) != nil
	return a
}

func (f *OpenApiFunctionDiscovery) DetectType(ctx context.Context, baseurl *url.URL) (*plugins.ServiceSpec, error) {
	var spec *plugins.ServiceSpec

	err := contextutils.NewExponentialBackoff(contextutils.ExponentioalBackoff{MaxDuration: &f.detectionTimeout}).Backoff(ctx, func(ctx context.Context) error {
		var err error
		spec, err = f.detectUpstreamTypeOnce(ctx, baseurl)
		return err
	})

	return spec, err
}

func (f *OpenApiFunctionDiscovery) detectUpstreamTypeOnce(ctx context.Context, baseUrl *url.URL) (*plugins.ServiceSpec, error) {
	// run detection and get functions
	var errs error
	logger := contextutils.LoggerFrom(ctx)

	logger.Debugf("attempting to detect openApi base url %v", baseUrl)

	switch baseUrl.Scheme {
	case "http":
		fallthrough
	case "https":
		// nothing to do as this baseurl already has an http address.
	case "tcp":
		// if it is a tcp address, assume it is plain http
		baseUrl.Scheme = "http"
	default:
		return nil, fmt.Errorf("unsupported baseurl for openApi discovery %v", baseUrl)
	}
	for _, uri := range f.openApiUrisToTry {
		Url := baseUrl.ResolveReference(&url.URL{Path: uri}).String()
		req, err := http.NewRequest("GET", Url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "invalid url for request")
		}
		req.Header.Set("X-Gloo-Discovery", "OpenApi-Discovery")

		req = req.WithContext(ctx)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			errs = multierror.Append(errs, errors.Wrapf(err, "could not perform HTTP GET on resolved addr: %v", Url))
			continue
		}
		// might have found a openApi service
		if res.StatusCode == http.StatusOK {
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to read response body")
			}
			if _, err := GetOpenApi3Doc(b); err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			// definitely found openApi
			svcInfo := &plugins.ServiceSpec{
				PluginType: &plugins.ServiceSpec_Rest{
					Rest: &rest_plugins.ServiceSpec{
						SwaggerInfo: &rest_plugins.ServiceSpec_SwaggerInfo{
							SwaggerSpec: &rest_plugins.ServiceSpec_SwaggerInfo_Url{
								Url: Url,
							},
						},
					},
				},
			}
			return svcInfo, nil
		}
		errs = multierror.Append(errs, errors.Errorf("path: %v response code: %v headers: %v", uri, res.Status, res.Header))
	}
	// not a openApi upstream
	return nil, errors.Wrapf(errs, "service at %s does not implement the openApi spec at a known endpoint, "+
		"or was unreachable", baseUrl.String())

}

func GetOpenApi3Doc(body []byte) (*openapi3.T, error) {
	doc, err := openapi3.NewLoader().LoadFromData(body)
	if err != nil {
		return nil, err
	}
	return doc, nil

}

func (f *OpenApiFunctionDiscovery) DetectFunctions(ctx context.Context, url *url.URL, _ func() fds.Dependencies, updatecb func(fds.UpstreamMutator) error) error {
	in := f.upstream
	fmt.Println("keerthan reached here")
	spec := getswagspec(in)
	if spec == nil || spec.GetSwaggerSpec() == nil {
		return errors.New("upstream doesn't have a openApi spec")
	}
	switch document := spec.GetSwaggerSpec().(type) {
	case *rest_plugins.ServiceSpec_SwaggerInfo_Url:
		fmt.Println("keerthan detecting functions from url", document.Url)
		return f.detectFunctionsFromUrl(ctx, document.Url, in)
	case *rest_plugins.ServiceSpec_SwaggerInfo_Inline:
		return f.detectFunctionsFromInline(ctx, document.Inline, in, updatecb)
	}
	fmt.Println("upstream doens't have an opeanpi source")
	return errors.New("upstream doesn't have a openApi source")
}

func (f *OpenApiFunctionDiscovery) detectFunctionsFromUrl(ctx context.Context, url string, in *v1.Upstream) error {
	for {
		err := contextutils.NewExponentioalBackoff(contextutils.ExponentioalBackoff{}).Backoff(ctx, func(ctx context.Context) error {

			spec, err := RetrieveOpenApiDocFromUrl(ctx, url)
			if err != nil {
				return err
			}
			err = f.detectFunctionsFromSpec(ctx, spec)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// ignore other erros as we would like to continue forever.
		}

		if err := contextutils.Sleep(ctx, f.functionPollTime); err != nil {
			return err
		}
	}

}

func (f *OpenApiFunctionDiscovery) detectFunctionsFromInline(ctx context.Context, document string, in *v1.Upstream, updatecb func(fds.UpstreamMutator) error) error {
	spec, err := GetOpenApi3Doc([]byte(document))
	if err != nil {
		return err
	}
	return f.detectFunctionsFromSpec(ctx, spec)
}

func (f *OpenApiFunctionDiscovery) detectFunctionsFromSpec(ctx context.Context, openApiSpec *openapi3.T) error {
	translator := graphql.NewOasToGqlTranslator(f.upstream)
	schemaAst, _, resolutions, err := translator.CreateGraphqlSchema([]*openapi3.T{openApiSpec})
	if err != nil {
		return err
	}
	schema := ast.Node(schemaAst)
	printedSchema := printer.Print(schema).(string)
	schemaCrd := &v1alpha1.GraphQLSchema{
		Metadata: &core.Metadata{
			Name:      f.upstream.GetMetadata().GetName(),
			Namespace: f.upstream.GetMetadata().GetNamespace(),
		},
		ExecutableSchema: &v1alpha1.ExecutableSchema{
			Executor: &v1alpha1.Executor{
				Executor: &v1alpha1.Executor_Local_{
					Local: &v1alpha1.Executor_Local{
						Resolutions:         resolutions,
						EnableIntrospection: true,
					},
				},
			},
			SchemaDefinition: printedSchema,
		},
	}
	_, err = f.graphqlClient.Write(schemaCrd, clients.WriteOpts{})
	if err != nil {
		return errors.Wrapf(err, "unable to write schema %+v", schemaCrd)
	}
	return nil
}

func RetrieveOpenApiDocFromUrl(ctx context.Context, url string) (*openapi3.T, error) {
	docBytes, err := LoadFromFileOrHTTP(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "loading swagger doc from url")
	}
	return GetOpenApi3Doc(docBytes)
}

func LoadFromFileOrHTTP(ctx context.Context, url string) ([]byte, error) {
	return swag.LoadStrategy(url, ioutil.ReadFile, loadHTTPBytes(ctx))(url)
}

func loadHTTPBytes(ctx context.Context) func(path string) ([]byte, error) {
	return func(path string) ([]byte, error) {
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}
		req = req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req)
		defer func() {
			if resp != nil {
				if e := resp.Body.Close(); e != nil {
					contextutils.LoggerFrom(ctx).Debug(e)
				}
			}
		}()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("could not access document at %q [%s] ", path, resp.Status)
		}

		return ioutil.ReadAll(resp.Body)
	}
}
