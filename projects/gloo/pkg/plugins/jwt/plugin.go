package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"sort"

	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/jwt_authn/v3"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-multierror"
	errors "github.com/rotisserie/eris"
	envoycore "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/jwt"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"gopkg.in/square/go-jose.v2"
)

const (
	ExtensionName     = "jwt"
	DisableName       = "-any:cf7a7de2-83ff-45ce-b697-f57d6a4775b5-"
	StateName         = "filterState"
	PayloadInMetadata = "principal"

	RemoteJwksTimeoutSecs = 5
	SoloJwtFilterName     = "io.solo.filters.http.solo_jwt_authn_staged"

	AfterExtAuthStage  = uint32(0)
	BeforeExtAuthStage = uint32(1)
)

var (
	_ plugins.Plugin            = new(Plugin)
	_ plugins.VirtualHostPlugin = new(Plugin)
	_ plugins.RoutePlugin       = new(Plugin)
	_ plugins.HttpFilterPlugin  = new(Plugin)

	beforeExtauthFilterStage = plugins.BeforeStage(plugins.AuthNStage)
	afterExtauthFilterStage  = plugins.DuringStage(plugins.AuthNStage)
)

// gather all the configurations from all the vhosts and place them in the filter config
// place a per filter config on the vhost
// that's it!

// as for rbac:
// convert config to per filter config
// thats it!

type Plugin struct {
	requireJwtBeforeExtauthFilter bool
	uniqProviders                 map[uint32]map[string]*v3.JwtProvider
	perVhostProviders             map[uint32]map[*v1.VirtualHost][]string
	perRouteJwtRequirements       map[uint32]map[string]*v3.JwtRequirement
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.perVhostProviders = map[uint32]map[*v1.VirtualHost][]string{
		BeforeExtAuthStage: make(map[*v1.VirtualHost][]string),
		AfterExtAuthStage:  make(map[*v1.VirtualHost][]string),
	}
	p.uniqProviders = map[uint32]map[string]*v3.JwtProvider{
		BeforeExtAuthStage: make(map[string]*v3.JwtProvider),
		AfterExtAuthStage:  make(map[string]*v3.JwtProvider),
	}
	p.perRouteJwtRequirements = map[uint32]map[string]*v3.JwtRequirement{
		BeforeExtAuthStage: make(map[string]*v3.JwtRequirement),
		AfterExtAuthStage:  make(map[string]*v3.JwtRequirement),
	}
	p.requireJwtBeforeExtauthFilter = false
	return nil
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	jwtRoute := p.convertRouteJwtConfig(in.GetOptions())
	// If no config for jwt exists on route, do not create any routePerFilter config
	if jwtRoute == nil || (jwtRoute.BeforeExtAuth == nil && jwtRoute.AfterExtAuth == nil) {
		// no config found, nothing to do here
		return nil
	}
	stagedCfg := &StagedJwtAuthnPerRoute{
		JwtConfigs: make(map[uint32]*SoloJwtAuthnPerRoute),
	}

	if routeCfg := p.getPerRouteFilterConfig(in, jwtRoute.BeforeExtAuth); routeCfg != nil {
		stagedCfg.JwtConfigs[BeforeExtAuthStage] = routeCfg
	}

	if routeCfg := p.getPerRouteFilterConfig(in, jwtRoute.AfterExtAuth); routeCfg != nil {
		stagedCfg.JwtConfigs[AfterExtAuthStage] = routeCfg
	}
	if len(stagedCfg.JwtConfigs) == 0 {
		return nil
	}

	return pluginutils.SetRoutePerFilterConfig(out, SoloJwtFilterName, stagedCfg)
}

