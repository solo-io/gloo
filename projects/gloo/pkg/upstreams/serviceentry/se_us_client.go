package serviceentry

import (
	"context"
	"fmt"
	"strings"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/NoOpUpstreamClient"
	"github.com/solo-io/go-utils/contextutils"
	skclients "github.com/solo-io/solo-kit/pkg/api/v1/clients"

	"istio.io/istio/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func NewServiceEntryUpstreamClient(client kube.Client) v1.UpstreamClient {
	return &serviceEntryClient{client: client}
}

type serviceEntryClient struct {
	client kube.Client
}

// List implements v1.UpstreamClient.
func (s *serviceEntryClient) List(namespace string, opts skclients.ListOpts) (v1.UpstreamList, error) {
	listOpts := metav1.ListOptions{}
	if opts.ExpressionSelector != "" {
		listOpts.LabelSelector = opts.ExpressionSelector
	} else {
		sel := labels.NewSelector()
		for k, v := range opts.Selector {
			req, _ := labels.NewRequirement(k, selection.Equals, strings.Split(v, ","))
			sel = sel.Add(*req)
		}
		listOpts.LabelSelector = sel.String()

	}
	s.client.Istio().NetworkingV1beta1().ServiceEntries(namespace).List(opts.Ctx, listOpts)
}

// Watch implements v1.UpstreamClient.
func (s *serviceEntryClient) Watch(namespace string, opts skclients.WatchOpts) (<-chan v1.UpstreamList, <-chan error, error) {
	panic("unimplemented")
}

// We don't actually use the following in thi shim, but we must satisfy the UpstreamClient interface.

const notImplementedErrMsg = "this operation is not supported by this client"

func (c *serviceEntryClient) BaseClient() skclients.ResourceClient {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return &NoOpUpstreamClient.NoOpUpstreamClient{}
}

func (c *serviceEntryClient) Register() error {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return fmt.Errorf(notImplementedErrMsg)
}

func (c *serviceEntryClient) Read(namespace, name string, opts skclients.ReadOpts) (*v1.Upstream, error) {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return nil, fmt.Errorf(notImplementedErrMsg)
}

func (c *serviceEntryClient) Write(resource *v1.Upstream, opts skclients.WriteOpts) (*v1.Upstream, error) {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return nil, fmt.Errorf(notImplementedErrMsg)
}

func (c *serviceEntryClient) Delete(namespace, name string, opts skclients.DeleteOpts) error {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return fmt.Errorf(notImplementedErrMsg)
}
