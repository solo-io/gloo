package utils

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/selectors"
	"k8s.io/apimachinery/pkg/selection"
)

var (
	SelectorExpressionOperatorValues = map[selectors.Selector_Expression_Operator]selection.Operator{
		selectors.Selector_Expression_Equals:       selection.Equals,
		selectors.Selector_Expression_DoubleEquals: selection.DoubleEquals,
		selectors.Selector_Expression_NotEquals:    selection.NotEquals,
		selectors.Selector_Expression_In:           selection.In,
		selectors.Selector_Expression_NotIn:        selection.NotIn,
		selectors.Selector_Expression_Exists:       selection.Exists,
		selectors.Selector_Expression_DoesNotExist: selection.DoesNotExist,
		selectors.Selector_Expression_GreaterThan:  selection.GreaterThan,
		selectors.Selector_Expression_LessThan:     selection.LessThan,
	}
)

func AllNamespaces(watchNamespaces []string) bool {

	if len(watchNamespaces) == 0 {
		return true
	}
	if len(watchNamespaces) == 1 && watchNamespaces[0] == "" {
		return true
	}
	return false
}

func ProcessWatchNamespaces(watchNamespaces []string, writeNamespace string) []string {
	if AllNamespaces(watchNamespaces) {
		return watchNamespaces
	}

	var writeNamespaceProvided bool
	for _, ns := range watchNamespaces {
		if ns == writeNamespace {
			writeNamespaceProvided = true
			break
		}
	}
	if !writeNamespaceProvided {
		watchNamespaces = append(watchNamespaces, writeNamespace)
	}

	return watchNamespaces
}

// ConvertExpressionSelectorToString will create a string representation of the Selector
// Expression data struct.
func ConvertExpressionSelectorToString(expressionSelectors []*selectors.Selector_Expression) (string, error) {
	if len(expressionSelectors) == 0 {
		return "", nil
	}
	var buffer bytes.Buffer
	endOfSelectors := len(expressionSelectors) - 1
	for i, sel := range expressionSelectors {
		op := sel.GetOperator()
		key := sel.GetKey()
		values := sel.GetValues()
		// TODO-JAKE might need to get rid of these selectors, and change the hybrid selectors used.
		if op == selectors.Selector_Expression_DoesNotExist || op == selectors.Selector_Expression_GreaterThan || op == selectors.Selector_Expression_LessThan {
			return "", errors.New("cannot select !, <, or > as operators for expression selectors")
		}
		switch op {
		case selectors.Selector_Expression_Exists, selectors.Selector_Expression_In,
			selectors.Selector_Expression_NotIn:
			buffer.WriteString(key)
			buffer.WriteByte(' ')
			buffer.WriteString(string(SelectorExpressionOperatorValues[op]))
			buffer.WriteByte(' ')
			buffer.WriteByte('(')
			buffer.WriteString(strings.Join(values, ","))
			buffer.WriteByte(')')
		case selectors.Selector_Expression_Equals, selectors.Selector_Expression_DoubleEquals,
			selectors.Selector_Expression_NotEquals:
			buffer.WriteString(key)
			buffer.WriteString(string(SelectorExpressionOperatorValues[op]))
			if len(values) > 1 {
				return "", errors.New("each expression selector operator must have a value associated with it")
			}
			buffer.WriteString(values[0])
		}
		if i < endOfSelectors {
			buffer.WriteByte(',')
		}
	}
	return buffer.String(), nil
}

func GetPodNamespace() string {
	if podNamespace := os.Getenv("POD_NAMESPACE"); podNamespace != "" {
		return podNamespace
	}
	return "gloo-system"
}
