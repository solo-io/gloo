package listeneroptions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwquery "github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/listenerutils"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	lisquery "github.com/solo-io/gloo/projects/gateway2/translator/plugins/listeneroptions/query"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"go.uber.org/zap"
	"istio.io/istio/pkg/kube/krt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	_ plugins.ListenerPlugin = &plugin{}
	_ plugins.StatusPlugin   = &plugin{}

	ReadingListenerOptionErrStr = "error reading ListenerOption"
)

// legacyStatus holds the structures needed to derive and report classic GE status
type legacyStatus struct {
	subresourceStatus map[string]*core.Status
	errors            []*validation.HttpListenerReport_Error
	warnings          []string
}

func newLegacyStatus() *legacyStatus {
	return &legacyStatus{
		subresourceStatus: map[string]*core.Status{},
		errors:            []*validation.HttpListenerReport_Error{},
		warnings:          []string{},
	}
}

type legacyStatusCache map[types.NamespacedName]*legacyStatus

type plugin struct {
	gwQueries                gwquery.GatewayQueries
	lisOptQueries            lisquery.ListenerOptionQueries
	legacyStatusCache        legacyStatusCache
	listenerOptionCollection krt.Collection[*solokubev1.ListenerOption]
	statusReporter           reporter.StatusReporter
}

func NewPlugin(
	gwQueries gwquery.GatewayQueries,
	client client.Client,
	listenerOptionCollection krt.Collection[*solokubev1.ListenerOption],
	statusReporter reporter.StatusReporter,
) *plugin {
	return &plugin{
		gwQueries:                gwQueries,
		lisOptQueries:            lisquery.NewQuery(client),
		legacyStatusCache:        make(legacyStatusCache),
		listenerOptionCollection: listenerOptionCollection,
		statusReporter:           statusReporter,
	}
}

func (p *plugin) ApplyListenerPlugin(
	ctx context.Context,
	listenerCtx *plugins.ListenerContext,
	outListener *v1.Listener,
) error {
	// attachedOption represents the ListenerOptions targeting the Gateway on which this listener resides, and/or
	// the ListenerOptions which specifies this listener in section name
	attachedOptions, err := p.lisOptQueries.GetAttachedListenerOptions(ctx, listenerCtx.GwListener, listenerCtx.Gateway, listenerCtx.ListenerSet)
	if err != nil {
		return err
	}

	if len(attachedOptions) == 0 {
		return nil
	}

	// use the first option (highest in priority)
	// see for more context: https://github.com/solo-io/solo-projects/issues/6313
	optionUsed := attachedOptions[0]
	optionsUsed := []*solokubev1.ListenerOption{attachedOptions[0]}

	if outListener.GetOptions() != nil {
		outListener.Options, _ = glooutils.ShallowMergeListenerOptions(outListener.GetOptions(), optionUsed.Spec.GetOptions())
	} else {
		outListener.Options = optionUsed.Spec.GetOptions()
	}

	listenerutils.AppendSourceToListener(outListener, optionUsed)

	nn := client.ObjectKeyFromObject(optionUsed)
	p.legacyStatusCache[nn] = newLegacyStatus()

	// set a warning on any unused ListenerOptions
	for _, opt := range attachedOptions[1:] {
		nn := client.ObjectKeyFromObject(opt)
		cacheEntry := p.legacyStatusCache[nn]
		if cacheEntry == nil {
			cacheEntry = newLegacyStatus()
		}

		warning := fmt.Sprintf("ListenerOption '%s' not attached to Gateway '%s/%s' "+
			"due to higher priority ListenerOption '%s'", nn, listenerCtx.Gateway.Namespace,
			listenerCtx.Gateway.Name, optionsToStr(optionsUsed))
		cacheEntry.warnings = append(cacheEntry.warnings, warning)

		p.legacyStatusCache[nn] = cacheEntry
	}

	return nil
}

func (p *plugin) InitStatusPlugin(
	ctx context.Context,
	statusCtx *plugins.StatusContext,
) error {
	for _, proxyWithReports := range statusCtx.ProxiesWithReports {
		fmt.Printf("Processing proxy reports: %s\n", proxyWithReports.Reports.ProxyReport)

		listenerOptionsErrors := extractListenerOptionsErrors(proxyWithReports.Reports.ProxyReport)
		fmt.Printf("ListenerOptions errors: %v\n", listenerOptionsErrors)

		for loKey := range listenerOptionsErrors {
			p.legacyStatusCache[loKey] = newLegacyStatus()
		}
	}

	return nil
}

func (p *plugin) MergeStatusPlugin(ctx context.Context, source any) error {
	sourceStatusPlugin, ok := source.(*plugin)
	if !ok {
		return nil
	}

	for cacheKey, status := range sourceStatusPlugin.legacyStatusCache {
		destStatus, ok := p.legacyStatusCache[cacheKey]
		if !ok {
			destStatus = newLegacyStatus()
		}

		destStatus.errors = append(destStatus.errors, status.errors...)
		destStatus.warnings = append(destStatus.warnings, status.warnings...)

		for subresourceKey, subresourceStatus := range status.subresourceStatus {
			destStatus.subresourceStatus[subresourceKey] = subresourceStatus
		}

		p.legacyStatusCache[cacheKey] = destStatus
	}

	return nil
}

