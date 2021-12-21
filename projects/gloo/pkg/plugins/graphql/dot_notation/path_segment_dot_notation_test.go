package dot_notation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/dot_notation"
)

func Key(key string) *v2.PathSegment {
	return &v2.PathSegment{Segment: &v2.PathSegment_Key{Key: key}}
}

func All() *v2.PathSegment {
	return &v2.PathSegment{Segment: &v2.PathSegment_All{All: true}}
}

func Index(index int) *v2.PathSegment {
	return &v2.PathSegment{Segment: &v2.PathSegment_Index{Index: uint32(index)}}
}

var _ = Describe("Dot Notation Test", func() {
	testDotNotationTranslation := func(dot string, pathSegments []*v2.PathSegment, expectedError error) {

		path, err := dot_notation.DotNotationToPathSegments(dot)
		if expectedError == nil {
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
		} else {
			ExpectWithOffset(1, err).To(MatchError(expectedError))
		}
		ExpectWithOffset(1, path).To(HaveLen(len(pathSegments)))
		for idx, pathSegment := range path {
			ExpectWithOffset(1, pathSegments[idx]).To(matchers.MatchProto(pathSegment))
		}
	}

	Context("Test Dot Notation translation", func() {
		It("translates chained keys correctly", func() {

			testDotNotationTranslation("a", []*v2.PathSegment{
				Key("a"),
			}, nil)

			testDotNotationTranslation("a.b.c", []*v2.PathSegment{
				Key("a"),
				Key("b"),
				Key("c"),
			}, nil)

		})
		It("catches trailing dot correctly", func() {
			testDotNotationTranslation("a.b.c.", nil, errors.New("Unable to parse 'a.b.c.' due to trailing dot!"))
			testDotNotationTranslation(".", nil, errors.New("Unexpected char: ., index: 0, key: ."))
		})
		It("translates chained indices correctly", func() {

			testDotNotationTranslation("[9]", []*v2.PathSegment{
				Index(9),
			}, nil)

			testDotNotationTranslation("a[9].c[0]", []*v2.PathSegment{
				Key("a"),
				Index(9),
				Key("c"),
				Index(0),
			}, nil)

			testDotNotationTranslation("a[*][*].c[9]", []*v2.PathSegment{
				Key("a"),
				All(),
				All(),
				Key("c"),
				Index(9),
			}, nil)

		})
		It("catches misused brackets correctly", func() {
			testDotNotationTranslation("a[]", nil, errors.New("Unexpected char: ], index: 2, key: a[]"))
			testDotNotationTranslation("a.[*]", nil, errors.New("Unexpected char: [, index: 2, key: a.[*]"))
		})

		It("translates special cases correctly", func() {

			testDotNotationTranslation("headers.:method", []*v2.PathSegment{
				Key("headers"),
				Key(":method"),
			}, nil)

		})
	})
})
