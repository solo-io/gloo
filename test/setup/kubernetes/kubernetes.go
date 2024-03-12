package kubernetes

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	glooinstancev1 "github.com/solo-io/solo-apis/pkg/api/fed.solo.io/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func getExternalAddress(
	ctx context.Context,
	namespace, portName, serviceSelector string,
	isKindCluster bool,
	client kubernetes.Interface,
) (string, error) {
	services, err := client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{LabelSelector: serviceSelector})
	if err != nil {
		return "", err
	}

	if len(services.Items) == 0 {
		return "", fmt.Errorf("no services found for given selector %q in namespace %s", serviceSelector, namespace)
	}

	var (
		nodePort int32
		nodeAddr string
	)

	for _, svc := range services.Items {
		for _, port := range svc.Spec.Ports {
			if port.Name != portName {
				continue
			}

			if isKindCluster {
				nodePort = port.NodePort
				break
			}

			for _, ingress := range svc.Status.LoadBalancer.Ingress {
				return fmt.Sprintf("%s:%d", ingress.IP, port.Port), nil
			}
		}
	}

	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP && isKindCluster {
				nodeAddr = addr.Address
				break
			}
		}
	}

	return fmt.Sprintf("%s:%d", nodeAddr, nodePort), nil
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
	copy := *serviceAccount

	copy.Name = versionName(copy.Name, version)

	return &copy
}

func versionService(version string, service *corev1.Service) *corev1.Service {
	if service == nil {
		return nil
	}
	copy := *service

	copy.Name = versionName(copy.Name, version)

	copy.Spec.Selector = versionedLabels(version, copy.Spec.Selector)

	return &copy
}

func versionDeployment(version string, deployment *appsv1.Deployment) *appsv1.Deployment {
	if deployment == nil {
		return nil
	}

	copy := *deployment

	name := versionName(copy.Name, version)

	copy.Name = name
	copy.Labels = versionedLabels(version, copy.Labels)
	copy.Spec.Template.Labels = versionedLabels(version, copy.Spec.Template.Labels)
	copy.Spec.Selector.MatchLabels = versionedLabels(version, copy.Spec.Selector.MatchLabels)
	copy.Spec.Template.Spec.ServiceAccountName = name

	if _, ok := copy.Spec.Template.ObjectMeta.Labels["security.policy.gloo.solo.io/service_account"]; ok {
		copy.Spec.Template.ObjectMeta.Labels["security.policy.gloo.solo.io/service_account"] = name
	}

	// Update the SERVICE_VERSION environment variable.
	for i, container := range copy.Spec.Template.Spec.Containers {
		for j, env := range container.Env {
			if env.Name == "SERVICE_VERSION" {
				copy.Spec.Template.Spec.Containers[i].Env[j].Value = version
			}
		}
	}
	return &copy
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

func WaitUntilPodsRunning(
	ctx context.Context,
	kubeClient kubernetes.Interface,
	timeout time.Duration,
	namespace string,
	podPrefixes ...string,
) error {
	pods := kubeClient.CoreV1().Pods(namespace)
	podsWithPrefixReady := func(prefix string) (bool, error) {
		list, err := pods.List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		var podsWithPrefix []corev1.Pod
		for _, pod := range list.Items {
			if strings.HasPrefix(pod.Name, prefix) {
				podsWithPrefix = append(podsWithPrefix, pod)
			}
		}
		if len(podsWithPrefix) == 0 {
			return false, fmt.Errorf("no pods found with prefix %v", prefix)
		}
		for _, pod := range podsWithPrefix {
			var podReady bool
			for _, cond := range pod.Status.Conditions {
				if cond.Type == corev1.ContainersReady && cond.Status == corev1.ConditionTrue {
					podReady = true
					break
				}
			}
			if !podReady {
				return false, nil
			}
		}
		return true, nil
	}
	failed := time.After(timeout)
	notYetRunning := make(map[string]struct{})
	for {
		select {
		case <-failed:
			return fmt.Errorf("timed out waiting for pods to come online: %v", notYetRunning)
		case <-time.After(time.Second / 2):
			notYetRunning = make(map[string]struct{})
			for _, prefix := range podPrefixes {
				ready, err := podsWithPrefixReady(prefix)
				if err != nil {
					fmt.Printf("failed to get pod status: %v\n", err)
					notYetRunning[prefix] = struct{}{}
				}
				if !ready {
					notYetRunning[prefix] = struct{}{}
				}
			}
			if len(notYetRunning) == 0 {
				return nil
			}
		}

	}
}
