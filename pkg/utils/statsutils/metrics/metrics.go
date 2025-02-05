package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/statsutils"
	gwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/jsonpath"
)

type MetricLabels = gloov1.Settings_ObservabilityOptions_MetricLabels

var Names = map[schema.GroupVersionKind]string{
	gwv1.GatewayGVK:         "validation.gateway.solo.io/gateway_config_status",
	gwv1.RouteTableGVK:      "validation.gateway.solo.io/route_table_config_status",
	gwv1.VirtualServiceGVK:  "validation.gateway.solo.io/virtual_service_config_status",
	gloov1.ProxyGVK:         "validation.gateway.solo.io/proxy_config_status",
	gloov1.SecretGVK:        "validation.gateway.solo.io/secret_config_status",
	gloov1.UpstreamGVK:      "validation.gateway.solo.io/upstream_config_status",
	gloov1.UpstreamGroupGVK: "validation.gateway.solo.io/upsteam_group_config_status",
}

var descriptions = map[schema.GroupVersionKind]string{
	gwv1.GatewayGVK:         "The health status of gateway resources in the system. A value of 0 indicates that there are no issues.",
	gwv1.RouteTableGVK:      "The health status of route table resources in the system. A value of 0 indicates that there are no issues.",
	gwv1.VirtualServiceGVK:  "The health status of virtual service resources in the system. A value of 0 indicates that there are no issues.",
	gloov1.ProxyGVK:         "The health status of proxy resources in the system. A value of 0 indicates that there are no issues.",
	gloov1.SecretGVK:        "The health status of secret resources in the system. A value of 0 indicates that there are no issues.",
	gloov1.UpstreamGVK:      "The health status of upstream resources in the system. A value of 0 indicates that there are no issues.",
	gloov1.UpstreamGroupGVK: "The health status of upstream group resources in the system. A value of 0 indicates that there are no issues.",
}

// ConfigStatusMetrics is a collection of metrics, each of which records if the configuration for
// a particular resource type is valid
type ConfigStatusMetrics struct {
	opts    map[string]*MetricLabels
	metrics map[schema.GroupVersionKind]*resourceMetric
}

// resourceMetric is a wrapper around a gauge. In addition to a gauge itself, it stores information
// regarding which labels should get applied to it, and how to obtain the values for those labels.
type resourceMetric struct {
	gauge       *stats.Int64Measure
	labelToPath map[string]string
}

func GetDefaultConfigStatusOptions() map[string]*MetricLabels {
	return make(map[string]*MetricLabels)
}

// NewConfigStatusMetrics creates and returns a ConfigStatusMetrics from the specified options.
// If the options are invalid, an error is returned.
func NewConfigStatusMetrics(opts map[string]*MetricLabels) (ConfigStatusMetrics, error) {
	metrics, err := prepareMetrics(opts)
	if err != nil {
		return ConfigStatusMetrics{}, err
	}

	configMetrics := ConfigStatusMetrics{
		opts:    opts,
		metrics: metrics,
	}

	return configMetrics, nil
}

func prepareMetrics(opts map[string]*MetricLabels) (map[schema.GroupVersionKind]*resourceMetric, error) {
	metrics := make(map[schema.GroupVersionKind]*resourceMetric)
	for gvkString, labels := range opts {
		gvk, err := parseGroupVersionKind(gvkString)
		if err != nil {
			return metrics, err
		}
		metric, err := newResourceMetric(gvk, labels.GetLabelToPath())
		if err != nil {
			return metrics, err
		}
		metrics[gvk] = metric
	}
	return metrics, nil
}

func parseGroupVersionKind(arg string) (schema.GroupVersionKind, error) {
	gvk, _ := schema.ParseKindArg(arg)
	if gvk == nil {
		return schema.GroupVersionKind{}, errors.Errorf("unable to parse GVK from string '%s'", arg)
	}
	if _, ok := Names[*gvk]; !ok {
		return schema.GroupVersionKind{}, errors.Errorf("config status metric reporting is not supported for resource type '%s'", arg)
	}
	return *gvk, nil
}

func resourceToGVK(resource resources.Resource) (schema.GroupVersionKind, error) {
	switch resource.(type) {
	// Gateway resources
	case *gwv1.Gateway:
		return gwv1.GatewayGVK, nil
	case *gwv1.RouteTable:
		return gwv1.RouteTableGVK, nil
	case *gwv1.VirtualService:
		return gwv1.VirtualServiceGVK, nil
	// Gloo resources
	case *gloov1.Proxy:
		return gloov1.ProxyGVK, nil
	case *gloov1.Secret:
		return gloov1.SecretGVK, nil
	case *gloov1.Upstream:
		return gloov1.UpstreamGVK, nil
	case *gloov1.UpstreamGroup:
		return gloov1.UpstreamGroupGVK, nil
	default:
		return schema.GroupVersionKind{}, errors.Errorf("config status metric reporting is not supported for resource type: %T", resource)
	}
}

