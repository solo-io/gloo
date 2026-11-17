package validator

import (
	"context"
	"errors"
	"fmt"
	"hash"
	"time"

	"github.com/solo-io/gloo/pkg/utils/statsutils"
	"github.com/solo-io/gloo/projects/envoyinit/pkg/runner"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/go-utils/contextutils"
	"go.opencensus.io/stats"
	"google.golang.org/protobuf/runtime/protoiface"
	"k8s.io/utils/lru"
)

// DefaultCacheSize defines the default size of the LRU cache used by the validator
const DefaultCacheSize int = 1024

// crashResultTTL is how long a crashed validation's result is reused. One config is validated once
// for each route and virtual host that uses it, and a crash is retried like any interruption, so
// reusing the crash keeps one config from starting envoy over and over. A crash can come from the
// pod's environment rather than the config, so after this long the config is validated again.
const crashResultTTL = time.Minute

// crashResult is the cache entry for a validation that crashed.
type crashResult struct {
	err       error
	crashedAt time.Time
}

// Validator validates an envoy config by running it by envoy in validate mode. This requires the envoy binary to be present at $ENVOY_BINARY (defaults to /usr/local/bin/envoy).
// Results are cached via an LRU cache for performance
type Validator interface {
	// ValidateConfig validates the given envoy config and returns any out and error from envoy. Returns nil if the envoy binary is not found.
	ValidateConfig(ctx context.Context, config HashableProtoMessage) error

	// CacheLength returns the returns the number of items in the cache
	CacheLength() int
}

type validator struct {
	filterName string
	// lruCache is a map of: (config hash) -> error state
	// this is usually a typed error but may be an untyped nil interface, or a crashResult
	lruCache *lru.Cache
	// Counter to increment on cache hits
	cacheHits *stats.Int64Measure
	// Counter to increment on cache misses
	cacheMisses *stats.Int64Measure
}

// New returns a new Validator
func New(name string, filterName string, opts ...Option) validator {
	cfg := processOptions(name, opts...)
	return validator{
		filterName:  filterName,
		lruCache:    lru.New(cfg.cacheSize),
		cacheHits:   cfg.cacheHits,
		cacheMisses: cfg.cacheMisses,
	}
}

// HashableProtoMessage defines a ProtoMessage that can be hashed. Useful when passing different ProtoMessages objects that need to be hashed.
type HashableProtoMessage interface {
	protoiface.MessageV1
	Hash(hasher hash.Hash64) (uint64, error)
}

func (v validator) ValidateConfig(ctx context.Context, config HashableProtoMessage) error {
	hash, err := config.Hash(nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).DPanicf("error hashing the config, should never happen: %v", err)
		return err
	}

	// This proto has already been validated, return the result. A crash result is used only until
	// it expires.
	if cached, ok := v.lruCache.Get(hash); ok {
		crash, isCrash := cached.(crashResult)
		if !isCrash || time.Since(crash.crashedAt) < crashResultTTL {
			statsutils.MeasureOne(
				ctx,
				v.cacheHits,
			)
			if isCrash {
				return crash.err
			}
			// Error may be nil here since it's just the cached result
			// so return it as a nil err after cast worst case.
			errCasted, _ := cached.(error)
			return errCasted
		}
	}
	statsutils.MeasureOne(
		ctx,
		v.cacheMisses,
	)

	err = bootstrap.ValidateBootstrap(ctx, v.filterName, config)

	// Check for a crash before checking for an interruption: a crash also matches
	// ErrValidationInterrupted, and is handled differently. The crash error is cached for
	// crashResultTTL. Until then, every validation of this
	// config, retries included, returns the cached error without starting envoy. After that, the
	// next validation starts envoy again.
	if errors.Is(err, runner.ErrValidationCrashed) {
		err = fmt.Errorf("validating %s config: %w", v.filterName, err)
		v.lruCache.Add(hash, crashResult{err: err, crashedAt: time.Now()})
		return err
	}

	// Any other interruption says nothing about the config and is not cached.
	if errors.Is(err, runner.ErrValidationInterrupted) {
		return fmt.Errorf("envoy validation of %s config was interrupted: %w", v.filterName, err)
	}

	v.lruCache.Add(hash, err)
	return err
}

func (v validator) CacheLength() int {
	return v.lruCache.Len()
}
