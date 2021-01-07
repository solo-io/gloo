package test_matchers

import (
	"fmt"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
)

func MatchesFailover(failScheme *fedv1.FailoverScheme) gomock.Matcher {
	return &failoverMatcher{
		failScheme: failScheme,
	}
}

type failoverMatcher struct {
	failScheme *fedv1.FailoverScheme
}

func (f *failoverMatcher) Matches(x interface{}) bool {
	failScheme, ok := x.(*fedv1.FailoverScheme)
	if !ok {
		return false
	}

	if !reflect.DeepEqual(failScheme.ObjectMeta, f.failScheme.ObjectMeta) ||
		!reflect.DeepEqual(failScheme.Spec.Primary, f.failScheme.Spec.Primary) ||
		!reflect.DeepEqual(failScheme.Spec.FailoverGroups, f.failScheme.Spec.FailoverGroups) {
		return false
	}
	// skip processing time as it will always be different
	return failScheme.Status.GetState() == f.failScheme.Status.GetState() &&
		failScheme.Status.GetObservedGeneration() == f.failScheme.Status.GetObservedGeneration() &&
		failScheme.Status.GetMessage() == f.failScheme.Status.GetMessage()
}

func (f *failoverMatcher) String() string {
	return fmt.Sprintf("%v", f.failScheme)
}

func MatchFailoverStatus(failScheme *fed_types.FailoverSchemeStatus) types.GomegaMatcher {
	return &failoverStatusMatcher{
		failSchemeStatus: failScheme,
	}
}

type failoverStatusMatcher struct {
	failSchemeStatus *fed_types.FailoverSchemeStatus
}

func (f *failoverStatusMatcher) Match(actual interface{}) (success bool, err error) {
	failScheme, ok := actual.(*fed_types.FailoverSchemeStatus)
	if !ok {
		return false, nil
	}
	// skip processing time as it will always be different
	return failScheme.GetState() == f.failSchemeStatus.GetState() &&
		failScheme.GetObservedGeneration() == f.failSchemeStatus.GetObservedGeneration() &&
		failScheme.GetMessage() == f.failSchemeStatus.GetMessage(), nil
}

func (f *failoverStatusMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "To be identical except the time to", f.failSchemeStatus)
}

func (f *failoverStatusMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "Not to be identical except the time to", f.failSchemeStatus)
}
