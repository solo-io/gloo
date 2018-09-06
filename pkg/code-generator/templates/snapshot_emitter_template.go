package templates

import (
	"text/template"
)

var ResourceGroupEmitterTemplate = template.Must(template.New("resource_group_emitter").Funcs(funcs).Parse(`package {{ .Project.PackageName }}

{{- $clients := new_str_slice }}
{{- range .Resources}}
{{- $clients := (append_str_slice $clients (printf "%vClient %vClient"  (lower_camel .Name) .Name)) }}
{{- end}}
{{- $clients := (join_str_slice $clients ", ") }}

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type {{ .GoName }}Emitter interface {
	Register() error
{{- range .Resources}}
	{{ .Name }}() {{ .Name }}Client
{{- end}}
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *{{ .GoName }}Snapshot, <-chan error, error)
}

func New{{ .GoName }}Emitter({{ $clients }}) {{ .GoName }}Emitter {
	return &{{ lower_camel .GoName }}Emitter{
{{- range .Resources}}
		{{ lower_camel .Name }}: {{ lower_camel .Name }}Client,
{{- end}}
	}
}

type {{ lower_camel .GoName }}Emitter struct {
{{- range .Resources}}
	{{ lower_camel .Name }} {{ .Name }}Client
{{- end}}
}

func (c *{{ lower_camel .GoName }}Emitter) Register() error {
{{- range .Resources}}
	if err := c.{{ lower_camel .Name }}.Register(); err != nil {
		return err
	}
{{- end}}
	return nil
}

{{- range .Resources}}

func (c *{{ lower_camel $.GoName }}Emitter) {{ .Name }}() {{ .Name }}Client {
	return c.{{ lower_camel .Name }}
}
{{- end}}

func (c *{{ lower_camel .GoName }}Emitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *{{ .GoName }}Snapshot, <-chan error, error) {
	snapshots := make(chan *{{ .GoName }}Snapshot)
	errs := make(chan error)
	var done sync.WaitGroup

	currentSnapshot := {{ .GoName }}Snapshot{}

	sync := func(newSnapshot {{ .GoName }}Snapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}

	for _, namespace := range watchNamespaces {
{{- range .Resources}}
		{{ lower_camel .Name }}Chan, {{ lower_camel .Name }}Errs, err := c.{{ lower_camel .Name }}.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting {{ .Name }} watch")
		}
		done.Add(1)
		go func() {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, {{ lower_camel .Name }}Errs, namespace+"-{{ lower_camel .PluralName }}")
		}()


		done.Add(1)
		go func(namespace string, {{ lower_camel .Name }}Chan  <- chan {{ .Name }}List) {
			defer done.Done()
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case {{ lower_camel .Name }}List := <-{{ lower_camel .Name }}Chan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.{{ .PluralName }}.Clear(namespace)
					newSnapshot.{{ .PluralName }}.Add({{ lower_camel .Name }}List...)
					sync(newSnapshot)
				}
			}
		}(namespace, {{ lower_camel .Name }}Chan)
{{- end}}
	}


	go func() {
		select {
		case <-opts.Ctx.Done():
			done.Wait()
			close(snapshots)
			close(errs)
		}
	}()
	return snapshots, errs, nil
}
`))
