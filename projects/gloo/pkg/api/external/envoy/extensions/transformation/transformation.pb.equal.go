// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation.proto

package transformation

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
func (m *FilterTransformations) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*FilterTransformations)
	if !ok {
		that2, ok := that.(FilterTransformations)
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

	if len(m.GetTransformations()) != len(target.GetTransformations()) {
		return false
	}
	for idx, v := range m.GetTransformations() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetTransformations()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetTransformations()[idx]) {
				return false
			}
		}

	}

	if m.GetStage() != target.GetStage() {
		return false
	}

	if m.GetLogRequestResponseInfo() != target.GetLogRequestResponseInfo() {
		return false
	}

	return true
}

// Equal function
func (m *TransformationRule) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*TransformationRule)
	if !ok {
		that2, ok := that.(TransformationRule)
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

	if h, ok := interface{}(m.GetMatch()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMatch()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMatch(), target.GetMatch()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetRouteTransformations()).(equality.Equalizer); ok {
		if !h.Equal(target.GetRouteTransformations()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetRouteTransformations(), target.GetRouteTransformations()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *RouteTransformations) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*RouteTransformations)
	if !ok {
		that2, ok := that.(RouteTransformations)
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

	if h, ok := interface{}(m.GetRequestTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetRequestTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetRequestTransformation(), target.GetRequestTransformation()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetResponseTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetResponseTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetResponseTransformation(), target.GetResponseTransformation()) {
			return false
		}
	}

	if m.GetClearRouteCache() != target.GetClearRouteCache() {
		return false
	}

	if len(m.GetTransformations()) != len(target.GetTransformations()) {
		return false
	}
	for idx, v := range m.GetTransformations() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetTransformations()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetTransformations()[idx]) {
				return false
			}
		}

	}

	return true
}

// Equal function
func (m *ResponseMatcher) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*ResponseMatcher)
	if !ok {
		that2, ok := that.(ResponseMatcher)
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

	if len(m.GetHeaders()) != len(target.GetHeaders()) {
		return false
	}
	for idx, v := range m.GetHeaders() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetHeaders()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetHeaders()[idx]) {
				return false
			}
		}

	}

	if h, ok := interface{}(m.GetResponseCodeDetails()).(equality.Equalizer); ok {
		if !h.Equal(target.GetResponseCodeDetails()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetResponseCodeDetails(), target.GetResponseCodeDetails()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *ResponseTransformationRule) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*ResponseTransformationRule)
	if !ok {
		that2, ok := that.(ResponseTransformationRule)
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

	if h, ok := interface{}(m.GetMatch()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMatch()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMatch(), target.GetMatch()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetResponseTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetResponseTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetResponseTransformation(), target.GetResponseTransformation()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *Transformation) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*Transformation)
	if !ok {
		that2, ok := that.(Transformation)
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

	if h, ok := interface{}(m.GetLogRequestResponseInfo()).(equality.Equalizer); ok {
		if !h.Equal(target.GetLogRequestResponseInfo()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetLogRequestResponseInfo(), target.GetLogRequestResponseInfo()) {
			return false
		}
	}

	switch m.TransformationType.(type) {

	case *Transformation_TransformationTemplate:
		if _, ok := target.TransformationType.(*Transformation_TransformationTemplate); !ok {
			return false
		}

		if h, ok := interface{}(m.GetTransformationTemplate()).(equality.Equalizer); ok {
			if !h.Equal(target.GetTransformationTemplate()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetTransformationTemplate(), target.GetTransformationTemplate()) {
				return false
			}
		}

	case *Transformation_HeaderBodyTransform:
		if _, ok := target.TransformationType.(*Transformation_HeaderBodyTransform); !ok {
			return false
		}

		if h, ok := interface{}(m.GetHeaderBodyTransform()).(equality.Equalizer); ok {
			if !h.Equal(target.GetHeaderBodyTransform()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetHeaderBodyTransform(), target.GetHeaderBodyTransform()) {
				return false
			}
		}

	case *Transformation_TransformerConfig:
		if _, ok := target.TransformationType.(*Transformation_TransformerConfig); !ok {
			return false
		}

		if h, ok := interface{}(m.GetTransformerConfig()).(equality.Equalizer); ok {
			if !h.Equal(target.GetTransformerConfig()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetTransformerConfig(), target.GetTransformerConfig()) {
				return false
			}
		}

	case *Transformation_AiTransformation:
		if _, ok := target.TransformationType.(*Transformation_AiTransformation); !ok {
			return false
		}

		if h, ok := interface{}(m.GetAiTransformation()).(equality.Equalizer); ok {
			if !h.Equal(target.GetAiTransformation()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetAiTransformation(), target.GetAiTransformation()) {
				return false
			}
		}

	default:
		// m is nil but target is not nil
		if m.TransformationType != target.TransformationType {
			return false
		}
	}

	return true
}

