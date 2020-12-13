package internal

import (
	"github.com/golang/protobuf/proto"
	rate_limiter "github.com/solo-io/rate-limiter/pkg/api/ratelimit.solo.io/v1alpha1"
	rate_limiter_types "github.com/solo-io/rate-limiter/pkg/api/ratelimit.solo.io/v1alpha1/types"
	solo_apis_rl "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
)

func ToRateLimiterResource(in *solo_apis_rl.RateLimitConfig) (*rate_limiter.RateLimitConfig, error) {
	if in == nil {
		return nil, nil
	}

	rlSpec, err := toRateLimiterResourceSpec(&in.Spec)
	if err != nil {
		return nil, err
	}
	rlStatus, err := toRateLimiterResourceStatus(&in.Status)
	if err != nil {
		return nil, err
	}

	out := &rate_limiter.RateLimitConfig{}
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = *in.ObjectMeta.DeepCopy()
	out.Spec = *rlSpec
	out.Status = *rlStatus

	return out, nil
}

func ToRateLimiterResourceSpec_Raw(in *solo_apis_rl.RateLimitConfigSpec_Raw) (*rate_limiter_types.RateLimitConfigSpec_Raw, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &rate_limiter_types.RateLimitConfigSpec_Raw{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func ToSoloAPIsResourceSpec_Raw(in *rate_limiter_types.RateLimitConfigSpec_Raw) (*solo_apis_rl.RateLimitConfigSpec_Raw, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &solo_apis_rl.RateLimitConfigSpec_Raw{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func ToRateLimiterSetDescriptor(in *solo_apis_rl.SetDescriptor) (*rate_limiter_types.SetDescriptor, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &rate_limiter_types.SetDescriptor{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func ToRateLimiterSetDescriptors(in []*solo_apis_rl.SetDescriptor) ([]*rate_limiter_types.SetDescriptor, error) {
	if in == nil {
		return nil, nil
	}

	out := make([]*rate_limiter_types.SetDescriptor, len(in))
	for i, descriptor := range in {
		convertedDescriptor, err := ToRateLimiterSetDescriptor(descriptor)
		if err != nil {
			return nil, err
		}
		out[i] = convertedDescriptor
	}
	return out, nil
}

func ToRateLimiterDescriptor(in *solo_apis_rl.Descriptor) (*rate_limiter_types.Descriptor, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &rate_limiter_types.Descriptor{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func ToRateLimiterDescriptors(in []*solo_apis_rl.Descriptor) ([]*rate_limiter_types.Descriptor, error) {
	if in == nil {
		return nil, nil
	}

	out := make([]*rate_limiter_types.Descriptor, len(in))
	for i, descriptor := range in {
		convertedDescriptor, err := ToRateLimiterDescriptor(descriptor)
		if err != nil {
			return nil, err
		}
		out[i] = convertedDescriptor
	}
	return out, nil
}

func ToSoloAPIsSetDescriptor(in *rate_limiter_types.SetDescriptor) (*solo_apis_rl.SetDescriptor, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &solo_apis_rl.SetDescriptor{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func ToSoloAPIsSetDescriptors(in []*rate_limiter_types.SetDescriptor) ([]*solo_apis_rl.SetDescriptor, error) {
	if in == nil {
		return nil, nil
	}

	out := make([]*solo_apis_rl.SetDescriptor, len(in))
	for i, descriptor := range in {
		convertedDescriptor, err := ToSoloAPIsSetDescriptor(descriptor)
		if err != nil {
			return nil, err
		}
		out[i] = convertedDescriptor
	}
	return out, nil
}

func ToSoloAPIsDescriptor(in *rate_limiter_types.Descriptor) (*solo_apis_rl.Descriptor, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &solo_apis_rl.Descriptor{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func ToRateLimiterActionsSlice(in []*solo_apis_rl.RateLimitActions) ([]*rate_limiter_types.RateLimitActions, error) {
	if in == nil {
		return nil, nil
	}

	out := make([]*rate_limiter_types.RateLimitActions, len(in))
	for i, actionsElement := range in {
		converted, err := ToRateLimiterActions(actionsElement)
		if err != nil {
			return nil, err
		}
		out[i] = converted
	}
	return out, nil
}

func ToRateLimiterActions(in *solo_apis_rl.RateLimitActions) (*rate_limiter_types.RateLimitActions, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &rate_limiter_types.RateLimitActions{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func ToSoloAPIsActionsSlice(in []*rate_limiter_types.RateLimitActions) ([]*solo_apis_rl.RateLimitActions, error) {
	if in == nil {
		return nil, nil
	}

	out := make([]*solo_apis_rl.RateLimitActions, len(in))
	for i, actionsElement := range in {
		converted, err := toSoloAPIsActions(actionsElement)
		if err != nil {
			return nil, err
		}
		out[i] = converted
	}
	return out, nil
}

func toSoloAPIsActions(in *rate_limiter_types.RateLimitActions) (*solo_apis_rl.RateLimitActions, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &solo_apis_rl.RateLimitActions{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func toRateLimiterResourceSpec(in *solo_apis_rl.RateLimitConfigSpec) (*rate_limiter_types.RateLimitConfigSpec, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &rate_limiter_types.RateLimitConfigSpec{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func toRateLimiterResourceStatus(in *solo_apis_rl.RateLimitConfigStatus) (*rate_limiter_types.RateLimitConfigStatus, error) {
	if in == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &rate_limiter_types.RateLimitConfigStatus{}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}
