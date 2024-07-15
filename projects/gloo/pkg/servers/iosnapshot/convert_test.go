package iosnapshot

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
)

var _ = Describe("Convert", func() {

	It("will fail if the ApiSnapshot has a new top level field", func() {
		// This test checks whether the ApiSnapshot has a new top-level field.
		// If the number of fields changes, it should be used as an indication that the
		// `snapshotToKubeResources` function might need to be updated to support the new field.
		Expect(reflect.TypeOf(v1snap.ApiSnapshot{}).NumField()).To(
			Equal(16),
			"wrong number of fields found in ApiSnapshot",
		)
	})
})