// Equal function
func (m *Extraction) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*Extraction)
	if !ok {
		that2, ok := that.(Extraction)
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

	if strings.Compare(m.GetRegex(), target.GetRegex()) != 0 {
		return false
	}

	if m.GetSubgroup() != target.GetSubgroup() {
		return false
	}

	if h, ok := interface{}(m.GetReplacementText()).(equality.Equalizer); ok {
		if !h.Equal(target.GetReplacementText()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetReplacementText(), target.GetReplacementText()) {
			return false
		}
	}

	if m.GetMode() != target.GetMode() {
		return false
	}

	switch m.Source.(type) {

	case *Extraction_Header:
		if _, ok := target.Source.(*Extraction_Header); !ok {
			return false
		}

		if strings.Compare(m.GetHeader(), target.GetHeader()) != 0 {
			return false
		}

	case *Extraction_Body:
		if _, ok := target.Source.(*Extraction_Body); !ok {
			return false
		}

		if h, ok := interface{}(m.GetBody()).(equality.Equalizer); ok {
			if !h.Equal(target.GetBody()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetBody(), target.GetBody()) {
				return false
			}
		}

	default:
		// m is nil but target is not nil
		if m.Source != target.Source {
			return false
		}
	}

	return true
}

// Equal function
func (m *TransformationTemplate) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*TransformationTemplate)
	if !ok {
		that2, ok := that.(TransformationTemplate)
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

	if m.GetAdvancedTemplates() != target.GetAdvancedTemplates() {
		return false
	}

	if len(m.GetExtractors()) != len(target.GetExtractors()) {
		return false
	}
	for k, v := range m.GetExtractors() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetExtractors()[k]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetExtractors()[k]) {
				return false
			}
		}

	}

	if len(m.GetHeaders()) != len(target.GetHeaders()) {
		return false
	}
	for k, v := range m.GetHeaders() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetHeaders()[k]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetHeaders()[k]) {
				return false
			}
		}

	}

	if len(m.GetHeadersToAppend()) != len(target.GetHeadersToAppend()) {
		return false
	}
	for idx, v := range m.GetHeadersToAppend() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetHeadersToAppend()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetHeadersToAppend()[idx]) {
				return false
			}
		}

	}

	if len(m.GetHeadersToRemove()) != len(target.GetHeadersToRemove()) {
		return false
	}
	for idx, v := range m.GetHeadersToRemove() {

		if strings.Compare(v, target.GetHeadersToRemove()[idx]) != 0 {
			return false
		}

	}

	if m.GetParseBodyBehavior() != target.GetParseBodyBehavior() {
		return false
	}

	if m.GetIgnoreErrorOnParse() != target.GetIgnoreErrorOnParse() {
		return false
	}

	if len(m.GetDynamicMetadataValues()) != len(target.GetDynamicMetadataValues()) {
		return false
	}
	for idx, v := range m.GetDynamicMetadataValues() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetDynamicMetadataValues()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetDynamicMetadataValues()[idx]) {
				return false
			}
		}

	}

	if m.GetEscapeCharacters() != target.GetEscapeCharacters() {
		return false
	}

	if h, ok := interface{}(m.GetSpanTransformer()).(equality.Equalizer); ok {
		if !h.Equal(target.GetSpanTransformer()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetSpanTransformer(), target.GetSpanTransformer()) {
			return false
		}
	}

	switch m.BodyTransformation.(type) {

	case *TransformationTemplate_Body:
		if _, ok := target.BodyTransformation.(*TransformationTemplate_Body); !ok {
			return false
		}

		if h, ok := interface{}(m.GetBody()).(equality.Equalizer); ok {
			if !h.Equal(target.GetBody()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetBody(), target.GetBody()) {
				return false
			}
		}

	case *TransformationTemplate_Passthrough:
		if _, ok := target.BodyTransformation.(*TransformationTemplate_Passthrough); !ok {
			return false
		}

		if h, ok := interface{}(m.GetPassthrough()).(equality.Equalizer); ok {
			if !h.Equal(target.GetPassthrough()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetPassthrough(), target.GetPassthrough()) {
				return false
			}
		}

	case *TransformationTemplate_MergeExtractorsToBody:
		if _, ok := target.BodyTransformation.(*TransformationTemplate_MergeExtractorsToBody); !ok {
			return false
		}

		if h, ok := interface{}(m.GetMergeExtractorsToBody()).(equality.Equalizer); ok {
			if !h.Equal(target.GetMergeExtractorsToBody()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetMergeExtractorsToBody(), target.GetMergeExtractorsToBody()) {
				return false
			}
		}

	case *TransformationTemplate_MergeJsonKeys:
		if _, ok := target.BodyTransformation.(*TransformationTemplate_MergeJsonKeys); !ok {
			return false
		}

		if h, ok := interface{}(m.GetMergeJsonKeys()).(equality.Equalizer); ok {
			if !h.Equal(target.GetMergeJsonKeys()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetMergeJsonKeys(), target.GetMergeJsonKeys()) {
				return false
			}
		}

	default:
		// m is nil but target is not nil
		if m.BodyTransformation != target.BodyTransformation {
			return false
		}
	}

	return true
}