func (p *Plugin) ProcessVirtualHost(
	params plugins.VirtualHostParams,
	in *v1.VirtualHost,
	out *envoy_config_route_v3.VirtualHost,
) error {
	var jwtExt = p.convertVhostJwtConfig(in.GetOptions())
	// If no config exists for vhost, do not create any vhost per filter config
	if jwtExt == nil || (jwtExt.BeforeExtAuth == nil && jwtExt.AfterExtAuth == nil) {
		// no config found, nothing to do here
		return nil
	}
	if jwtExt.BeforeExtAuth != nil {
		p.requireJwtBeforeExtauthFilter = true
	}
	stagedCfg := &StagedJwtAuthnPerRoute{
		JwtConfigs: make(map[uint32]*SoloJwtAuthnPerRoute),
	}
	if jwtExt.BeforeExtAuth != nil {
		cfg, err := p.getVhostFilterConfig(in, jwtExt.BeforeExtAuth, BeforeExtAuthStage)
		if err != nil {
			return err
		}
		stagedCfg.JwtConfigs[BeforeExtAuthStage] = cfg
	}
	if jwtExt.AfterExtAuth != nil {
		cfg, err := p.getVhostFilterConfig(in, jwtExt.AfterExtAuth, AfterExtAuthStage)
		if err != nil {
			return err
		}
		stagedCfg.JwtConfigs[AfterExtAuthStage] = cfg
	}
	return pluginutils.SetVhostPerFilterConfig(out, SoloJwtFilterName, stagedCfg)
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {

	var filters []plugins.StagedHttpFilter
	// Get filter config for after extauth (default)
	stagedFilter, err := p.getFilterForStage(AfterExtAuthStage, afterExtauthFilterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, *stagedFilter)

	if p.requireJwtBeforeExtauthFilter {
		stagedFilter, err = p.getFilterForStage(BeforeExtAuthStage, beforeExtauthFilterStage)
		if err != nil {
			return nil, err
		}
		filters = append(filters, *stagedFilter)
	}
	return filters, nil
}

func (p *Plugin) getFilterForStage(stage uint32, filterStage plugins.FilterStage) (*plugins.StagedHttpFilter, error) {
	cfg := p.getJwtAuthnCfgForStage(stage)
	filterConfig := &JwtWithStage{
		Stage:    stage,
		JwtAuthn: cfg,
	}
	stagedFilter, err := plugins.NewStagedFilterWithConfig(SoloJwtFilterName, filterConfig, filterStage)
	if err != nil {
		return nil, err
	}
	return &stagedFilter, nil
}

func (p *Plugin) getJwtAuthnCfgForStage(stage uint32) *v3.JwtAuthentication {
	cfg := &v3.JwtAuthentication{
		Providers: make(map[string]*v3.JwtProvider),
		FilterStateRules: &v3.FilterStateRule{
			Name:     getFilterStateNameForStage(stage),
			Requires: make(map[string]*v3.JwtRequirement),
		},
	}
	for providerName, provider := range p.uniqProviders[stage] {
		cfg.Providers[providerName] = provider
	}
	for virtualHost, jwtReq := range p.perVhostProviders[stage] {
		cfg.FilterStateRules.Requires[virtualHost.Name] = p.getRequirement(virtualHost, jwtReq, stage)
	}
	for route, jwtReq := range p.perRouteJwtRequirements[stage] {
		cfg.FilterStateRules.Requires[route] = jwtReq
	}

	// this should never happen, but let's make sure
	if _, ok := cfg.FilterStateRules.Requires[DisableName]; ok {
		// DisableName is reserved for a nil verifier, which will cause the JWT filter
		// do become a NOP
		panic("DisableName already in use")
	}

	return cfg
}

func (p *Plugin) getRequirement(vHost *v1.VirtualHost, providers []string, stage uint32) *v3.JwtRequirement {
	var jwtReq *v3.JwtRequirement
	// Get requirements from virtual host providers
	if len(providers) == 1 {
		jwtReq = &v3.JwtRequirement{
			RequiresType: &v3.JwtRequirement_ProviderName{
				ProviderName: providers[0],
			},
		}
	} else {
		var reqs []*v3.JwtRequirement
		for _, provider := range providers {
			r := &v3.JwtRequirement{
				RequiresType: &v3.JwtRequirement_ProviderName{
					ProviderName: provider,
				},
			}
			reqs = append(reqs, r)
		}

		// sort for idempotency
		sort.Slice(reqs, func(i, j int) bool { return reqs[i].GetProviderName() < reqs[j].GetProviderName() })
		jwtReq = &v3.JwtRequirement{
			RequiresType: &v3.JwtRequirement_RequiresAny{
				RequiresAny: &v3.JwtRequirementOrList{
					Requirements: reqs,
				},
			},
		}
	}

	// If the jwt can fail or can be missing, but will still pass jwt requirement,
	// OR the current JWT provider reqs with a allow_missing_or_failed jwt requirement
	var allowMissingOrFailed bool
	switch stage {
	case BeforeExtAuthStage:
		allowMissingOrFailed = vHost.GetOptions().GetJwtStaged().GetBeforeExtAuth().GetAllowMissingOrFailedJwt()
	case AfterExtAuthStage:
		// Deprecated jwt config defaults to AfterExtAuthStage, look there as well for allow_missing_or_failed
		missingOrFailed := vHost.GetOptions().GetJwtStaged().GetAfterExtAuth().GetAllowMissingOrFailedJwt()
		missingOrFailedDeprecated := vHost.GetOptions().GetJwt().GetAllowMissingOrFailedJwt()
		allowMissingOrFailed = missingOrFailed || missingOrFailedDeprecated

	}
	if allowMissingOrFailed {
		missingOrFailedReq := &v3.JwtRequirement{
			RequiresType: &v3.JwtRequirement_AllowMissingOrFailed{
				AllowMissingOrFailed: &empty.Empty{},
			},
		}
		jwtReq = &v3.JwtRequirement{
			RequiresType: &v3.JwtRequirement_RequiresAny{
				// Requires Any will OR the two requirements
				RequiresAny: &v3.JwtRequirementOrList{
					Requirements: []*v3.JwtRequirement{
						jwtReq,
						missingOrFailedReq,
					},
				},
			},
		}
	}
	return jwtReq
}

func translateProvider(j *jwt.Provider) (*v3.JwtProvider, error) {
	provider := &v3.JwtProvider{
		Issuer:    j.Issuer,
		Audiences: j.Audiences,
		Forward:   j.KeepToken,
	}
	translateTokenSource(j, provider)

	err := translateJwks(j, provider)
	return provider, err
}

func translateTokenSource(j *jwt.Provider, provider *v3.JwtProvider) {
	if headers := j.GetTokenSource().GetHeaders(); len(headers) != 0 {
		for _, header := range headers {
			provider.FromHeaders = append(provider.FromHeaders, &v3.JwtHeader{
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

func routeJwtRequirementName(routeName string) string {
	return fmt.Sprintf("route_%s", routeName)
}

func (p *Plugin) translateProviders(in *v1.VirtualHost, j jwt.VhostExtension, claimsToHeader map[string]*SoloJwtAuthnPerRoute_ClaimToHeaders, stage uint32) error {
	for name, provider := range j.GetProviders() {
		envoyProvider, err := translateProvider(provider)
		if err != nil {
			return err
		}
		name := ProviderName(in.GetName(), name)
		envoyProvider.PayloadInMetadata = name
		p.uniqProviders[stage][name] = envoyProvider
		claimsToHeader[name] = translateClaimsToHeader(provider.ClaimsToHeaders)
		p.perVhostProviders[stage][in] = append(p.perVhostProviders[stage][in], name)
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

func (p *Plugin) convertVhostJwtConfig(opts *v1.VirtualHostOptions) *jwt.JwtStagedVhostExtension {
	ret := &jwt.JwtStagedVhostExtension{}

	switch opts.GetJwtConfig().(type) {
	case *v1.VirtualHostOptions_JwtStaged:
		{
			ret = opts.GetJwtStaged()
		}
	case *v1.VirtualHostOptions_Jwt:
		{
			ret.AfterExtAuth = opts.GetJwt()
		}
	}

	return ret
}

func (p *Plugin) convertRouteJwtConfig(opts *v1.RouteOptions) *jwt.JwtStagedRouteExtension {
	ret := &jwt.JwtStagedRouteExtension{}

	switch opts.GetJwtConfig().(type) {
	case *v1.RouteOptions_JwtStaged:
		{
			ret = opts.GetJwtStaged()
		}
	case *v1.RouteOptions_Jwt:
		{
			ret.AfterExtAuth = opts.GetJwt()
		}
	}

	return ret
}

func (p *Plugin) getPerRouteFilterConfig(in *v1.Route, routeCfg *jwt.RouteExtension) *SoloJwtAuthnPerRoute {
	if routeCfg != nil && routeCfg.Disable {
		return &SoloJwtAuthnPerRoute{Requirement: DisableName}
	}
	return nil
}

func (p *Plugin) getVhostFilterConfig(in *v1.VirtualHost, vhostCfg *jwt.VhostExtension, stage uint32) (*SoloJwtAuthnPerRoute, error) {
	claimsToHeader := make(map[string]*SoloJwtAuthnPerRoute_ClaimToHeaders)
	err := p.translateProviders(in, *vhostCfg, claimsToHeader, stage)
	if err != nil {
		return nil, errors.Wrapf(err, "Error translating provider for "+in.Name)
	}
	clearRouteCache := len(claimsToHeader) != 0
	routeCfg := &SoloJwtAuthnPerRoute{
		Requirement:       in.Name,
		PayloadInMetadata: PayloadInMetadata,
		ClaimsToHeaders:   claimsToHeader,
		ClearRouteCache:   clearRouteCache,
	}
	return routeCfg, nil
}

func (p *Plugin) requireEarlyJwtFilter() {
	p.requireJwtBeforeExtauthFilter = true
}

// Avoid collisions between multiple stages of jwt with same filter state name
func getFilterStateNameForStage(stage uint32) string {
	return fmt.Sprintf("stage%d-%s", stage, StateName)
}

type jwksSource interface {
	GetJwks() *jwt.Jwks
}

func translateJwks(j jwksSource, out *v3.JwtProvider) error {
	if j.GetJwks() == nil {
		return errors.New("no jwks source provided")
	}
	switch jwks := j.GetJwks().Jwks.(type) {
	case *jwt.Jwks_Remote:
		if jwks.Remote.UpstreamRef == nil {
			return errors.New("upstream ref must not be empty in jwks source")
		}
		out.JwksSourceSpecifier = &v3.JwtProvider_RemoteJwks{
			RemoteJwks: &v3.RemoteJwks{
				CacheDuration: jwks.Remote.GetCacheDuration(),
				HttpUri: &envoycore.HttpUri{
					Timeout: &duration.Duration{Seconds: RemoteJwksTimeoutSecs},
					Uri:     jwks.Remote.Url,
					HttpUpstreamType: &envoycore.HttpUri_Cluster{
						Cluster: translator.UpstreamToClusterName(jwks.Remote.UpstreamRef),
					},
				},
			},
		}
	case *jwt.Jwks_Local:

		keyset, err := TranslateKey(jwks.Local.Key)
		if err != nil {
			return errors.Wrapf(err, "failed to parse inline jwks")
		}

		keysetJson, err := json.Marshal(keyset)
		if err != nil {
			return errors.Wrapf(err, "failed to serialize inline jwks")
		}

		out.JwksSourceSpecifier = &v3.JwtProvider_LocalJwks{
			LocalJwks: &envoycore.DataSource{
				Specifier: &envoycore.DataSource_InlineString{
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
	multierr = multierror.Append(multierr, errors.Wrapf(err, "PEM"))

	ks, err = parseKeySet(key)
	if err == nil {
		if len(ks.Keys) != 0 {
			return ks, nil
		}
		err = errors.New("no keys in set")
	}
	multierr = multierror.Append(multierr, errors.Wrapf(err, "JWKS"))

	ks, err = parseKey(key)
	if err == nil {
		return ks, nil
	}
	multierr = multierror.Append(multierr, errors.Wrapf(err, "JWK"))

	return nil, errors.Wrapf(multierr, "cannot parse local jwks")
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
