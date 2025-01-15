package utils

import (
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

const wildcardField = "*"

// OptionsMergeResult is an enum indicating the result of the merge of options from src to dst.
type OptionsMergeResult int

const (
	// OptionsMergedNone indicates that no fields were merged from src to dst.
	OptionsMergedNone OptionsMergeResult = iota

	// OptionsMergedPartial indicates that some but not all fields were merged from src to dst.
	OptionsMergedPartial

	// OptionsMergedFull indicates that all fields set in dst were overwritten by src.
	OptionsMergedFull
)

// ShallowMerge sets dst to the value of src, if src is non-zero and dst is zero-valued
// It returns a boolean indicating whether src overwrote dst.
func ShallowMerge(dst, src reflect.Value) bool {
	if !src.IsValid() {
		return false
	}

	if dst.CanSet() && !isEmptyValue(src) && isEmptyValue(dst) {
		dst.Set(src)
		return true
	}

	return false
}

// maySetValue sets dst to the value of src if:
// - src is set (has a non-nil value) and
// - dst is nil or overwrite is true
//
// It returns a boolean indicating whether src overwrote dst.
func maySetValue(dst, src reflect.Value, overwrite bool) bool {
	if src.CanSet() && !src.IsNil() && dst.CanSet() && (overwrite || dst.IsNil()) {
		dst.Set(src)
		return true
	}

	return false
}

// Copied from some previous version of https://github.com/golang/go/blob/68a32ced0f7b1b9abf9fd948db53c668ef6b1c66/src/encoding/json/encode.go#L304
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return isEmptyValue(v.Elem())
	case reflect.Func:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
}

// ShallowMergeRouteOptions merges the top-level fields of src into dst.
// The fields in dst that have non-zero values will not be overwritten.
// It performs a shallow merge of top-level fields only.
// It returns a boolean indicating whether any fields in src overwrote dst.
func ShallowMergeRouteOptions(dst, src *v1.RouteOptions) (*v1.RouteOptions, bool) {
	if src == nil {
		return dst, false
	}

	if dst == nil {
		return src.Clone().(*v1.RouteOptions), true
	}

	dstValue, srcValue := reflect.ValueOf(dst).Elem(), reflect.ValueOf(src).Elem()

	overwrote := false
	for i := range dstValue.NumField() {
		dstField, srcField := dstValue.Field(i), srcValue.Field(i)
		if srcOverride := ShallowMerge(dstField, srcField); srcOverride {
			overwrote = true
		}
	}

	return dst, overwrote
}

// ShallowMergeVirtualHostOptions merges the top-level fields of src into dst.
// The fields in dst that have non-zero values will not be overwritten.
// It performs a shallow merge of top-level fields only.
// It returns a boolean indicating whether any fields in src overwrote dst.
func ShallowMergeVirtualHostOptions(dst, src *v1.VirtualHostOptions) (*v1.VirtualHostOptions, bool) {
	if src == nil {
		return dst, false
	}

	if dst == nil {
		return src.Clone().(*v1.VirtualHostOptions), true
	}

	dstValue, srcValue := reflect.ValueOf(dst).Elem(), reflect.ValueOf(src).Elem()

	overwrote := false
	for i := range dstValue.NumField() {
		dstField, srcField := dstValue.Field(i), srcValue.Field(i)
		if srcOverride := ShallowMerge(dstField, srcField); srcOverride {
			overwrote = true
		}
	}

	return dst, overwrote
}

// MergeRouteOptionsWithOverrides merges the top-level fields of src that are allowed to be merged into dst.
// allowedOverrides is a set of field names in lowercase that are allowed to be merged, and may contain a wildcard field "*"
// to allow all fields to be merged.
// When allowedOverrides is empty, only fields unset in dst and set in src will be merged into dst.
//
// It returns an Enum indicating the result of the merge.
func MergeRouteOptionsWithOverrides(dst, src *v1.RouteOptions, allowedOverrides sets.Set[string]) (*v1.RouteOptions, OptionsMergeResult) {
	if src == nil {
		return dst, OptionsMergedNone
	}

	if dst == nil {
		return src.Clone().(*v1.RouteOptions), OptionsMergedFull
	}

	dstValue, srcValue := reflect.ValueOf(dst).Elem(), reflect.ValueOf(src).Elem()

	// By default, do not allow fields in src to overwrite fields in dst.
	// If allowedOverrides is non-empty, enable overwrites for the allowed fields.
	overwriteByDefault := false
	if allowedOverrides.Len() > 0 {
		overwriteByDefault = true
	}

	var srcFieldsUsed int
	var dstFieldsSet int
	var dstFieldsOverwritten int
	for i := range dstValue.NumField() {
		dstField, srcField := dstValue.Field(i), srcValue.Field(i)

		// NOTE: important to pre-compute this because dstFieldsOverwritten must be
		// incremented based on the original value of dstField and not the overwritten value
		dstIsSet := dstField.CanSet() && !dstField.IsNil()
		if dstIsSet {
			dstFieldsSet++
		}

		// Allow overrides if the field in dst is unset as merging from src into dst by default
		// allows src to augment dst by setting fields unset in dst, hence the override check only
		// applies when the field in dst is set (dstIsSet=true).
		if dstIsSet && overwriteByDefault && !(allowedOverrides.Has(wildcardField) ||
			allowedOverrides.Has(strings.ToLower(dstValue.Type().Field(i).Name))) {
			continue
		}

		if srcOverride := maySetValue(dstField, srcField, overwriteByDefault); srcOverride {
			srcFieldsUsed++
			if dstIsSet {
				dstFieldsOverwritten++
			}
		}
	}

	var overrideState OptionsMergeResult
	if srcFieldsUsed == 0 {
		overrideState = OptionsMergedNone
	} else if dstFieldsSet == dstFieldsOverwritten {
		overrideState = OptionsMergedFull
	} else {
		overrideState = OptionsMergedPartial
	}

	return dst, overrideState
}
