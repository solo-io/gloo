package grafana

import (
	"bytes"
	"strings"
	"text/template"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

type upstreamStats struct {
	ClusterName  string
	NameTemplate string
	Overwrite    bool
}

func NameToUid(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

func GenerateDashboard(upstream *gloov1.Upstream) ([]byte, error) {
	stats := upstreamStats{
		ClusterName:  NameToUid(upstream.Metadata.GetName()),
		NameTemplate: "{{zone}} ({{host}})",
		Overwrite:    true,
	}
	return tmplExec(dashboardTemplate, stats)
}

func GenerateSnapshot(upstream *gloov1.Upstream) ([]byte, error) {
	stats := upstreamStats{
		ClusterName:  NameToUid(upstream.Metadata.GetName()),
		NameTemplate: "{{zone}} ({{host}})",
		Overwrite:    true,
	}
	return tmplExec(snapshotTemlpate, stats)
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
