package templates

import (
	"text/template"
)

var ResourceGroupSnapshotTemplate = template.Must(template.New("resource_group_snapshot").Funcs(funcs).Parse(
	`package {{ .Project.Version }}

import (
	{{ .Imports }}
	"go.uber.org/zap"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type {{ .GoName }}Snapshot struct {
{{- range .Resources}}
	{{ upper_camel .PluralName }} {{ .ImportPrefix }}{{ upper_camel .PluralName }}ByNamespace
{{- end}}
}

func (s {{ .GoName }}Snapshot) Clone() {{ .GoName }}Snapshot {
	return {{ .GoName }}Snapshot{
{{- range .Resources}}
		{{ upper_camel .PluralName }}: s.{{ upper_camel .PluralName }}.Clone(),
{{- end}}
	}
}

func (s {{ .GoName }}Snapshot) snapshotToHash() {{ .GoName }}Snapshot {
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

	return snapshotForHashing
}

func (s {{ .GoName }}Snapshot) Hash() uint64 {
	return s.hashStruct(s.snapshotToHash())
 }

 func (s {{ .GoName }}Snapshot) HashFields() []zap.Field {
	snapshotForHashing := s.snapshotToHash()
	var fields []zap.Field


{{- range .Resources}}
	{{ lower_camel .PluralName }} := s.hashStruct(snapshotForHashing.{{ upper_camel .PluralName }}.List())
	fields = append(fields, zap.Uint64("{{ lower_camel .PluralName }}", {{ lower_camel .PluralName }} ))
{{- end}}

	return append(fields, zap.Uint64("snapshotHash",  s.hashStruct(snapshotForHashing)))
 }
 
func (s {{ .GoName }}Snapshot) hashStruct(v interface{}) uint64 {
	h, err := hashstructure.Hash(v, nil)
	 if err != nil {
		 panic(err)
	 }
	 return h
 }


`))
