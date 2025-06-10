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
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
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
		listenerOptionsErrors := extractListenerOptionsErrors(proxyWithReports.Reports.ProxyReport)
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
	for _, proxyWithReport := range statusCtx.ProxiesWithReports {
		loErrors := extractListenerOptionsErrors(proxyWithReport.Reports.ProxyReport)
		for loKey, loerrs := range loErrors {
			statusForLO, ok := p.legacyStatusCache[loKey]
			if !ok {
				continue
			}

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
		for _, loerr := range status.errors {
			listenerOptionReport.AddError(loObjSk, errors.New(loerr.GetReason()))
		}

		// add any warnings
		listenerOptionReport.AddWarnings(loObjSk, status.warnings...)

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
	httpListenerReports := getAllHttpListenerReports(proxyReport.GetListenerReports())

	for _, hlr := range httpListenerReports {
		for _, hlerr := range hlr.GetErrors() {
			if loKey, ok := extractListenerOptionSourceKeys(hlerr); ok {
				listenerErrors[loKey] = append(listenerErrors[loKey], hlerr)
			}
		}
	}

	return listenerErrors
}

func getAllHttpListenerReports(listenerReports []*validation.ListenerReport) []*validation.HttpListenerReport {
	allReports := make([]*validation.HttpListenerReport, 0)
	for _, lr := range listenerReports {
		for _, hlr := range lr.GetAggregateListenerReport().GetHttpListenerReports() {
			allReports = append(allReports, hlr)
		}
	}
	return allReports
}

func extractListenerOptionSourceKeys(hlre *validation.HttpListenerReport_Error) (types.NamespacedName, bool) {
	metadata := hlre.GetMetadata()
	if metadata == nil {
		return types.NamespacedName{}, false
	}

	for _, src := range metadata.GetSources() {
		resourceKind := src.GetResourceKind()
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
