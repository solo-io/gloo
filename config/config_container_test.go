package config_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/glue/config"
)

var _ = Describe("Config", func() {
	Describe("Update", func() {
		var (
			called  bool
			handler UpdateHandler
			cfg     *Config
		)
		Context("handler returns an error", func() {
			called = false
			cfg = NewConfig()
			handler = func(blobs []byte) error {
				called = true
				return errors.New("intended error")
			}
			cfg.RegisterHandler(handler)
			err := cfg.Update([]byte("ignored"))
			It("returns an error", func() {
				Expect(called).To(BeTrue())
				Expect(err).To(HaveOccurred())
			})
		})
		Context("handler succeeds", func() {
			called = false
			cfg = NewConfig()
			handler = func(blobs []byte) error {
				called = true
				return nil
			}
			cfg.RegisterHandler(handler)
			err := cfg.Update([]byte("ignored"))
			It("calls the handler", func() {
				Expect(called).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
