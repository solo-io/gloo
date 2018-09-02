package templates

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var funcs = template.FuncMap{
	"join":        strings.Join,
	"lowercase":   strings.ToLower,
	"lower_camel": strcase.ToLowerCamel,
	"upper_camel": strcase.ToCamel,
	"snake":       strcase.ToSnake,
	"new_str_slice": func() *[]string {
		var v []string
		return &v
	},
	"append_str_slice": func(to *[]string, str string) *[]string {
		*to = append(*to, str)
		return to
	},
	"join_str_slice": func(slc *[]string, sep string) string {
		return strings.Join(*slc, sep)
	},
	"new_bool": func() *bool {
		var v bool
		return &v
	},
	"set_bool": func(v *bool, val bool) *bool {
		*v = val
		return v
	},
}