func (m *ConfigStatusMetrics) SetResourceStatus(ctx context.Context, resource resources.Resource, status *core.Status) {
	fmt.Printf("SetResourceStatus: %s\n, %s\n", resource.GetMetadata().Ref(), status.GetState())

	if status.GetState() == core.Status_Warning || status.GetState() == core.Status_Rejected {
		m.SetResourceInvalid(ctx, resource)
		return
	}
	// Don't bother setting the metric while pending, we'll set it momentarily when it transitions
	if status.GetState() == core.Status_Accepted {
		m.SetResourceValid(ctx, resource)
	}
}

func (m *ConfigStatusMetrics) SetResourceValid(ctx context.Context, resource resources.Resource) {
	log := contextutils.LoggerFrom(ctx)
	gvk, err := resourceToGVK(resource)
	if err != nil {
		log.Debugf(err.Error())
		return
	}
	if m.metrics[gvk] != nil {
		log.Debugf("Setting '%s' config metric valid", resource.GetMetadata().Ref())
		mutators, err := getMutators(m.metrics[gvk], resource)
		if err != nil {
			log.Errorf("Error setting labels on %s: %s", Names[gvk], err.Error())
		}
		statsutils.MeasureZero(ctx, m.metrics[gvk].gauge, mutators...)
	} else {
		log.Warnf("Skipping setting valid metric for resource %s, no metric configured", resource.GetMetadata().Ref())
	}
}

func (m *ConfigStatusMetrics) SetResourceInvalid(ctx context.Context, resource resources.Resource) {
	log := contextutils.LoggerFrom(ctx)
	gvk, err := resourceToGVK(resource)
	if err != nil {
		log.Debugf(err.Error())
		return
	}
	if m.metrics[gvk] != nil {
		log.Debugf("Setting '%s' config metric invalid", resource.GetMetadata().Ref())
		mutators, err := getMutators(m.metrics[gvk], resource)
		if err != nil {
			log.Errorf("Error setting labels on %s: %s", Names[gvk], err.Error())
		}
		statsutils.MeasureOne(ctx, m.metrics[gvk].gauge, mutators...)
	} else {
		log.Warnf("Skipping setting invalid metric for resource %s, no metric configured", resource.GetMetadata().Ref())
	}
}

// ClearMetrics removes all metrics from the ConfigStatusMetrics
func (m *ConfigStatusMetrics) ClearMetrics(ctx context.Context) {
	log := contextutils.LoggerFrom(ctx)

	someViewsUnregistered := false

	// Iterate through the resource metrics and unregister them
	for _, metric := range m.metrics {
		v := view.Find(metric.gauge.Name())
		if v != nil {
			view.Unregister(v)
			someViewsUnregistered = true
		}
	}

	// Only sleep some metrics were unregistered
	if someViewsUnregistered {
		// Wait for the view to be unregistered (a channel is used)
		// This is necessary because the view is unregistered asynchronously.
		// We may not need this after we upgrade to an newer metrics package
		time.Sleep(1 * time.Second)
	}

	// Add fresh metrics
	var err error
	m.metrics, err = prepareMetrics(m.opts)
	if err != nil {
		log.Errorf("Error clearing resource metrics: %s", err.Error())
	}
}

func getMutators(metric *resourceMetric, resource resources.Resource) ([]tag.Mutator, error) {
	numLabels := len(metric.labelToPath)
	mutators := make([]tag.Mutator, numLabels)
	i := 0
	for k, v := range metric.labelToPath {
		key, err := tag.NewKey(k)
		if err != nil {
			return nil, err
		}
		value, err := extractValueFromResource(resource, v)
		if err != nil {
			return nil, err
		}
		mutators[i] = tag.Upsert(key, value)
		i++
	}
	return mutators, nil
}

// Grab the value at the specified json path from the resource
func extractValueFromResource(resource resources.Resource, jsonPath string) (string, error) {
	j := jsonpath.New("ConfigStatusMetrics")
	// Parse the template
	err := j.Parse(jsonPath)
	if err != nil {
		return "", err
	}
	// grab the result from the resource
	values, err := j.FindResults(resource)
	if err != nil {
		return "", nil
	}

	var valueStrings []string
	if len(values) == 0 || len(values[0]) == 0 {
		valueStrings = append(valueStrings, "<none>")
	}
	for i := range values {
		for j := range values[i] {
			valueStrings = append(valueStrings, fmt.Sprintf("%v", values[i][j].Interface()))
		}
	}
	output := strings.Join(valueStrings, ",")
	return output, nil
}

// Returns a resourceMetric, or nil if labelToPath is nil or empty. An error is returned if the
// labelToPath configuration is invalid (for example, specifies an invalid label key).
func newResourceMetric(gvk schema.GroupVersionKind, labelToPath map[string]string) (*resourceMetric, error) {
	numLabels := len(labelToPath)
	if numLabels == 0 {
		return nil, nil
	}
	tagKeys := make([]tag.Key, numLabels)
	i := 0
	for k := range labelToPath {
		var err error
		tagKeys[i], err = tag.NewKey(k)
		if err != nil {
			return nil, errors.Wrapf(err, "Error creating resourceMetric for %s", Names[gvk])
		}
		i++
	}
	return &resourceMetric{
		gauge:       statsutils.MakeGauge(Names[gvk], descriptions[gvk], tagKeys...),
		labelToPath: labelToPath,
	}, nil
}
