package codegen

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/solo-io/skv2/codegen/render"
)

func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"base_group": func(grp render.Group) string {
			// e.g. enterprise_gloo, gloo, gateway
			return strings.ReplaceAll(baseGroupImportName(grp), fmt.Sprintf("_solo_io_%s", grp.Version), "")
		},
		"base_Group": func(grp render.Group) string {
			// e.g. EnterpriseGloo
			return strcase.ToCamel(strings.ReplaceAll(baseGroupImportName(grp), fmt.Sprintf("_solo_io_%s", grp.Version), ""))
		},
		"base_group_import_name": baseGroupImportName,
		"base_kind": func(fedKind string) string {
			return strings.ReplaceAll(fedKind, "Federated", "")
		},
		// Proto group minus `.solo.io`, as used in our directory names, e.g. enterprise.gloo
		// Replace ratelimit.api with ratelimit
		"base_proto_group_shorthand": func(grp render.Group) string {
			output := strings.ReplaceAll(grp.Group, "fed.", "")
			output = strings.ReplaceAll(output, "ratelimit.api", "ratelimit")
			return strings.ReplaceAll(output, ".solo.io", "")
		},
		"base_proto_group": func(grp render.Group) string { // e.g. enterprise.gloo.solo.io
			return strings.ReplaceAll(grp.Group, "fed.", "")
		},
		"base_group_version": func(grp render.Group) string { // e.g. v1, v1alpha1
			return grp.Version
		},
		// Helper functions
		"uppercase": func(str string) string {
			return strings.ToUpper(str)
		},
		"lowercase": func(str string) string {
			return strings.ToLower(str)
		},
	}
}

// e.g. fed_gateway_solo_io_v1
func groupImportName(grp render.Group) string {
	name := strings.ReplaceAll(grp.GroupVersion.String(), "/", "_")
	name = strings.ReplaceAll(name, "ratelimit.api", "ratelimit")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

// e.g. enterprise_gloo_solo_io_v1
func baseGroupImportName(grp render.Group) string {
	return strings.ReplaceAll(groupImportName(grp), "fed_", "")
}
