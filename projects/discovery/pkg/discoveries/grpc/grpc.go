package grpc

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/pkg/errors"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	discovery "github.com/solo-io/solo-kit/projects/discovery/pkg"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
	grpc_plugins "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/grpc"
)

func getgrpcspec(u *v1.Upstream) *grpc_plugins.ServiceSpec {
	spec, ok := u.UpstreamSpec.UpstreamType.(v1.ServiceSpecGetter)
	if !ok {
		return nil
	}
	grpcwrapper, ok := spec.GetServiceSpec().PluginType.(*plugins.ServiceSpec_Grpc)
	if !ok {
		return nil
	}
	grpc := grpcwrapper.Grpc
	return grpc
}

type FunctionDiscoveryFactory struct {
	detectionTimeout   time.Duration
	detectionRetryBase time.Duration
	functionPollTime   time.Duration
	fileclient         v1.ArtifactClient
}

func NewFunctionDiscovery(u *v1.Upstream) discovery.UpstreamFunctionDiscovery {
	return &UpstreamFunctionDiscovery{
		upstream: u,
	}
}

type UpstreamFunctionDiscovery struct {
	upstream   *v1.Upstream
	fileclient v1.ArtifactClient
}

func (f *UpstreamFunctionDiscovery) IsFunctional() bool {
	return getgrpcspec(f.upstream) != nil
}

func (f *UpstreamFunctionDiscovery) DetectType(ctx context.Context, url *url.URL) (*plugins.ServiceSpec, error) {
	log := contextutils.LoggerFrom(ctx)
	log.Debugf("attempting to detect GRPC for %s", f.upstream.Metadata.Name)

	var dialopts []grpc.DialOption
	if url.Scheme != "https" {
		dialopts = append(dialopts, grpc.WithInsecure())
	}

	addr := fmt.Sprintf("%s:%d", url.Host, url.Port)

	cc, err := grpc.Dial(addr, dialopts...)
	if err != nil {
		return nil, errors.Wrapf(err, "dialing grpc on %v", addr)
	}
	refClient := grpcreflect.NewClient(context.Background(), reflectpb.NewServerReflectionClient(cc))

	services, err := refClient.ListServices()
	if err != nil {
		return nil, errors.Wrapf(err, "listing services. are you sure %v implements reflection?", addr)
	}
	log.Infof("%v discovered as a gRPC service", addr)
	var (
		serviceNames []string
	)
	descriptors := &descriptor.FileDescriptorSet{}
	for _, s := range services {
		// ignore the reflection descriptor
		if s == "grpc.reflection.v1alpha.ServerReflection" {
			continue
		}
		files, err := getAllDescriptors(refClient, s)
		if err != nil {
			return nil, errors.Wrapf(err, "getting descriptors for svc %s", s)
		}
		descriptors.File = append(descriptors.File, files...)

		parts := strings.Split(s, ".")
		serviceSuffix := parts[len(parts)-1]
		serviceNames = append(serviceNames, serviceSuffix)
	}

	b, err := proto.Marshal(descriptors)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling proto descriptors")
	}

	// the : is a necessary separator for kube file storage
	// otherwise it's just ignored
	fileRef := fmt.Sprintf("%v:%v.descriptors",
		"grpc-discovery",
		strings.Join(serviceNames, "."))

	file := &v1.Artifact{
		Metadata: core.Metadata{
			Name:      fileRef,
			Namespace: f.upstream.Metadata.Namespace,
		},
		Data: map[string]string{
			// TODO(yuval-k): this is not a string
			"descriptors": string(b),
		},
	}

	wo := clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	}

	if _, err := f.fileclient.Write(file, wo); err != nil {
		return nil, errors.Wrap(err, "creating file for discovered descriptors")
	}

	svcInfo := &plugins.ServiceSpec{
		PluginType: &plugins.ServiceSpec_Grpc{
			Grpc: &grpc_plugins.ServiceSpec{},
		},
	}
	return svcInfo, nil

}

func (f *UpstreamFunctionDiscovery) DetectFunctions(ctx context.Context, secrets func() v1.SecretList, out func(discovery.UpstreamMutator) error) error {
	panic("TODO")
}

func getAllDescriptors(refClient *grpcreflect.Client, s string) ([]*descriptor.FileDescriptorProto, error) {
	root, err := refClient.FileContainingSymbol(s)
	if err != nil {
		return nil, errors.Wrapf(err, "getting file for symbol %s", s)
	}
	files := getDepTree(root)
	return files, nil
}

func getDepTree(root *desc.FileDescriptor) []*descriptor.FileDescriptorProto {
	var deps []*descriptor.FileDescriptorProto
	for _, dep := range root.GetDependencies() {
		deps = append(deps, getDepTree(dep)...)
	}
	deps = append(deps, root.AsFileDescriptorProto())
	return deps
}
