package settingsutil

import (
	"context"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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

func SetNamespacesToWatch(s *v1.Settings, namespaces []string) error {

	if len(s.GetWatchNamespaces()) != 0 {
		namespacesToWatch = s.GetWatchNamespaces()
		return nil
	}

	if len(s.GetWatchNamespaceSelectors()) == 0 {
		namespacesToWatch = []string{""}
	}

	var selectors []labels.Selector
	selectedNamespaces := sets.NewString()

	for _, selector := range s.GetWatchNamespaceSelectors() {
		ls, err := metav1.LabelSelectorAsSelector(selector)
		if err != nil {
			return err
		}
		selectors = append(selectors, ls)
	}

	namespacesToWatch = selectedNamespaces.List()
	return nil
}

func GetNamespaces(s *v1.Settings) []string {
	return namespacesToWatch
}
