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
type Transition{{ .ResourceType }}Func func(original, desired *{{ .ResourceType }})

type {{ .ResourceType }}Reconciler interface {
	Reconcile(namespace string, desiredResources []*{{ .ResourceType }}, opts clients.ListOpts) error
}

func {{ lowercase .ResourceType }}sToResources(list {{ .ResourceType }}List) []resources.Resource {
	var resourceList []resources.Resource
	for _, {{ lowercase .ResourceType }} := range list {
		resourceList = append(resourceList, {{ lowercase .ResourceType }})
	}
	return resourceList
}

func New{{ .ResourceType }}Reconciler(client {{ .ResourceType }}Client, transition Transition{{ .ResourceType }}Func) {{ .ResourceType }}Reconciler {
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) {
			transition(original.(*{{ .ResourceType }}), desired.(*{{ .ResourceType }}))
		}
	}
	return &{{ lowercase .ResourceType }}Reconciler{
		base: reconcile.NewReconciler(client.BaseClient(), transitionResources),
	}
}

type {{ lowercase .ResourceType }}Reconciler struct {
	base reconcile.Reconciler
}

func (r *{{ lowercase .ResourceType }}Reconciler) Reconcile(namespace string, desiredResources []*{{ .ResourceType }}, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "{{ lowercase .ResourceType }}_reconciler")
	return r.base.Reconcile(namespace, {{ lowercase .ResourceType }}sToResources(desiredResources), opts)
}
`
