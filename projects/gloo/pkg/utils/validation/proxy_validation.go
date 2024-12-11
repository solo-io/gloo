package validation

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gloo/pkg/utils"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"go.uber.org/multierr"
)

var (
	RouteErrorMsg      = "Route Error"
	RouteIdentifierTxt = "Route Name"
)

const (
	ErrorLevels_WARNING string = "WARNING"
	ErrorLevels_ERROR   string = "ERROR"
)

// ErrorLevelContext carries additional information to be passed in the
// ErrorWithKnownLevel object
type ErrorLevelContext struct {
	// A given TCP listener can have multiple TCP hosts; HostNum provides the
	// numerical index of the host on the listener associated with the error to
	// be reported
	HostNum *int
}

// Types implementing ErrorWithKnownLevel are able to report on the severity of
// the error they describe by evaluating ErrorLevel()
type ErrorWithKnownLevel interface {
	error
	// The severity of the error - should return either ErrorLevels_WARNING or
	// ErrorLevels_ERROR
	ErrorLevel() string
	// Additional contextual information required to report the error/warning
	GetContext() ErrorLevelContext
	// The instance of the error itself
	GetError() error
}

// TcpHostWarning implements ErrorWithKnownLevel; it is intended to allow
// reporting certain errors as warnings and others as errors, and provides the
// necessary context required for reporting errors as such
type TcpHostWarning struct {
	Err      error
	ErrLevel string
	Context  ErrorLevelContext
}

func (tcpHostWarning *TcpHostWarning) ErrorLevel() string {
	return tcpHostWarning.ErrLevel
}

func (tcpHostWarning *TcpHostWarning) Error() string {
	return fmt.Sprintf("TcpHost error: %v", tcpHostWarning.Err)
}

func (tcpHostWarning *TcpHostWarning) GetContext() ErrorLevelContext {
	return tcpHostWarning.Context
}

// return the instance of the Error this object is wrapping
func (tcpHostWarning *TcpHostWarning) GetError() error {
	return tcpHostWarning.Err
}