// Equal function
func (m *InjaTemplate) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*InjaTemplate)
	if !ok {
		that2, ok := that.(InjaTemplate)
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

	if strings.Compare(m.GetText(), target.GetText()) != 0 {
		return false
	}

	return true
}

// Equal function
func (m *Passthrough) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*Passthrough)
	if !ok {
		that2, ok := that.(Passthrough)
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

	return true
}

// Equal function
func (m *MergeExtractorsToBody) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*MergeExtractorsToBody)
	if !ok {
		that2, ok := that.(MergeExtractorsToBody)
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

	return true
}

// Equal function
func (m *MergeJsonKeys) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*MergeJsonKeys)
	if !ok {
		that2, ok := that.(MergeJsonKeys)
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

	if len(m.GetJsonKeys()) != len(target.GetJsonKeys()) {
		return false
	}
	for k, v := range m.GetJsonKeys() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetJsonKeys()[k]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetJsonKeys()[k]) {
				return false
			}
		}

	}

	return true
}

// Equal function
func (m *HeaderBodyTransform) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*HeaderBodyTransform)
	if !ok {
		that2, ok := that.(HeaderBodyTransform)
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

	if m.GetAddRequestMetadata() != target.GetAddRequestMetadata() {
		return false
	}

	return true
}

// Equal function
func (m *FieldDefault) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*FieldDefault)
	if !ok {
		that2, ok := that.(FieldDefault)
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

	if strings.Compare(m.GetField(), target.GetField()) != 0 {
		return false
	}

	if h, ok := interface{}(m.GetValue()).(equality.Equalizer); ok {
		if !h.Equal(target.GetValue()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetValue(), target.GetValue()) {
			return false
		}
	}

	if m.GetOverride() != target.GetOverride() {
		return false
	}

	return true
}

// Equal function
func (m *PromptEnrichment) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*PromptEnrichment)
	if !ok {
		that2, ok := that.(PromptEnrichment)
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

	if len(m.GetPrepend()) != len(target.GetPrepend()) {
		return false
	}
	for idx, v := range m.GetPrepend() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetPrepend()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetPrepend()[idx]) {
				return false
			}
		}

	}

	if len(m.GetAppend()) != len(target.GetAppend()) {
		return false
	}
	for idx, v := range m.GetAppend() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetAppend()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetAppend()[idx]) {
				return false
			}
		}

	}

	return true
}

// Equal function
func (m *AiTransformation) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*AiTransformation)
	if !ok {
		that2, ok := that.(AiTransformation)
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

	if m.GetEnableChatStreaming() != target.GetEnableChatStreaming() {
		return false
	}

	if len(m.GetFieldDefaults()) != len(target.GetFieldDefaults()) {
		return false
	}
	for idx, v := range m.GetFieldDefaults() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetFieldDefaults()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetFieldDefaults()[idx]) {
				return false
			}
		}

	}

	if h, ok := interface{}(m.GetPromptEnrichment()).(equality.Equalizer); ok {
		if !h.Equal(target.GetPromptEnrichment()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetPromptEnrichment(), target.GetPromptEnrichment()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *TransformationRule_Transformations) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*TransformationRule_Transformations)
	if !ok {
		that2, ok := that.(TransformationRule_Transformations)
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

	if h, ok := interface{}(m.GetRequestTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetRequestTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetRequestTransformation(), target.GetRequestTransformation()) {
			return false
		}
	}

	if m.GetClearRouteCache() != target.GetClearRouteCache() {
		return false
	}

	if h, ok := interface{}(m.GetResponseTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetResponseTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetResponseTransformation(), target.GetResponseTransformation()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetOnStreamCompletionTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetOnStreamCompletionTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetOnStreamCompletionTransformation(), target.GetOnStreamCompletionTransformation()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *RouteTransformations_RouteTransformation) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*RouteTransformations_RouteTransformation)
	if !ok {
		that2, ok := that.(RouteTransformations_RouteTransformation)
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

	if m.GetStage() != target.GetStage() {
		return false
	}

	switch m.Match.(type) {

	case *RouteTransformations_RouteTransformation_RequestMatch_:
		if _, ok := target.Match.(*RouteTransformations_RouteTransformation_RequestMatch_); !ok {
			return false
		}

		if h, ok := interface{}(m.GetRequestMatch()).(equality.Equalizer); ok {
			if !h.Equal(target.GetRequestMatch()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetRequestMatch(), target.GetRequestMatch()) {
				return false
			}
		}

	case *RouteTransformations_RouteTransformation_ResponseMatch_:
		if _, ok := target.Match.(*RouteTransformations_RouteTransformation_ResponseMatch_); !ok {
			return false
		}

		if h, ok := interface{}(m.GetResponseMatch()).(equality.Equalizer); ok {
			if !h.Equal(target.GetResponseMatch()) {
				return false
			}
		} else {
			if !proto.Equal(m.GetResponseMatch(), target.GetResponseMatch()) {
				return false
			}
		}

	default:
		// m is nil but target is not nil
		if m.Match != target.Match {
			return false
		}
	}

	return true
}

