package setup

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/pkg/utils/setuputils"
	gloostatusutils "github.com/solo-io/gloo/pkg/utils/statusutils"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway2/controller"
	"github.com/solo-io/gloo/projects/gateway2/deployer"
	"github.com/solo-io/gloo/projects/gateway2/extensions"
	"github.com/solo-io/gloo/projects/gateway2/krtcollections"
	"github.com/solo-io/gloo/projects/gateway2/proxy_syncer"
	"github.com/solo-io/gloo/projects/gloo/constants"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glookubev1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/shared"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

var settingsGVR = glookubev1.SchemeGroupVersion.WithResource("settings")

func createKubeClient(restConfig *rest.Config) (istiokube.Client, error) {
	restCfg := istiokube.NewClientConfigForRestConfig(restConfig)
	client, err := istiokube.NewClient(restCfg, "")
	if err != nil {
		return nil, err
	}
	istiokube.EnableCrdWatcher(client)
	return client, nil
}

func getInitialSettings(ctx context.Context, c istiokube.Client, nns types.NamespacedName) *glookubev1.Settings {
	// get initial settings
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("getting initial settings. gvr: %v", settingsGVR)

	i, err := c.Dynamic().Resource(settingsGVR).Namespace(nns.Namespace).Get(ctx, nns.Name, metav1.GetOptions{})
	if err != nil {
		logger.Panicf("failed to get initial settings: %v", err)
		return nil
	}
	logger.Infof("got initial settings")

	var empty glookubev1.Settings
	out := &empty
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(i.UnstructuredContent(), out)
	if err != nil {
		logger.Panicf("failed converting unstructured into settings: %v", i)
		return nil
	}
	return out
}

func StartGGv2(ctx context.Context,
	setupOpts *bootstrap.SetupOpts,
	uccBuilder krtcollections.UniquelyConnectedClientsBulider,
	extensionsFactory extensions.K8sGatewayExtensionsFactory,
	pluginRegistryFactory func(opts registry.PluginOpts) plugins.PluginRegistryFactory,
) error {
	restConfig := ctrl.GetConfigOrDie()

	return StartGGv2WithConfig(ctx, setupOpts, restConfig, uccBuilder, extensionsFactory, pluginRegistryFactory, setuputils.SetupNamespaceName())
}

