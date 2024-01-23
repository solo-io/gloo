package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
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

// Source: https://github.com/grafana/grafana/blob/fc73bc1161c2801a027ad76a7022408d845a73df/pkg/util/shortid_generator.go#L11
// We additionally want to replace -
var uidInvalidCharacters = regexp.MustCompile(`[^a-zA-Z0-9\_]`)

const (
	// when the observability pod updates a dashboard, use this pre-canned message to indicate it was automated
	DefaultCommitMessage = "__gloo-auto-gen-dashboard__"
	// Ref: https://grafana.com/docs/grafana/latest/developers/http_api/dashboard/#identifier-id-vs-unique-identifier-uid
	GrafanaUIDMaxLength = 40
)

type upstreamTemplateGenerator struct {
	upstream                   *gloov1.Upstream
	dashboardJsonTemplate      string
	dashboardPrefix            string
	extraMetricQueryParameters string
	extraDashboardTags         []string
}

var _ TemplateGenerator = &upstreamTemplateGenerator{}

func NewUpstreamTemplateGenerator(upstream *gloov1.Upstream, dashboardJsonTemplate string, opts ...Option) TemplateGenerator {
	options := processOptions(opts...)
	return &upstreamTemplateGenerator{
		upstream:                   upstream,
		dashboardJsonTemplate:      dashboardJsonTemplate,
		dashboardPrefix:            options.dashboardPrefix,
		extraMetricQueryParameters: options.extraMetricQueryParameters,
		extraDashboardTags:         options.extraDashboardTags,
	}
}

func (t *upstreamTemplateGenerator) generateUidWithoutPrefix() string {
	// return trailing chars because they are more likely to be distinct
	maxNameLength := GrafanaUIDMaxLength - len(t.dashboardPrefix)
	name := t.upstream.Metadata.Name
	if len(name) > maxNameLength {
		name = name[len(name)-maxNameLength:]
	}
	return uidInvalidCharacters.ReplaceAllLiteralString(name, "_") // replace invalid characters with _
}

func (t *upstreamTemplateGenerator) GenerateUid() string {
	var sb strings.Builder
	sb.WriteString(t.dashboardPrefix)
	sb.WriteString(t.generateUidWithoutPrefix())
	return uidInvalidCharacters.ReplaceAllLiteralString(sb.String(), "_") // replace invalid characters with _
}

// ToUID converts the string to a Grafana compatible UID.
// It truncates it to 40 char, the max UID length and
// replaces all non UID characters based off https://github.com/grafana/grafana/blob/fc73bc1161c2801a027ad76a7022408d845a73df/pkg/util/shortid_generator.go#L11
func ToUID(str string) string {
	if len(str) > GrafanaUIDMaxLength {
		str = str[:GrafanaUIDMaxLength]
	}
	return uidInvalidCharacters.ReplaceAllLiteralString(str, "_")
}

