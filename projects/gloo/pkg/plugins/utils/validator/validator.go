package validator

import (
	"context"
	"encoding/json"
	"hash"
	"hash/fnv"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/bootstrap"
	envoyvalidation "github.com/solo-io/gloo/pkg/utils/envoyutils/validation"
	"github.com/solo-io/gloo/pkg/utils/statsutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"

	"github.com/solo-io/go-utils/contextutils"

	"go.opencensus.io/stats"
	"google.golang.org/protobuf/runtime/protoiface"
	"k8s.io/utils/lru"
)

// DefaultCacheSize defines the default size of the LRU cache used by the validator
const DefaultCacheSize int = 1024

// Validator validates an envoy config by running it by envoy in validate mode. This requires the envoy binary to be present at $ENVOY_BINARY (defaults to /usr/local/bin/envoy).
// Results are cached via an LRU cache for performance
type Validator interface {
	// ValidateConfig validates the given envoy config and returns any out and error from envoy. Returns nil if the envoy binary is not found.
	ValidateConfig(ctx context.Context, config HashableProtoMessage) error

	// ValidateSnapshot validates the given snapshot and returns any out and error from envoy. Returns nil if the envoy binary is not found.
	ValidateSnapshot(ctx context.Context, snap HashableSnapshot) error

	// CacheLength returns the returns the number of items in the cache
	CacheLength() int
}

var _ Validator = new(validator)

type validator struct {
	// filterName to be used if validating a specific filter config
	// e.g. for transformations or waf.
	filterName string
	// lruCache is a map of: (config hash) -> error state
	// this is usually a typed error but may be an untyped nil interface
	lruCache *lru.Cache
	// Counter to increment on cache hits
	cacheHits *stats.Int64Measure
	// Counter to increment on cache misses
	cacheMisses *stats.Int64Measure
	// Hasher to use for caching
	hasher hash.Hash64
}

// New returns a new Validator
func New(name string, opts ...Option) validator {
	cfg := processOptions(name, opts...)
	return validator{
		filterName:  cfg.filterName,
		lruCache:    lru.New(cfg.cacheSize),
		cacheHits:   cfg.cacheHits,
		cacheMisses: cfg.cacheMisses,
		hasher:      cfg.hasher,
	}
}

// HashableProtoMessage defines a ProtoMessage that can be hashed. Useful when passing different ProtoMessages objects that need to be hashed.
type HashableProtoMessage interface {
	protoiface.MessageV1
	Hash(hasher hash.Hash64) (uint64, error)
}

func (v validator) ValidateConfig(ctx context.Context, config HashableProtoMessage) error {
	// Always use the proto's generated hasher.
	hash, err := config.Hash(nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).DPanicf("error hashing the config, should never happen: %v", err)
		return err
	}

	// This proto has already been validated, return the result
	if err, ok := v.lruCache.Get(hash); ok {
		statsutils.MeasureOne(
			ctx,
			v.cacheHits,
		)
		// Error may be nil here since it's just the cached result
		// so return it as a nil err after cast worst case.
		errCasted, _ := err.(error)
		return errCasted
	}
	statsutils.MeasureOne(
		ctx,
		v.cacheMisses,
	)
	filterBootstrap, err := bootstrap.FromFilter(v.filterName, config)
	if err != nil {
		contextutils.LoggerFrom(ctx).DPanicf("error constructing valid bootstrap from the config, should never happen: %v", err)
		return err
	}

	err = envoyvalidation.ValidateBootstrap(ctx, filterBootstrap)
	v.lruCache.Add(hash, err)
	return err
}

func (v validator) ValidateSnapshot(ctx context.Context, snap HashableSnapshot) error {
	hash, err := snap.Hash(v.hasher)
	if err != nil {
		contextutils.LoggerFrom(ctx).DPanicf("error hashing the snapshot, should never happen: %v", err)
		return err
	}

	// This proto has already been validated, return the result
	if err, ok := v.lruCache.Get(hash); ok {
		statsutils.MeasureOne(
			ctx,
			v.cacheHits,
		)
		// Error may be nil here since it's just the cached result
		// so return it as a nil err after cast worst case.
		errCasted, _ := err.(error)
		return errCasted
	}
	statsutils.MeasureOne(
		ctx,
		v.cacheMisses,
	)

	err = envoyvalidation.ValidateSnapshot(ctx, snap)
	v.lruCache.Add(hash, err)
	return err
}

func (v validator) CacheLength() int {
	return v.lruCache.Len()
}

func NewHashableSnapshot(snap envoycache.Snapshot) *hashableSnapshot {
	return &hashableSnapshot{
		Snapshot: snap,
	}
}

// HashableSnapshot defines a snapshot that can be hashed.
type HashableSnapshot interface {
	envoycache.Snapshot
	Hash(hasher hash.Hash64) (uint64, error)
}

// hashableSnapshot implements HashableSnapshot.
type hashableSnapshot struct {
	envoycache.Snapshot
}

func (h *hashableSnapshot) Hash(hasher hash.Hash64) (uint64, error) {
	if hasher == nil {
		hasher = fnv.New64()
	}
	b, err := json.Marshal(h)
	if err != nil {
		return 0, err
	}

	_, err = hasher.Write(b)
	if err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}
