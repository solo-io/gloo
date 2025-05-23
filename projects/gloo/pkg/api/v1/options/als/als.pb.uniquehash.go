// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/als/als.proto

package als

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"strconv"

	safe_hasher "github.com/solo-io/protoc-gen-ext/pkg/hasher"
	"github.com/solo-io/protoc-gen-ext/pkg/hasher/hashstructure"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = new(hash.Hash64)
	_ = fnv.New64
	_ = strconv.Itoa
	_ = hashstructure.Hash
	_ = new(safe_hasher.SafeHasher)
)

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *AccessLoggingService) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.AccessLoggingService")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("AccessLog")); err != nil {
		return 0, err
	}
	for i, v := range m.GetAccessLog() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if h, ok := interface{}(v).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("v")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(v, nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("v")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *AccessLog) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.AccessLog")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetFilter()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Filter")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetFilter(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Filter")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	switch m.OutputDestination.(type) {

	case *AccessLog_FileSink:

		if h, ok := interface{}(m.GetFileSink()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("FileSink")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetFileSink(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("FileSink")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLog_GrpcService:

		if h, ok := interface{}(m.GetGrpcService()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("GrpcService")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetGrpcService(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("GrpcService")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLog_OpenTelemetryService:

		if h, ok := interface{}(m.GetOpenTelemetryService()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("OpenTelemetryService")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetOpenTelemetryService(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("OpenTelemetryService")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *FileSink) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.FileSink")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Path")); err != nil {
		return 0, err
	}
	if _, err = hasher.Write([]byte(m.GetPath())); err != nil {
		return 0, err
	}

	switch m.OutputFormat.(type) {

	case *FileSink_StringFormat:

		if _, err = hasher.Write([]byte("StringFormat")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(m.GetStringFormat())); err != nil {
			return 0, err
		}

	case *FileSink_JsonFormat:

		if h, ok := interface{}(m.GetJsonFormat()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("JsonFormat")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetJsonFormat(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("JsonFormat")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *GrpcService) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.GrpcService")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("LogName")); err != nil {
		return 0, err
	}
	if _, err = hasher.Write([]byte(m.GetLogName())); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("AdditionalRequestHeadersToLog")); err != nil {
		return 0, err
	}
	for i, v := range m.GetAdditionalRequestHeadersToLog() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if _, err = hasher.Write([]byte("v")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	if _, err = hasher.Write([]byte("AdditionalResponseHeadersToLog")); err != nil {
		return 0, err
	}
	for i, v := range m.GetAdditionalResponseHeadersToLog() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if _, err = hasher.Write([]byte("v")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	if _, err = hasher.Write([]byte("AdditionalResponseTrailersToLog")); err != nil {
		return 0, err
	}
	for i, v := range m.GetAdditionalResponseTrailersToLog() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if _, err = hasher.Write([]byte("v")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	if _, err = hasher.Write([]byte("FilterStateObjectsToLog")); err != nil {
		return 0, err
	}
	for i, v := range m.GetFilterStateObjectsToLog() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if _, err = hasher.Write([]byte("v")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	switch m.ServiceRef.(type) {

	case *GrpcService_StaticClusterName:

		if _, err = hasher.Write([]byte("StaticClusterName")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(m.GetStaticClusterName())); err != nil {
			return 0, err
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *OpenTelemetryGrpcCollector) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.OpenTelemetryGrpcCollector")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Endpoint")); err != nil {
		return 0, err
	}
	if _, err = hasher.Write([]byte(m.GetEndpoint())); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Authority")); err != nil {
		return 0, err
	}
	if _, err = hasher.Write([]byte(m.GetAuthority())); err != nil {
		return 0, err
	}

	{
		var result uint64
		innerHash := fnv.New64()
		for k, v := range m.GetHeaders() {
			innerHash.Reset()

			if _, err = innerHash.Write([]byte("v")); err != nil {
				return 0, err
			}
			if _, err = innerHash.Write([]byte(v)); err != nil {
				return 0, err
			}

			if _, err = innerHash.Write([]byte("k")); err != nil {
				return 0, err
			}
			if _, err = innerHash.Write([]byte(k)); err != nil {
				return 0, err
			}

			result = result ^ innerHash.Sum64()
		}
		err = binary.Write(hasher, binary.LittleEndian, result)
		if err != nil {
			return 0, err
		}

	}

	if _, err = hasher.Write([]byte("Insecure")); err != nil {
		return 0, err
	}
	err = binary.Write(hasher, binary.LittleEndian, m.GetInsecure())
	if err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetSslConfig()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("SslConfig")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetSslConfig(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("SslConfig")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if h, ok := interface{}(m.GetTimeout()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Timeout")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetTimeout(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Timeout")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *OpenTelemetryService) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.OpenTelemetryService")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("LogName")); err != nil {
		return 0, err
	}
	if _, err = hasher.Write([]byte(m.GetLogName())); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("FilterStateObjectsToLog")); err != nil {
		return 0, err
	}
	for i, v := range m.GetFilterStateObjectsToLog() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if _, err = hasher.Write([]byte("v")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	if _, err = hasher.Write([]byte("DisableBuiltinLabels")); err != nil {
		return 0, err
	}
	err = binary.Write(hasher, binary.LittleEndian, m.GetDisableBuiltinLabels())
	if err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetBody()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Body")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetBody(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Body")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if h, ok := interface{}(m.GetAttributes()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Attributes")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetAttributes(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Attributes")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	switch m.Destination.(type) {

	case *OpenTelemetryService_Collector:

		if h, ok := interface{}(m.GetCollector()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("Collector")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetCollector(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("Collector")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *AccessLogFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.AccessLogFilter")); err != nil {
		return 0, err
	}

	switch m.FilterSpecifier.(type) {

	case *AccessLogFilter_StatusCodeFilter:

		if h, ok := interface{}(m.GetStatusCodeFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("StatusCodeFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetStatusCodeFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("StatusCodeFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_DurationFilter:

		if h, ok := interface{}(m.GetDurationFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("DurationFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetDurationFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("DurationFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_NotHealthCheckFilter:

		if h, ok := interface{}(m.GetNotHealthCheckFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("NotHealthCheckFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetNotHealthCheckFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("NotHealthCheckFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_TraceableFilter:

		if h, ok := interface{}(m.GetTraceableFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("TraceableFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetTraceableFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("TraceableFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_RuntimeFilter:

		if h, ok := interface{}(m.GetRuntimeFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("RuntimeFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetRuntimeFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("RuntimeFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_AndFilter:

		if h, ok := interface{}(m.GetAndFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("AndFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetAndFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("AndFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_OrFilter:

		if h, ok := interface{}(m.GetOrFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("OrFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetOrFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("OrFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_HeaderFilter:

		if h, ok := interface{}(m.GetHeaderFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("HeaderFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetHeaderFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("HeaderFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_ResponseFlagFilter:

		if h, ok := interface{}(m.GetResponseFlagFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("ResponseFlagFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetResponseFlagFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("ResponseFlagFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	case *AccessLogFilter_GrpcStatusFilter:

		if h, ok := interface{}(m.GetGrpcStatusFilter()).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("GrpcStatusFilter")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(m.GetGrpcStatusFilter(), nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("GrpcStatusFilter")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *ComparisonFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.ComparisonFilter")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Op")); err != nil {
		return 0, err
	}
	err = binary.Write(hasher, binary.LittleEndian, m.GetOp())
	if err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetValue()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Value")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetValue(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Value")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *StatusCodeFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.StatusCodeFilter")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetComparison()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Comparison")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetComparison(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Comparison")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *DurationFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.DurationFilter")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetComparison()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Comparison")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetComparison(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Comparison")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *NotHealthCheckFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.NotHealthCheckFilter")); err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *TraceableFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.TraceableFilter")); err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *RuntimeFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.RuntimeFilter")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("RuntimeKey")); err != nil {
		return 0, err
	}
	if _, err = hasher.Write([]byte(m.GetRuntimeKey())); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetPercentSampled()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("PercentSampled")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetPercentSampled(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("PercentSampled")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if _, err = hasher.Write([]byte("UseIndependentRandomness")); err != nil {
		return 0, err
	}
	err = binary.Write(hasher, binary.LittleEndian, m.GetUseIndependentRandomness())
	if err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *AndFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.AndFilter")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Filters")); err != nil {
		return 0, err
	}
	for i, v := range m.GetFilters() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if h, ok := interface{}(v).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("v")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(v, nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("v")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *OrFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.OrFilter")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Filters")); err != nil {
		return 0, err
	}
	for i, v := range m.GetFilters() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if h, ok := interface{}(v).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("v")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(v, nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("v")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *HeaderFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.HeaderFilter")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetHeader()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Header")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetHeader(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Header")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *ResponseFlagFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.ResponseFlagFilter")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Flags")); err != nil {
		return 0, err
	}
	for i, v := range m.GetFlags() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if _, err = hasher.Write([]byte("v")); err != nil {
			return 0, err
		}
		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	return hasher.Sum64(), nil
}

// HashUnique function generates a hash of the object that is unique to the object by
// hashing field name and value pairs.
// Replaces Hash due to original hashing implemention only using field values. The omission
// of the field name in the hash calculation can lead to hash collisions.
func (m *GrpcStatusFilter) HashUnique(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("als.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als.GrpcStatusFilter")); err != nil {
		return 0, err
	}

	if _, err = hasher.Write([]byte("Statuses")); err != nil {
		return 0, err
	}
	for i, v := range m.GetStatuses() {
		if _, err = hasher.Write([]byte(strconv.Itoa(i))); err != nil {
			return 0, err
		}

		if _, err = hasher.Write([]byte("v")); err != nil {
			return 0, err
		}
		err = binary.Write(hasher, binary.LittleEndian, v)
		if err != nil {
			return 0, err
		}

	}

	if _, err = hasher.Write([]byte("Exclude")); err != nil {
		return 0, err
	}
	err = binary.Write(hasher, binary.LittleEndian, m.GetExclude())
	if err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}
