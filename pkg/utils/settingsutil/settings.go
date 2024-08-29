package settingsutil

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/solo-io/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/external/kubernetes/namespace"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/lru"
)

type settingsKeyStruct struct{}

var (
	mu sync.Mutex

	// Setting a cache size of 2 should suffice for :
	// - The current translation loop
	// - The new loop about to run when the settings have changed
	namespacesToWatchCache = lru.New(2)

	settingsKey = settingsKeyStruct{}
)

func WithSettings(ctx context.Context, settings *v1.Settings) context.Context {
	return context.WithValue(ctx, settingsKey, settings)
}

// Deprecated: potentially unsafe panic
func FromContext(ctx context.Context) *v1.Settings {
	if ctx == nil {
		return nil
	}
	settings := MaybeFromContext(ctx)
	if settings != nil {
		return settings
	}
	// we should always have settings when this method is called.
	panic("no settings on context")
}

func MaybeFromContext(ctx context.Context) *v1.Settings {
	if ctx == nil {
		return nil
	}
	if settings, ok := ctx.Value(settingsKey).(*v1.Settings); ok {
		return settings
	}
	return nil
}

func IsAllNamespacesFromSettings(settings *v1.Settings) bool {
	if settings == nil {
		return false
	}
	return IsAllNamespaces(GetNamespacesToWatch(settings))
}

func IsAllNamespaces(watchNs []string) bool {
	switch {
	case len(watchNs) == 0:
		return true
	case len(watchNs) == 1 && watchNs[0] == "":
		return true
	default:
		return false
	}
}

func GenerateNamespacesToWatch(settings *v1.Settings, namespaces kubernetes.KubeNamespaceList) ([]string, error) {
	writeNamespace := settings.GetDiscoveryNamespace()
	if writeNamespace == "" {
		writeNamespace = defaults.GlooSystem
	}

	if len(settings.GetWatchNamespaces()) != 0 {
		return utils.ProcessWatchNamespaces(settings.GetWatchNamespaces(), writeNamespace), nil
	}

	if len(settings.GetWatchNamespaceSelectors()) == 0 {
		return []string{""}, nil
	}

	var selectors []labels.Selector
	selectedNamespaces := sets.NewString()

	fmt.Println("--------------------- Selectors : ")
	for _, selector := range settings.GetWatchNamespaceSelectors() {
		ls, err := LabelSelectorAsSelector(selector)
		fmt.Println(ls.String())
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, ls)
	}

	for _, ns := range namespaces {
		fmt.Println(ns.Name)
		for _, selector := range selectors {
			fmt.Println(ns.Labels)
			if selector.Matches(labels.Set(ns.Labels)) {
				fmt.Println("--------------------- Adding : ", ns.Name)
				selectedNamespaces.Insert(ns.Name)
				break
			}
		}
	}

	fmt.Println("--------------------- Adding : ", writeNamespace)
	selectedNamespaces.Insert(writeNamespace)

	return selectedNamespaces.List(), nil
}

func setNamespacesToWatch(settings *v1.Settings, namespaces []string) {
	// fmt.Println("--------------------- FROM : ", namespacesToWatch)
	// namespacesToWatch = namespaces
	// fmt.Println("--------------------- TO : ", namespacesToWatch)

	namespacesToWatchCache.Add(settings.MustHash(), namespaces)
	fmt.Println("--------------------- TO : ", namespaces)
}

func UpdateNamespacesToWatch(settings *v1.Settings, namespaces kubernetes.KubeNamespaceList) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	newNamespacesToWatch, err := GenerateNamespacesToWatch(settings, namespaces)
	if err != nil {
		return false, err
	}
	fmt.Println("--------------------- NS : ", newNamespacesToWatch)

	ns, ok := namespacesToWatchCache.Get(settings.MustHash())
	if ok {
		currentNamespacesToWatch, ok := ns.([]string)
		if ok && slices.Equal(newNamespacesToWatch, currentNamespacesToWatch) {
			return false, nil
		}
	}

	setNamespacesToWatch(settings, newNamespacesToWatch)

	return true, nil
}

func GetAllNamespaces() (kubernetes.KubeNamespaceList, error) {
	kubeClient := helpers.MustKubeClient()
	kubeCache, _ := cache.NewKubeCoreCache(context.TODO(), kubeClient)
	nsClient := namespace.NewNamespaceClient(kubeClient, kubeCache)

	return nsClient.List(clients.ListOpts{})
}

func GetNamespacesToWatch(settings *v1.Settings) []string {
	ns, ok := namespacesToWatchCache.Get(settings.MustHash())
	if ok {
		currentNamespacesToWatch, ok := ns.([]string)
		if ok {
			return currentNamespacesToWatch
		}
	}

	// Fallback to fetching all namespaces and updating the cache if not found
	allNamespaces, err := GetAllNamespaces()
	if err != nil {
		panic("Unable to fetch namespaces")
	}
	UpdateNamespacesToWatch(settings, allNamespaces)
	ns, _ = namespacesToWatchCache.Get(settings.MustHash())
	return ns.([]string)
}

func LabelSelectorAsSelector(ps *v1.LabelSelector) (labels.Selector, error) {
	if ps == nil {
		return labels.Nothing(), nil
	}
	if len(ps.GetMatchLabels())+len(ps.GetMatchExpressions()) == 0 {
		return labels.Everything(), nil
	}
	requirements := make([]labels.Requirement, 0, len(ps.GetMatchLabels())+len(ps.GetMatchExpressions()))
	for k, v := range ps.GetMatchLabels() {
		r, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	for _, expr := range ps.GetMatchExpressions() {
		var op selection.Operator
		switch metav1.LabelSelectorOperator(expr.GetOperator()) {
		case metav1.LabelSelectorOpIn:
			op = selection.In
		case metav1.LabelSelectorOpNotIn:
			op = selection.NotIn
		case metav1.LabelSelectorOpExists:
			op = selection.Exists
		case metav1.LabelSelectorOpDoesNotExist:
			op = selection.DoesNotExist
		default:
			return nil, fmt.Errorf("%q is not a valid label selector operator", expr.GetOperator())
		}
		r, err := labels.NewRequirement(expr.GetKey(), op, append([]string(nil), expr.GetValues()...))
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	selector := labels.NewSelector()
	selector = selector.Add(requirements...)
	return selector, nil
}
