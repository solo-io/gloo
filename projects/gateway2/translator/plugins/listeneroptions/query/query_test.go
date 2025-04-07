package query_test

import (
	"context"

	"github.com/solo-io/gloo/pkg/schemes"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/listeneroptions/query"
	plugintestutils "github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils/testutils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ListenerOptionsBuilder struct{}

func (b ListenerOptionsBuilder) Build(def *plugintestutils.OptionsDef) client.Object {
	now := metav1.Now()
	lo := &solokubev1.ListenerOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:              def.Name,
			Namespace:         def.Namespace,
			CreationTimestamp: now,
		},
		Spec: sologatewayv1.ListenerOption{
			Options: &v1.ListenerOptions{},
		},
	}

	if def.TargetRefs != nil {
		lo.Spec.TargetRefs = def.TargetRefs
	}

	return lo
}

func getQuery(ctx context.Context, deps []client.Object, listener *gwv1.Listener, gw *gwv1.Gateway, listenerSet *apixv1a1.XListenerSet) ([]client.Object, error) {

	builder := fake.NewClientBuilder().WithScheme(schemes.GatewayScheme())
	query.IterateIndices(func(o client.Object, f string, fun client.IndexerFunc) error {
		builder.WithIndex(o, f, fun)
		return nil
	})
	fakeClient := builder.WithObjects(deps...).Build()
	qry := query.NewQuery(fakeClient)

	listenerOptions, err := qry.GetAttachedListenerOptions(ctx, listener, gw, listenerSet)
	if err != nil {
		return nil, err
	}

	// Convert []*solokubev1.HttpListenerOption to []client.Object
	objects := make([]client.Object, len(listenerOptions))
	for i, vho := range listenerOptions {
		objects[i] = vho // vho is a pointer to a struct that implements client.Object
	}

	return objects, nil

}

var _ = plugintestutils.TestListenerOptionPlugin(getQuery, ListenerOptionsBuilder{})
