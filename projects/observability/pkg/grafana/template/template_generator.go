package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

//go:generate mockgen -destination mocks/mock_template_generator.go -package mocks github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template TemplateGenerator

// Render a template based on the given upstream
type TemplateGenerator interface {
	GenerateSnapshot() ([]byte, error)
	GenerateDashboardPost(dashboardFolderId uint) (*grafana.DashboardPostRequest, error)
	GenerateUid() string
}

// when the observability pod updates a dashboard, use this pre-canned message to indicate it was automated
const DefaultCommitMessage = "__gloo-auto-gen-dashboard__"

type upstreamTemplateGenerator struct {
	upstream              *gloov1.Upstream
	dashboardJsonTemplate string
}

var _ TemplateGenerator = &upstreamTemplateGenerator{}

func NewUpstreamTemplateGenerator(upstream *gloov1.Upstream, dashboardJsonTemplate string) TemplateGenerator {
	return &upstreamTemplateGenerator{
		upstream:              upstream,
		dashboardJsonTemplate: dashboardJsonTemplate,
	}
}

func (t *upstreamTemplateGenerator) GenerateUid() string {
	// Uid has max 40 chars
	// return trailing chars because they are more likely to be distinct
	name := t.upstream.Metadata.Name
	if len(name) > 40 {
		name = name[len(name)-41 : len(name)-1]
	}
	return strings.Replace(name, "-", "_", -1)
}

func (t *upstreamTemplateGenerator) GenerateDashboardPost(dashboardFolderId uint) (*grafana.DashboardPostRequest, error) {
	stats := upstreamStats{
		Uid:              t.GenerateUid(),
		EnvoyClusterName: t.buildEnvoyClusterName(),
		NameTemplate:     "{{zone}} ({{host}})",
		Overwrite:        true,
	}

	renderedDash, err := tmplExec(t.dashboardJsonTemplate, "upstream.json", stats)
	if err != nil {
		return nil, err
	}

	return &grafana.DashboardPostRequest{
		Dashboard: renderedDash,
		FolderId:  dashboardFolderId,
		Message:   DefaultCommitMessage,
		Overwrite: true,
	}, nil
}

func (t *upstreamTemplateGenerator) GenerateSnapshot() ([]byte, error) {
	stats := upstreamStats{
		EnvoyClusterName: t.buildEnvoyClusterName(),
		Uid:              t.GenerateUid(),
		NameTemplate:     "{{zone}} ({{host}})",
		Overwrite:        true,
	}

	snapshotPayload := buildSnapshotPayloadTemplate(t.dashboardJsonTemplate)
	return tmplExec(snapshotPayload, "upstream.json", stats)
}

type defaultJsonGenerator struct {
	uid, dashboardJson string
}

func NewDefaultJsonGenerator(uid, dashboardJson string) TemplateGenerator {
	return &defaultJsonGenerator{
		uid:           uid,
		dashboardJson: dashboardJson,
	}
}

func (t *defaultJsonGenerator) GenerateUid() string {
	return t.uid
}

func (t *defaultJsonGenerator) GenerateDashboardPost(dashboardFolderId uint) (*grafana.DashboardPostRequest, error) {
	return &grafana.DashboardPostRequest{
		Dashboard: []byte(t.dashboardJson),
		Message:   DefaultCommitMessage,
		FolderId:  dashboardFolderId,
		Overwrite: false,
	}, nil
}

func (t *defaultJsonGenerator) GenerateSnapshot() ([]byte, error) {
	return []byte{}, fmt.Errorf("GenerateSnapshot not implemented for defaultJsonGenerator")
}

func tmplExec(tmplStr, filename string, us upstreamStats) ([]byte, error) {
	tmpl, err := template.New(filename).Parse(tmplStr)
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

func (t *upstreamTemplateGenerator) buildEnvoyClusterName() string {
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