func (p *plugin) ApplyStatusPlugin(
	ctx context.Context,
	statusCtx *plugins.StatusContext,
) error {
	logger := contextutils.LoggerFrom(ctx).Desugar()

	for _, proxyWithReport := range statusCtx.ProxiesWithReports {
		proxyStatus := p.statusReporter.StatusFromReport(proxyWithReport.Reports.ResourceReports[proxyWithReport.Proxy], nil)

		loErrors := extractListenerOptionsErrors(proxyWithReport.Reports.ProxyReport)
		for loKey, loerrs := range loErrors {
			statusForLO, ok := p.legacyStatusCache[loKey]
			if !ok {
				continue
			}

			// set the subresource status for this specific proxy on the RO
			thisSubresourceStatus := statusForLO.subresourceStatus
			thisSubresourceStatus[xds.SnapshotCacheKey(proxyWithReport.Proxy)] = proxyStatus
			statusForLO.subresourceStatus = thisSubresourceStatus

			// add any routeErrors from this Proxy translation
			statusForLO.errors = append(statusForLO.errors, loerrs...)

			// update the cache
			p.legacyStatusCache[loKey] = statusForLO
		}
	}

	listenerOptionReport := make(reporter.ResourceReports)
	var multierr *multierror.Error
	for loKey, status := range p.legacyStatusCache {
		// get the object
		key := krt.Named{Namespace: loKey.Namespace, Name: loKey.Name}.ResourceName()
		maybeloObj := p.listenerOptionCollection.GetKey(key)
		if maybeloObj == nil {
			err := errors.New("ListenerOption not found")
			multierr = multierror.Append(multierr, eris.Wrapf(err, "%s %s in namespace %s",
				ReadingListenerOptionErrStr, loKey.Name, loKey.Namespace))
			continue
		}

		loObj := **maybeloObj
		loObj.Spec.Metadata = &core.Metadata{}
		loObj.Spec.GetMetadata().Name = loObj.GetName()
		loObj.Spec.GetMetadata().Namespace = loObj.GetNamespace()
		loObjSk := &loObj.Spec
		loObjSk.NamespacedStatuses = &loObj.Status

		// mark the object as accepted
		listenerOptionReport.Accept(loObjSk)

		// add any errors
		for i, loerr := range status.errors {
			loErr := errors.New(loerr.GetReason())
			logger.Debug("adding error to ListenerOption status", zap.Stringer("ListenerOption", loKey),
				zap.Error(loErr), zap.Int("errorIndex", i))
			listenerOptionReport.AddError(loObjSk, loErr)
		}

		// actually write out the reports
		err := p.statusReporter.WriteReports(ctx, listenerOptionReport, status.subresourceStatus)
		if err != nil {
			err = eris.Wrapf(err, "failed to write status for ListenerOption %s", loKey)
			multierr = multierror.Append(multierr, err)
			continue
		}
	}

	return multierr.ErrorOrNil()
}

func extractListenerOptionsErrors(
	proxyReport *validation.ProxyReport,
) map[types.NamespacedName][]*validation.HttpListenerReport_Error {
	listenerErrors := make(map[types.NamespacedName][]*validation.HttpListenerReport_Error)
	listenerReports := getAllListenerReports(proxyReport.GetListenerReports())
	fmt.Printf("Processing all listener reports: %v\n", listenerReports)

	for _, lr := range listenerReports {
		fmt.Printf("Processing listener report: %s\n", lr)

		alr := lr.GetAggregateListenerReport()
		if alr == nil {
			fmt.Printf("No aggregate listener report found\n")
			continue
		}

		fmt.Printf("Processing aggregate listener report: %s\n", alr)

		for _, hlr := range alr.GetHttpListenerReports() {
			for _, hlerr := range hlr.GetErrors() {
				fmt.Printf("Processing HTTP listener report: %s\n", hlerr)

				if loKey, ok := extractListenerOptionSourceKeys(hlerr); ok {
					errors := listenerErrors[loKey]
					errors = append(errors, hlerr)
					listenerErrors[loKey] = errors
				}
			}
		}
	}

	return listenerErrors
}

func getAllListenerReports(listenerReports []*validation.ListenerReport) []*validation.ListenerReport {
	allReports := make([]*validation.ListenerReport, 0)
	for _, lr := range listenerReports {
		allReports = append(allReports, lr)
	}
	return allReports
}

func extractListenerOptionSourceKeys(listenerError *validation.HttpListenerReport_Error) (types.NamespacedName, bool) {
	metadata := listenerError.GetMetadata()
	if metadata == nil {
		fmt.Printf("No metadata found: %v\n", listenerError)
		return types.NamespacedName{}, false
	}

	for _, src := range metadata.GetSources() {
		resourceKind := src.GetResourceKind()
		fmt.Printf("Resource kind: %s\n", resourceKind)

		if resourceKind == sologatewayv1.ListenerOptionGVK.Kind {
			resRef := src.GetResourceRef()
			return types.NamespacedName{
				Namespace: resRef.GetNamespace(),
				Name:      resRef.GetName(),
			}, true
		}
	}

	return types.NamespacedName{}, false
}

func optionsToStr(opts []*solokubev1.ListenerOption) string {
	resourceNames := make([]string, len(opts))
	for i, opt := range opts {
		resourceNames[i] = client.ObjectKeyFromObject(opt).String()
	}
	return strings.Join(resourceNames, ", ")
}
