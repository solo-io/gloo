package helpers

import . "github.com/onsi/gomega"

func Must(err error) {
	Expect(err).NotTo(HaveOccurred())
}
