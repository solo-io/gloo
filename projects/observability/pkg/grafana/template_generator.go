package grafana

import (
	"bytes"
	"strings"
	"text/template"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

//go:generate mockgen -destination mocks/mock_template_generator.go -package mocks github.com/solo-io/solo-projects/projects/observability/pkg/grafana TemplateGenerator

type TemplateGenerator interface {
	GenerateSnapshot() ([]byte, error)
	GenerateDashboard() ([]byte, error)
	GenerateUid() string
}

type templateGenerator struct {
	upstream *gloov1.Upstream
}

var _ TemplateGenerator = &templateGenerator{}

func NewTemplateGenerator(upstream *gloov1.Upstream) TemplateGenerator {
	return &templateGenerator{upstream: upstream}
}

func (t *templateGenerator) GenerateUid() string {
	// Uid has max 40 chars
	// return trailing chars because they are more likely to be distinct
	name := t.upstream.Metadata.Name
	if len(name) > 40 {
		name = name[len(name)-41 : len(name)-1]
	}
	return nameToEnvoyStats(name)
}

func (t *templateGenerator) GenerateDashboard() ([]byte, error) {
	upstream := t.upstream
	stats := upstreamStats{
		ClusterName:  nameToEnvoyStats(upstream.Metadata.GetName()),
		Uid:          t.GenerateUid(),
		NameTemplate: "{{zone}} ({{host}})",
		Overwrite:    true,
	}
	return tmplExec(dashboardTemplate, stats)
}

func (t *templateGenerator) GenerateSnapshot() ([]byte, error) {
	upstream := t.upstream
	stats := upstreamStats{
		ClusterName:  nameToEnvoyStats(upstream.Metadata.GetName()),
		Uid:          t.GenerateUid(),
		NameTemplate: "{{zone}} ({{host}})",
		Overwrite:    true,
	}
	return tmplExec(snapshotTemplate, stats)
}

func tmplExec(tmplStr string, us upstreamStats) ([]byte, error) {
	tmpl, err := template.New("upstream.json").Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, us)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func nameToEnvoyStats(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

type upstreamStats struct {
	Uid          string
	ClusterName  string
	NameTemplate string
	Overwrite    bool
}
