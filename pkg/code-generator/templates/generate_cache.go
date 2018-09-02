package templates

const emitterTemplateContents = `package {{ .PackageName }}

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type Emitter interface {
	Register() error
{{- range .ResourceTypes}}
	{{ . }}() {{ . }}Client
{{- end}}
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewEmitter({{ clients . true }}) Emitter {
	return &emitter{
{{- range .ResourceTypes}}
		{{ lower_camel . }}: {{ lower_camel . }}Client,
{{- end}}
	}
}

type emitter struct {
{{- range .ResourceTypes}}
	{{ lower_camel . }} {{ . }}Client
{{- end}}
}

func (c *emitter) Register() error {
{{- range .ResourceTypes}}
	if err := c.{{ lower_camel . }}.Register(); err != nil {
		return err
	}
{{- end}}
	return nil
}

{{- range .ResourceTypes}}

func (c *emitter) {{ . }}() {{ . }}Client {
	return c.{{ lower_camel . }}
}
{{- end}}

func (c *emitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
	snapshots := make(chan *Snapshot)
	errs := make(chan error)

	currentSnapshot := Snapshot{}

	sync := func(newSnapshot Snapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}

	for _, namespace := range watchNamespaces {
{{- range .ResourceTypes}}
		{{ lower_camel . }}Chan, {{ lower_camel . }}Errs, err := c.{{ lower_camel . }}.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting {{ . }} watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, {{ lower_camel . }}Errs, namespace+"-{{ lower_camel (resource . $).PluralName }}")
		go func(namespace string, {{ lower_camel . }}Chan  <- chan {{ . }}List) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case {{ lower_camel . }}List := <-{{ lower_camel . }}Chan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.{{ (resource . $).PluralName }}.Clear(namespace)
					newSnapshot.{{ (resource . $).PluralName }}.Add({{ lower_camel . }}List...)
					sync(newSnapshot)
				}
			}
		}(namespace, {{ lower_camel . }}Chan)
{{- end}}
	}


	go func() {
		select {
		case <-opts.Ctx.Done():
			close(snapshots)
			close(errs)
		}
	}()
	return snapshots, errs, nil
}
`

const cacheTestTemplateContents = `package {{ .PackageName }}
`
