package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

//go:generate mockgen -destination mocks/mock_template_generator.go -package mocks github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template TemplateGenerator

// Render a template based on the given upstream
type TemplateGenerator interface {
	GenerateSnapshot() ([]byte, error)
	GenerateDashboard(dashboardFolderId uint) ([]byte, error)
	GenerateUid() string
}

// when the observability pod updates a dashboard, use this pre-canned message to indicate it was automated
const DefaultCommitMessage = "__gloo-auto-gen-dashboard__"

type templateGenerator struct {
	upstream              *gloov1.Upstream
	dashboardJsonTemplate string
}

var _ TemplateGenerator = &templateGenerator{}

func NewTemplateGenerator(upstream *gloov1.Upstream, dashboardJsonTemplate string) TemplateGenerator {
	return &templateGenerator{
		upstream:              upstream,
		dashboardJsonTemplate: dashboardJsonTemplate,
	}
}

func (t *templateGenerator) GenerateUid() string {
	// Uid has max 40 chars
	// return trailing chars because they are more likely to be distinct
	name := t.upstream.Metadata.Name
	if len(name) > 40 {
		name = name[len(name)-41 : len(name)-1]
	}
	return strings.Replace(name, "-", "_", -1)
}

func (t *templateGenerator) GenerateDashboard(dashboardFolderId uint) ([]byte, error) {
	stats := upstreamStats{
		Uid:              t.GenerateUid(),
		EnvoyClusterName: t.buildEnvoyClusterName(),
		NameTemplate:     "{{zone}} ({{host}})",
		Overwrite:        true,
	}

	dashboardPayload := buildDashboardPayloadTemplate(t.dashboardJsonTemplate, dashboardFolderId)
	return tmplExec(dashboardPayload, stats)
}

func (t *templateGenerator) GenerateSnapshot() ([]byte, error) {
	stats := upstreamStats{
		EnvoyClusterName: t.buildEnvoyClusterName(),
		Uid:              t.GenerateUid(),
		NameTemplate:     "{{zone}} ({{host}})",
		Overwrite:        true,
	}

	snapshotPayload := buildSnapshotPayloadTemplate(t.dashboardJsonTemplate)
	return tmplExec(snapshotPayload, stats)
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

func (t *templateGenerator) buildEnvoyClusterName() string {
	return translator.UpstreamToClusterName(t.upstream.Metadata.Ref())
}

type upstreamStats struct {
	Uid              string
	EnvoyClusterName string
	NameTemplate     string
	Overwrite        bool
}

func buildDashboardPayloadTemplate(dashboardTemplate string, folderId uint) string {
	return fmt.Sprintf(`
{
  "dashboard": 
	%s,
  "overwrite": {{.Overwrite}},
  "message": "%s",
  "folderId": %d
}
`, dashboardTemplate, DefaultCommitMessage, folderId)
}

func buildSnapshotPayloadTemplate(dashboardTemplate string) string {
	return fmt.Sprintf(`
{
  "dashboard": 
	%s,
  "name": "{{.EnvoyClusterName}}"
}
`, dashboardTemplate)
}
