package utils

import (
	"fmt"
	envoyal "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	"google.golang.org/protobuf/proto"
)

var (
	NoValueError = func(filterName string, fieldName string) error {
		return eris.Errorf("No value found for field %s of %s", fieldName, filterName)
	}
	InvalidEnumValueError = func(filterName string, fieldName string, value string) error {
		return eris.Errorf("Invalid value of %s in Enum field %s of %s", value, fieldName, filterName)
	}
	WrapInvalidEnumValueError = func(filterName string, err error) error {
		return eris.Wrap(err, fmt.Sprintf("Invalid subfilter in %s", filterName))
	}
)

// Since we are using the same proto def, marshal out of gloo format and unmarshal into envoy format
func TranslateFilter(accessLog *envoyal.AccessLog, inFilter *als.AccessLogFilter) error {
	if inFilter == nil {
		return nil
	}

	// We need to validate the enums in the filter manually because the protobuf libraries
	// do not validate them, for "compatibilty reasons". It's nicer to catch them here instead
	// of sending bad configs to Envoy.
	if err := validateFilterEnums(inFilter); err != nil {
		return err
	}

	bytes, err := proto.Marshal(inFilter)
	if err != nil {
		return err
	}

	outFilter := &envoyal.AccessLogFilter{}
	if err := proto.Unmarshal(bytes, outFilter); err != nil {
		return err
	}

	accessLog.Filter = outFilter
	return nil
}

func validateFilterEnums(filter *als.AccessLogFilter) error {
	switch filter := filter.GetFilterSpecifier().(type) {
	case *als.AccessLogFilter_RuntimeFilter:
		denominator := filter.RuntimeFilter.GetPercentSampled().GetDenominator()
		name := v3.FractionalPercent_DenominatorType_name[int32(denominator.Number())]
		if name == "" {
			return InvalidEnumValueError("RuntimeFilter", "FractionalPercent.Denominator", denominator.String())
		}
		runtimeKey := filter.RuntimeFilter.GetRuntimeKey()
		if len(runtimeKey) == 0 {
			return NoValueError("RuntimeFilter", "FractionalPercent.RuntimeKey")
		}
	case *als.AccessLogFilter_StatusCodeFilter:
		op := filter.StatusCodeFilter.GetComparison().GetOp()
		name := als.ComparisonFilter_Op_name[int32(op.Number())]
		if name == "" {
			return InvalidEnumValueError("StatusCodeFilter", "ComparisonFilter.Op", op.String())
		}
		value := filter.StatusCodeFilter.GetComparison().GetValue()
		if value == nil {
			return NoValueError("StatusCodeFilter", "ComparisonFilter.Value")
		}
		if len(value.GetRuntimeKey()) == 0 {
			return NoValueError("StatusCodeFilter", "ComparisonFilter.Value.RuntimeKey")
		}
	case *als.AccessLogFilter_DurationFilter:
		op := filter.DurationFilter.GetComparison().GetOp()
		name := als.ComparisonFilter_Op_name[int32(op.Number())]
		if name == "" {
			return InvalidEnumValueError("DurationFilter", "ComparisonFilter.Op", op.String())
		}
		value := filter.DurationFilter.GetComparison().GetValue()
		if value == nil {
			return NoValueError("DurationFilter", "ComparisonFilter.Value")
		}
		if len(value.GetRuntimeKey()) == 0 {
			return NoValueError("DurationFilter", "ComparisonFilter.Value.RuntimeKey")
		}
	case *als.AccessLogFilter_AndFilter:
		subfilters := filter.AndFilter.GetFilters()
		for _, f := range subfilters {
			err := validateFilterEnums(f)
			if err != nil {
				return WrapInvalidEnumValueError("AndFilter", err)
			}
		}
	case *als.AccessLogFilter_OrFilter:
		subfilters := filter.OrFilter.GetFilters()
		for _, f := range subfilters {
			err := validateFilterEnums(f)
			if err != nil {
				return WrapInvalidEnumValueError("OrFilter", err)
			}
		}
	case *als.AccessLogFilter_GrpcStatusFilter:
		statuses := filter.GrpcStatusFilter.GetStatuses()
		for _, status := range statuses {
			name := als.GrpcStatusFilter_Status_name[int32(status.Number())]
			if name == "" {
				return InvalidEnumValueError("GrpcStatusFilter", "Status", status.String())
			}
		}
	}

	return nil
}