func MakeReport(proxy *v1.Proxy) *validation.ProxyReport {
	listeners := proxy.GetListeners()
	listenerReports := make([]*validation.ListenerReport, len(listeners))

	for i, listener := range listeners {
		switch listenerType := listener.GetListenerType().(type) {
		case *v1.Listener_HttpListener:

			listenerReports[i] = &validation.ListenerReport{
				ListenerTypeReport: &validation.ListenerReport_HttpListenerReport{
					HttpListenerReport: &validation.HttpListenerReport{
						VirtualHostReports: makeVhostReports(listenerType.HttpListener.GetVirtualHosts()),
					},
				},
			}
		case *v1.Listener_TcpListener:
			listenerReports[i] = &validation.ListenerReport{
				ListenerTypeReport: &validation.ListenerReport_TcpListenerReport{
					TcpListenerReport: &validation.TcpListenerReport{
						TcpHostReports: makeTcpHostReports(listenerType.TcpListener.GetTcpHosts()),
					},
				},
			}
		case *v1.Listener_HybridListener:
			matchedListenerReports := make(map[string]*validation.MatchedListenerReport)
			for _, matchedListener := range listenerType.HybridListener.GetMatchedListeners() {
				matchedListenerReport := &validation.MatchedListenerReport{}
				switch matchedListenerType := matchedListener.GetListenerType().(type) {
				case *v1.MatchedListener_HttpListener:
					matchedListenerReport.ListenerReportType = &validation.MatchedListenerReport_HttpListenerReport{
						HttpListenerReport: &validation.HttpListenerReport{
							VirtualHostReports: makeVhostReports(matchedListenerType.HttpListener.GetVirtualHosts()),
						},
					}
				case *v1.MatchedListener_TcpListener:
					matchedListenerReport.ListenerReportType = &validation.MatchedListenerReport_TcpListenerReport{
						TcpListenerReport: &validation.TcpListenerReport{
							TcpHostReports: makeTcpHostReports(matchedListenerType.TcpListener.GetTcpHosts()),
						},
					}
				}
				matchedListenerReports[utils.MatchedRouteConfigName(listener, matchedListener.GetMatcher())] = matchedListenerReport
			}

			listenerReports[i] = &validation.ListenerReport{
				ListenerTypeReport: &validation.ListenerReport_HybridListenerReport{
					HybridListenerReport: &validation.HybridListenerReport{
						MatchedListenerReports: matchedListenerReports,
					},
				},
			}
		case *v1.Listener_AggregateListener:
			httpListenerReports := make(map[string]*validation.HttpListenerReport)
			tcpListenerReports := make(map[string]*validation.TcpListenerReport)
			httpResources := listenerType.AggregateListener.GetHttpResources()
			for _, httpFilterChain := range listenerType.AggregateListener.GetHttpFilterChains() {
				var virtualHosts []*v1.VirtualHost
				for _, vhostRef := range httpFilterChain.GetVirtualHostRefs() {
					virtualHosts = append(virtualHosts, httpResources.GetVirtualHosts()[vhostRef])
				}

				httpListenerReport := &validation.HttpListenerReport{
					VirtualHostReports: makeVhostReports(virtualHosts),
				}
				httpListenerReports[utils.MatchedRouteConfigName(listener, httpFilterChain.GetMatcher())] = httpListenerReport
			}

			for _, tcpListener := range listenerType.AggregateListener.GetTcpListeners() {
				tcpListenerReport := &validation.TcpListenerReport{
					TcpHostReports: makeTcpHostReports(tcpListener.GetTcpListener().GetTcpHosts()),
				}
				tcpListenerReports[utils.MatchedRouteConfigName(listener, tcpListener.GetMatcher())] = tcpListenerReport
			}

			listenerReports[i] = &validation.ListenerReport{
				ListenerTypeReport: &validation.ListenerReport_AggregateListenerReport{
					AggregateListenerReport: &validation.AggregateListenerReport{
						HttpListenerReports: httpListenerReports,
						TcpListenerReports:  tcpListenerReports,
					},
				},
			}

		}
	}

	return &validation.ProxyReport{
		ListenerReports: listenerReports,
	}
}

func makeVhostReports(virtualHosts []*v1.VirtualHost) []*validation.VirtualHostReport {
	vhostReports := make([]*validation.VirtualHostReport, len(virtualHosts))

	for j, vh := range virtualHosts {
		routeReports := make([]*validation.RouteReport, len(vh.GetRoutes()))
		for k := range vh.GetRoutes() {
			routeReports[k] = &validation.RouteReport{}
		}

		vhostReports[j] = &validation.VirtualHostReport{
			RouteReports: routeReports,
		}
	}

	return vhostReports
}

func makeTcpHostReports(tcpHosts []*v1.TcpHost) []*validation.TcpHostReport {
	tcpHostReports := make([]*validation.TcpHostReport, len(tcpHosts))
	for j := range tcpHosts {
		tcpHostReports[j] = &validation.TcpHostReport{}
	}
	return tcpHostReports
}

func formattedError(level, errType, reason string) error {
	return errors.Errorf("%v Error: %v. Reason: %v", level, errType, reason)
}

func formattedWarning(level, errType, reason string) string {
	return fmt.Sprintf("%v Warning: %v. Reason: %v", level, errType, reason)
}

// GetListenerErr returns the formatted errors on the ListenerReport only. It does NOT recursively aggregate errors in its sub-reports such as its host report, route report, etc.
func GetListenerErr(listener *validation.ListenerReport) []error {
	errs := make([]error, len(listener.GetErrors()))
	for _, errReport := range listener.GetErrors() {
		errs = append(errs, formattedError("Listener", errReport.GetType().String(), errReport.GetReason()))
	}
	return errs
}

