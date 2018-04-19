package consul

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/helpers/local"
)

var _ = Describe("Consul UpstreamController Units", func() {
	Describe("uniqueTagSets", func() {
		It("takes a list of unsorted, duplicated tags and sorts and de-dupes them", func() {
			input := [][]string{
				{"foo", "bar", "baz"},
				{"foo", "baz"},
				{"baz", "foo"},
				{"baz", "fooa"},
			}
			out := uniqueTagSets(input)
			Expect(out).To(Equal([][]string{
				{"baz", "foo"},
				{"baz", "fooa"},
				{"bar", "baz", "foo"},
			}))
		})
	})
})

var _ = Describe("Consul UpstreamController Units", func() {
	if os.Getenv("RUN_CONSUL_TESTS") != "1" {
		log.Printf("This test downloads and runs consul and is disabled by default. To enable, set RUN_CONSUL_TESTS=1 in your env.")
		return
	}

	var (
		consulFactory  *localhelpers.ConsulFactory
		consulInstance *localhelpers.ConsulInstance
		err            error
	)

	var _ = BeforeSuite(func() {
		consulFactory, err = localhelpers.NewConsulFactory()
		helpers.Must(err)
		consulInstance, err = consulFactory.NewConsulInstance()
		helpers.Must(err)
		err = consulInstance.Run()
		helpers.Must(err)
	})

	var _ = AfterSuite(func() {
		consulInstance.Clean()
		consulFactory.Clean()
	})

	Describe("uniqueTagSets", func() {
		It("takes a list of unsorted, duplicated tags and sorts and de-dupes them", func() {
			input := [][]string{
				{"foo", "bar", "baz"},
				{"foo", "baz"},
				{"baz", "foo"},
				{"baz", "fooa"},
			}
			out := uniqueTagSets(input)
			Expect(out).To(Equal(nil))
		})
	})
})
