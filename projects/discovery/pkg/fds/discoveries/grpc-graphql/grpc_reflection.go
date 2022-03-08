package grpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/golang/protobuf/proto"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	grpc_plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/grpc"
)

func getGrpcSpec(u *v1.Upstream) *grpc_plugins.ServiceSpec {
	upstreamType, ok := u.GetUpstreamType().(v1.ServiceSpecGetter)
	if !ok {
		return nil
	}

	if upstreamType.GetServiceSpec() == nil {
		return nil
	}

	grpcWrapper, ok := upstreamType.GetServiceSpec().GetPluginType().(*plugins.ServiceSpec_Grpc)
	if !ok {
		return nil
	}
	return grpcWrapper.Grpc
}

func NewFunctionDiscoveryFactory() fds.FunctionDiscoveryFactory {
	return &FunctionDiscoveryFactory{
		DetectionTimeout: time.Second * 15,
		FunctionPollTime: time.Second * 15,
	}
}

var _ fds.FunctionDiscoveryFactory = new(FunctionDiscoveryFactory)

func (f *FunctionDiscoveryFactory) FunctionDiscoveryFactoryName() string {
	return "GrpcGraphqlFunctionDiscoveryFactory"
}

// FunctionDiscoveryFactory returns a FunctionDiscovery that can be used to discover functions
type FunctionDiscoveryFactory struct {
	DetectionTimeout   time.Duration
	DetectionRetryBase time.Duration
	FunctionPollTime   time.Duration
	Artifacts          v1.ArtifactClient
}

// NewFunctionDiscovery returns a FunctionDiscovery that can be used to discover functions
func (f *FunctionDiscoveryFactory) NewFunctionDiscovery(u *v1.Upstream, clients fds.AdditionalClients) fds.UpstreamFunctionDiscovery {
	return &GraphqlSchemaDiscovery{
		upstream:         u,
		graphqlClient:    clients.GraphqlClient,
		detectionTimeout: f.DetectionTimeout,
	}
}

// GraphqlSchemaDiscovery represents a function discovery for upstream
type GraphqlSchemaDiscovery struct {
	upstream         *v1.Upstream
	graphqlClient    v1alpha1.GraphQLSchemaClient
	detectionTimeout time.Duration
}

// IsFunctional returns true if the upstream has already been discovered as a
// gRPC upstream with reflection.
func (f *GraphqlSchemaDiscovery) IsFunctional() bool {
	if getGrpcSpec(f.upstream) != nil {
		return true
	}
	return false
}

func (f *GraphqlSchemaDiscovery) DetectType(ctx context.Context, url *url.URL) (*plugins.ServiceSpec, error) {
	log := contextutils.LoggerFrom(ctx)
	log.Debugf("attempting to detect GRPC for %s", f.upstream.GetMetadata().GetName())

	refClient, closeConn, err := getClient(ctx, url)
	if err != nil {
		return nil, err
	}

	defer closeConn()

	_, err = refClient.ListServices()
	if err != nil {
		return nil, errors.Wrapf(err, "listing services. are you sure %v implements reflection?", url)
	}

	svcInfo := &plugins.ServiceSpec{
		PluginType: &plugins.ServiceSpec_Grpc{
			Grpc: &grpc_plugins.ServiceSpec{},
		},
	}
	return svcInfo, nil
}

func (f *GraphqlSchemaDiscovery) DetectFunctions(ctx context.Context, url *url.URL, _ func() fds.Dependencies, updatecb func(fds.UpstreamMutator) error) error {
	err := contextutils.NewExponentialBackoff(contextutils.ExponentialBackoff{
		MaxDuration: &f.detectionTimeout,
	}).Backoff(ctx, func(ctx context.Context) error {
		err := f.DetectFunctionsOnce(ctx, url, updatecb)
		if err != nil {
			contextutils.LoggerFrom(ctx).Warnf("Unable to create GraphQLSchemas from gRPC reflection for upstream %s in namespace %s: %s",
				f.upstream.GetMetadata().GetName(),
				f.upstream.GetMetadata().GetNamespace(),
				err)
		}
		return err
	})

	if err != nil {
		if ctx.Err() != nil {
			return multierror.Append(err, ctx.Err())
		}
		// ignore other errors as we would like to continue forever.
	}
	if err := contextutils.Sleep(ctx, 30*time.Second); err != nil {
		return err
	}
	return nil
}

// Creating an interface here for the sake of mocking in tests
type GrpcReflectionClient interface {
	ListServices() ([]string, error)
	FileContainingSymbol(symbol string) (*desc.FileDescriptor, error)
}

