package aws

import (
	"log"
	"testing"
	"time"
)

func TestDiffWithOneNilArray(t *testing.T) {
	l := []Lambda{Lambda{Name: "f", Qualifier: "q"}}
	var r []Lambda
	if !diff(l, r) {
		t.Error("expected nil array to be different")
	}
}

func TestDiffWithDifferentLength(t *testing.T) {
	l := []Lambda{Lambda{Name: "f", Qualifier: "q"}}
	r := []Lambda{Lambda{Name: "f", Qualifier: "q"},
		Lambda{Name: "f", Qualifier: "q1"}}
	if !diff(l, r) {
		t.Error("expected arrays to be different")
	}
}

func TestDiffWithSimilarArray(t *testing.T) {
	l := []Lambda{Lambda{Name: "f", Qualifier: "q1"},
		Lambda{Name: "f", Qualifier: "q"}}
	r := []Lambda{Lambda{Name: "f", Qualifier: "q"},
		Lambda{Name: "f", Qualifier: "q1"}}
	if diff(l, r) {
		t.Error("expected arrays to be same")
	}
}

func TestUpdaterIsCalledIntially(t *testing.T) {
	f := func(r string, a AccessToken) ([]Lambda, error) {
		log.Println("fetching for " + r)
		return []Lambda{Lambda{Name: r + "_func1",
			Qualifier: "$LATEST"}}, nil
	}
	stop := make(chan struct{})
	called := false
	timer := time.NewTimer(1 * time.Second)
	go func() {
		<-timer.C
		if !called {
			t.Error("timed out and update still not called")
		}
	}()

	u := func(r Region) error {
		log.Println("updating ", r)
		called = true
		close(stop)
		return nil
	}
	awsPoller := NewAWSPoller(f, u)
	awsPoller.AddUpdateRegion(Region{
		ID:   "test",
		Name: "us-east-1",
	})

	awsPoller.Start(1*time.Millisecond, stop)
	<-stop
	timer.Stop()
}

func TestUpdaterIsNotCalledIfLambdasDonotChange(t *testing.T) {
	lambdas := []Lambda{Lambda{Name: "func", Qualifier: "v1"}}
	f := func(r string, a AccessToken) ([]Lambda, error) {
		log.Println("fetching for " + r)
		return lambdas, nil
	}
	stop := make(chan struct{})
	timer := time.NewTimer(500 * time.Millisecond)
	go func() {
		<-timer.C
		close(stop)
	}()

	u := func(r Region) error {
		log.Println("updating ", r)
		t.Error("lambdas didn't change should not be called")
		return nil
	}
	awsPoller := NewAWSPoller(f, u)
	region := Region{
		ID:      "test",
		Name:    "us-east-1",
		Lambdas: lambdas}
	awsPoller.AddUpdateRegion(region)

	awsPoller.Start(1*time.Millisecond, stop)
	<-stop
}
