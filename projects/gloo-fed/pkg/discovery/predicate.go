package discovery

import (
	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	AppLabel = "app"
	Gloo     = "gloo"
)

var (
	GlooAppLabels = map[string]string{
		AppLabel: Gloo,
	}
)

var GlooWorkloadPredicate = predicate.Funcs{
	CreateFunc: func(e event.CreateEvent) bool {
		return isGlooResource(e.Object)
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		return isGlooResource(e.Object)
	},
	UpdateFunc: func(e event.UpdateEvent) bool {
		return isGlooResource(e.ObjectOld) || isGlooResource(e.ObjectNew)
	},
	GenericFunc: func(e event.GenericEvent) bool {
		return isGlooResource(e.Object)
	},
}

func isGlooResource(obj runtime.Object) bool {
	switch typedObject := obj.(type) {
	case *apps_v1.Deployment:
		return labels.SelectorFromSet(GlooAppLabels).Matches(labels.Set(typedObject.GetLabels()))
	case *apps_v1.DaemonSet:
		return labels.SelectorFromSet(GlooAppLabels).Matches(labels.Set(typedObject.GetLabels()))
	case *batch_v1.Job:
		return labels.SelectorFromSet(GlooAppLabels).Matches(labels.Set(typedObject.GetLabels()))
	default:
		return false
	}
}
