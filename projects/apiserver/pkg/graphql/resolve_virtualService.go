package graphql

import (
	"context"

	"github.com/gogo/protobuf/types"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
	ratelimitapi "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

func (r *namespaceResolver) VirtualServices(ctx context.Context, obj *customtypes.Namespace) ([]*models.VirtualService, error) {
	list, err := r.VirtualServiceClient.List(obj.Name, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualServices(list)
}

func (r *namespaceResolver) VirtualService(ctx context.Context, obj *customtypes.Namespace, name string) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServiceClient.Read(obj.Name, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(virtualService)
}

type virtualServiceMutationResolver struct{ *ApiResolver }

func (r *virtualServiceMutationResolver) Create(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	v1VirtualService, err := NewConverter(r.ApiResolver, ctx).ConvertInputVirtualService(virtualService)
	if err != nil {
		return nil, err
	}
	out, err := r.VirtualServiceClient.Write(v1VirtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: false,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

// Reads the virtual service identified for update from storage
// Steps through the update object and applies only the requested updates
func (r *virtualServiceMutationResolver) Update(ctx context.Context, obj *customtypes.VirtualServiceMutation, guid string, resourceVersion string, updates models.InputUpdateVirtualService) (*models.VirtualService, error) {
	_, namespace, name, err := resources.SplitKey(guid)
	if err != nil {
		return &models.VirtualService{}, err
	}
	virtualService, err := r.VirtualServiceClient.Read(namespace, name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if updates.Domains != nil {
		virtualService.VirtualHost.Domains = updates.Domains
	}

	if updates.Metadata != nil {
		if updates.Metadata.Name != nil {
			return nil, errors.Errorf("Changing name is not a valid operation. Please delete and recreate this resource.")
			// return an error for now.
			// We could delete and recreate the resource under a new name with:
			// virtualService.Metadata.Name = *updates.Metadata.Name
			// virtualService.Metadata.ResourceVersion = ""
			// ...But that's probably not what the user wants
			// Consider adding a "Reference Name" field to allow the user to update the displayed name without changing the CRD ID
		}
	}

	if updates.SslConfig != nil {
		return nil, errors.Errorf("SSLConfig updates are not yet supported.")
	}

	if updates.RateLimitConfig != nil {
		if err := applyRateLimitConfigUpdate(updates.RateLimitConfig, virtualService.VirtualHost); err != nil {
			return nil, err
		}
	}

	out, err := r.VirtualServiceClient.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func applyRateLimitConfigUpdate(in *models.InputRateLimitConfig, out *v1.VirtualHost) error {
	currentRlc := &ratelimitapi.IngressRateLimit{}
	if out.VirtualHostPlugins == nil {
		out.VirtualHostPlugins = &v1.VirtualHostPlugins{}
	}
	p := out.VirtualHostPlugins
	err := utils.UnmarshalExtension(p, ratelimit.ExtensionName, currentRlc)
	if err != nil && err != utils.NotFoundError {
		return errors.Wrapf(err, "failed to unmarshal proto message to %v plugin", ratelimit.ExtensionName)
	}

	if err := applyInputRateLimitSpec(in.AuthorizedLimits, currentRlc, true); err != nil {
		return err
	}
	if err := applyInputRateLimitSpec(in.AnonymousLimits, currentRlc, false); err != nil {
		return err
	}
	rlStruct, err := util.MessageToStruct(currentRlc)
	if err != nil {
		return err
	}
	if p.Extensions == nil {
		p.Extensions = &v1.Extensions{}
	}
	if p.Extensions.Configs == nil {
		p.Extensions.Configs = map[string]*types.Struct{ratelimit.ExtensionName: rlStruct}
	} else {
		p.Extensions.Configs[ratelimit.ExtensionName] = rlStruct
	}
	return nil
}

func (r *virtualServiceMutationResolver) Delete(ctx context.Context, obj *customtypes.VirtualServiceMutation, guid string) (*models.VirtualService, error) {
	_, namespace, name, err := resources.SplitKey(guid)
	if err != nil {
		return &models.VirtualService{}, err
	}
	virtualService, err := r.VirtualServiceClient.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.VirtualServiceClient.Delete(namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(virtualService)
}

func (r *virtualServiceMutationResolver) AddRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceId string, resourceVersion string, index int, route models.InputRoute) (*models.VirtualService, error) {
	_, namespace, name, err := resources.SplitKey(virtualServiceId)
	if err != nil {
		return &models.VirtualService{}, err
	}
	v1Route, err := NewConverter(r.ApiResolver, ctx).ConvertInputRoute(route)
	if err != nil {
		return nil, err
	}

	virtualService, err := r.VirtualServiceClient.Read(namespace, name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index > len(virtualService.VirtualHost.Routes) {
		index = len(virtualService.VirtualHost.Routes)
	}
	virtualService.VirtualHost.Routes = append(virtualService.VirtualHost.Routes, nil)
	copy(virtualService.VirtualHost.Routes[index+1:], virtualService.VirtualHost.Routes[index:])
	virtualService.VirtualHost.Routes[index] = v1Route

	out, err := r.VirtualServiceClient.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func (r *virtualServiceMutationResolver) UpdateRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceId string, resourceVersion string, index int, route models.InputRoute) (*models.VirtualService, error) {
	_, namespace, name, err := resources.SplitKey(virtualServiceId)
	if err != nil {
		return &models.VirtualService{}, err
	}
	v1Route, err := NewConverter(r.ApiResolver, ctx).ConvertInputRoute(route)
	if err != nil {
		return nil, err
	}

	virtualService, err := r.VirtualServiceClient.Read(namespace, name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index > len(virtualService.VirtualHost.Routes) {
		return nil, errors.Errorf("index out of bounds")
	}

	virtualService.VirtualHost.Routes[index] = v1Route

	out, err := r.VirtualServiceClient.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func (r *virtualServiceMutationResolver) DeleteRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceId string, resourceVersion string, index int) (*models.VirtualService, error) {
	_, namespace, name, err := resources.SplitKey(virtualServiceId)
	if err != nil {
		return &models.VirtualService{}, err
	}
	virtualService, err := r.VirtualServiceClient.Read(namespace, name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index > len(virtualService.VirtualHost.Routes) {
		return nil, errors.Errorf("index out of bounds")
	}

	virtualService.VirtualHost.Routes = append(virtualService.VirtualHost.Routes[:index], virtualService.VirtualHost.Routes[index+1:]...)

	out, err := r.VirtualServiceClient.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func (r *virtualServiceMutationResolver) SwapRoutes(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceId string, resourceVersion string, index1 int, index2 int) (*models.VirtualService, error) {
	_, namespace, name, err := resources.SplitKey(virtualServiceId)
	if err != nil {
		return &models.VirtualService{}, err
	}
	virtualService, err := r.VirtualServiceClient.Read(namespace, name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index1 > len(virtualService.VirtualHost.Routes) || index2 > len(virtualService.VirtualHost.Routes) || index1 < 0 || index2 < 0 {
		return nil, errors.Errorf("index out of bounds")
	}

	virtualService.VirtualHost.Routes[index1], virtualService.VirtualHost.Routes[index2] = virtualService.VirtualHost.Routes[index2], virtualService.VirtualHost.Routes[index1]
	out, err := r.VirtualServiceClient.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)

}

// Removes the route at fromIndex and inserts it at toIndex.
// Any routes in between shift to fill the hole or to make room.
func (r *virtualServiceMutationResolver) ShiftRoutes(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceId string, resourceVersion string, fromIndex int, toIndex int) (*models.VirtualService, error) {
	_, namespace, name, err := resources.SplitKey(virtualServiceId)
	if err != nil {
		return &models.VirtualService{}, err
	}
	virtualService, err := r.VirtualServiceClient.Read(namespace, name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if fromIndex > len(virtualService.VirtualHost.Routes) || toIndex > len(virtualService.VirtualHost.Routes) || fromIndex < 0 || toIndex < 0 {
		return nil, errors.Errorf("index out of bounds")
	}

	if toIndex < fromIndex {
		// anchor on the fromIndex and swap until all updated
		for i := toIndex; i < fromIndex; i++ {
			virtualService.VirtualHost.Routes[fromIndex], virtualService.VirtualHost.Routes[i] = virtualService.VirtualHost.Routes[i], virtualService.VirtualHost.Routes[fromIndex]
		}
	} else {
		// anchor on the toIndex and swap until all updated
		for i := toIndex; i > fromIndex; i-- {
			virtualService.VirtualHost.Routes[fromIndex], virtualService.VirtualHost.Routes[i] = virtualService.VirtualHost.Routes[i], virtualService.VirtualHost.Routes[fromIndex]
		}
	}

	out, err := r.VirtualServiceClient.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)

}
