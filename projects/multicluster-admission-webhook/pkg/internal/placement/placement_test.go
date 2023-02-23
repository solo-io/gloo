package placement_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/placement"
)

var _ = Describe("Placement", func() {
	var (
		ctx     = context.TODO()
		matcher placement.Matcher
	)

	BeforeEach(func() {
		matcher = placement.NewMatcher()
	})

	It("will return false if resource placement is invalid", func() {
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{},
			&multicluster_types.Placement{
				Clusters: []string{"cluster-1", "cluster-3"},
			},
		)
		Expect(matches).To(BeFalse())
	})

	It("will fail if rule placement is invalid", func() {
		matcher := placement.NewMatcher()
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1", "cluster-2"},
			},
			&multicluster_types.Placement{},
		)
		Expect(matches).To(BeFalse())
	})

	It("will not match if rule clusters are non-wildcard, and resource clusters are wildcard", func() {
		matcher := placement.NewMatcher()
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"*"},
			},
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1", "cluster-3"},
			},
		)
		Expect(matches).To(BeFalse())
	})

	It("will not match if rule does not contain all clusters specified in resource", func() {
		matcher := placement.NewMatcher()
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1", "cluster-3"},
			},
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1"},
			},
		)
		Expect(matches).To(BeFalse())
	})

	It("will not match if rule namespaces are non-wildcard, and resource namespaces are wildcard", func() {
		matcher := placement.NewMatcher()
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{
				Namespaces: []string{"*"},
				Clusters:   []string{"cluster-1", "cluster-3"},
			},
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1", "cluster-3"},
			},
		)
		Expect(matches).To(BeFalse())
	})

	It("will not match if rule does not contain all namespaces specified in resource", func() {
		matcher := placement.NewMatcher()
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1", "cluster-3"},
			},
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1"},
			},
		)
		Expect(matches).To(BeFalse())
	})

	It("will match if rule is wildcard", func() {
		matcher := placement.NewMatcher()
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1", "cluster-3"},
			},
			&multicluster_types.Placement{
				Namespaces: []string{"*"},
				Clusters:   []string{"*"},
			},
		)
		Expect(matches).To(BeTrue())
	})

	It("will match if rule is superset of resource", func() {
		matcher := placement.NewMatcher()
		matches := matcher.Matches(
			ctx,
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1"},
				Clusters:   []string{"cluster-1"},
			},
			&multicluster_types.Placement{
				Namespaces: []string{"namespace-1", "namespace-2"},
				Clusters:   []string{"cluster-1", "cluster-3"},
			},
		)
		Expect(matches).To(BeTrue())
	})
})
