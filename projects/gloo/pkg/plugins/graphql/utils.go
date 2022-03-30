package graphql

import (
	"regexp"
	"strings"
	"unicode"
)

var regex = regexp.MustCompile(`\n[ ]+\n`)

// Does 3 things to pretty print a kube string (allow the string to display in multiline format in yaml)
func PrettyPrintKubeString(input string) string {
	input = strings.ReplaceAll(input, "\t", "  ")
	input = regex.ReplaceAllString(input, "\n\n")
	input = strings.TrimRightFunc(input, unicode.IsSpace)
	return input
}
