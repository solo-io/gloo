package grafana

import (
	"bytes"
	"fmt"
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
	return strings.Replace(name, "-", "_", -1)
}

func (t *templateGenerator) GenerateDashboard() ([]byte, error) {
	stats := upstreamStats{
		Uid:              t.GenerateUid(),
		EnvoyClusterName: t.buildEnvoyClusterName(),
		NameTemplate:     "{{zone}} ({{host}})",
		Overwrite:        true,
	}
	return tmplExec(dashboardTemplate, stats)
}

func (t *templateGenerator) GenerateSnapshot() ([]byte, error) {
	stats := upstreamStats{
		EnvoyClusterName: t.buildEnvoyClusterName(),
		Uid:              t.GenerateUid(),
		NameTemplate:     "{{zone}} ({{host}})",
		Overwrite:        true,
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

func (t *templateGenerator) buildEnvoyClusterName() string {
	us := t.upstream.GetUpstreamSpec()
	switch us.GetUpstreamType().(type) {

	// kubernetes upstreams have their prometheus statistics built using metadata about the service being represented
	case *gloov1.UpstreamSpec_Kube:
		kube := us.GetKube()
		return fmt.Sprintf("%s-%s-%d_%s", kube.ServiceNamespace, kube.ServiceName, kube.ServicePort, t.upstream.Metadata.Namespace)

	// all other types just use the name/namespace of the upstream itself
	default:
		return fmt.Sprintf("%s_%s", t.upstream.Metadata.Name, t.upstream.Metadata.Namespace)
	}
}

type upstreamStats struct {
	Uid              string
	EnvoyClusterName string
	NameTemplate     string
	Overwrite        bool
}
