package federation

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/pkg/multicluster/watch"
)

// There are 2 places we currently have configurable retries in Gloo Fed:
//  1. When the skv2 cluster watcher fails to create/start a new manager, in response to a cluster
//     event (i.e. new cluster registered/detected)
//  2. When gloo fed fails to reset the status for federated resources, in response to a cluster event.
//     Note, this only refers to status updates triggered by cluster-related changes. For other changes
//     like updates to federated resources, we rely on controller-runtime to requeue events when
//     reconcile fails.
const (
	// these correspond the the options we will pass to skv2 cluster watcher to retry manager
	// creation failures
	ClusterWatcherRemoteRetryTypeEnv      = "CW_REMOTE_RETRY_TYPE"
	ClusterWatcherRemoteRetryDelayEnv     = "CW_REMOTE_RETRY_DELAY"
	ClusterWatcherRemoteRetryMaxDelayEnv  = "CW_REMOTE_RETRY_MAX_DELAY"
	ClusterWatcherRemoteRetryMaxJitterEnv = "CW_REMOTE_RETRY_MAX_JITTER"
	ClusterWatcherRemoteRetryAttemptsEnv  = "CW_REMOTE_RETRY_ATTEMPTS"

	// these are the options we use in the gloo fed cluster handlers to retry fed status
	// update failures
	ClusterWatcherLocalRetryTypeEnv      = "CW_LOCAL_RETRY_TYPE"
	ClusterWatcherLocalRetryDelayEnv     = "CW_LOCAL_RETRY_DELAY"
	ClusterWatcherLocalRetryMaxDelayEnv  = "CW_LOCAL_RETRY_MAX_DELAY"
	ClusterWatcherLocalRetryMaxJitterEnv = "CW_LOCAL_RETRY_MAX_JITTER"
	ClusterWatcherLocalRetryAttemptsEnv  = "CW_LOCAL_RETRY_ATTEMPTS"

	// the retry type env variable values are strings which we'll convert to the
	// appropriate types to pass to skv2 or retry-go
	RetryTypeBackoff = "backoff"
	RetryTypeFixed   = "fixed"

	// reasonable defaults for the various options
	DefaultTypeRemote      = RetryTypeBackoff
	DefaultDelayRemote     = time.Second
	DefaultMaxDelayRemote  = 0 * time.Second
	DefaultMaxJitterRemote = 100 * time.Millisecond
	DefaultAttemptsRemote  = uint(0)

	DefaultTypeLocal      = RetryTypeBackoff
	DefaultDelayLocal     = 100 * time.Millisecond
	DefaultMaxDelayLocal  = 0 * time.Second
	DefaultMaxJitterLocal = 100 * time.Millisecond
	DefaultAttemptsLocal  = uint(5)
)

// GetClusterWatcherRemoteRetryOptions gets retry options to pass to the skv2 cluster watcher.
//
// Whenever any changes are made to kubeconfig secrets (which contain remote cluster config),
// Gloo Fed starts a new manager using the remote cluster config (via a skv2 clusterWatcher).
// In the case that the manager creation/startup fails (e.g. remote cluster is not reachable),
// the cluster watcher will retry using the policy defined by the CW_REMOTE_* env variables.
//
// For any variables that are not set or not parseable, we will fall back to reasonable defaults.
func GetClusterWatcherRemoteRetryOptions(ctx context.Context) watch.RetryOptions {
	var delayType watch.RetryDelayType
	if getRetryTypeFromEnv(ctx, ClusterWatcherRemoteRetryTypeEnv, DefaultTypeRemote) == RetryTypeFixed {
		delayType = watch.RetryDelayType_Fixed
	} else {
		delayType = watch.RetryDelayType_Backoff
	}

	delay := getDurationFromEnv(ctx, ClusterWatcherRemoteRetryDelayEnv, DefaultDelayRemote)
	maxDelay := getDurationFromEnv(ctx, ClusterWatcherRemoteRetryMaxDelayEnv, DefaultMaxDelayRemote)
	maxJitter := getDurationFromEnv(ctx, ClusterWatcherRemoteRetryMaxJitterEnv, DefaultMaxJitterRemote)
	attempts := getUintFromEnv(ctx, ClusterWatcherRemoteRetryAttemptsEnv, DefaultAttemptsRemote)

	return watch.RetryOptions{
		DelayType: delayType,
		Delay:     &delay,
		MaxDelay:  &maxDelay,
		MaxJitter: &maxJitter,
		Attempts:  &attempts,
	}
}

