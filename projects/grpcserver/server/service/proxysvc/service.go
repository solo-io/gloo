package proxysvc

import (
	"compress/zlib"
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/compress"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubecrd "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/solo.io/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"

	"github.com/ghodss/yaml"
	"go.uber.org/zap"
)

// ErrNotCompressed is thrown if the caller tries to decompress a spec which is not compressed
var ErrNotCompressed = errors.New("could not decompress spec 'compressedSpec' was not found on the resource")

type proxyGrpcService struct {
	ctx             context.Context
	clientCache     client.ClientCache
	rawGetter       rawgetter.RawGetter
	statusConverter status.InputResourceStatusGetter
	settingsValues  settings.ValuesClient
}

func NewProxyGrpcService(ctx context.Context, clientCache client.ClientCache, rawGetter rawgetter.RawGetter, statusConverter status.InputResourceStatusGetter, settingsValues settings.ValuesClient) v1.ProxyApiServer {
	return &proxyGrpcService{
		ctx:             ctx,
		clientCache:     clientCache,
		rawGetter:       rawGetter,
		statusConverter: statusConverter,
		settingsValues:  settingsValues,
	}
}

func (s *proxyGrpcService) GetProxy(ctx context.Context, request *v1.GetProxyRequest) (*v1.GetProxyResponse, error) {
	proxy, err := s.clientCache.GetProxyClient().Read(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToGetProxyError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	proxyDetails, err := s.getDetails(proxy)
	if err != nil {
		wrapped := FailedToGetProxyError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.GetProxyResponse{ProxyDetails: proxyDetails}, nil
}

func (s *proxyGrpcService) ListProxies(ctx context.Context, request *v1.ListProxiesRequest) (*v1.ListProxiesResponse, error) {
	var proxyDetailsList []*v1.ProxyDetails
	for _, ns := range s.settingsValues.GetWatchNamespaces() {
		proxiesInNamespace, err := s.clientCache.GetProxyClient().List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListProxiesError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, p := range proxiesInNamespace {
			pd, err := s.getDetails(p)
			if err != nil {
				wrapped := FailedToGetProxyError(err, &core.ResourceRef{Name: p.Metadata.Name, Namespace: p.Metadata.Namespace})
				contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
				return nil, wrapped
			}
			proxyDetailsList = append(proxyDetailsList, pd)
		}
	}
	return &v1.ListProxiesResponse{ProxyDetails: proxyDetailsList}, nil
}

func (s *proxyGrpcService) getDetails(proxy *gloov1.Proxy) (*v1.ProxyDetails, error) {
	raw := s.rawGetter.GetRaw(s.ctx, proxy, gloov1.ProxyCrd)

	// If spec compression is enabled, decompress it in the raw content before sending to the UI
	if proxy.Metadata.Annotations[compress.CompressedKey] == "true" && strings.Contains(raw.Content, "compressedSpec") {
		err := DecompressSpec(raw)
		if err != nil {
			return nil, err
		}
	}

	return &v1.ProxyDetails{
		Proxy:  proxy,
		Raw:    raw,
		Status: s.statusConverter.GetApiStatusFromResource(proxy),
	}, nil
}

// DecompressSpec takes a raw representation of a kube resource, and
// decompresses the Spec.compressedSpec.
func DecompressSpec(raw *v1.Raw) error {
	resourceFromYaml := &kubecrd.Resource{}
	err := yaml.Unmarshal([]byte(raw.Content), resourceFromYaml)
	if err != nil {
		return err
	}

	compressedSpec, ok := (*resourceFromYaml.Spec)["compressedSpec"].(string)
	if !ok {
		// Not compressed
		return ErrNotCompressed
	}

	// Base64 decode
	decodedSpec, err := base64.StdEncoding.DecodeString(compressedSpec)
	if err != nil {
		return err
	}

	// Zlib decompress
	r, err := zlib.NewReader(strings.NewReader(string(decodedSpec)))
	if err != nil {
		return err
	}
	defer r.Close()

	ds, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// Unmarshal bytes into kubecrd.Spec
	decompressedSpec := &kubecrd.Spec{}
	err = yaml.Unmarshal(ds, decompressedSpec)
	if err != nil {
		return err
	}

	// Replace compressedSpec with the uncompressed spec
	resourceFromYaml.Spec = decompressedSpec
	rawContent, err := yaml.Marshal(resourceFromYaml)
	if err != nil {
		return err
	}

	raw.Content = string(rawContent)

	return nil
}
