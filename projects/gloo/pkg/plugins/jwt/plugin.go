package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/jwt"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/jwt_authn/v2alpha"
	"github.com/hashicorp/go-multierror"
	jose "gopkg.in/square/go-jose.v2"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

//go:generate protoc -I$GOPATH/src/github.com/envoyproxy/protoc-gen-validate -I. -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=jwt,paths=source_relative:. solo_jwt_authn.proto

const (
	JwtFilterName     = "io.solo.filters.http.solo_jwt_authn"
	ExtensionName     = "jwt"
	DisableName       = "-any:cf7a7de2-83ff-45ce-b697-f57d6a4775b5-"
	StateName         = "filterState"
	filterStage       = plugins.InAuth
	PayloadInMetadata = "principal"
)

// gather all the configurations from all the vhosts and place them in the filter config
// place a per filter config on the vhost
// that's it!

// as for rbac:
// convert config to per filter config
// thats it!

type Plugin struct {
	uniqProviders map[string]*envoyauth.JwtProvider

	perVhostProviders map[string][]string
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.perVhostProviders = make(map[string][]string)
	p.uniqProviders = make(map[string]*envoyauth.JwtProvider)
	return nil
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {
	var jwtRoute jwt.RouteExtension
	err := utils.UnmarshalExtension(in.RoutePlugins, ExtensionName, &jwtRoute)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto any to jwt plugin")
	}

	if jwtRoute.Disable {
		return pluginutils.SetRoutePerFilterConfig(out, JwtFilterName, &SoloJwtAuthnPerRoute{Requirement: DisableName})
	}
	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	// get the jwt config from the vhost
	var jwtExt jwt.VhostExtension
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &jwtExt)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto any to jwt plugin")
	}

	cfgName := in.Name

	claimsToHeader := make(map[string]*SoloJwtAuthnPerRoute_ClaimToHeaders)

	err = p.translateProviders(in.Name, jwtExt, claimsToHeader)
	if err != nil {
		return errors.Wrapf(err, "Error translating provider for "+cfgName)
	}
	clearRouteCache := len(claimsToHeader) != 0
	cfg := &SoloJwtAuthnPerRoute{
		Requirement:       cfgName,
		PayloadInMetadata: PayloadInMetadata,
		ClaimsToHeaders:   claimsToHeader,
		ClearRouteCache:   clearRouteCache,
	}
	return pluginutils.SetVhostPerFilterConfig(out, JwtFilterName, cfg)
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {

	cfg := &envoyauth.JwtAuthentication{
		Providers: make(map[string]*envoyauth.JwtProvider),
		FilterStateRules: &envoyauth.FilterStateRule{
			Name:     StateName,
			Requires: make(map[string]*envoyauth.JwtRequirement),
		},
	}
	for k, v := range p.uniqProviders {
		cfg.Providers[k] = v
	}
	for k, v := range p.perVhostProviders {
		cfg.FilterStateRules.Requires[k] = p.getRequirement(k, v)
	}

	// this should never happen, but let's make sure
	if _, ok := cfg.FilterStateRules.Requires[DisableName]; ok {
		// DisableName is reserved for a nil verifier, which will cause the JWT filter
		// do become a NOP
		panic("DisableName already in use")
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(JwtFilterName, cfg, filterStage)
	if err != nil {
		return nil, err
	}
	var filters []plugins.StagedHttpFilter
	filters = append(filters, stagedFilter)
	return filters, nil
}

func (p *Plugin) getRequirement(name string, providers []string) *envoyauth.JwtRequirement {

	if len(providers) == 1 {
		return &envoyauth.JwtRequirement{
			RequiresType: &envoyauth.JwtRequirement_ProviderName{
				ProviderName: providers[0],
			},
		}
	}
	var reqs []*envoyauth.JwtRequirement
	for _, provider := range providers {
		r := &envoyauth.JwtRequirement{
			RequiresType: &envoyauth.JwtRequirement_ProviderName{
				ProviderName: provider,
			},
		}
		reqs = append(reqs, r)
	}

	// sort for idempotancy
	sort.Slice(reqs, func(i, j int) bool { return reqs[i].GetProviderName() < reqs[j].GetProviderName() })
	return &envoyauth.JwtRequirement{
		RequiresType: &envoyauth.JwtRequirement_RequiresAny{
			RequiresAny: &envoyauth.JwtRequirementOrList{
				Requirements: reqs,
			},
		},
	}

	// TODO: add OR for all providers in the vhost name

}

func translateDefaultProvider(j jwt.VhostExtension) (*envoyauth.JwtProvider, error) {
	if j.Jwks == nil {
		return nil, errors.New("JWKS source not provided")
	}

	provider := &envoyauth.JwtProvider{
		Issuer:    j.Issuer,
		Audiences: j.Audiences,
	}
	provider.PayloadInMetadata = PayloadInMetadata

	err := translateJwks(&j, provider)
	return provider, err
}
func translateProvider(j *jwt.Provider) (*envoyauth.JwtProvider, error) {
	provider := &envoyauth.JwtProvider{
		Issuer:    j.Issuer,
		Audiences: j.Audiences,
		Forward:   j.KeepToken,
	}
	translateTokenSource(j, provider)

	err := translateJwks(j, provider)
	return provider, err
}

func translateTokenSource(j *jwt.Provider, provider *envoyauth.JwtProvider) {
	if headers := j.GetTokenSource().GetHeaders(); len(headers) != 0 {
		for _, header := range headers {
			provider.FromHeaders = append(provider.FromHeaders, &envoyauth.JwtHeader{
				Name:        header.Header,
				ValuePrefix: header.Prefix,
			})
		}
	}
	provider.FromParams = j.GetTokenSource().GetQueryParams()
}

func ProviderName(vhostname, providername string) string {
	return fmt.Sprintf("%s_%s", vhostname, providername)
}

func (p *Plugin) translateProviders(vhostname string, j jwt.VhostExtension, claimsToHeader map[string]*SoloJwtAuthnPerRoute_ClaimToHeaders) error {
	if len(j.GetProviders()) == 0 {
		// TODO(yuval-k): This is to support fieldsw on the VhostExtensions that were there before
		// providers were added. remove the fields in the VhostExtensions and then remove this code,
		// when they are not in used any more.
		provider, err := translateDefaultProvider(j)
		if err != nil {
			return err
		}
		name := ProviderName(vhostname, "default")
		p.uniqProviders[name] = provider
		p.perVhostProviders[vhostname] = append(p.perVhostProviders[vhostname], name)
		return nil
	}

	for name, provider := range j.GetProviders() {
		envoyProvider, err := translateProvider(provider)
		if err != nil {
			return err
		}
		name := ProviderName(vhostname, name)
		envoyProvider.PayloadInMetadata = name
		p.uniqProviders[name] = envoyProvider
		claimsToHeader[name] = translateClaimsToHeader(provider.ClaimsToHeaders)
		p.perVhostProviders[vhostname] = append(p.perVhostProviders[vhostname], name)
	}
	return nil
}

func translateClaimsToHeader(c2hs []*jwt.ClaimToHeader) *SoloJwtAuthnPerRoute_ClaimToHeaders {
	var ret []*SoloJwtAuthnPerRoute_ClaimToHeader
	for _, c2h := range c2hs {
		ret = append(ret, &SoloJwtAuthnPerRoute_ClaimToHeader{
			Claim:  c2h.Claim,
			Header: c2h.Header,
			Append: c2h.Append,
		})
	}
	if ret == nil {
		return nil
	}
	return &SoloJwtAuthnPerRoute_ClaimToHeaders{
		Claims: ret,
	}
}

type jwksSource interface {
	GetJwks() *jwt.Jwks
}

func translateJwks(j jwksSource, out *envoyauth.JwtProvider) error {
	if j.GetJwks() == nil {
		return errors.New("no jwks source provided")
	}
	switch jwks := j.GetJwks().Jwks.(type) {
	case *jwt.Jwks_Remote:
		if jwks.Remote.UpstreamRef == nil {
			return errors.New("upstream ref must not be empty in jwks source")
		}
		out.JwksSourceSpecifier = &envoyauth.JwtProvider_RemoteJwks{
			RemoteJwks: &envoyauth.RemoteJwks{
				HttpUri: &envoycore.HttpUri{
					Uri: jwks.Remote.Url,
					HttpUpstreamType: &envoycore.HttpUri_Cluster{
						Cluster: translator.UpstreamToClusterName(*jwks.Remote.UpstreamRef),
					},
				},
			},
		}
	case (*jwt.Jwks_Local):

		keyset, err := TranslateKey(jwks.Local.Key)
		if err != nil {
			return errors.Wrap(err, "failed to parse inline jwks")
		}

		keysetJson, err := json.Marshal(keyset)
		if err != nil {
			return errors.Wrap(err, "failed to serialize inline jwks")
		}

		out.JwksSourceSpecifier = &envoyauth.JwtProvider_LocalJwks{
			LocalJwks: &core.DataSource{
				Specifier: &core.DataSource_InlineString{
					InlineString: string(keysetJson),
				},
			},
		}
	default:
		return errors.New("unknown jwks source")
	}
	return nil
}

func TranslateKey(key string) (*jose.JSONWebKeySet, error) {
	// key can be an individual key, a key set or a pem block public key:
	// is it a pem block?
	var multierr error
	ks, err := parsePem(key)
	if err == nil {
		return ks, nil
	}
	multierr = multierror.Append(multierr, errors.Wrap(err, "PEM"))

	ks, err = parseKeySet(key)
	if err == nil {
		if len(ks.Keys) != 0 {
			return ks, nil
		}
		err = errors.New("no keys in set")
	}
	multierr = multierror.Append(multierr, errors.Wrap(err, "JWKS"))

	ks, err = parseKey(key)
	if err == nil {
		return ks, nil
	}
	multierr = multierror.Append(multierr, errors.Wrap(err, "JWK"))

	return nil, errors.Wrap(multierr, "cannot parse local jwks")
}

func parseKeySet(key string) (*jose.JSONWebKeySet, error) {
	var keyset jose.JSONWebKeySet
	err := json.Unmarshal([]byte(key), &keyset)
	return &keyset, err
}

func parseKey(key string) (*jose.JSONWebKeySet, error) {
	var jwk jose.JSONWebKey
	err := json.Unmarshal([]byte(key), &jwk)
	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk},
	}, err
}

func parsePem(key string) (*jose.JSONWebKeySet, error) {

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	var err error
	var publicKey interface{}
	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		publicKey, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	alg := ""
	switch publicKey.(type) {
	// RS256 implied for hash
	case *rsa.PublicKey:
		alg = "RS256"

	// envoy doesn't support this; uncomment when it does.
	// case *ecdsa.PublicKey:
	// 	alg = "ES256"

	default:
		return nil, errors.New("unsupported public key. only RSA public key is supported in PEM format")
	}

	jwk := jose.JSONWebKey{
		Key:       publicKey,
		Algorithm: alg,
		Use:       "sig",
	}
	keySet := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk},
	}
	return keySet, nil
}
