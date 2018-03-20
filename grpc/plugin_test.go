package grpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"

	. "github.com/solo-io/gloo-plugins/grpc"
	"github.com/solo-io/gloo/pkg/log"
)

var _ = Describe("Plugin", func() {
	It("unmarshal file descriptor", func() {
		b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
		Expect(err).NotTo(HaveOccurred())
		descriptor, err := ConvertProto(b)
		Expect(err).NotTo(HaveOccurred())
		log.Printf("%v", FuncsForProto("Bookstore", descriptor))
	})
	FIt("unmarshal file descriptor", func() {
		b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
		Expect(err).NotTo(HaveOccurred())
		descriptor, err := ConvertProto(b)
		Expect(err).NotTo(HaveOccurred())
		AddHttpRulesToProto("my-upstream", "Bookstore", descriptor)
		log.Printf("%v", FuncsForProto("Bookstore", descriptor))
	})
})
