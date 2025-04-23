package als

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	envoy_al "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_al_file_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	envoy_al_grpc "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/grpc/v3"
	envoy_al_otel "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/open_telemetry/v3"
	envoy_metadata_formatter "github.com/envoyproxy/go-control-plane/envoy/extensions/formatter/metadata/v3"
	envoy_req_without_query "github.com/envoyproxy/go-control-plane/envoy/extensions/formatter/req_without_query/v3"
	envoy_tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/constants"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"google.golang.org/protobuf/proto"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	translatorutil "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

const (
	// OpenTelemetryAccessLog sink for the OpenTelemetry access log service
	OpenTelemetryAccessLog = "envoy.access_loggers.open_telemetry"
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

type aclType string

var (
	Hcm          aclType = "hcm"
	HttpListener aclType = "http-listener"
	Tcp          aclType = "tcp"
	Udp          aclType = "udp"
)

// DetectUnusefulCmds will detect commands that are not useful in the current configuration
// It returns errors that may some day be bubbled up arbitrarly high.
// See https://github.com/envoyproxy/envoy/blob/313b6fb7cf0f806e74a2d42c93e7c1fcccce2391/docs/root/configuration/observability/access_log/usage.rst?plain=1#L114-L123
func DetectUnusefulCmds(filterLocationType aclType, proposedLogFormats []*envoy_al.AccessLog) error {

	// TODO: programatically make sure we cover all command operators as found
	// in https://github.com/envoyproxy/envoy/blob/0f3e4aa373db6bbb7643b1bb60b0cb60d5b39df8/source/common/formatter/stream_info_formatter.cc#L1443
	// This could take place in ci especially if envoy version has been changed.

	var unusefulCmds []string
	switch filterLocationType {
	case Hcm:
		unusefulCmds = []string{"DOWNSTREAM_TRANSPORT_FAILURE_REASON"}
	case HttpListener:
	case Tcp:
		unusefulCmds = []string{"DOWNSTREAM_TRANSPORT_FAILURE_REASON", "REQ", "RESP", "TRAILER", "METADATA"}
	case Udp:
		unusefulCmds = []string{"UPSTREAM_TRANSPORT_FAILURE_REASON", "DOWNSTREAM_TRANSPORT_FAILURE_REASON", "REQ", "RESP", "TRAILER"}
	default:
		return errors.New("unknown accesslog level, cannot detect unuseful commands")
	}

	// TODO: Warn against deprecated commands for future proofing
	// deprecatedCmds := []struct{ cur, new string }{{cur: "DYNAMIC_METADATA", new: "METADATA"}}

	var errs []error
	for _, plf := range proposedLogFormats {
		// Dont put command operators in Names
		// Therefore we can extract their usage via strings
		dumpedACLSTR := fmt.Sprintf("%v", plf)
		issues := []string{}
		for _, cmd := range unusefulCmds {
			if strings.Contains(dumpedACLSTR, cmd) {
				issues = append(issues, cmd)
			}
		}
		// For a given proposed format this will include all the detected bad operators.
		if len(issues) > 0 {
			errs = append(errs, fmt.Errorf("unuseful command operators found in access log %s: %v", plf.GetName(), issues))
		}

	}

	return errors.Join(errs...)
}

// ProcessAccessLogPlugins will configure access logging for envoy, regardless of whether it will be applied to
// an HttpConnectionManager, http listener, TcpProxy NetworkFilter or perhaps someday a UdpProxy NetworkFilter.
// We have exposed plugins to allow configuration of http listeners and filters across multiple plugins.
// However, the TCP proxy is still configured by the TCP plugin only.
// To keep our access logging translation in a single place, we expose this function
// and the TCP plugin calls out to it.
func ProcessAccessLogPlugins(params plugins.Params, service *als.AccessLoggingService,
	logCfg []*envoy_al.AccessLog) ([]*envoy_al.AccessLog, error) {
	results := make([]*envoy_al.AccessLog, 0, len(service.GetAccessLog()))
	for _, al := range service.GetAccessLog() {
		var newAlsCfg envoy_al.AccessLog
		var err error

		// Make the "base" config with output destination
		switch cfgType := al.GetOutputDestination().(type) {
		case *als.AccessLog_FileSink:
			var cfg envoy_al_file_v3.FileAccessLog
			if err = copyFileSettings(&cfg, cfgType); err != nil {
				return nil, err
			}

			if newAlsCfg, err = translatorutil.NewAccessLogWithConfig(wellknown.FileAccessLog, &cfg); err != nil {
				return nil, err
			}

		case *als.AccessLog_GrpcService:
			var cfg envoy_al_grpc.HttpGrpcAccessLogConfig
			err := copyGrpcSettings(&cfg, cfgType)
			if err != nil {
				return nil, err
			}

			newAlsCfg, err = translatorutil.NewAccessLogWithConfig(wellknown.HTTPGRPCAccessLog, &cfg)
			if err != nil {
				return nil, err
			}
		case *als.AccessLog_OpenTelemetryService:
			var cfg envoy_al_otel.OpenTelemetryAccessLogConfig
			if err = copyOtelSettings(params, &cfg, cfgType); err != nil {
				return nil, err
			}

			newAlsCfg, err = translatorutil.NewAccessLogWithConfig(OpenTelemetryAccessLog, &cfg)
			if err != nil {
				return nil, err
			}
		}

		// Create and add the filter
		filter := al.GetFilter()
		err = translateFilter(&newAlsCfg, filter)
		if err != nil {
			return nil, err
		}

		results = append(results, &newAlsCfg)

	}

	logCfg = append(logCfg, results...)
	return logCfg, nil
}

// Since we are using the same proto def, marshal out of gloo format and unmarshal into envoy format
func translateFilter(accessLog *envoy_al.AccessLog, inFilter *als.AccessLogFilter) error {
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

	outFilter := &envoy_al.AccessLogFilter{}
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

func copyGrpcSettings(cfg *envoy_al_grpc.HttpGrpcAccessLogConfig, alsSettings *als.AccessLog_GrpcService) error {
	if alsSettings.GrpcService == nil {
		return eris.New("grpc service object cannot be nil")
	}

	svc := &envoy_core_v3.GrpcService{
		TargetSpecifier: &envoy_core_v3.GrpcService_EnvoyGrpc_{
			EnvoyGrpc: &envoy_core_v3.GrpcService_EnvoyGrpc{
				ClusterName: alsSettings.GrpcService.GetStaticClusterName(),
			},
		},
	}
	cfg.AdditionalRequestHeadersToLog = alsSettings.GrpcService.GetAdditionalRequestHeadersToLog()
	cfg.AdditionalResponseHeadersToLog = alsSettings.GrpcService.GetAdditionalResponseHeadersToLog()
	cfg.AdditionalResponseTrailersToLog = alsSettings.GrpcService.GetAdditionalResponseTrailersToLog()
	cfg.CommonConfig = &envoy_al_grpc.CommonGrpcAccessLogConfig{
		LogName:                 alsSettings.GrpcService.GetLogName(),
		GrpcService:             svc,
		TransportApiVersion:     envoy_core_v3.ApiVersion_V3,
		FilterStateObjectsToLog: alsSettings.GrpcService.GetFilterStateObjectsToLog(),
	}
	return cfg.Validate()
}

// getClustersForAccessLogs returns the clusters for the access loggers (if needed)
func getClustersForAccessLogs(
	params plugins.Params,
	proxy *v1.Proxy,
	reports reporter.ResourceReports,
	service *als.AccessLoggingService,
) []*envoy_config_cluster_v3.Cluster {
	clusters := []*envoy_config_cluster_v3.Cluster{}

	for _, al := range service.GetAccessLog() {
		switch cfgType := al.GetOutputDestination().(type) {
		case *als.AccessLog_OpenTelemetryService:
			// Create the cluster for the OpenTelemetry access log service
			cluster := createOtelCollectorCluster(params, proxy, reports, cfgType)
			if cluster == nil {
				continue
			}

			// Add the cluster to the list of clusters
			clusters = append(clusters, cluster)
		}
	}

	return clusters
}

func otelCollectorName(logName string) string {
	return fmt.Sprintf("%sotel_logs_%s", constants.SoloGeneratedClusterPrefix, logName)
}

// createOtelCollectorCluster creates a cluster for the OpenTelemetry collector
// that will receive the access logs.
func createOtelCollectorCluster(
	params plugins.Params,
	proxy *v1.Proxy,
	reports reporter.ResourceReports,
	alsSettings *als.AccessLog_OpenTelemetryService,
) *envoy_config_cluster_v3.Cluster {
	if alsSettings.OpenTelemetryService == nil {
		return nil
	}

	collector := alsSettings.OpenTelemetryService.GetCollector()
	if collector == nil {
		return nil
	}

	clusterName := otelCollectorName(alsSettings.OpenTelemetryService.GetLogName())

	host, port, err := net.SplitHostPort(collector.GetEndpoint())
	if err != nil {
		reports.AddError(proxy, fmt.Errorf("invalid OTEL log endpoint (%s): %v", collector.GetEndpoint(), err))
		return nil
	}

	if host == "" {
		reports.AddError(proxy, fmt.Errorf("invalid OTEL log endpoint (%s) missing host: %v",
			collector.GetEndpoint(), err))
		return nil
	}

	if port == "" {
		reports.AddError(proxy, fmt.Errorf("invalid OTEL log endpoint (%s) missing port: %v",
			collector.GetEndpoint(), err))
		return nil
	}

	discoveryType := envoy_config_cluster_v3.Cluster_STRICT_DNS
	addr := net.ParseIP(host)
	if addr != nil {
		discoveryType = envoy_config_cluster_v3.Cluster_STATIC
	}

	portValue, err := strconv.Atoi(port)
	if err != nil {
		reports.AddError(proxy, fmt.Errorf("invalid OTEL log endpoint port (%s): %v", port, err))
		return nil
	}

	cluster := &envoy_config_cluster_v3.Cluster{
		Name:           clusterName,
		ConnectTimeout: collector.GetTimeout(),
		// required to force envoy to use http2
		Http2ProtocolOptions: &envoy_core_v3.Http2ProtocolOptions{},
		ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{
			Type: discoveryType,
		},
		LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
				{
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Address: &envoy_core_v3.Address{
										Address: &envoy_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_core_v3.SocketAddress{
												Address: host,
												PortSpecifier: &envoy_core_v3.SocketAddress_PortValue{
													PortValue: uint32(portValue),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// if the collector is not insecure, we need to add the TLS context
	if !collector.GetInsecure() {
		cfg := &envoy_tls_v3.UpstreamTlsContext{
			Sni: host,
			CommonTlsContext: &envoy_tls_v3.CommonTlsContext{
				// default params
				TlsParams: &envoy_tls_v3.TlsParameters{},
			},
		}

		if sslConfig := collector.GetSslConfig(); sslConfig != nil {
			cfg, err = utils.NewSslConfigTranslator().ResolveUpstreamSslConfig(params.Snapshot.Secrets, sslConfig)
			if err != nil {
				// if we are configured to warn on missing tls secret and we match that error, add a
				// warning instead of error to the report.
				if params.Settings.GetGateway().GetValidation().GetWarnMissingTlsSecret().GetValue() &&
					errors.Is(err, utils.SslSecretNotFoundError) {
					reports.AddWarning(proxy, err.Error())
				} else {
					reports.AddError(proxy, err)
					return nil
				}
			}
		}

		typedConfig, err := utils.MessageToAny(cfg)
		if err != nil {
			reports.AddError(proxy, err)
			return nil
		}

		cluster.TransportSocket = &envoy_config_core_v3.TransportSocket{
			Name:       wellknown.TransportSocketTls,
			ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{TypedConfig: typedConfig},
		}
	}

	return cluster
}

func copyOtelSettings(params plugins.Params, cfg *envoy_al_otel.OpenTelemetryAccessLogConfig,
	alsSettings *als.AccessLog_OpenTelemetryService) error {
	if alsSettings.OpenTelemetryService == nil {
		return eris.New("OpenTelemetry service object cannot be empty")
	}

	if alsSettings.OpenTelemetryService.GetLogName() == "" {
		return eris.New("OpenTelemetry service log name cannot be empty")
	}

	collector := alsSettings.OpenTelemetryService.GetCollector()
	if collector == nil {
		return eris.New("OpenTelemetry service collector must be unset")
	}

	// check the ssl config and return error if a problem
	var sslWarning *translator.Warning
	if sslConfig := collector.GetSslConfig(); sslConfig != nil {
		_, err := utils.NewSslConfigTranslator().ResolveUpstreamSslConfig(params.Snapshot.Secrets, sslConfig)
		if err != nil {
			if params.Settings.GetGateway().GetValidation().GetWarnMissingTlsSecret().GetValue() &&
				errors.Is(err, utils.SslSecretNotFoundError) {
				sslWarning = &translator.Warning{
					Message: err.Error(),
				}
			} else {
				return err
			}
		}
	}

	cfg.CommonConfig = &envoy_al_grpc.CommonGrpcAccessLogConfig{
		LogName: alsSettings.OpenTelemetryService.GetLogName(),
		GrpcService: &envoy_core_v3.GrpcService{
			TargetSpecifier: &envoy_core_v3.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &envoy_core_v3.GrpcService_EnvoyGrpc{
					ClusterName: otelCollectorName(alsSettings.OpenTelemetryService.GetLogName()),
					Authority:   collector.GetAuthority(),
				},
			},
			InitialMetadata: convertHeaders(collector.GetHeaders()),
			Timeout:         collector.GetTimeout(),
		},
		FilterStateObjectsToLog: alsSettings.OpenTelemetryService.GetFilterStateObjectsToLog(),
	}
	cfg.DisableBuiltinLabels = alsSettings.OpenTelemetryService.GetDisableBuiltinLabels()
	cfg.Body = alsSettings.OpenTelemetryService.GetBody()
	cfg.Attributes = alsSettings.OpenTelemetryService.GetAttributes()

	err := cfg.Validate()
	if err != nil {
		return err
	}

	if sslWarning != nil {
		return sslWarning
	}

	return nil
}

func copyFileSettings(cfg *envoy_al_file_v3.FileAccessLog, alsSettings *als.AccessLog_FileSink) error {
	cfg.Path = alsSettings.FileSink.GetPath()

	formatterExtensions, err := getFormatterExtensions()
	if err != nil {
		return err
	}

	switch fileSinkType := alsSettings.FileSink.GetOutputFormat().(type) {
	case *als.FileSink_StringFormat:
		if fileSinkType.StringFormat != "" {
			cfg.AccessLogFormat = &envoy_al_file_v3.FileAccessLog_LogFormat{
				LogFormat: &envoy_core_v3.SubstitutionFormatString{
					Format: &envoy_core_v3.SubstitutionFormatString_TextFormat{
						TextFormat: fileSinkType.StringFormat,
					},
					Formatters: formatterExtensions,
				},
			}
		}
	case *als.FileSink_JsonFormat:
		cfg.AccessLogFormat = &envoy_al_file_v3.FileAccessLog_LogFormat{
			LogFormat: &envoy_core_v3.SubstitutionFormatString{
				Format: &envoy_core_v3.SubstitutionFormatString_JsonFormat{
					JsonFormat: fileSinkType.JsonFormat,
				},
				Formatters: formatterExtensions,
			},
		}
	}
	return cfg.Validate()
}

func getFormatterExtensions() ([]*envoy_core_v3.TypedExtensionConfig, error) {
	reqWithoutQueryFormatter := &envoy_req_without_query.ReqWithoutQuery{}
	reqWithoutQueryFormatterTc, err := utils.MessageToAny(reqWithoutQueryFormatter)
	if err != nil {
		return nil, err
	}

	mdFormatter := &envoy_metadata_formatter.Metadata{}
	mdFormatterTc, err := utils.MessageToAny(mdFormatter)
	if err != nil {
		return nil, err
	}

	return []*envoy_core_v3.TypedExtensionConfig{
		{
			Name:        "envoy.formatter.req_without_query",
			TypedConfig: reqWithoutQueryFormatterTc,
		},
		{
			Name:        "envoy.formatter.metadata",
			TypedConfig: mdFormatterTc,
		},
	}, nil
}

// convertHeaders converts a map of headers to add to a slice of Envoy HeaderValues
func convertHeaders(headersToAddMap map[string]string) []*envoy_config_core_v3.HeaderValue {
	var headersToAdd []*envoy_config_core_v3.HeaderValue
	for k, v := range headersToAddMap {
		headersToAdd = append(headersToAdd, &envoy_config_core_v3.HeaderValue{
			Key:   k,
			Value: v,
		})
	}
	return headersToAdd
}
