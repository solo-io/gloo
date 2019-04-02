package queries

import (
	"fmt"
	"strings"
)

func ReplaceNamespace(raw, targetNamespace string) string {
	const nsSpec = `"namespace":"default"`
	newSpec := fmt.Sprintf("\"namespace\":\"%v\"", targetNamespace)
	return strings.Replace(raw, nsSpec, newSpec, -1)
}

func ReplaceNamespaces(input []string, targetNamespace string) []string {
	out := []string{}
	for _, s := range input {
		converted := ReplaceNamespace(s, targetNamespace)
		out = append(out, converted)
	}
	return out
}
