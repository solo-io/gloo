package matchers

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/solo-io/solo-projects/test/gomega/transforms"
)

type CookieDataMatcher struct {
	c *transforms.CookieData
}

func (m *CookieDataMatcher) Match(actual interface{}) (success bool, err error) {
	actualData, ok := actual.(*transforms.CookieData)
	if !ok {
		return false, fmt.Errorf("actual value must be of type *extauth_test.CookieData")
	}

	if m.c.Value == nil {
		// Allows for partially matching in the case that we only care about other field(s),
		// or the value is generated/unknown.
		m.c.Value = actualData.Value
	}

	var valueMatcher types.GomegaMatcher
	switch m.c.Value.(type) {
	case string:
		valueMatcher = Equal(m.c.Value)
	case types.GomegaMatcher:
		valueMatcher = m.c.Value.(types.GomegaMatcher)
	default:
		return false, fmt.Errorf("expected value must be either a string or GomegaMatcher")
	}

	ok, err = valueMatcher.Match(actualData.Value)
	if err != nil || !ok {
		return ok, err
	}

	matcher := HaveField("HttpOnly", m.c.HttpOnly)
	ok, err = matcher.Match(actualData)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (m *CookieDataMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("%s \n", m.FailureMessage(actual))
}

func (m *CookieDataMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("%s \n", m.NegatedFailureMessage(actual))
}

func MatchCookieData(c *transforms.CookieData) *CookieDataMatcher {
	return &CookieDataMatcher{
		c: c,
	}
}
