// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/proxy_protocol.proto

package v3

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
func (m *ProxyProtocolPassThroughTLVs) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*ProxyProtocolPassThroughTLVs)
	if !ok {
		that2, ok := that.(ProxyProtocolPassThroughTLVs)
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

	if m.GetMatchType() != target.GetMatchType() {
		return false
	}

	if len(m.GetTlvType()) != len(target.GetTlvType()) {
		return false
	}
	for idx, v := range m.GetTlvType() {

		if v != target.GetTlvType()[idx] {
			return false
		}

	}

	return true
}

// Equal function
func (m *ProxyProtocolConfig) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*ProxyProtocolConfig)
	if !ok {
		that2, ok := that.(ProxyProtocolConfig)
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

	if m.GetVersion() != target.GetVersion() {
		return false
	}

	if h, ok := interface{}(m.GetPassThroughTlvs()).(equality.Equalizer); ok {
		if !h.Equal(target.GetPassThroughTlvs()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetPassThroughTlvs(), target.GetPassThroughTlvs()) {
			return false
		}
	}

	return true
}