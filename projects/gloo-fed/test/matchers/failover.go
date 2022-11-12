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

	return StatusesAreEqual(&failScheme.Status, &f.failScheme.Status)
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

	return StatusesAreEqual(failScheme, f.failSchemeStatus), nil
}

func (f *failoverStatusMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "To be identical except the time to", f.failSchemeStatus)
}

func (f *failoverStatusMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "Not to be identical except the time to", f.failSchemeStatus)
}

// StatusesAreEqual returns true if the two StatusesAreEqual objects are identical, false otherwise
func StatusesAreEqual(left, right *fed_types.FailoverSchemeStatus) bool {
	leftStatuses := left.GetNamespacedStatuses()
	rightStatuses := right.GetNamespacedStatuses()

	if len(leftStatuses) != len(rightStatuses) {
		return false
	}

	for ns, leftStatus := range leftStatuses {
		if !SingleStatusesAreEqual(leftStatus, rightStatuses[ns]) {
			return false
		}
	}
	return true
}

// SingleStatusesAreEqual returns true if the two FailoverSchemeStatus_Status objects are identical, false otherwise
func SingleStatusesAreEqual(left, right *fed_types.FailoverSchemeStatus_Status) bool {
	// we skip processing time as it will always be different
	return left.GetState() == right.GetState() &&
		left.GetObservedGeneration() == right.GetObservedGeneration() &&
		left.GetMessage() == right.GetMessage()
}
