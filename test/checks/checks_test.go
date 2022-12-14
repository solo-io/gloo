package checks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checks mod", func() {

	// Ginkgo v2.0.0 contains describetable so this could be wr
	// 	DescribeTable("needed replaces are present",

	// 	func(name, path, oldPath, version string){
	// 		allPackages, err := modfile.Parse()
	// 		Expect(err).NotTo(HaveOccurred())

	// 		replacedPackages := allPackages.Replace
	// 		Expect(replacedPackages).NotTo(BeNil())

	// 		var replace modfile.ReplacedGoPackage
	// 		for _, replacedGoPkg := range replacedPackages {
	// 			if replacedGoPkg.Old.Path == oldPath {
	// 				replace = replacedGoPkg
	// 				break
	// 			}
	// 		}

	// 		Expect(replace).NotTo(BeNil())
	// 		Expect(replace.New.Path).To(Equal(path))
	// 		Expect(replace.New.Version).To(Equal(version))
	// 	},

	// 	// we no longer need to replace klog but this was a long time example Entry("k8s.io/klog", "k8s.io/klog", "github.com/stefanprodan/klog", "v0.0.0-20181102134211-b9b56d5dfc92"),

	// )

})
