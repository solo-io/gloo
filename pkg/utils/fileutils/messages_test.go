package fileutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"

	"github.com/gogo/protobuf/types"
	. "github.com/solo-io/solo-kit/pkg/utils/fileutils"
)

var _ = Describe("Messages", func() {
	var filename string
	BeforeEach(func() {
		f, err := ioutil.TempFile("", "messages_test")
		Expect(err).NotTo(HaveOccurred())
		filename = f.Name()
	})
	AfterEach(func() {
		os.RemoveAll(filename)
	})
	It("Writes and reads proto messages into files", func() {
		input := &types.Struct{
			Fields: map[string]*types.Value{
				"foo": {
					Kind: &types.Value_StringValue{StringValue: "bar"},
				},
			},
		}
		err := WriteToFile(filename, input)
		Expect(err).NotTo(HaveOccurred())

		b, err := ioutil.ReadFile(filename)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(b)).To(Equal("foo: bar\n"))

		var output types.Struct
		err = ReadFileInto(filename, &output)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(Equal(*input))
	})
})
