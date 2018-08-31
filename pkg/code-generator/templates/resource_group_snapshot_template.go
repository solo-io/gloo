package templates

import (
	"text/template"
)

var ResourceGroupSnapshotTemplate = template.Must(template.New("resource_group_snapshot").Funcs(funcs).Parse(`package {{ .PackageName }}

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type {{ .Name }}Snapshot struct {
{{- range .Resources}}
	{{ upper_camel .PluralName }} {{ upper_camel .PluralName }}ByNamespace
{{- end}}
}

func (s {{ .Name }}Snapshot) Clone() {{ .Name }}Snapshot {
	return {{ .Name }}Snapshot{
{{- range .ResourceTypes}}
		{{ upper_camel .PluralName }}: s.{{ upper_camel .PluralName }}.Clone(),
{{- end}}
	}
}

func (s {{ .Name }}Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
{{- range .ResourceTypes}}
	for _, {{ lower_camel . }} := range snapshotForHashing.{{ upper_camel .PluralName }}.List() {
		resources.UpdateMetadata({{ lower_camel . }}, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
{{- if (index $.ResourceLevelParams .).IsInputType }}
		{{ lower_camel . }}.SetStatus(core.Status{})
{{- end }}
	}
{{- end}}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}
`))
