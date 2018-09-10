package apiclient

import (
	"io"
	"sort"
	"strings"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/apiserver"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/labels"
)

type ResourceClient struct {
	grpc         apiserver.ApiServerClient
	resourceType resources.Resource
	token        string
}

func NewResourceClient(cc *grpc.ClientConn, token string, resourceType resources.Resource) *ResourceClient {
	return &ResourceClient{
		grpc:         apiserver.NewApiServerClient(cc),
		resourceType: resourceType,
		token:        token,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Kind() string {
	return resources.Kind(rc.resourceType)
}

func (rc *ResourceClient) NewResource() resources.Resource {
	return resources.Clone(rc.resourceType)
}

func (rc *ResourceClient) Register() error {
	return nil
}

func (rc *ResourceClient) Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	opts.Ctx = metadata.AppendToOutgoingContext(opts.Ctx, "authorization", "bearer "+rc.token)
	resp, err := rc.grpc.Read(opts.Ctx, &apiserver.ReadRequest{
		Name:      name,
		Namespace: namespace,
		Kind:      rc.Kind(),
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok && strings.Contains(stat.Message(), "does not exist") {
			return nil, errors.NewNotExistErr(namespace, name)
		}
		return nil, errors.Wrapf(err, "performing grpc request")
	}
	resource := rc.NewResource()
	if err := protoutils.UnmarshalStruct(resp.Resource.Data, resource); err != nil {
		return nil, errors.Wrapf(err, "reading proto struct into %v", rc.Kind())
	}
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	if err := resources.Validate(resource); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	opts.Ctx = metadata.AppendToOutgoingContext(opts.Ctx, "authorization", "bearer "+rc.token)
	data, err := protoutils.MarshalStruct(resource)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal resource")
	}
	resp, err := rc.grpc.Write(opts.Ctx, &apiserver.WriteRequest{
		Resource: &apiserver.Resource{
			Data: data,
			Kind: rc.Kind(),
		},
		OverwriteExisting: opts.OverwriteExisting,
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok && strings.Contains(stat.Message(), "exists") {
			return nil, errors.NewExistErr(resource.GetMetadata())
		}
		return nil, errors.Wrapf(err, "performing grpc request")
	}
	written := rc.NewResource()
	if err := protoutils.UnmarshalStruct(resp.Resource.Data, written); err != nil {
		return nil, errors.Wrapf(err, "reading proto struct into %v", rc.Kind())
	}
	return written, nil
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = metadata.AppendToOutgoingContext(opts.Ctx, "authorization", "bearer "+rc.token)
	_, err := rc.grpc.Delete(opts.Ctx, &apiserver.DeleteRequest{
		Name:           name,
		Namespace:      namespace,
		Kind:           rc.Kind(),
		IgnoreNotExist: opts.IgnoreNotExist,
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok && strings.Contains(stat.Message(), "does not exist") {
			return errors.NewNotExistErr(namespace, name)
		}
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) (resources.ResourceList, error) {
	opts = opts.WithDefaults()
	opts.Ctx = metadata.AppendToOutgoingContext(opts.Ctx, "authorization", "bearer "+rc.token)
	resp, err := rc.grpc.List(opts.Ctx, &apiserver.ListRequest{
		Namespace: namespace,
		Kind:      rc.Kind(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "performing grpc request")
	}

	var resourceList resources.ResourceList
	for _, resourceData := range resp.ResourceList {
		resource := rc.NewResource()
		if err := protoutils.UnmarshalStruct(resourceData.Data, resource); err != nil {
			return nil, errors.Wrapf(err, "reading proto struct into %v", rc.Kind())
		}
		if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(resource.GetMetadata().Labels)) {
			resourceList = append(resourceList, resource)
		}
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan resources.ResourceList, <-chan error, error) {
	opts = opts.WithDefaults()
	opts.Ctx = metadata.AppendToOutgoingContext(opts.Ctx, "authorization", "bearer "+rc.token)
	secs := opts.RefreshRate.Seconds()
	nanos := int64(opts.RefreshRate.Seconds()) % int64(time.Duration(opts.RefreshRate.Seconds())*time.Second)
	resp, err := rc.grpc.Watch(opts.Ctx, &apiserver.WatchRequest{
		SyncFrequency: &types.Duration{
			Seconds: int64(secs),
			Nanos:   int32(nanos),
		},
		Namespace: namespace,
		Kind:      rc.Kind(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "performing grpc request")
	}

	resourcesChan := make(chan resources.ResourceList)
	errs := make(chan error)
	go func() {
		list, err := rc.List(namespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}()
	go func() {
		for {
			select {
			case <-opts.Ctx.Done():
				close(resourcesChan)
				close(errs)
				return
			default:
				resourceDataList, err := resp.Recv()
				if err == io.EOF {
					errs <- errors.Wrapf(err, "grpc stream closed")
					return
				}
				if err != nil {
					errs <- err
					continue
				}
				var resourceList resources.ResourceList
				for _, resourceData := range resourceDataList.ResourceList {
					resource := rc.NewResource()
					if err := protoutils.UnmarshalStruct(resourceData.Data, resource); err != nil {
						errs <- errors.Wrapf(err, "reading proto struct into %v", rc.Kind())
						continue
					}
					if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(resource.GetMetadata().Labels)) {
						resourceList = append(resourceList, resource)
					}
				}

				sort.SliceStable(resourceList, func(i, j int) bool {
					return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
				})
				resourcesChan <- resourceList
			}
		}
	}()

	return resourcesChan, errs, nil
}
