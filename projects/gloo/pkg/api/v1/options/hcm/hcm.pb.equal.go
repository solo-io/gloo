// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/hcm/hcm.proto

package hcm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	equality "github.com/solo-io/protoc-gen-ext/pkg/equality"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = bytes.Compare
	_ = strings.Compare
	_ = equality.Equalizer(nil)
	_ = proto.Message(nil)
)

// Equal function
func (m *HttpConnectionManagerSettings) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*HttpConnectionManagerSettings)
	if !ok {
		that2, ok := that.(HttpConnectionManagerSettings)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if h, ok := interface{}(m.GetSkipXffAppend()).(equality.Equalizer); ok {
		if !h.Equal(target.GetSkipXffAppend()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetSkipXffAppend(), target.GetSkipXffAppend()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetVia()).(equality.Equalizer); ok {
		if !h.Equal(target.GetVia()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetVia(), target.GetVia()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetXffNumTrustedHops()).(equality.Equalizer); ok {
		if !h.Equal(target.GetXffNumTrustedHops()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetXffNumTrustedHops(), target.GetXffNumTrustedHops()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetUseRemoteAddress()).(equality.Equalizer); ok {
		if !h.Equal(target.GetUseRemoteAddress()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetUseRemoteAddress(), target.GetUseRemoteAddress()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetGenerateRequestId()).(equality.Equalizer); ok {
		if !h.Equal(target.GetGenerateRequestId()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetGenerateRequestId(), target.GetGenerateRequestId()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetProxy_100Continue()).(equality.Equalizer); ok {
		if !h.Equal(target.GetProxy_100Continue()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetProxy_100Continue(), target.GetProxy_100Continue()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetStreamIdleTimeout()).(equality.Equalizer); ok {
		if !h.Equal(target.GetStreamIdleTimeout()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetStreamIdleTimeout(), target.GetStreamIdleTimeout()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetIdleTimeout()).(equality.Equalizer); ok {
		if !h.Equal(target.GetIdleTimeout()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetIdleTimeout(), target.GetIdleTimeout()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetMaxRequestHeadersKb()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMaxRequestHeadersKb()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMaxRequestHeadersKb(), target.GetMaxRequestHeadersKb()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetRequestTimeout()).(equality.Equalizer); ok {
		if !h.Equal(target.GetRequestTimeout()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetRequestTimeout(), target.GetRequestTimeout()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetRequestHeadersTimeout()).(equality.Equalizer); ok {
		if !h.Equal(target.GetRequestHeadersTimeout()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetRequestHeadersTimeout(), target.GetRequestHeadersTimeout()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetDrainTimeout()).(equality.Equalizer); ok {
		if !h.Equal(target.GetDrainTimeout()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetDrainTimeout(), target.GetDrainTimeout()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetDelayedCloseTimeout()).(equality.Equalizer); ok {
		if !h.Equal(target.GetDelayedCloseTimeout()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetDelayedCloseTimeout(), target.GetDelayedCloseTimeout()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetServerName()).(equality.Equalizer); ok {
		if !h.Equal(target.GetServerName()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetServerName(), target.GetServerName()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetStripAnyHostPort()).(equality.Equalizer); ok {
		if !h.Equal(target.GetStripAnyHostPort()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetStripAnyHostPort(), target.GetStripAnyHostPort()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetAcceptHttp_10()).(equality.Equalizer); ok {
		if !h.Equal(target.GetAcceptHttp_10()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetAcceptHttp_10(), target.GetAcceptHttp_10()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetDefaultHostForHttp_10()).(equality.Equalizer); ok {
		if !h.Equal(target.GetDefaultHostForHttp_10()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetDefaultHostForHttp_10(), target.GetDefaultHostForHttp_10()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetAllowChunkedLength()).(equality.Equalizer); ok {
		if !h.Equal(target.GetAllowChunkedLength()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetAllowChunkedLength(), target.GetAllowChunkedLength()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetEnableTrailers()).(equality.Equalizer); ok {
		if !h.Equal(target.GetEnableTrailers()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetEnableTrailers(), target.GetEnableTrailers()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetTracing()).(equality.Equalizer); ok {
		if !h.Equal(target.GetTracing()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetTracing(), target.GetTracing()) {
			return false
		}
	}

	if m.GetForwardClientCertDetails() != target.GetForwardClientCertDetails() {
		return false
	}

	if h, ok := interface{}(m.GetSetCurrentClientCertDetails()).(equality.Equalizer); ok {
		if !h.Equal(target.GetSetCurrentClientCertDetails()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetSetCurrentClientCertDetails(), target.GetSetCurrentClientCertDetails()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetPreserveExternalRequestId()).(equality.Equalizer); ok {
		if !h.Equal(target.GetPreserveExternalRequestId()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetPreserveExternalRequestId(), target.GetPreserveExternalRequestId()) {
			return false
		}
	}

	if len(m.GetUpgrades()) != len(target.GetUpgrades()) {
		return false
	}
	for idx, v := range m.GetUpgrades() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetUpgrades()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetUpgrades()[idx]) {
				return false
			}
		}

	}

	if h, ok := interface{}(m.GetMaxConnectionDuration()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMaxConnectionDuration()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMaxConnectionDuration(), target.GetMaxConnectionDuration()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetMaxStreamDuration()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMaxStreamDuration()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMaxStreamDuration(), target.GetMaxStreamDuration()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetMaxHeadersCount()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMaxHeadersCount()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMaxHeadersCount(), target.GetMaxHeadersCount()) {
			return false
		}
	}

	if m.GetHeadersWithUnderscoresAction() != target.GetHeadersWithUnderscoresAction() {
		return false
	}

	if h, ok := interface{}(m.GetMaxRequestsPerConnection()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMaxRequestsPerConnection()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMaxRequestsPerConnection(), target.GetMaxRequestsPerConnection()) {
			return false
		}
	}

	if m.GetServerHeaderTransformation() != target.GetServerHeaderTransformation() {
		return false
	}

	if m.GetPathWithEscapedSlashesAction() != target.GetPathWithEscapedSlashesAction() {
		return false
	}

	if m.GetCodecType() != target.GetCodecType() {
		return false
	}

	if h, ok := interface{}(m.GetMergeSlashes()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMergeSlashes()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMergeSlashes(), target.GetMergeSlashes()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetNormalizePath()).(equality.Equalizer); ok {
		if !h.Equal(target.GetNormalizePath()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetNormalizePath(), target.GetNormalizePath()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetUuidRequestIdConfig()).(equality.Equalizer); ok {
		if !h.Equal(target.GetUuidRequestIdConfig()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetUuidRequestIdConfig(), target.GetUuidRequestIdConfig()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetHttp2ProtocolOptions()).(equality.Equalizer); ok {
		if !h.Equal(target.GetHttp2ProtocolOptions()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetHttp2ProtocolOptions(), target.GetHttp2ProtocolOptions()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetInternalAddressConfig()).(equality.Equalizer); ok {
		if !h.Equal(target.GetInternalAddressConfig()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetInternalAddressConfig(), target.GetInternalAddressConfig()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetAppendXForwardedPort()).(equality.Equalizer); ok {
		if !h.Equal(target.GetAppendXForwardedPort()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetAppendXForwardedPort(), target.GetAppendXForwardedPort()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetEarlyHeaderManipulation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetEarlyHeaderManipulation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetEarlyHeaderManipulation(), target.GetEarlyHeaderManipulation()) {
			return false
		}
	}

	switch m.HeaderFormat.(type) {

	case *HttpConnectionManagerSettings_ProperCaseHeaderKeyFormat:
		if _, ok := target.HeaderFormat.(*HttpConnectionManagerSettings_ProperCaseHeaderKeyFormat); !ok {
			return false
		}

		if h, ok := interface{}(m.GetProperCaseHeaderKeyFormat()).(equality.Equalizer); ok {
			if !h.Equal(target.GetProperCaseHeaderKeyFormat()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetProperCaseHeaderKeyFormat(), target.GetProperCaseHeaderKeyFormat()) {
				return false
			}
		}

	case *HttpConnectionManagerSettings_PreserveCaseHeaderKeyFormat:
		if _, ok := target.HeaderFormat.(*HttpConnectionManagerSettings_PreserveCaseHeaderKeyFormat); !ok {
			return false
		}

		if h, ok := interface{}(m.GetPreserveCaseHeaderKeyFormat()).(equality.Equalizer); ok {
			if !h.Equal(target.GetPreserveCaseHeaderKeyFormat()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetPreserveCaseHeaderKeyFormat(), target.GetPreserveCaseHeaderKeyFormat()) {
				return false
			}
		}

	default:
		// m is nil but target is not nil
		if m.HeaderFormat != target.HeaderFormat {
			return false
		}
	}

	return true
}

// Equal function
func (m *HttpConnectionManagerSettings_SetCurrentClientCertDetails) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*HttpConnectionManagerSettings_SetCurrentClientCertDetails)
	if !ok {
		that2, ok := that.(HttpConnectionManagerSettings_SetCurrentClientCertDetails)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if h, ok := interface{}(m.GetSubject()).(equality.Equalizer); ok {
		if !h.Equal(target.GetSubject()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetSubject(), target.GetSubject()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetCert()).(equality.Equalizer); ok {
		if !h.Equal(target.GetCert()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetCert(), target.GetCert()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetChain()).(equality.Equalizer); ok {
		if !h.Equal(target.GetChain()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetChain(), target.GetChain()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetDns()).(equality.Equalizer); ok {
		if !h.Equal(target.GetDns()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetDns(), target.GetDns()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetUri()).(equality.Equalizer); ok {
		if !h.Equal(target.GetUri()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetUri(), target.GetUri()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *HttpConnectionManagerSettings_UuidRequestIdConfigSettings) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*HttpConnectionManagerSettings_UuidRequestIdConfigSettings)
	if !ok {
		that2, ok := that.(HttpConnectionManagerSettings_UuidRequestIdConfigSettings)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if h, ok := interface{}(m.GetPackTraceReason()).(equality.Equalizer); ok {
		if !h.Equal(target.GetPackTraceReason()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetPackTraceReason(), target.GetPackTraceReason()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetUseRequestIdForTraceSampling()).(equality.Equalizer); ok {
		if !h.Equal(target.GetUseRequestIdForTraceSampling()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetUseRequestIdForTraceSampling(), target.GetUseRequestIdForTraceSampling()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *HttpConnectionManagerSettings_CidrRange) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*HttpConnectionManagerSettings_CidrRange)
	if !ok {
		that2, ok := that.(HttpConnectionManagerSettings_CidrRange)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if strings.Compare(m.GetAddressPrefix(), target.GetAddressPrefix()) != 0 {
		return false
	}

	if h, ok := interface{}(m.GetPrefixLen()).(equality.Equalizer); ok {
		if !h.Equal(target.GetPrefixLen()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetPrefixLen(), target.GetPrefixLen()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *HttpConnectionManagerSettings_InternalAddressConfig) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*HttpConnectionManagerSettings_InternalAddressConfig)
	if !ok {
		that2, ok := that.(HttpConnectionManagerSettings_InternalAddressConfig)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if h, ok := interface{}(m.GetUnixSockets()).(equality.Equalizer); ok {
		if !h.Equal(target.GetUnixSockets()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetUnixSockets(), target.GetUnixSockets()) {
			return false
		}
	}

	if len(m.GetCidrRanges()) != len(target.GetCidrRanges()) {
		return false
	}
	for idx, v := range m.GetCidrRanges() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetCidrRanges()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetCidrRanges()[idx]) {
				return false
			}
		}

	}

	return true
}
