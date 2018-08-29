package typed

import (
	"bytes"
	"text/template"
)

func GenerateReconcilerCode(params ResourceLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := typedReconcilerTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var typedReconcilerTemplate = template.Must(template.New("typed_client").Funcs(funcs).Parse(typedReconcilerTemplateContents))

const typedReconcilerTemplateContents = `package {{ .PackageName }}

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type Transition{{ .ResourceType }}Func func(original, desired *{{ .ResourceType }}) error

type {{ .ResourceType }}Reconciler interface {
	Reconcile(namespace string, desiredResources []*{{ .ResourceType }}, transition Transition{{ .ResourceType }}Func, opts clients.ListOpts) error
}

func {{ lower_camel .ResourceType }}sToResources(list {{ .ResourceType }}List) resources.ResourceList {
	var resourceList resources.ResourceList
	for _, {{ lower_camel .ResourceType }} := range list {
		resourceList = append(resourceList, {{ lower_camel .ResourceType }})
	}
	return resourceList
}

func New{{ .ResourceType }}Reconciler(client {{ .ResourceType }}Client) {{ .ResourceType }}Reconciler {
	return &{{ lower_camel .ResourceType }}Reconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type {{ lower_camel .ResourceType }}Reconciler struct {
	base reconcile.Reconciler
}

func (r *{{ lower_camel .ResourceType }}Reconciler) Reconcile(namespace string, desiredResources []*{{ .ResourceType }}, transition Transition{{ .ResourceType }}Func, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "{{ lower_camel .ResourceType }}_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) error {
			return transition(original.(*{{ .ResourceType }}), desired.(*{{ .ResourceType }}))
		}
	}
	return r.base.Reconcile(namespace, {{ lower_camel .ResourceType }}sToResources(desiredResources), transitionResources, opts)
}
`