// Equal function
func (m *RouteTransformations_RouteTransformation_RequestMatch) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*RouteTransformations_RouteTransformation_RequestMatch)
	if !ok {
		that2, ok := that.(RouteTransformations_RouteTransformation_RequestMatch)
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

	if h, ok := interface{}(m.GetMatch()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMatch()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMatch(), target.GetMatch()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetRequestTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetRequestTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetRequestTransformation(), target.GetRequestTransformation()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetResponseTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetResponseTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetResponseTransformation(), target.GetResponseTransformation()) {
			return false
		}
	}

	if m.GetClearRouteCache() != target.GetClearRouteCache() {
		return false
	}

	return true
}

// Equal function
func (m *RouteTransformations_RouteTransformation_ResponseMatch) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*RouteTransformations_RouteTransformation_ResponseMatch)
	if !ok {
		that2, ok := that.(RouteTransformations_RouteTransformation_ResponseMatch)
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

	if h, ok := interface{}(m.GetMatch()).(equality.Equalizer); ok {
		if !h.Equal(target.GetMatch()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetMatch(), target.GetMatch()) {
			return false
		}
	}

	if h, ok := interface{}(m.GetResponseTransformation()).(equality.Equalizer); ok {
		if !h.Equal(target.GetResponseTransformation()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetResponseTransformation(), target.GetResponseTransformation()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *TransformationTemplate_HeaderToAppend) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*TransformationTemplate_HeaderToAppend)
	if !ok {
		that2, ok := that.(TransformationTemplate_HeaderToAppend)
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

	if strings.Compare(m.GetKey(), target.GetKey()) != 0 {
		return false
	}

	if h, ok := interface{}(m.GetValue()).(equality.Equalizer); ok {
		if !h.Equal(target.GetValue()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetValue(), target.GetValue()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *TransformationTemplate_DynamicMetadataValue) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*TransformationTemplate_DynamicMetadataValue)
	if !ok {
		that2, ok := that.(TransformationTemplate_DynamicMetadataValue)
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

	if strings.Compare(m.GetMetadataNamespace(), target.GetMetadataNamespace()) != 0 {
		return false
	}

	if strings.Compare(m.GetKey(), target.GetKey()) != 0 {
		return false
	}

	if h, ok := interface{}(m.GetValue()).(equality.Equalizer); ok {
		if !h.Equal(target.GetValue()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetValue(), target.GetValue()) {
			return false
		}
	}

	if m.GetJsonToProto() != target.GetJsonToProto() {
		return false
	}

	return true
}

// Equal function
func (m *TransformationTemplate_SpanTransformer) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*TransformationTemplate_SpanTransformer)
	if !ok {
		that2, ok := that.(TransformationTemplate_SpanTransformer)
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

	if h, ok := interface{}(m.GetName()).(equality.Equalizer); ok {
		if !h.Equal(target.GetName()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetName(), target.GetName()) {
			return false
		}
	}

	return true
}

// Equal function
func (m *MergeJsonKeys_OverridableTemplate) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*MergeJsonKeys_OverridableTemplate)
	if !ok {
		that2, ok := that.(MergeJsonKeys_OverridableTemplate)
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

	if h, ok := interface{}(m.GetTmpl()).(equality.Equalizer); ok {
		if !h.Equal(target.GetTmpl()) {
			return false
		}
	} else {
		if !proto.Equal(m.GetTmpl(), target.GetTmpl()) {
			return false
		}
	}

	if m.GetOverrideEmpty() != target.GetOverrideEmpty() {
		return false
	}

	return true
}

// Equal function
func (m *PromptEnrichment_Message) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*PromptEnrichment_Message)
	if !ok {
		that2, ok := that.(PromptEnrichment_Message)
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

	if strings.Compare(m.GetRole(), target.GetRole()) != 0 {
		return false
	}

	if strings.Compare(m.GetContent(), target.GetContent()) != 0 {
		return false
	}

	return true
}
