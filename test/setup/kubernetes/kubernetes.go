package kubernetes

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	glooinstancev1 "github.com/solo-io/solo-apis/pkg/api/fed.solo.io/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/cmd/rollout"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
)

func NewClient(ctx context.Context, dir, name string) (string, string, kubernetes.Interface, client.Client, error) {
	baseCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", "", nil, nil, errors.Wrap(err, "failed to load kube config")
	}

	kubeContext, kubeConfig, config, err := splitKubeconfig(dir, name, baseCfg)
	if err != nil {
		return "", "", nil, nil, errors.Wrap(err, "failed to split kube config")
	}

	scheme := runtime.NewScheme()

	// k8s resources
	if err = corev1.AddToScheme(scheme); err != nil {
		return "", "", nil, nil, err
	}
	if err = appsv1.AddToScheme(scheme); err != nil {
		return "", "", nil, nil, err
	}
	// k8s gateway resources
	if err = v1alpha2.AddToScheme(scheme); err != nil {
		return "", "", nil, nil, err
	}
	if err = v1beta1.AddToScheme(scheme); err != nil {
		return "", "", nil, nil, err
	}
	if err = v1.AddToScheme(scheme); err != nil {
		return "", "", nil, nil, err
	}
	// gloo resources
	if err = glooinstancev1.AddToScheme(scheme); err != nil {
		return "", "", nil, nil, err
	}

	ctrllog.SetLogger(ctrlzap.New(ctrlzap.WriteTo(io.Discard)))

	ctrl, err := manager.New(config, manager.Options{
		Metrics: metricsserver.Options{BindAddress: "0"},
		Scheme:  scheme,
		Logger:  ctrlzap.New(ctrlzap.WriteTo(io.Discard)),
	})
	if err != nil {
		return "", "", nil, nil, errors.Wrap(err, "failed to create controller manager")
	}

	go func() {
		if err := ctrl.Start(ctx); err != nil {
			contextutils.LoggerFrom(ctx).Error(err)
		}
	}()

	return kubeContext, kubeConfig, kubernetes.NewForConfigOrDie(config), ctrl.GetClient(), nil
}

func splitKubeconfig(dir, name string, base *api.Config) (string, string, *rest.Config, error) {
	var context, cluster string
	for k := range base.Clusters {
		if !strings.Contains(k, name) {
			continue
		}
		cluster = k
	}

	for k := range base.Contexts {
		if !strings.Contains(k, name) {
			continue
		}
		context = k
	}

	if context == "" {
		return "", "", nil, errors.New(fmt.Sprintf("failed to find context for %s", name))
	}

	authInfo := base.Contexts[context].AuthInfo

	config := api.Config{
		CurrentContext: context,
		Clusters:       map[string]*api.Cluster{cluster: base.Clusters[cluster]},
		Contexts:       map[string]*api.Context{context: base.Contexts[context]},
		AuthInfos:      map[string]*api.AuthInfo{authInfo: base.AuthInfos[authInfo]},
	}

	if len(config.Clusters) == 0 || len(config.Contexts) == 0 || len(config.AuthInfos) == 0 {
		return "", "", nil, errors.New(fmt.Sprintf("failed to find cluster, context, or authInfo for %s", name))
	}

	path := filepath.Join(dir, name)

	if err := clientcmd.WriteToFile(config, path); err != nil {
		return "", "", nil, errors.Wrap(err, "failed to encode client config")
	}

	rest, err := clientcmd.NewNonInteractiveClientConfig(config, config.CurrentContext, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to get client config")
	}

	return context, path, rest, nil

}

func RolloutStatus(namespace string, cluster *Cluster) error {
	kubeConfig := cluster.GetKubeConfig()

	flags := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeConfig,
		Namespace:  &namespace,
	}
	cmd := rollout.NewCmdRolloutStatus(util.NewFactory(flags), genericiooptions.IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stdout,
	})

	cmd.SetArgs([]string{"deployment,daemonset,statefulset"})

	return cmd.Execute()
}

type TypedObject[T client.Object] interface {
	DeepCopy() T
}

func CreateOrUpdate[T client.Object](ctx context.Context, object client.Object, cluster client.Client) error {
	if object == nil {
		return nil
	}

	if err := cluster.Create(ctx, object); err == nil && !k8errors.IsAlreadyExists(err) {
		return nil
	}

	typedObj, ok := object.(TypedObject[T])
	if !ok {
		// We can't update it
		return nil
	}

	key := types.NamespacedName{
		Name:      object.GetName(),
		Namespace: object.GetNamespace(),
	}

	existing := typedObj.DeepCopy()

	if err := cluster.Get(ctx, key, existing); err != nil {
		return err
	}

	// We'll run into resource version errors if we don't apply existing version to our new object
	object.SetResourceVersion(existing.GetResourceVersion())

	if err := cluster.Update(ctx, object); err != nil && !k8errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func versionServiceAccount(version string, serviceAccount *corev1.ServiceAccount) *corev1.ServiceAccount {
	if serviceAccount == nil {
		return nil
	}
	copySA := *serviceAccount

	copySA.Name = versionName(copySA.Name, version)

	return &copySA
}

func versionService(version string, service *corev1.Service) *corev1.Service {
	if service == nil {
		return nil
	}
	copyService := *service

	copyService.Name = versionName(copyService.Name, version)

	copyService.Spec.Selector = versionedLabels(version, copyService.Spec.Selector)

	return &copyService
}

func versionDeployment(version string, deployment *appsv1.Deployment) *appsv1.Deployment {
	if deployment == nil {
		return nil
	}

	copyDeployment := *deployment

	name := versionName(copyDeployment.Name, version)

	copyDeployment.Name = name
	copyDeployment.Labels = versionedLabels(version, copyDeployment.Labels)
	copyDeployment.Spec.Template.Labels = versionedLabels(version, copyDeployment.Spec.Template.Labels)
	copyDeployment.Spec.Selector.MatchLabels = versionedLabels(version, copyDeployment.Spec.Selector.MatchLabels)
	copyDeployment.Spec.Template.Spec.ServiceAccountName = name

	if _, ok := copyDeployment.Spec.Template.ObjectMeta.Labels["security.policy.gloo.solo.io/service_account"]; ok {
		copyDeployment.Spec.Template.ObjectMeta.Labels["security.policy.gloo.solo.io/service_account"] = name
	}

	// Update the SERVICE_VERSION environment variable.
	for i, container := range copyDeployment.Spec.Template.Spec.Containers {
		for j, env := range container.Env {
			if env.Name == "SERVICE_VERSION" {
				copyDeployment.Spec.Template.Spec.Containers[i].Env[j].Value = version
			}
		}
	}
	return &copyDeployment
}

func versionName(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}

func versionedLabels(version string, labels map[string]string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}

	labels["version"] = version

	return labels
}