// GetListenerWarn returns the formatted warnings on the ListenerReport only. It does NOT recursively aggregate errors in its sub-reports such as its host report, route report, etc.
func GetListenerWarn(listener *validation.ListenerReport) []string {
	warnings := make([]string, len(listener.GetWarnings()))
	for _, warning := range listener.GetWarnings() {
		warnings = append(warnings, formattedWarning("Listener", warning.GetType().String(), warning.GetReason()))
	}
	return warnings
}

// GetHttpListenerErr returns the formatted errors on the HttpListenerReport only. It does NOT recursively aggregate errors in its sub-reports such as VirtualHostReports, RouteReports, etc.
func GetHttpListenerErr(httpListener *validation.HttpListenerReport) []error {
	errs := make([]error, len(httpListener.GetErrors()))
	for _, errReport := range httpListener.GetErrors() {
		errs = append(errs, formattedError("HttpListener", errReport.GetType().String(), errReport.GetReason()))
	}
	return errs
}

// GetHttpListenerWarning returns the formatted warnings on the HttpListenerReport only. It does NOT recursively aggregate warnings in its sub-reports such as VirtualHostReports, RouteReports, etc.
func GetHttpListenerWarning(httpListener *validation.HttpListenerReport) []string {
	warnings := make([]string, len(httpListener.GetWarnings()))
	for _, warning := range httpListener.GetWarnings() {
		warnings = append(warnings, formattedWarning("HttpListener", warning.GetType().String(), warning.GetReason()))
	}
	return warnings
}

// GetVirtualHostErr returns the formatted errors on the VirtualHostReport only. It does NOT recursively aggregate errors in its sub-reports such as Routes, etc.
func GetVirtualHostErr(virtualHost *validation.VirtualHostReport) []error {
	errs := make([]error, len(virtualHost.GetErrors()))
	for _, errReport := range virtualHost.GetErrors() {
		errs = append(errs, formattedError("VirtualHost", errReport.GetType().String(), errReport.GetReason()))
	}
	return errs
}

// GetRouteErr returns the formatted errors on the RouteReport.
func GetRouteErr(route *validation.RouteReport) []error {
	errs := make([]error, len(route.GetErrors()))
	for _, errReport := range route.GetErrors() {
		routeError := errors.Errorf("%v. Reason: %v", errReport.GetType().String(), errReport.GetReason())
		errs = append(errs, errors.Wrap(routeError, RouteErrorMsg))
	}
	return errs
}

// GetRouteWarning returns the formatted warnings on the RouteReport.
func GetRouteWarning(route *validation.RouteReport) []string {
	warnings := make([]string, len(route.GetWarnings()))
	for _, warning := range route.GetWarnings() {
		warnings = append(warnings, formattedWarning("Route", warning.GetType().String(), warning.GetReason()))
	}
	return warnings
}

// GetTcpListenerErr returns the formatted errors on the TcpListenerReport only. It does NOT recursively aggregate errors in its sub-reports such as TcpHostReports, etc.
func GetTcpListenerErr(tcpListener *validation.TcpListenerReport) []error {
	errs := make([]error, len(tcpListener.GetErrors()))
	for _, errReport := range tcpListener.GetErrors() {
		errs = append(errs, formattedError("TcpListener", errReport.GetType().String(), errReport.GetReason()))
	}
	return errs
}

// GetTcpListenerWarning returns the formatted warnings on the TcpListenerReport only. It does NOT recursively aggregate warnings in its sub-reports such as TcpHostReports, etc.
func GetTcpListenerWarning(tcpListener *validation.TcpListenerReport) []string {
	warnings := make([]string, len(tcpListener.GetWarnings()))
	for _, warning := range tcpListener.GetWarnings() {
		warnings = append(warnings, formattedWarning("TcpListener", warning.GetType().String(), warning.GetReason()))
	}
	return warnings
}

// GetTcpHostErr returns the formatted errors on the TcpHostReport.
func GetTcpHostErr(tcpHost *validation.TcpHostReport) []error {
	errs := make([]error, len(tcpHost.GetErrors()))
	for _, errReport := range tcpHost.GetErrors() {
		errs = append(errs, formattedError("TcpHost", errReport.GetType().String(), errReport.GetReason()))
	}
	return errs
}

