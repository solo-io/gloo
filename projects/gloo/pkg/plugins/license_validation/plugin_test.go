package license_validation_test

import (
	"bytes"
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/pkg/license"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/license_validation"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// This license is set to expire on 2/26/2051
	// It was generated from solo-io/licensing using
	// go run pkg/cmd/genlicense/main.go -days 10000 -enterprise -product gloo
	validLicense   = "eyJleHAiOjI1NjEwNDIzNjQsImlhdCI6MTY5NzA0MjM2NCwiayI6ImVRc240USIsImx0IjoiZW50IiwicHJvZHVjdCI6Imdsb28ifQ.ITDsn8vi3n10zJXaa5E0bWPTj0VG4jEpCXjbdIO9qWg"
	invalidLicense = "bad license"
)

var _ = Describe("Plugin", func() {

	var (
		p      plugins.Plugin
		params plugins.InitParams

		buf    *bytes.Buffer
		bws    *zapcore.BufferedWriteSyncer
		logger *zap.SugaredLogger
	)

	BeforeEach(func() {
		buf = &bytes.Buffer{}
		bws = &zapcore.BufferedWriteSyncer{WS: zapcore.AddSync(buf)}
		encoderCfg := zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "logger",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), bws, zap.DebugLevel)
		logger = zap.New(core).Sugar()
		contextutils.SetFallbackLogger(logger)

		params = plugins.InitParams{Ctx: context.Background()}
	})

	When("invalid license is set on startup", func() {
		BeforeEach(func() {
			err := os.Setenv(license.EnvName, invalidLicense)
			Expect(err).NotTo(HaveOccurred())

			// In practice the LicensedFeatureProvider is created on startup and
			// ValidateAndSetLicense is called exactly once at that time
			lfp := license.NewLicensedFeatureProvider()
			lfp.ValidateAndSetLicense(context.Background())
			p = NewPlugin(lfp)
		})

		AfterEach(func() {
			err := os.Unsetenv(license.EnvName)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should error log a notice about invalid license even if license becomes valid", func() {
			By("testing without modifying invalid license")
			p.Init(params)

			err := logger.Sync()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("LICENSE ERROR"))

			By("testing after setting valid license")
			buf.Reset()
			err = os.Setenv(license.EnvName, validLicense)
			Expect(err).NotTo(HaveOccurred())

			p.Init(params)

			err = logger.Sync()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("LICENSE ERROR"))
		})
	})

	When("valid license is set on startup", func() {
		BeforeEach(func() {
			err := os.Setenv(license.EnvName, validLicense)
			Expect(err).NotTo(HaveOccurred())

			// In practice the LicensedFeatureProvider is created on startup and
			// ValidateAndSetLicense is called exactly once at that time
			lfp := license.NewLicensedFeatureProvider()
			lfp.ValidateAndSetLicense(context.Background())
			p = NewPlugin(lfp)
		})

		AfterEach(func() {
			err := os.Unsetenv(license.EnvName)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not error log a notice about invalid license, even if license becomes invalid", func() {
			By("testing without modifying valid license")
			p.Init(params)

			err := logger.Sync()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).NotTo(ContainSubstring("LICENSE ERROR"))

			By("testing after modifying valid license")
			buf.Reset()
			err = os.Setenv(license.EnvName, invalidLicense)
			Expect(err).NotTo(HaveOccurred())

			p.Init(params)

			err = logger.Sync()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).NotTo(ContainSubstring("LICENSE ERROR"))
		})
	})
})