func (t *upstreamTemplateGenerator) GenerateDashboardPost(dashboardFolderId uint) (*grafana.DashboardPostRequest, error) {
	extraDashboardTags, err := generateExtraTags(t.extraDashboardTags)
	if err != nil {
		return nil, err
	}
	stats := upstreamStats{
		DashboardPrefix:            t.dashboardPrefix,
		ExtraDashboardTags:         extraDashboardTags,
		ExtraMetricQueryParameters: generateExtraMetricsQueryParameters(t.extraMetricQueryParameters),
		Uid:                        t.generateUidWithoutPrefix(),
		EnvoyClusterName:           t.buildEnvoyClusterName(),
		NameTemplate:               "{{zone}} ({{host}})",
		Overwrite:                  true,
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
	extraDashboardTags, err := generateExtraTags(t.extraDashboardTags)
	if err != nil {
		return nil, err
	}
	stats := upstreamStats{
		DashboardPrefix:            t.dashboardPrefix,
		ExtraDashboardTags:         extraDashboardTags,
		ExtraMetricQueryParameters: generateExtraMetricsQueryParameters(t.extraMetricQueryParameters),
		EnvoyClusterName:           t.buildEnvoyClusterName(),
		Uid:                        t.generateUidWithoutPrefix(),
		NameTemplate:               "{{zone}} ({{host}})",
		Overwrite:                  true,
	}

	snapshotPayload := buildSnapshotPayloadTemplate(t.dashboardJsonTemplate)
	return tmplExec(snapshotPayload, "upstream.json", stats)
}

type defaultDashboardTemplateGenerator struct {
	uid, dashboardJson, dashboardPrefix, extraMetricQueryParameters string
	extraDashboardTags                                              []string
}

func NewDefaultDashboardTemplateGenerator(uid, dashboardJson string, opts ...Option) TemplateGenerator {
	options := processOptions(opts...)
	return &defaultDashboardTemplateGenerator{
		uid:                        uid,
		dashboardJson:              dashboardJson,
		dashboardPrefix:            options.dashboardPrefix,
		extraMetricQueryParameters: options.extraMetricQueryParameters,
		extraDashboardTags:         options.extraDashboardTags,
	}
}

func (t *defaultDashboardTemplateGenerator) GenerateUid() string {
	return t.uid
}

func (t *defaultDashboardTemplateGenerator) GenerateDashboardPost(dashboardFolderId uint) (*grafana.DashboardPostRequest, error) {
	extraDashboardTags, err := generateExtraTags(t.extraDashboardTags)
	if err != nil {
		return nil, err
	}
	// Since the default dashboard contains templates pipelines, just do a direct string replace
	dashboardJson := strings.ReplaceAll(t.dashboardJson, "{{.DashboardPrefix}}", t.dashboardPrefix)
	dashboardJson = strings.ReplaceAll(dashboardJson, "{{.ExtraDashboardTags}}", extraDashboardTags)
	dashboardJson = strings.ReplaceAll(dashboardJson, "{{.ExtraMetricQueryParameters}}", generateExtraMetricsQueryParameters(t.extraMetricQueryParameters))
	return &grafana.DashboardPostRequest{
		Dashboard: []byte(dashboardJson),
		Message:   DefaultCommitMessage,
		FolderId:  dashboardFolderId,
		Overwrite: false,
	}, nil
}

func (t *defaultDashboardTemplateGenerator) GenerateSnapshot() ([]byte, error) {
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
	Uid                        string
	EnvoyClusterName           string
	NameTemplate               string
	Overwrite                  bool
	DashboardPrefix            string
	ExtraMetricQueryParameters string
	ExtraDashboardTags         string
}

func generateExtraTags(tags []string) (string, error) {
	if len(tags) > 0 {
		strTags, err := json.Marshal(tags)
		if err != nil {
			return "", err
		}
		// Since this has to be appended to an existing list of tags, prefix it with a comma and remove the `[]`
		return "," + string(strTags[1:len(strTags)-1]), nil
	}
	return "", nil
}

func generateExtraMetricsQueryParameters(extraMetricQueryParameters string) string {
	if extraMetricQueryParameters != "" {
		// Since the extra query params are inside of a string, the quotes need to be escaped if not already done
		// Unfortunately Go doesn't support negative look behind matching, so we do it the fun way
		// 1. Replace all " with \" (now the pre-exiting escaped quotes are \\")
		// 2. Replace all \\" with \"
		// Eg: `cluster=\"clue",proxy="pro\"` should be converted into `cluster=\"clue\",proxy=\"pro\"`
		// 0. extraMetricQueryParameters = `cluster=\"clue",proxy="pro\"`
		// 1. params = cluster=\\"clue\",proxy=\"pro\\"
		params := strings.ReplaceAll(extraMetricQueryParameters, "\"", "\\\"")
		// 2. params = cluster=\"clue\",proxy=\"pro\"
		params = strings.ReplaceAll(params, "\\\\\"", "\\\"")
		// Prepend a comma since it needs to join a list of existing params
		params = "," + params
		return params
	}
	return extraMetricQueryParameters
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
  "name": "{{.DashboardPrefix}}{{.EnvoyClusterName}}"
}
`, dashboardTemplate)
}