// GetTcpHostWarning returns the formatted warnings on the TcpHostReport.
func GetTcpHostWarning(tcpHost *validation.TcpHostReport) []string {
	warnings := make([]string, len(tcpHost.GetWarnings()))
	for _, warning := range tcpHost.GetWarnings() {
		warnings = append(warnings, formattedWarning("TcpHost", warning.GetType().String(), warning.GetReason()))
	}
	return warnings
}

// GetListenerError returns the aggregated errors on the ListenerReport. Based on the type of listener,
// it aggregates the errors in its sub-reports such as its host report, route report, etc.
func GetListenerError(listener *validation.ListenerReport) []error {
	var errs []error

	if err := GetListenerErr(listener); err != nil {
		errs = append(errs, err...)
	}
	switch listenerType := listener.GetListenerTypeReport().(type) {
	case *validation.ListenerReport_HttpListenerReport:
		errs = append(errs, getHttpListenerReportErrs(listenerType.HttpListenerReport)...)

	case *validation.ListenerReport_TcpListenerReport:
		errs = append(errs, getTcpListenerReportErrs(listenerType.TcpListenerReport)...)

	case *validation.ListenerReport_HybridListenerReport:
		for _, mr := range listenerType.HybridListenerReport.GetMatchedListenerReports() {
			switch lrt := mr.GetListenerReportType().(type) {
			case *validation.MatchedListenerReport_HttpListenerReport:
				errs = append(errs, getHttpListenerReportErrs(lrt.HttpListenerReport)...)
			case *validation.MatchedListenerReport_TcpListenerReport:
				errs = append(errs, getTcpListenerReportErrs(lrt.TcpListenerReport)...)
			}
		}

	case *validation.ListenerReport_AggregateListenerReport:
		for _, httpListenerReport := range listenerType.AggregateListenerReport.GetHttpListenerReports() {
			errs = append(errs, getHttpListenerReportErrs(httpListenerReport)...)
		}
		for _, tcpListenerReport := range listenerType.AggregateListenerReport.GetTcpListenerReports() {
			errs = append(errs, getTcpListenerReportErrs(tcpListenerReport)...)
		}
	}

	return errs
}

// GetListenerWarning returns the aggregated warnings on the ListenerReport. Based on the type of listener,
// it aggregates the warnings in its sub-reports such as its host report, route report, etc.
func GetListenerWarning(listener *validation.ListenerReport) []string {
	var warnings []string

	if warning := GetListenerWarn(listener); warning != nil {
		warnings = append(warnings, warning...)
	}
	switch listenerType := listener.GetListenerTypeReport().(type) {
	case *validation.ListenerReport_HttpListenerReport:
		warnings = append(warnings, getHttpListenerReportWarns(listenerType.HttpListenerReport)...)

	case *validation.ListenerReport_TcpListenerReport:
		warnings = append(warnings, getTcpListenerReportWarns(listenerType.TcpListenerReport)...)

	case *validation.ListenerReport_HybridListenerReport:
		for _, mr := range listenerType.HybridListenerReport.GetMatchedListenerReports() {
			switch lrt := mr.GetListenerReportType().(type) {
			case *validation.MatchedListenerReport_HttpListenerReport:
				warnings = append(warnings, getHttpListenerReportWarns(lrt.HttpListenerReport)...)
			case *validation.MatchedListenerReport_TcpListenerReport:
				warnings = append(warnings, getTcpListenerReportWarns(lrt.TcpListenerReport)...)
			}
		}

	case *validation.ListenerReport_AggregateListenerReport:
		for _, httpListenerReport := range listenerType.AggregateListenerReport.GetHttpListenerReports() {
			warnings = append(warnings, getHttpListenerReportWarns(httpListenerReport)...)
		}
		for _, tcpListenerReport := range listenerType.AggregateListenerReport.GetTcpListenerReports() {
			warnings = append(warnings, getTcpListenerReportWarns(tcpListenerReport)...)
		}
	}

	return warnings
}

