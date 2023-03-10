package federation_test

import (
	"context"
	"errors"
	"math"
	"os"
	"time"

	"github.com/avast/retry-go/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/skv2/pkg/multicluster/watch"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
)

var _ = Describe("Retries", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		retryEnvVars map[string]string
	)

	setEnv := func(vars map[string]string) {
		retryEnvVars = vars

		var err error
		for key, val := range retryEnvVars {
			err = os.Setenv(key, val)
			Expect(err).NotTo(HaveOccurred())
		}
	}

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		var err error
		for key := range retryEnvVars {
			err = os.Unsetenv(key)
			Expect(err).NotTo(HaveOccurred())
		}

		cancel()
	})

	Context("GetClusterWatcherRemoteRetryOptions", func() {
		// In these tests we just verify that we pass the expected retry options to skv2
		It("returns default values when no env vars are set", func() {
			opt := federation.GetClusterWatcherRemoteRetryOptions(ctx)
			expectedDelay := federation.DefaultDelayRemote
			expectedMaxDelay := federation.DefaultMaxDelayRemote
			expectedMaxJitter := federation.DefaultMaxJitterRemote
			expectedAttempts := federation.DefaultAttemptsRemote
			Expect(opt).To(Equal(watch.RetryOptions{
				DelayType: watch.RetryDelayType_Backoff,
				Delay:     &expectedDelay,
				MaxDelay:  &expectedMaxDelay,
				MaxJitter: &expectedMaxJitter,
				Attempts:  &expectedAttempts,
			}))
		})

		It("can set values via env variables", func() {
			setEnv(map[string]string{
				federation.ClusterWatcherRemoteRetryTypeEnv:      "fixed",
				federation.ClusterWatcherRemoteRetryDelayEnv:     "700ms",
				federation.ClusterWatcherRemoteRetryMaxDelayEnv:  "0",
				federation.ClusterWatcherRemoteRetryMaxJitterEnv: "50ms",
				federation.ClusterWatcherRemoteRetryAttemptsEnv:  "73",
			})
			opt := federation.GetClusterWatcherRemoteRetryOptions(ctx)
			expectedDelay := 700 * time.Millisecond
			expectedMaxDelay := 0 * time.Second
			expectedMaxJitter := 50 * time.Millisecond
			expectedAttempts := uint(73)
			Expect(opt).To(Equal(watch.RetryOptions{
				DelayType: watch.RetryDelayType_Fixed,
				Delay:     &expectedDelay,
				MaxDelay:  &expectedMaxDelay,
				MaxJitter: &expectedMaxJitter,
				Attempts:  &expectedAttempts,
			}))
		})

		It("uses default value if an env variable is invalid", func() {
			setEnv(map[string]string{
				federation.ClusterWatcherRemoteRetryTypeEnv:      "blah",   // invalid
				federation.ClusterWatcherRemoteRetryDelayEnv:     "123num", // invalid
				federation.ClusterWatcherRemoteRetryMaxDelayEnv:  "1m",
				federation.ClusterWatcherRemoteRetryMaxJitterEnv: "0",
				federation.ClusterWatcherRemoteRetryAttemptsEnv:  "-1", // invalid
			})
			opt := federation.GetClusterWatcherRemoteRetryOptions(ctx)
			expectedDelay := federation.DefaultDelayRemote
			expectedMaxDelay := 1 * time.Minute
			expectedMaxJitter := 0 * time.Second
			expectedAttempts := federation.DefaultAttemptsRemote
			Expect(opt).To(Equal(watch.RetryOptions{
				DelayType: watch.RetryDelayType_Backoff,
				Delay:     &expectedDelay,
				MaxDelay:  &expectedMaxDelay,
				MaxJitter: &expectedMaxJitter,
				Attempts:  &expectedAttempts,
			}))
		})
	})

	Context("GetClusterWatcherLocalRetryOptions", func() {
		It("uses default values when no env vars are set", func() {
			opts := federation.GetClusterWatcherLocalRetryOptions(ctx)

			// call a func that always returns an error, with the given retry options.
			// keep track of total time taken, and total attempts, and compare to expected values.
			start := time.Now()
			count := uint(0)
			err := retry.Do(func() error {
				count++
				return errors.New("this is an error")
			}, opts...)
			totalTime := time.Since(start)

			Expect(err).To(HaveOccurred())

			// we expect all the attempts to be exhausted
			Expect(count).To(Equal(federation.DefaultAttemptsLocal), "incorrect number of attempts")

			// this is approximately the total time the retries should take, with exponential backoff
			expectedTotalTime := federation.DefaultDelayLocal *
				time.Duration(math.Exp2(float64(federation.DefaultAttemptsLocal)-1))
			// compare to the actual time, taking into account some randomness/jitter
			Expect(totalTime).To(BeNumerically("~", expectedTotalTime, 500*time.Millisecond))
		})

		It("can set values via env variables", func() {
			setEnv(map[string]string{
				federation.ClusterWatcherLocalRetryTypeEnv:      "fixed",
				federation.ClusterWatcherLocalRetryDelayEnv:     "70ms",
				federation.ClusterWatcherLocalRetryMaxDelayEnv:  "1m",
				federation.ClusterWatcherLocalRetryMaxJitterEnv: "0",
				federation.ClusterWatcherLocalRetryAttemptsEnv:  "7",
			})

			opts := federation.GetClusterWatcherLocalRetryOptions(ctx)

			start := time.Now()
			count := uint(0)
			err := retry.Do(func() error {
				count++
				return errors.New("this is an error")
			}, opts...)
			totalTime := time.Since(start)

			Expect(err).To(HaveOccurred())
			Expect(count).To(Equal(uint(7)), "incorrect number of attempts")

			// 70ms fixed retry intervals and 6 retries (attempts minus 1)
			expectedTotalTime := 70 * time.Millisecond * 6

			// compare to the actual time. since we set jitter to 0, the threshold can be smaller
			Expect(totalTime).To(BeNumerically("~", expectedTotalTime, 50*time.Millisecond))
		})
	})

})