// GetClusterWatcherLocalRetryOptions gets retry options for updating federated resource
// statuses when cluster changes are detected.
//
// When changes to a managed cluster occur (such as registering or deregistering a cluster),
// Gloo Fed triggers a reconcile loop by resetting the status for that cluster on all the
// federated resources. Sometimes the status update can fail (e.g. due to write errors), which
// will cause Gloo Fed not to re-reconcile. In this case, the status updates will be retried
// using the policy defined by the CW_LOCAL_* env variables.
//
// For any variables that are not set or not parseable, we will fall back to reasonable defaults.
//
// Note, that while this function is similar to the one above, they have different return types.
// The function above constructs options to pass to skv2, whereas this function creates retry-go
// options that Gloo Fed uses directly.
func GetClusterWatcherLocalRetryOptions(ctx context.Context) []retry.Option {
	var delayType retry.DelayTypeFunc
	if getRetryTypeFromEnv(ctx, ClusterWatcherLocalRetryTypeEnv, DefaultTypeLocal) == RetryTypeFixed {
		delayType = retry.FixedDelay
	} else {
		delayType = retry.BackOffDelay
	}

	delay := getDurationFromEnv(ctx, ClusterWatcherLocalRetryDelayEnv, DefaultDelayLocal)
	maxDelay := getDurationFromEnv(ctx, ClusterWatcherLocalRetryMaxDelayEnv, DefaultMaxDelayLocal)
	maxJitter := getDurationFromEnv(ctx, ClusterWatcherLocalRetryMaxJitterEnv, DefaultMaxJitterLocal)
	attempts := getUintFromEnv(ctx, ClusterWatcherLocalRetryAttemptsEnv, DefaultAttemptsLocal)

	// construct the retry options with the above values
	retryOptions := []retry.Option{
		retry.Delay(delay),
		retry.MaxDelay(maxDelay),
		retry.Attempts(attempts),
	}
	if maxJitter > 0 {
		// add a random delay with max jitter to the specified delay type
		retryOptions = append(retryOptions,
			retry.DelayType(retry.CombineDelay(delayType, retry.RandomDelay)),
			retry.MaxJitter(maxJitter))
	} else {
		// if maxJitter was explicitly set to 0, don't add randomness or jitter
		retryOptions = append(retryOptions, retry.DelayType(delayType))
	}
	return retryOptions
}

// Gets the retry type ("backoff" or "fixed") from the given env variable key.
// If the value is anything other than "backoff" or "fixed", returns the given default value.
func getRetryTypeFromEnv(ctx context.Context, key string, defaultValue string) string {
	retryType := os.Getenv(key)
	if retryType != "" {
		if retryType == RetryTypeBackoff || retryType == RetryTypeFixed {
			return retryType
		}
		contextutils.LoggerFrom(ctx).Warnf("invalid retry type %s found for environment variable %s. supported values are [backoff, fixed].", retryType, key)
	}
	return defaultValue
}

// Gets the duration specified by the given env variable key.
// If the env variable is not set, or the value can't be parsed, returns the given default value.
func getDurationFromEnv(ctx context.Context, key string, defaultValue time.Duration) time.Duration {
	if durationStr := os.Getenv(key); durationStr != "" {
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			contextutils.LoggerFrom(ctx).Warnf("invalid duration string %s found for environment variable %s: %v", durationStr, key, err)
			return defaultValue
		}
		return duration

	}
	return defaultValue
}

// Gets the uint specified by the given env variable key.
// If the env variable is not set, or the value can't be parsed, returns the given default value.
func getUintFromEnv(ctx context.Context, key string, defaultValue uint) uint {
	if numStr := os.Getenv(key); numStr != "" {
		num, err := strconv.ParseUint(numStr, 10, 32)
		if err != nil {
			contextutils.LoggerFrom(ctx).Warnf("invalid uint string %s found for environment variable %s: %v", numStr, key, err)
			return defaultValue
		}
		ptr := uint(num)
		return ptr
	}
	return defaultValue
}