// GetProxyError returns the aggregated errors on the ProxyReport. It iterates through the listeners and
// based on the type of listener, it aggregates the errors in its sub-reports such as its host report, route report, etc.
func GetProxyError(proxyRpt *validation.ProxyReport) error {
	var errs []error
	for _, listener := range proxyRpt.GetListenerReports() {
		listenerErrs := GetListenerError(listener)

		if listenerErrs != nil {
			errs = append(errs, listenerErrs...)
		}
	}

	var combinedErr error
	for _, err := range errs {
		combinedErr = multierr.Append(combinedErr, err)
	}

	return combinedErr
}

// getTcpListenerReportErrs returns the aggregated errors on the TcpListenerReport.
// It aggregates the errors in its sub-reports such as TcpHostReports, etc.
func getTcpListenerReportErrs(tcpListenerReport *validation.TcpListenerReport) []error {
	var errs []error

	if err := GetTcpListenerErr(tcpListenerReport); err != nil {
		errs = append(errs, err...)
	}
	for _, hostReport := range tcpListenerReport.GetTcpHostReports() {
		if err := GetTcpHostErr(hostReport); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}

// getTcpListenerReportWarns returns the aggregated warnings on the TcpListenerReport.
// It aggregates the warnings in its sub-reports such as TcpHostReports, etc.
func getTcpListenerReportWarns(tcpListenerReport *validation.TcpListenerReport) []string {
	var warnings []string

	if warning := GetTcpListenerWarning(tcpListenerReport); warning != nil {
		warnings = append(warnings, warning...)
	}
	for _, hostReport := range tcpListenerReport.GetTcpHostReports() {
		if err := GetTcpHostWarning(hostReport); err != nil {
			warnings = append(warnings, err...)
		}
	}

	return warnings
}

// getHttpListenerReportErrs returns the aggregated errors on the HttpListenerReport.
// It aggregates the errors in its sub-reports such as VirtualHostReports, RouteReports, etc.
func getHttpListenerReportErrs(httpListenerReport *validation.HttpListenerReport) []error {
	var errs []error

	if err := GetHttpListenerErr(httpListenerReport); err != nil {
		errs = append(errs, err...)
	}
	for _, vhReport := range httpListenerReport.GetVirtualHostReports() {
		if err := GetVirtualHostErr(vhReport); err != nil {
			errs = append(errs, err...)
		}
		for _, routeReport := range vhReport.GetRouteReports() {
			if err := GetRouteErr(routeReport); err != nil {
				errs = append(errs, err...)
			}
		}
	}

	return errs
}

// getHttpListenerReportWarns returns the aggregated warnings on the HttpListenerReport.
// It aggregates the warnings in its sub-reports such as Routes, etc.
func getHttpListenerReportWarns(httpListenerReport *validation.HttpListenerReport) []string {
	var warnings []string

	if err := GetHttpListenerWarning(httpListenerReport); err != nil {
		warnings = append(warnings, err...)
	}
	for _, vhReport := range httpListenerReport.GetVirtualHostReports() {
		for _, routeReport := range vhReport.GetRouteReports() {
			if err := GetRouteWarning(routeReport); err != nil {
				warnings = append(warnings, err...)
			}
		}
	}

	return warnings
}

// GetProxyWarning returns the aggregated warnings on the ProxyReport. It iterates through the listeners and
// based on the type of listener, it aggregates the warnings in its sub-reports such as its host report, route report, etc.
func GetProxyWarning(proxyRpt *validation.ProxyReport) []string {
	var warnings []string

	for _, listenerReport := range proxyRpt.GetListenerReports() {
		warnings = append(warnings, GetListenerWarning(listenerReport)...)
	}

	return warnings
}

func AppendListenerError(listenerReport *validation.ListenerReport, errType validation.ListenerReport_Error_Type, reason string) {
	listenerReport.Errors = append(listenerReport.GetErrors(), &validation.ListenerReport_Error{
		Type:   errType,
		Reason: reason,
	})
}

func AppendListenerWarning(listenerReport *validation.ListenerReport, errType validation.ListenerReport_Warning_Type, reason string) {
	listenerReport.Warnings = append(listenerReport.GetWarnings(), &validation.ListenerReport_Warning{
		Type:   errType,
		Reason: reason,
	})
}

func AppendVirtualHostError(virtualHostReport *validation.VirtualHostReport, errType validation.VirtualHostReport_Error_Type, reason string) {
	virtualHostReport.Errors = append(virtualHostReport.GetErrors(), &validation.VirtualHostReport_Error{
		Type:   errType,
		Reason: reason,
	})
}

func AppendVirtualHostErrorWithMetadata(virtualHostReport *validation.VirtualHostReport, errType validation.VirtualHostReport_Error_Type, reason string, metadata *v1.SourceMetadata) {
	virtualHostReport.Errors = append(virtualHostReport.GetErrors(), &validation.VirtualHostReport_Error{
		Type:     errType,
		Reason:   reason,
		Metadata: metadata,
	})
}

func AppendHTTPListenerError(httpListenerReport *validation.HttpListenerReport, errType validation.HttpListenerReport_Error_Type, reason string) {
	httpListenerReport.Errors = append(httpListenerReport.GetErrors(), &validation.HttpListenerReport_Error{
		Type:   errType,
		Reason: reason,
	})
}

func AppendHTTPListenerWarning(httpListenerReport *validation.HttpListenerReport, errType validation.HttpListenerReport_Warning_Type, reason string) {
	httpListenerReport.Warnings = append(httpListenerReport.GetWarnings(), &validation.HttpListenerReport_Warning{
		Type:   errType,
		Reason: reason,
	})
}

func AppendTCPListenerError(tcpListenerReport *validation.TcpListenerReport, errType validation.TcpListenerReport_Error_Type, reason string) {
	tcpListenerReport.Errors = append(tcpListenerReport.GetErrors(), &validation.TcpListenerReport_Error{
		Type:   errType,
		Reason: reason,
	})
}

func AppendTCPListenerWarning(httpListenerReport *validation.TcpListenerReport, errType validation.TcpListenerReport_Warning_Type, reason string) {
	httpListenerReport.Warnings = append(httpListenerReport.GetWarnings(), &validation.TcpListenerReport_Warning{
		Type:   errType,
		Reason: reason,
	})
}

func AppendRouteError(routeReport *validation.RouteReport, errType validation.RouteReport_Error_Type, reason string, routeName string) {
	routeReport.Errors = append(routeReport.GetErrors(), &validation.RouteReport_Error{
		Type:   errType,
		Reason: fmt.Sprintf("%s. %s: %s", reason, RouteIdentifierTxt, routeName),
	})
}

func AppendRouteErrorWithMetadata(routeReport *validation.RouteReport, errType validation.RouteReport_Error_Type, reason string, routeName string, metadata *v1.SourceMetadata) {
	routeReport.Errors = append(routeReport.GetErrors(), &validation.RouteReport_Error{
		Type:     errType,
		Reason:   fmt.Sprintf("%s. %s: %s", reason, RouteIdentifierTxt, routeName),
		Metadata: metadata,
	})
}

func AppendRouteWarning(routeReport *validation.RouteReport, errType validation.RouteReport_Warning_Type, reason string) {
	routeReport.Warnings = append(routeReport.GetWarnings(), &validation.RouteReport_Warning{
		Type:   errType,
		Reason: reason,
	})
}

func AppendTcpHostWarning(tcpHostReport *validation.TcpHostReport, errType validation.TcpHostReport_Warning_Type, reason string) {
	tcpHostReport.Warnings = append(tcpHostReport.GetWarnings(), &validation.TcpHostReport_Warning{
		Type:   errType,
		Reason: reason,
	})
}
