package types_test

import (
	"testing"

	"github.com/solo-io/gloo/projects/gateway2/translator/types"
	"github.com/solo-io/gloo/projects/gateway2/utils"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

func TestGetListeners(t *testing.T) {
	gw := gw()
	ls1 := ls("ls1")
	ls2 := ls("ls2")

	cgw := &types.ConsolidatedGateway{
		Gateway:             gw,
		AllowedListenerSets: []*apixv1a1.XListenerSet{ls1},
		DeniedListenerSets:  []*apixv1a1.XListenerSet{ls2},
	}

	assert.Len(t, cgw.GetListeners(ls1), 1)
	assert.Len(t, cgw.GetListeners(ls2), 0)
}

func TestGetConsolidatedListeners(t *testing.T) {
	gw := gw()
	ls1 := ls("ls1")
	ls2 := ls("ls2")

	cgw := &types.ConsolidatedGateway{
		Gateway:             gw,
		AllowedListenerSets: []*apixv1a1.XListenerSet{ls1},
		DeniedListenerSets:  []*apixv1a1.XListenerSet{ls2},
	}

	cl := cgw.GetConsolidatedListeners()
	assert.Len(t, cl, 2)

	lis1 := utils.ToListener(ls1.Spec.Listeners[0])
	assert.Equal(t, cl, []types.ConsolidatedListener{
		{
			Gateway:     gw,
			ListenerSet: nil,
			Listener:    &gw.Spec.Listeners[0],
		},
		{
			Gateway:     gw,
			ListenerSet: ls1,
			Listener:    &lis1,
		},
	})
}

func gw() *apiv1.Gateway {
	return &apiv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "gw",
		},
		Spec: apiv1.GatewaySpec{
			Listeners: []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
				},
			},
		},
	}
}

func ls(name string) *apixv1a1.XListenerSet {
	return &apixv1a1.XListenerSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
		},
		Spec: apixv1a1.ListenerSetSpec{
			Listeners: []apixv1a1.ListenerEntry{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
				},
			},
		},
	}
}
