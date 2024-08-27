package settingsutil

import (
	"context"
	"fmt"
	"slices"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/sets"
)

type settingsKeyStruct struct{}

var (
	settingsKey = settingsKeyStruct{}

	namespacesToWatch = []string{}
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

func IsAllNamespacesFromSettings(s *v1.Settings) bool {
	if s == nil {
		return false
	}
	return IsAllNamespaces(GetNamespaces(s))
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

func GenerateNamespacesToWatch(s *v1.Settings, namespaces kubernetes.KubeNamespaceList) ([]string, error) {
	if len(s.GetWatchNamespaces()) != 0 {
		return s.GetWatchNamespaces(), nil
	}

	if len(s.GetWatchNamespaceSelectors()) == 0 {
		return []string{""}, nil
	}

	var selectors []labels.Selector
	selectedNamespaces := sets.NewString()

	fmt.Println("--------------------- Selectors : ", selectedNamespaces)
	for _, selector := range s.GetWatchNamespaceSelectors() {
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
				selectedNamespaces.Insert(ns.Name)
				break
			}
		}
	}

	return selectedNamespaces.List(), nil
}

func setNamespacesToWatch(namespaces []string) {
	namespacesToWatch = namespaces
}

func UpdateNamespacesToWatch(s *v1.Settings, namespaces kubernetes.KubeNamespaceList) (bool, error) {
	ns, err := GenerateNamespacesToWatch(s, namespaces)
	if err != nil {
		return false, err
	}

	if slices.Equal(ns, namespacesToWatch) {
		return false, nil
	}

	setNamespacesToWatch(ns)
	return true, nil
}

func GetNamespaces(s *v1.Settings) []string {
	return namespacesToWatch
}

func LabelSelectorAsSelector(ps *v1.LabelSelector) (labels.Selector, error) {
	if ps == nil {
		return labels.Nothing(), nil
	}
	if len(ps.MatchLabels)+len(ps.MatchExpressions) == 0 {
		return labels.Everything(), nil
	}
	requirements := make([]labels.Requirement, 0, len(ps.MatchLabels)+len(ps.MatchExpressions))
	for k, v := range ps.MatchLabels {
		r, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	for _, expr := range ps.MatchExpressions {
		var op selection.Operator
		switch metav1.LabelSelectorOperator(expr.Operator) {
		case metav1.LabelSelectorOpIn:
			op = selection.In
		case metav1.LabelSelectorOpNotIn:
			op = selection.NotIn
		case metav1.LabelSelectorOpExists:
			op = selection.Exists
		case metav1.LabelSelectorOpDoesNotExist:
			op = selection.DoesNotExist
		default:
			return nil, fmt.Errorf("%q is not a valid label selector operator", expr.Operator)
		}
		r, err := labels.NewRequirement(expr.Key, op, append([]string(nil), expr.Values...))
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	selector := labels.NewSelector()
	selector = selector.Add(requirements...)
	return selector, nil
}
