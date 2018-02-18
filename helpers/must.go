package helpers

import "github.com/onsi/ginkgo"

func Must(err error) {
	if err != nil {
		ginkgo.Fail(err.Error())
	}
}
