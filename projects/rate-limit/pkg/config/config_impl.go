package configproto

import (
	"errors"
	"fmt"

	"github.com/solo-io/rate-limiter/pkg/config"

	pb_struct "github.com/envoyproxy/go-control-plane/envoy/api/v2/ratelimit"

	pb_rls "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v2"
	glooee "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	solorl "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"

	"go.uber.org/zap"
)

type rateLimitDescriptor struct {
	descriptors map[string]*rateLimitDescriptor
	limit       *config.RateLimit
}

type rateLimitDomain struct {
	rateLimitDescriptor
}

type rateLimitConfig struct {
	domains map[string]*rateLimitDomain
	logger  *zap.SugaredLogger
}

// map of keys generated from the user config
type keyedRateLimits map[string]*activeRateLimit

// just the information needed to compare against the request and count
type activeRateLimit struct {
	requestsPerUnit uint32
	unit            solorl.RateLimit_Unit
	// unit            pb.RateLimitResponse_RateLimit_Unit
}

// Create a new rate limit config entry.
// @param requestsPerUnit supplies the requests per unit of time for the entry.
// @param unit supplies the unit of time for the entry.
// @param key supplies the fully resolved key name of the entry.
// @param scope supplies the owning scope.
// @return the new config entry.
func NewRateLimit(
	requestsPerUnit uint32, unit pb_rls.RateLimitResponse_RateLimit_Unit, key string) *config.RateLimit {

	return &config.RateLimit{Key: key, RequestsPerUnit: requestsPerUnit, Unit: unit}
}

// Dump an individual descriptor for debugging purposes.
func (this *rateLimitDescriptor) dump() string {
	ret := ""
	if this.limit != nil {
		ret += fmt.Sprintf(
			"%s: unit=%s requests_per_unit=%d\n", this.limit.Key,
			this.limit.Unit.String(), this.limit.RequestsPerUnit)
	}
	for _, descriptor := range this.descriptors {
		ret += descriptor.dump()
	}
	return ret
}

func (this rateLimitConfig) Dump() string {
	ret := ""
	for _, domain := range this.domains {
		ret += domain.dump()
	}

	return ret
}

func (this rateLimitConfig) GetLimit(
	domain string, descriptor *pb_struct.RateLimitDescriptor) *config.RateLimit {

	this.logger.Debugf("starting get limit lookup")
	var rateLimit *config.RateLimit = nil
	value := this.domains[domain]
	if value == nil {
		this.logger.Debugf("unknown domain '%s'", domain)
		return rateLimit
	}

	descriptorsMap := value.descriptors
	for i, entry := range descriptor.Entries {
		// First see if key_value is in the map. If that isn't in the map we look for just key
		// to check for a default value.
		finalKey := entry.Key + "_" + entry.Value
		this.logger.Debugf("looking up key: %s", finalKey)
		nextDescriptor := descriptorsMap[finalKey]
		if nextDescriptor == nil {
			finalKey = entry.Key
			this.logger.Debugf("looking up key: %s", finalKey)
			nextDescriptor = descriptorsMap[finalKey]
		}

		if nextDescriptor != nil && nextDescriptor.limit != nil {
			this.logger.Debugf("found rate limit: %s", finalKey)
			if i == len(descriptor.Entries)-1 {
				rateLimit = nextDescriptor.limit
			} else {
				this.logger.Debugf("request depth does not match config depth, there are more entries in the request's descriptor")
			}
		}

		if nextDescriptor != nil && len(nextDescriptor.descriptors) > 0 {
			this.logger.Debugf("iterating to next level")
			descriptorsMap = nextDescriptor.descriptors
		} else {
			break
		}
	}

	return rateLimit
}

func (this *rateLimitDescriptor) loadDescriptors(logger *zap.SugaredLogger, parentKey string, descriptors []*glooee.Constraint) error {

	for _, descriptorConfig := range descriptors {
		if descriptorConfig.Key == "" {
			return errors.New("descriptor has empty key")
		}

		// Value is optional, so the final key for the map is either the key only or key_value.
		finalKey := descriptorConfig.Key
		if descriptorConfig.Value != "" {
			finalKey += "_" + descriptorConfig.Value
		}

		newParentKey := parentKey + finalKey
		if _, present := this.descriptors[finalKey]; present {
			return fmt.Errorf("duplicate descriptor composite key '%s'", newParentKey)
		}

		var rateLimit *config.RateLimit = nil
		var rateLimitDebugString string = ""
		if descriptorConfig.RateLimit != nil {
			value, present :=
				pb_rls.RateLimitResponse_RateLimit_Unit_value[solorl.RateLimit_Unit_name[int32(descriptorConfig.RateLimit.Unit)]]
			if !present || value == int32(pb_rls.RateLimitResponse_RateLimit_UNKNOWN) {
				return fmt.Errorf("invalid rate limit unit '%s'", descriptorConfig.RateLimit.Unit)
			}

			rateLimit = NewRateLimit(
				descriptorConfig.RateLimit.RequestsPerUnit, pb_rls.RateLimitResponse_RateLimit_Unit(value), newParentKey)
			rateLimitDebugString = fmt.Sprintf(
				" ratelimit={requests_per_unit=%d, unit=%s}", rateLimit.RequestsPerUnit,
				rateLimit.Unit.String())
		}

		logger.Debugf(
			"loading descriptor: key=%s%s", newParentKey, rateLimitDebugString)
		newDescriptor := &rateLimitDescriptor{map[string]*rateLimitDescriptor{}, rateLimit}
		err := newDescriptor.loadDescriptors(logger, newParentKey+".", descriptorConfig.Constraints)
		if err != nil {
			return err
		}
		this.descriptors[finalKey] = newDescriptor
	}
	return nil
}

type rateLimitConfigGenerator struct {
	logger *zap.SugaredLogger
}

func NewConfigGenerator(logger *zap.SugaredLogger) RateLimitConfigGenerator {
	return &rateLimitConfigGenerator{
		logger: logger,
	}
}

func (this *rateLimitConfigGenerator) GenerateConfig(configs []*glooee.RateLimitConfig) (config.RateLimitConfig, error) {
	limits := rateLimitConfig{
		domains: make(map[string]*rateLimitDomain),
		logger:  this.logger,
	}
	for _, cfg := range configs {
		if cfg.Domain == "" {
			return nil, errors.New("config file cannot have empty domain")
		}

		if _, present := limits.domains[cfg.Domain]; present {
			return nil, fmt.Errorf("duplicate domain '%s' in config", cfg.Domain)
		}

		domaincfg, err := this.makeConfig(cfg)
		if err != nil {
			return nil, err
		}
		limits.domains[cfg.Domain] = domaincfg
	}
	return limits, nil

}

func (this *rateLimitConfigGenerator) makeConfig(rc *glooee.RateLimitConfig) (*rateLimitDomain, error) {

	newDomain := &rateLimitDomain{rateLimitDescriptor{map[string]*rateLimitDescriptor{}, nil}}
	newDomain.loadDescriptors(this.logger, rc.Domain+".", rc.Constraints)
	return newDomain, nil
}