func StartGGv2WithConfig(ctx context.Context,
	setupOpts *bootstrap.SetupOpts,
	restConfig *rest.Config,
	uccBuilder krtcollections.UniquelyConnectedClientsBulider,
	extensionsFactory extensions.K8sGatewayExtensionsFactory,
	pluginRegistryFactory func(opts registry.PluginOpts) plugins.PluginRegistryFactory,
	settingsNns types.NamespacedName,
) error {
	ctx = contextutils.WithLogger(ctx, "k8s")

	logger := contextutils.LoggerFrom(ctx)
	logger.Info("starting gloo gateway")

	kubeClient, err := createKubeClient(restConfig)
	if err != nil {
		return err
	}

	initialSettings := getInitialSettings(ctx, kubeClient, settingsNns)
	if initialSettings == nil {
		return fmt.Errorf("initial settings not found")
	}

	logger.Info("creating krt collections")
	augmentedPods := krtcollections.NewPodsCollection(ctx, kubeClient, setupOpts.KrtDebugger)
	setting := proxy_syncer.SetupCollectionDynamic[glookubev1.Settings](
		ctx,
		kubeClient,
		settingsGVR,
		krt.WithName("GlooSettings"))

	augmentedPodsForUcc := augmentedPods
	if envutils.IsEnvTruthy("DISABLE_POD_LOCALITY_XDS") {
		augmentedPodsForUcc = nil
	}

	ucc := uccBuilder(ctx, setupOpts.KrtDebugger, augmentedPodsForUcc)

	settingsSingle := krt.NewSingleton(func(ctx krt.HandlerContext) *glookubev1.Settings {
		s := krt.FetchOne(ctx, setting,
			krt.FilterObjectName(settingsNns))
		if s != nil {
			return *s
		}
		return nil
	}, krt.WithName("GlooSettingsSingleton"))

	// Hardcode for now
	glooMtls := &deployer.GlooMtlsInfo{
		Enabled: true,
		TlsCert: &deployer.TlsCertInfo{
			CaCert:  []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURNakNDQWhxZ0F3SUJBZ0lJU05GVDN2MXpoVmN3RFFZSktvWklodmNOQVFFTEJRQXdKREVpTUNBR0ExVUUKQXhNWmMzVndaWEpuYkc5dkxYZGxZbWh2YjJzdFkyVnlkQzFqWVRBZUZ3MHlOVEF4TWpJeU1qTXlNamxhRncwegpOVEF4TWpBeU1qTXlNamxhTUNReElqQWdCZ05WQkFNVEdYTjFjR1Z5WjJ4dmJ5MTNaV0pvYjI5ckxXTmxjblF0ClkyRXdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEQjU2eW9zRFNnamxtam5jakQKRm53UkJwL0hKdHBGYzJTOW14bFF6QVRlc1Mzd1JhbS8wV1lXNjZkVGpPQUVSblRlZFc5dW94bisyTC94clBRaApWWmJXQ0tzTGJUZW1GWmc4bXpLeGhuM1BFditOL0xoWDdSQXAxUzEzdFRYMmVnM2hqQUI0U3d2RFFvTHBlbEl6CmRkSU1YeVRKM3FYZkN3c3JPSEFMalpaNXNyMFZTY0UxZXFqZ2VGT0NDTEJXSVV4V1dtWElyaTNZTGhvUDI2M1IKc0t0MkljVi9IaGRqSFdhclUxOFo1WXg4ZkNFZFlmaTNRYWtCL1NCclYvSFRWNVpoNkNFbHVrV3haWjJBZFZjbApLbVg3RDZJdE9MZ3VLbTRNR1o2MVhPa0Z4Ujgvak5Idm45MWJhK2Q4ZHVHNWxkVWZWYm5kVlV0ZEF1eFZUTHNXCm0vdmpBZ01CQUFHamFEQm1NQTRHQTFVZER3RUIvd1FFQXdJQ3BEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01CMEcKQTFVZERnUVdCQlRhMlJvcEhKT0NEZFR1Y2VhR3NRVHZHTUFhcHpBa0JnTlZIUkVFSFRBYmdobHpkWEJsY21kcwpiMjh0ZDJWaWFHOXZheTFqWlhKMExXTmhNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUErbndzNE8zNGRIU1YvCmFqNW9icDRFa0U1S0ZVdmh2dDBHZmg3amRxaWl6ajd6OFk4SUd0cmxtdFRVdVlCalJZY3ZhNnp5M3VWTW9Bb3kKeHZ1ZkwyWFN6aGg3c2VsaFlWaUlXdUlna0ozOWl2SUYwU2ZzYnBjWmR3WW4rVlEyNG0zWnN5Rjd2V2o0SGdKUwpyazNOZ3RVVkFsTzM0UDEvL2xieEpHRFNYTFoyMjFYajBaWVdTN25CNDluU01WSTZDSkhWUzlLUWxFY2N3Q0pXCjc2dmdCc3hVQXk1QVd2NmRyOWtuNWhtRWlBWTFRS3FmaEF4aVpEWDdLWklBNVprS1VXMS90dFpFYnphU3BhYVIKTnp2ZVR5YzdQcmNrVDNYMTRyL29DS2lWTWFWSUYvelVtVjV0d1hGMDQrTFlIaU5BcGFTSERmejRvTlJIQTFwNApKMTNLNlFVQwotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
			TlsCert: []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURpRENDQW5DZ0F3SUJBZ0lJWVpCc2FmRTJmN2d3RFFZSktvWklodmNOQVFFTEJRQXdKREVpTUNBR0ExVUUKQXhNWmMzVndaWEpuYkc5dkxYZGxZbWh2YjJzdFkyVnlkQzFqWVRBZUZ3MHlOVEF4TWpJeU1qTXlNamxhRncweQpOakF4TWpJeU1qTXlNamxhTURFeEVEQU9CZ05WQkFvVEIzTnZiRzh1YVc4eEhUQWJCZ05WQkFNVEZHZHNiMjh1CloyeHZieTF6ZVhOMFpXMHVjM1pqTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEKeHBXdVZVbkYybm5SZWhRM1lXMUpXWTBoYzd4RFNUaWpUbGdlSjJITDRBNkV4aEpxSG1EWkRaNEt1OTZJYlp3WgpnSERRRitUS1UrcE5DRThqN3VnNVAybFc0Y2EzQTQzQkFsSlpvZ3pPUEpvYW1yeVBhTEN1NCtOL1hlU1VsM2FICmNEaS9OYzF5YkkzNmRDN29JcTRBaUl1VjQ4M3U5ZG56QWRVb05tc3BQbmtndURPUzgrUFN5SFNpbmx6S0YzZ1UKZk4yYVhBRndXczNnMGk5M1piRHZlNzIyZ2FYTEQrSm9EMWdUbEhORmJRN1FDL1hzZU1hZjdMYllmend5K3h5Qgo0ZWxZYTlDQmdrR0VWbXRJcGxnRU8vdU5lNXhHVCtlLzVPczBXSW9zc09PMzJtK0kvcitDUkZ3ZVN4NnBiWWhrCk9vZjVRdWVDUDREclA2SDJROUtuS1FJREFRQUJvNEd3TUlHdE1BNEdBMVVkRHdFQi93UUVBd0lGb0RBZEJnTlYKSFNVRUZqQVVCZ2dyQmdFRkJRY0RBUVlJS3dZQkJRVUhBd0l3SHdZRFZSMGpCQmd3Rm9BVTJ0a2FLUnlUZ2czVQo3bkhtaHJFRTd4akFHcWN3V3dZRFZSMFJCRlF3VW9JRVoyeHZiNElRWjJ4dmJ5NW5iRzl2TFhONWMzUmxiWUlVCloyeHZieTVuYkc5dkxYTjVjM1JsYlM1emRtT0NJbWRzYjI4dVoyeHZieTF6ZVhOMFpXMHVjM1pqTG1Oc2RYTjAKWlhJdWJHOWpZV3d3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUhrMC9sc1lyc0tCcFQ5WFprQXdxV0lHeVM4LwpmVnhLR3JLbS9WSnR4akFFRk8vVWJFRWV4a3Q2N3FOa3lnWFFVWGpIa3FQTUxrSjBucFZ3bm05Yi9hL3c5K3BoCkprdWZSdHdYb2JleFRRMzRKb05yUUZQeVg0U3EvTnU4K2dlT1FJSERGYjI5bWE3OW5xTU9Oa2lzaThROFFsVmQKQVA0ckR3TGlFTFA1ajRHWXdaMmltTVp5MkxNY0tnYUFyVTBtTnJ5Wkk2N0xiNm9vVjVrNmpOaUpjTWtKSzY1WgpGYTBCZFVHZGNOUWo2UnNRL2Y5a25LM1dOVTkwYW1mcmNPbG1paGI3cDdwWjFMOXZId3hPRlAxQ1lVak9yZXRjCmUvUFE4amFwVU9RTytVVVRjL3EyYmtMQlhMa1hNY0hLMGdCQm9xTk05VmtXSWVrN29oZlBxb3AzYmRzPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
			TlsKey:  []byte("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBeHBXdVZVbkYybm5SZWhRM1lXMUpXWTBoYzd4RFNUaWpUbGdlSjJITDRBNkV4aEpxCkhtRFpEWjRLdTk2SWJad1pnSERRRitUS1UrcE5DRThqN3VnNVAybFc0Y2EzQTQzQkFsSlpvZ3pPUEpvYW1yeVAKYUxDdTQrTi9YZVNVbDNhSGNEaS9OYzF5YkkzNmRDN29JcTRBaUl1VjQ4M3U5ZG56QWRVb05tc3BQbmtndURPUwo4K1BTeUhTaW5sektGM2dVZk4yYVhBRndXczNnMGk5M1piRHZlNzIyZ2FYTEQrSm9EMWdUbEhORmJRN1FDL1hzCmVNYWY3TGJZZnp3eSt4eUI0ZWxZYTlDQmdrR0VWbXRJcGxnRU8vdU5lNXhHVCtlLzVPczBXSW9zc09PMzJtK0kKL3IrQ1JGd2VTeDZwYlloa09vZjVRdWVDUDREclA2SDJROUtuS1FJREFRQUJBb0lCQVFDY1pMTklMZkpvN2psQQpHSDNJOThXbGhoVkxUWC84UVdPelJvaHc0WDhyZEtPeVRqeE9zbDBlY1ZIb3hRZlNzdllPaGtvTUZ6NFV1bGh5CmE4bFQxSVdKWUE5eGZnc1IvR2g5eUpjWW5WY3F1UDZzMEVWczRJRExycFhYUHphYTFsa3gweThiVXpLRE5ZbUMKU0pLL0JTUWNaVG5sajRCYkZJdGg3UnBmU1NQUk16UnZHRFcrUlNlNzkzdHVIRjc2dVRvS29Mdkczc0dFOTh6ZQovNmhXNTFiTW1PbXkwbTFxU2lKUUFyYlR1ams3b2l3NERiYjU0T2QyaGJtaDhwMVFqMUZUUnhDVkp6RStQY1RyCnNWRExucWZJZHhSYWI4UnhTT1podGZZTUJTeHE4cUtuQzR4cHhkTllYeUpqYllzZFNqei9McDBpOW1CNjNFNTMKSUJYeG5aMkZBb0dCQU1tOGd3NmJqM3MzSzhVUUovMmdMdHVZT2xaTUg2eVVaZ3gzMkZWemVWcmJkeDE2Y1hFVApiS3puNHhad3BLc0NlUk5sOEFuWk90S0JUVk1NejdvellYd0hqMXpmamoxcURmOGlHSUYxM01qM3EyaTkwam5kCklrcU9pK0JlaUtDT1lvTE5qYjhlMktMOUdTOVY4SW1JL3I5YUVSK1BvU1BYbTM0T0c3YlRQOFRYQW9HQkFQd0EKSlRZamM1WEQ0KzE3Nmk5ak5rQmVXTjB2RXBnTFRhQnRnTTRGTko3MVF2Mko2TmdxTHUzZWk0U1Z5WkdwcE91UQpYM28zcWlhaUFEd2E3dWREWERMRnlsZUxjays3eFpKSUZuWlJLbHlnakxzUHJaZ0Ftd01zMndWLzdMa2kzSS8wCklyWnJnU3EvcEp5WmFTZmIrdDFGRXh1OVhheHZhK3VHK082d0IzUC9Bb0dCQUlJYzJMWnFOSzkyMVA4anZYZEMKZ2haVjU0SlRWTFo5ZkJnY2orWUZOVWNaZDRrR3VQWUNYanhpeno5ZzVZUDZjMWJFajMzNm9vcTBwTWVrNHJHbwpnLzUzN2NvcjBkVGdleWlMdUJ4L2hTZ0ZQWU92c2xCcHhMMHJsU0hnTnVTL0VPQm1iVDdRU1U5T1NKa1VKN0M5Ckwva0F3VHlHNlpweGJETndMQVhOMkRvL0FvR0JBTVJVaktsQU40WWdCd3ozOUwwVXE2aThtTGxDT2xkUXZ4clYKRlh2dEhGRVh2aWh2OElPeFlieWJITkdnTFZtWjlNNCtQZFNuVjU0ZnF0VXBHcVg4bWZGSW5kdFUzaXQybkhmYQpSLzNJUUp2SHpielRleWlvbUJ5Q0x1VjdCQUE5UTkrM2tlL1RrOSt0VFY5Z09rZituOVVTUXMvaTJmOUZFNng5CkRLWlJhSTBiQW9HQUx4Z0Y5RFZDbzVQMFB5MzBsM3cvR0dGTit3RkpONGx0RWtVcFo4NTI3b29TZ0kyZy9Mb1kKRGg0QThFNXNXbDVRWXptdXVUQ29WSkJNUUl0S1VWMlhXZjJySlN5ejZuWXY5ZFh4eVpqaVhyZW5nNENtaTU1YQo3MC9RQ29mQnp1aDF0dGVMMUIxZ2hsOXA5cXBNS3E5cFREZWtQQ1V2bktObkJWam1HbXYvUlVvPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="),
		},
	}

	logger.Info("Hardcoding glooMtls for now", "glooMtls", glooMtls)

	serviceClient := kclient.New[*corev1.Service](kubeClient)
	services := krt.WrapClient(serviceClient, krt.WithName("Services"))

	logger.Info("creating reporter")
	kubeGwStatusReporter := NewGenericStatusReporter(kubeClient, defaults.KubeGatewayReporter)

	glooReporter := NewGenericStatusReporter(kubeClient, defaults.GlooReporter)
	pluginOpts := registry.PluginOpts{
		Ctx:                     ctx,
		SidecarOnGatewayEnabled: envutils.IsEnvTruthy(constants.IstioInjectionEnabled),
		SvcCollection:           services,
	}
	logger.Info("initializing controller")
	c, err := controller.NewControllerBuilder(ctx, controller.StartConfig{
		ExtensionsFactory:    extensionsFactory,
		RestConfig:           restConfig,
		SetupOpts:            setupOpts,
		KubeGwStatusReporter: kubeGwStatusReporter,
		Translator:           setup.TranslatorFactory{PluginRegistry: pluginRegistryFactory(pluginOpts)},
		GlooStatusReporter:   glooReporter,
		Client:               kubeClient,
		AugmentedPods:        augmentedPods,
		UniqueClients:        ucc,

		InitialSettings: initialSettings,
		Settings:        settingsSingle,
		// Dev flag may be useful for development purposes; not currently tied to any user-facing API
		Dev:      false,
		GlooMtls: glooMtls,
		Debugger: setupOpts.KrtDebugger,
	})
	if err != nil {
		logger.Error("failed initializing controller: ", err)
		return err
	}

	/// no collections after this point

	logger.Info("waiting for cache sync")
	kubeClient.RunAndWait(ctx.Done())
	setting.Synced().WaitUntilSynced(ctx.Done())

	logger.Info("starting controller")
	return c.Start(ctx)
}

