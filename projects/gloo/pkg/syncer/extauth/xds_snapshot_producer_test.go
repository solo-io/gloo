package extauth_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("xDS Snapshot Producer", func() {

	Context("proxySourcedXdsSnapshotProducer", func() {
		// We intend to remove this instance of the XdsSnapshotProducer
		// All existing tests live in ../extauth_translator_syncer_test.go
	})

	Context("snapshotSourcedXdsSnapshotProducer", func() {
		// TODO test the behavior of the memoized xds snapshot producer
	})
})
