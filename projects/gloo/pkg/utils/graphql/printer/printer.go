package printer

import (
	"regexp"
	"strings"
	"unicode"
)

// used to get rid of the spaces between new lines
var betweenNewLinesRegex = regexp.MustCompile(`\n[ ]+\n`)

// used to get rid of all the spaces at the end of a line
var whiteSpacesAtTheEndOfLineRegex = regexp.MustCompile(`(.*)[ ]{1,}\n`)

// Does 3 things to pretty print a kube string (allow the string to display in multiline format in yaml)
func PrettyPrintKubeString(input string) string {
	input = strings.ReplaceAll(input, "\t", "  ")
	input = betweenNewLinesRegex.ReplaceAllString(input, "\n\n")
	input = strings.TrimRightFunc(input, unicode.IsSpace)
	input = whiteSpacesAtTheEndOfLineRegex.ReplaceAllString(input, "$1\n")
	return input
}
