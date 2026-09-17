package validator

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/constants"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
)

// Regression tests for interrupted validation
//
// A fake envoy script replaces the real binary so this runs anywhere. It drains stdin and reports
// the config as OK, optionally sleeping first so the test can cancel mid-run.
var _ = Describe("Validator with a cancelled context", func() {

	var (
		validConfig = &transformation.TransformationStages{InheritTransformation: true}
	)

	writeFakeEnvoy := func(sleepSeconds string) {
		dir := GinkgoT().TempDir()
		script := filepath.Join(dir, "envoy")
		// `exec sleep` replaces the shell so the SIGKILL lands on the process holding the stdout pipe.
		// Otherwise Wait blocks until the orphaned sleep exits.
		content := "#!/bin/sh\ncat >/dev/null\necho \"configuration '/dev/fd/0' OK\"\nexec sleep " + sleepSeconds + "\n"
		Expect(os.WriteFile(script, []byte(content), 0o755)).To(Succeed())
		GinkgoT().Setenv(constants.EnvoyBinaryEnv, script)
	}

	// writeFakeEnvoyRejecting rejects every config with an error report and exit status 1.
	writeFakeEnvoyRejecting := func() {
		dir := GinkgoT().TempDir()
		script := filepath.Join(dir, "envoy")
		content := "#!/bin/sh\ncat >/dev/null\necho \"error initializing configuration '/dev/fd/0': not a valid config\" >&2\nexit 1\n"
		Expect(os.WriteFile(script, []byte(content), 0o755)).To(Succeed())
		GinkgoT().Setenv(constants.EnvoyBinaryEnv, script)
	}

	It("validates a good config with a live context (sanity check of the fake envoy)", func() {
		writeFakeEnvoy("0")
		v := New("cancel-sanity", "cancel-sanity")

		Expect(v.ValidateConfig(context.Background(), validConfig)).To(Succeed())
		Expect(v.CacheLength()).To(Equal(1))
	})

	It("does not cache an interrupted validation (context cancelled before the fork starts)", func() {
		writeFakeEnvoy("0")
		v := New("cancel-before", "cancel-before")

		cancelled, cancel := context.WithCancel(context.Background())
		cancel()

		err := v.ValidateConfig(cancelled, validConfig)
		Expect(err).To(HaveOccurred(), "an interrupted validation should be reported to the caller")
		Expect(err.Error()).To(ContainSubstring("context canceled"))

		Expect(v.CacheLength()).To(Equal(0), "interrupted validation was cached as a config error")

		// A later validation of the same config on a healthy context must re-run envoy.
		Expect(v.ValidateConfig(context.Background(), validConfig)).To(Succeed(),
			"healthy translation was served the cached cancellation error")
	})

	// writeFakeEnvoySignalKilled dies on SIGKILL with the caller's context still live, like a fork
	// killed from outside (OOM killer, manual kill).
	writeFakeEnvoySignalKilled := func() {
		dir := GinkgoT().TempDir()
		script := filepath.Join(dir, "envoy")
		content := "#!/bin/sh\ncat >/dev/null\nkill -9 $$\necho \"configuration '/dev/fd/0' OK\"\n"
		Expect(os.WriteFile(script, []byte(content), 0o755)).To(Succeed())
		GinkgoT().Setenv(constants.EnvoyBinaryEnv, script)
	}

	It("does not cache a validation killed by an outside signal under a live context", func() {
		writeFakeEnvoySignalKilled()
		v := New("signal-killed", "signal-killed")

		err := v.ValidateConfig(context.Background(), validConfig)
		Expect(err).To(HaveOccurred(), "a killed validation should be reported to the caller")
		Expect(err.Error()).To(ContainSubstring("interrupted"))

		Expect(v.CacheLength()).To(Equal(0), "a signal-killed validation was cached as a config error")

		// With no outside kill, the same config must validate cleanly.
		writeFakeEnvoy("0")
		Expect(v.ValidateConfig(context.Background(), validConfig)).To(Succeed(),
			"healthy translation was served the stale signal-killed result")
	})

	It("still returns and caches a genuine rejection under a live context", func() {
		writeFakeEnvoyRejecting()
		v := New("reject", "reject")

		err := v.ValidateConfig(context.Background(), validConfig)
		Expect(err).To(HaveOccurred(), "a genuine rejection must be reported to the caller")
		Expect(err.Error()).NotTo(ContainSubstring("interrupted"), "a rejection must be distinguishable from an interruption")
		Expect(v.CacheLength()).To(Equal(1), "a genuine rejection is a property of the config and must be cached")

		// Swap in an envoy that accepts everything. The rejection must still be served, proving it
		// comes from the cache.
		writeFakeEnvoy("0")
		Expect(v.ValidateConfig(context.Background(), validConfig)).To(HaveOccurred(),
			"the cached rejection should be re-served without re-running envoy")
	})

	It("does not let a doomed context contaminate a concurrent validation of the same config", func() {
		// Neither call finds a cached result for the config, so both fork envoy. One context is
		// cancelled mid-run, the other completes. The cancelled run must not land in the shared cache.
		writeFakeEnvoy("2")
		v := New("concurrent", "concurrent")

		doomedCtx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		var doomedErr, healthyErr error
		wg.Add(2)
		go func() {
			defer GinkgoRecover()
			defer wg.Done()
			doomedErr = v.ValidateConfig(doomedCtx, validConfig)
		}()
		go func() {
			defer GinkgoRecover()
			defer wg.Done()
			healthyErr = v.ValidateConfig(context.Background(), validConfig)
		}()
		time.Sleep(300 * time.Millisecond)
		cancel()
		wg.Wait()

		Expect(doomedErr).To(HaveOccurred(), "the interrupted validation should be reported to its own caller")
		Expect(healthyErr).NotTo(HaveOccurred(), "the healthy validation must not be served the other caller's interruption")
		Expect(v.ValidateConfig(context.Background(), validConfig)).To(Succeed(),
			"the shared cache must not hold the interruption")
	})

	It("does not cache an interrupted validation (context cancelled while envoy is running)", func() {
		writeFakeEnvoy("5")
		v := New("cancel-during", "cancel-during")

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(300 * time.Millisecond)
			cancel()
		}()

		start := time.Now()
		err := v.ValidateConfig(ctx, validConfig)
		Expect(time.Since(start)).To(BeNumerically("<", 4*time.Second), "cancellation should have killed the child")
		Expect(err).To(HaveOccurred())

		Expect(v.CacheLength()).To(Equal(0), "interrupted validation was cached as a config error")

		// Use a fast fake envoy for the healthy retry so the test stays quick.
		writeFakeEnvoy("0")
		Expect(v.ValidateConfig(context.Background(), validConfig)).To(Succeed(),
			"healthy translation was served the cached cancellation error")
	})
})
