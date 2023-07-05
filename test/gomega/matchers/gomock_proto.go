package matchers

import (
	"strings"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/onsi/gomega/types"
)

func GomockMatchProto4(actual proto.Message) gomock.Matcher {
	return &protoMatcherImpl{
		actual: actual,
	}
}

type protoMatcherImpl struct {
	actual proto.Message
}

func (p *protoMatcherImpl) Matches(actual interface{}) bool {
	if protoMsg, ok := actual.(proto.Message); ok {
		return proto.Equal(p.actual, protoMsg)
	}
	return false
}

func (p *protoMatcherImpl) String() string {
	return "equals proto " + p.actual.(proto.Message).String()
}

func MatchesPublicFields(actual interface{}) types.GomegaMatcher {
	return &publicFieldMatcher{actual: actual}
}

type publicFieldMatcher struct {
	actual interface{}
	diff   []string
}

func (p *publicFieldMatcher) Match(actual interface{}) (success bool, err error) {
	diff := deep.Equal(p.actual, actual)
	p.diff = diff
	return len(diff) == 0, nil
}

func (p *publicFieldMatcher) FailureMessage(actual interface{}) (message string) {
	return strings.Join(p.diff, " ")
}

func (p *publicFieldMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return "not equal"
}