type genericStatusReporter struct {
	client               istiokube.Client
	kubeGwStatusReporter reporter.StatusReporter
	statusClient         resources.StatusClient
}

func NewGenericStatusReporter(client istiokube.Client, r string) reporter.StatusReporter {
	statusReporterNamespace := gloostatusutils.GetStatusReporterNamespaceOrDefault("gloo-system")
	statusClient := gloostatusutils.GetStatusClientForNamespace(statusReporterNamespace)

	kubeGwStatusReporter := reporter.NewReporter(
		r,
		statusClient,
	)
	return &genericStatusReporter{client: client, kubeGwStatusReporter: kubeGwStatusReporter, statusClient: statusClient}
}

// StatusFromReport implements reporter.StatusReporter.
func (g *genericStatusReporter) StatusFromReport(report reporter.Report, subresourceStatuses map[string]*core.Status) *core.Status {
	return g.kubeGwStatusReporter.StatusFromReport(report, subresourceStatuses)
}

// WriteReports implements reporter.StatusReporter.
func (g *genericStatusReporter) WriteReports(ctx context.Context, resourceErrs reporter.ResourceReports, subresourceStatuses map[string]*core.Status) error {
	ctx = contextutils.WithLogger(ctx, "reporter")
	logger := contextutils.LoggerFrom(ctx)

	var merr error

	// copy the map so we can iterate over the copy, deleting resources from
	// the original map if they are not found/no longer exist.
	resourceErrsCopy := make(reporter.ResourceReports, len(resourceErrs))
	for resource, report := range resourceErrs {
		resourceErrsCopy[resource] = report
	}

	for resource, report := range resourceErrsCopy {

		// check if resource is an internal upstream. if so skip it..
		if kubernetes.IsFakeKubeUpstream(resource.GetMetadata().GetName()) {
			continue
		}
		// check if resource is an internal upstream. Internal upstreams have ':' in their names so
		// the cannot be written to the cluster. if so skip it..
		if strings.IndexRune(resource.GetMetadata().GetName(), ':') >= 0 {
			continue
		}

		status := g.StatusFromReport(report, subresourceStatuses)
		status = trimStatus(status)

		resourceStatus := g.statusClient.GetStatus(resource)

		if status.Equal(resourceStatus) {
			// TODO: find a way to log this but it is noisy currently due to once per second status sync
			// see: projects/gateway2/proxy_syncer/kube_gw_translator_syncer.go#syncStatus(...)
			// and its call site in projects/gateway2/proxy_syncer/proxy_syncer.go
			// logger.Debugf("skipping report for %v as it has not changed", resource.GetMetadata().Ref())
			continue
		}

		resourceToWrite := resources.Clone(resource).(resources.InputResource)
		g.statusClient.SetStatus(resourceToWrite, status)
		writeErr := g.attemptUpdateStatus(ctx, resourceToWrite, status)

		if k8serrors.IsNotFound(writeErr) {
			logger.Debugf("did not write report for %v : %v because resource was not found", resourceToWrite.GetMetadata().Ref(), status)
			delete(resourceErrs, resource)
			continue
		}

		if writeErr != nil {
			err := fmt.Errorf("failed to write status %v for resource %v: %w", status, resource.GetMetadata().GetName(), writeErr)
			logger.Warn(err)
			merr = errors.Join(merr, err)
			continue
		}
		logger.Debugf("wrote report for %v : %v", resource.GetMetadata().Ref(), status)

	}
	return merr
}

