package graphql

import (
	"context"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

type subscriptionResolver struct{ *ApiResolver }

func (r subscriptionResolver) Upstreams(ctx context.Context, namespace string, selector *models.InputMapStringString) (<-chan []*models.Upstream, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	watch, errs, err := r.ApiResolver.UpstreamClient.Watch(namespace, clients.WatchOpts{
		// TODO(ilackarms): refresh rate
		RefreshRate: time.Minute * 10,
		Ctx:         ctx,
		Selector:    convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	upstreamsChan := make(chan []*models.Upstream)
	go func() {
		defer close(upstreamsChan)
		for {
			select {
			case list, ok := <-watch:
				if !ok {
					return
				}
				upstreams, err := NewConverter(r.ApiResolver, ctx).ConvertOutputUpstreams(list)
				if err != nil {
					// TODO(mitchdraft) log this
					return
				}
				select {
				case upstreamsChan <- upstreams:
				default:
					contextutils.LoggerFrom(ctx).Errorf("upstream channel full, cannot send list of %v upstreams", len(list))
				}
			case err, ok := <-errs:
				if !ok {
					return
				}
				contextutils.LoggerFrom(ctx).Errorf("error in upstream subscription: %v", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	return upstreamsChan, nil
}

func (r subscriptionResolver) VirtualServices(ctx context.Context, namespace string, selector *models.InputMapStringString) (<-chan []*models.VirtualService, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	watch, errs, err := r.ApiResolver.VirtualServiceClient.Watch(namespace, clients.WatchOpts{
		// TODO(ilackarms): refresh rate
		RefreshRate: time.Minute * 10,
		Ctx:         ctx,
		Selector:    convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	virtualServicesChan := make(chan []*models.VirtualService)
	go func() {
		defer close(virtualServicesChan)
		for {
			select {
			case list, ok := <-watch:
				if !ok {
					return
				}
				virtualServices, err := NewConverter(r.ApiResolver, ctx).ConvertOutputVirtualServices(list)
				if err != nil {
					// TODO(mitchdraft) log this
					return
				}
				select {
				case virtualServicesChan <- virtualServices:
				default:
					contextutils.LoggerFrom(ctx).Errorf("virtualService channel full, cannot send list of %v virtualServices", len(list))
				}
			case err, ok := <-errs:
				if !ok {
					return
				}
				contextutils.LoggerFrom(ctx).Errorf("error in virtualService subscription: %v", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	return virtualServicesChan, nil
}
