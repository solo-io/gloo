package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/models"
)

type virtualServiceQueryResolver struct{ *ApiResolver }

func (r *virtualServiceQueryResolver) List(ctx context.Context, obj *customtypes.VirtualServiceQuery, selector *models.InputMapStringString) ([]*models.VirtualService, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.VirtualServices.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualServices(list)
}

func (r *virtualServiceQueryResolver) Get(ctx context.Context, obj *customtypes.VirtualServiceQuery, name string) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(virtualService)
}

type virtualServiceMutationResolver struct{ *ApiResolver }

func (r *virtualServiceMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	v1VirtualService, err := NewConverter(r.ApiResolver, ctx).ConvertInputVirtualService(virtualService)
	if err != nil {
		return nil, err
	}
	out, err := r.VirtualServices.Write(v1VirtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func (r *virtualServiceMutationResolver) Create(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	return r.write(false, ctx, obj, virtualService)
}

// Reads the virtual service identifed for update from storage
// Steps through the update object and applies only the requested updates
func (r *virtualServiceMutationResolver) Update(ctx context.Context, obj *customtypes.VirtualServiceMutation, name string, resourceVersion string, updates models.InputUpdateVirtualService) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, name, clients.ReadOpts{Ctx: ctx})
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

	if updates.Plugins != nil {
		return nil, errors.Errorf("Plugin updates are not yet supported.")
	}

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}
func (r *virtualServiceMutationResolver) Delete(ctx context.Context, obj *customtypes.VirtualServiceMutation, name string) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.VirtualServices.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(virtualService)
}
func (r *virtualServiceMutationResolver) AddRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index int, route models.InputRoute) (*models.VirtualService, error) {
	v1Route, err := NewConverter(r.ApiResolver, ctx).ConvertInputRoute(route)
	if err != nil {
		return nil, err
	}

	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
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

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func (r *virtualServiceMutationResolver) UpdateRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index int, route models.InputRoute) (*models.VirtualService, error) {
	v1Route, err := NewConverter(r.ApiResolver, ctx).ConvertInputRoute(route)
	if err != nil {
		return nil, err
	}

	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
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

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func (r *virtualServiceMutationResolver) DeleteRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index int) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
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

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)
}

func (r *virtualServiceMutationResolver) SwapRoutes(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index1 int, index2 int) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
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
	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
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
func (r *virtualServiceMutationResolver) ShiftRoutes(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, fromIndex int, toIndex int) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
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

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualService(out)

}
