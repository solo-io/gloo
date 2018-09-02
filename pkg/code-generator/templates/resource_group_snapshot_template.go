package templates

import (
	"text/template"
)

var ResourceGroupSnapshotTemplate = template.Must(template.New("resource_group_snapshot").Funcs(funcs).Parse(`package {{ .Project.PackageName }}

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type {{ .GoName }}Snapshot struct {
{{- range .Resources}}
	{{ upper_camel .PluralName }} {{ upper_camel .PluralName }}ByNamespace
{{- end}}
}

func (s {{ .GoName }}Snapshot) Clone() {{ .GoName }}Snapshot {
	return {{ .GoName }}Snapshot{
{{- range .Resources}}
		{{ upper_camel .PluralName }}: s.{{ upper_camel .PluralName }}.Clone(),
{{- end}}
	}
}

func (s {{ .GoName }}Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
{{- range .Resources}}
	for _, {{ lower_camel .Name }} := range snapshotForHashing.{{ upper_camel .PluralName }}.List() {
		resources.UpdateMetadata({{ lower_camel .Name }}, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
{{- if .HasStatus }}
		{{ lower_camel .Name }}.SetStatus(core.Status{})
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