func (f *GraphqlSchemaDiscovery) GetSchemaBuilderForProtoFileDescriptor(refClient GrpcReflectionClient, descriptors *descriptor.FileDescriptorSet, services []string) (*SchemaBuilder, *v1alpha1.Executor_Local, error) {
	sb := NewSchemaBuilder()
	executor := &v1alpha1.Executor_Local{
		EnableIntrospection: true,
		Resolutions:         map[string]*v1alpha1.Resolution{},
	}
	for _, s := range services {
		// ignore the reflection descriptor
		if s == "grpc.reflection.v1alpha.ServerReflection" {
			continue
		}
		root, err := refClient.FileContainingSymbol(s)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "getting file for svc symbol %s", s)
		}
		files := getDepTree(root)

		descriptors.File = append(descriptors.GetFile(), files...)

		parts := strings.Split(s, ".")
		serviceName := parts[len(parts)-1]
		// find the service in the file and get its functions
		for _, svc := range root.GetServices() {
			if svc.GetName() == serviceName {
				methods := svc.GetMethods()
				for _, method := range methods {
					methodName := method.GetName()
					inputType := method.GetInputType()
					outputType := method.GetOutputType()
					resolverName := fmt.Sprintf("Query|%s.%s", serviceName, methodName)
					_, _, err = sb.CreateInputMessageType(inputType)
					if err != nil {
						return nil, nil, errors.Wrapf(err, "unable to translate input type %s of method %s for service %s",
							inputType.GetName(), method, svc.GetName(), f.upstream.GetMetadata().GetName(), f.upstream.GetMetadata().GetNamespace())
					}
					_, _, err = sb.CreateOutputMessageType(outputType)
					if err != nil {
						return nil, nil, errors.Wrapf(err, "unable to translate type %s of method %s for service %s",
							inputType.GetName(), method, svc.GetName(), f.upstream.GetMetadata().GetName(), f.upstream.GetMetadata().GetNamespace())
					}
					sb.AddQueryField(method.GetName(), inputType, outputType, resolverName)
					outgoingJsonBody := GenerateOutgoingJsonBodyForInputType(inputType, "{$args."+inputType.GetName())
					t := &v1alpha1.GrpcRequestTemplate{
						OutgoingMessageJson: outgoingJsonBody,
						ServiceName:         svc.GetFullyQualifiedName(),
						MethodName:          method.GetName(),
					}
					resolution := &v1alpha1.Resolution{
						Resolver: &v1alpha1.Resolution_GrpcResolver{
							GrpcResolver: &v1alpha1.GrpcResolver{
								UpstreamRef:      f.upstream.GetMetadata().Ref(),
								RequestTransform: t,
							},
						},
					}
					executor.Resolutions[resolverName] = resolution

				}
			}
		}
	}
	return sb, executor, nil
}

func (f *GraphqlSchemaDiscovery) BuildGraphQLSchemaFromGrpcReflection(refClient GrpcReflectionClient) (*v1alpha1.GraphQLSchema, error) {
	services, err := refClient.ListServices()
	if err != nil {
		return nil, errors.Wrapf(err, "listing services. are you sure upstream %s.%s implements reflection?", f.upstream.GetMetadata().GetName(), f.upstream.GetMetadata().GetNamespace())
	}
	descriptors := &descriptor.FileDescriptorSet{}

	schemaBuilder, executor, err := f.GetSchemaBuilderForProtoFileDescriptor(refClient, descriptors, services)
	if err != nil {
		return nil, errors.Wrapf(err, "error in translating gRPC reflection for upstream %s in namespace %s to GraphQLSchema",
			f.upstream.GetMetadata().GetName(), f.upstream.GetMetadata().GetNamespace())
	}

	doc := schemaBuilder.Build()
	schemaDef := fmt.Sprintf("%s", printer.Print(doc))
	d, err := proto.Marshal(descriptors)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling descriptors")
	}
	dest := make([]byte, base64.StdEncoding.EncodedLen(len(d)))
	base64.StdEncoding.Encode(dest, d)

	schema := &v1alpha1.GraphQLSchema{
		Metadata: &core.Metadata{
			Name:      f.upstream.GetMetadata().GetName(),
			Namespace: f.upstream.GetMetadata().GetNamespace(),
		},
		ExecutableSchema: &v1alpha1.ExecutableSchema{
			SchemaDefinition: schemaDef,
			Executor: &v1alpha1.Executor{
				Executor: &v1alpha1.Executor_Local_{
					Local: executor,
				},
			},
			GrpcDescriptorRegistry: &v1alpha1.GrpcDescriptorRegistry{
				DescriptorSet: &v1alpha1.GrpcDescriptorRegistry_ProtoDescriptorBin{
					ProtoDescriptorBin: d,
				},
			},
		},
	}
	return schema, nil
}

func (f *GraphqlSchemaDiscovery) DetectFunctionsOnce(ctx context.Context, url *url.URL, updatecb func(fds.UpstreamMutator) error) error {
	refClient, closeConn, err := getClient(ctx, url)
	if err != nil {
		return err
	}
	defer closeConn()

	schema, err := f.BuildGraphQLSchemaFromGrpcReflection(refClient)
	if err != nil {
		return errors.Wrap(err, "error creating schema from gRPC reflection")
	}
	_, err = f.graphqlClient.Write(schema, clients.WriteOpts{})
	return err
}

func getClient(ctx context.Context, url *url.URL) (*grpcreflect.Client, func() error, error) {
	var dialOpts []grpc.DialOption
	if url.Scheme != "https" {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}

	cc, err := grpc.Dial(url.Host, dialOpts...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "dialing grpc on %v", url.Host)
	}
	refClient := grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(cc))
	return refClient, cc.Close, nil
}

func getDepTree(root *desc.FileDescriptor) []*descriptor.FileDescriptorProto {
	var deps []*descriptor.FileDescriptorProto
	for _, dep := range root.GetDependencies() {
		deps = append(deps, getDepTree(dep)...)
	}
	deps = append(deps, root.AsFileDescriptorProto())
	return deps
}