func (g *genericStatusReporter) attemptUpdateStatus(ctx context.Context, resourceToWrite resources.InputResource, statusToWrite *core.Status) error {
	key := resources.Kind(resourceToWrite)
	crd, ok := kindToCrd[key]
	if !ok {
		err := fmt.Errorf("no crd found for kind %v", key)
		contextutils.LoggerFrom(ctx).DPanic(err)
		return err
	}
	ns := resourceToWrite.GetMetadata().GetNamespace()
	name := resourceToWrite.GetMetadata().GetName()

	data, err := shared.GetJsonPatchData(ctx, resourceToWrite)
	if err != nil {
		return fmt.Errorf("error getting status json patch data: %w", err)
	}

	_, err = g.client.Dynamic().Resource(crd.GroupVersion().WithResource(crd.CrdMeta.Plural)).Namespace(ns).Patch(ctx, name, types.JSONPatchType, data, metav1.PatchOptions{})
	return err
}

var _ reporter.StatusReporter = &genericStatusReporter{}

var kindToCrd = map[string]crd.Crd{}

func add(crd crd.Crd, resourceType resources.InputResource) {
	skKind := resources.Kind(resourceType)
	kindToCrd[skKind] = crd
}

func init() {
	add(gateway.RouteOptionCrd, new(gateway.RouteOption))
	add(gateway.VirtualHostOptionCrd, new(gateway.VirtualHostOption))
	add(gloov1.ProxyCrd, new(gloov1.Proxy))
	add(gloov1.UpstreamCrd, new(gloov1.Upstream))
	// add(rlv1alpha1.RateLimitCrd, new(rlv1alpha1.RateLimit))
	// add(rlv1alpha1.RateLimitCrd, new(rlv1alpha1.RateLimit))
}

func trimStatus(status *core.Status) *core.Status {
	// truncate status reason to a kilobyte, with max 100 keys in subresource statuses
	return trimStatusForMaxSize(status, reporter.MaxStatusBytes, reporter.MaxStatusKeys)
}

func trimStatusForMaxSize(status *core.Status, bytesPerKey, maxKeys int) *core.Status {
	if status == nil {
		return nil
	}
	if len(status.GetReason()) > bytesPerKey {
		status.Reason = status.GetReason()[:bytesPerKey]
	}

	if len(status.GetSubresourceStatuses()) > maxKeys {
		// sort for idempotency
		keys := make([]string, 0, len(status.GetSubresourceStatuses()))
		for key := range status.GetSubresourceStatuses() {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		trimmedSubresourceStatuses := make(map[string]*core.Status, maxKeys)
		for _, key := range keys[:maxKeys] {
			trimmedSubresourceStatuses[key] = status.GetSubresourceStatuses()[key]
		}
		status.SubresourceStatuses = trimmedSubresourceStatuses
	}

	for key, childStatus := range status.GetSubresourceStatuses() {
		// divide by two so total memory usage is bounded at: (num_keys * bytes_per_key) + (num_keys / 2 * bytes_per_key / 2) + ...
		// 100 * 1024b + 50 * 512b + 25 * 256b + 12 * 128b + 6 * 64b + 3 * 32b + 1 * 16b ~= 136 kilobytes
		//
		// 2147483647 bytes is k8s -> etcd limit in grpc connection. 2147483647 / 136 ~= 15788 resources at limit before we see an issue
		// https://github.com/solo-io/solo-projects/issues/4120
		status.GetSubresourceStatuses()[key] = trimStatusForMaxSize(childStatus, bytesPerKey/2, maxKeys/2)
	}
	return status
}
