package kube

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

var (
	MLists    = stats.Int64("kube/lists", "The number of lists", "1")
	MWatches  = stats.Int64("kube/watches", "The number of watches", "1")
	MWatchLen = stats.Float64("kube/watch", "The length of a watch session", "ms")

	KeyNamespace, _ = tag.NewKey("namespace")
	KeyKind, _      = tag.NewKey("kind")

	ListCountView = &view.View{
		Name:        "kube/lists-count",
		Measure:     MLists,
		Description: "The number of list calls",
		Aggregation: view.Count(),
		TagKeys: []tag.Key{
			KeyNamespace,
			KeyKind,
		},
	}
	WatchCountView = &view.View{
		Name:        "kube/watch-count",
		Measure:     MWatches,
		Description: "The number of watch calls",
		Aggregation: view.Count(),
		TagKeys: []tag.Key{
			KeyNamespace,
			KeyKind,
		},
	}
	WatchSeesionView = &view.View{
		Name:        "kube/watch-session",
		Description: "Watch session lengths in buckets",
		Measure:     MWatchLen,
		// [>=0ms, >=25ms, >=50ms, >=75ms, >=100ms, >=200ms, >=400ms, >=600ms, >=800ms, >=1s, >=2s, >=4s, >=6s]
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000, 10000),
		TagKeys: []tag.Key{
			KeyNamespace,
			KeyKind,
		},
	}
)

func init() {
	view.Register(ListCountView, WatchCountView, WatchSeesionView)
}

type ResourceClient struct {
	crd          crd.Crd
	apiexts      apiexts.Interface
	kube         versioned.Interface
	ownerLabel   string
	resourceName string
	resourceType resources.InputResource
}

func NewResourceClient(crd crd.Crd, cfg *rest.Config, resourceType resources.InputResource) (*ResourceClient, error) {
	apiExts, err := apiexts.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "creating api extensions client")
	}
	crdClient, err := versioned.NewForConfig(cfg, crd)
	if err != nil {
		return nil, errors.Wrapf(err, "creating crd client")
	}
	typeof := reflect.TypeOf(resourceType)
	resourceName := strings.Replace(typeof.String(), "*", "", -1)
	resourceName = strings.Replace(resourceName, ".", "", -1)
	return &ResourceClient{
		crd:          crd,
		apiexts:      apiExts,
		kube:         crdClient,
		resourceName: resourceName,
		resourceType: resourceType,
	}, nil
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Kind() string {
	return resources.Kind(rc.resourceType)
}

func (rc *ResourceClient) NewResource() resources.Resource {
	return resources.Clone(rc.resourceType)
}

func (rc *ResourceClient) Register() error {
	return rc.crd.Register(rc.apiexts)
}

func (rc *ResourceClient) Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	resourceCrd, err := rc.kube.ResourcesV1().Resources(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.NewNotExistErr(namespace, name, err)
		}
		return nil, errors.Wrapf(err, "reading resource from kubernetes")
	}
	resource := rc.NewResource()
	if resourceCrd.Spec != nil {
		if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
			return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
		}
	}
	resources.UpdateMetadata(resource, func(meta *core.Metadata) {
		meta.ResourceVersion = resourceCrd.ResourceVersion
	})
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	opts = opts.WithDefaults()
	if err := resources.Validate(resource); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	meta.Namespace = clients.DefaultNamespaceIfEmpty(meta.Namespace)

	// mutate and return clone
	clone := proto.Clone(resource).(resources.InputResource)
	clone.SetMetadata(meta)
	resourceCrd := rc.crd.KubeResource(clone)

	if rc.exist(meta.Namespace, meta.Name) {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(meta)
		}
		if _, updateerr := rc.kube.ResourcesV1().Resources(meta.Namespace).Update(resourceCrd); updateerr != nil {
			original, err := rc.kube.ResourcesV1().Resources(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			if err == nil {
				return nil, errors.Wrapf(updateerr, "updating kube resource %v:%v (want %v)", resourceCrd.Name, resourceCrd.ResourceVersion, original.ResourceVersion)
			}
			return nil, errors.Wrapf(updateerr, "updating kube resource %v", resourceCrd.Name)
		}
	} else {
		if _, err := rc.kube.ResourcesV1().Resources(meta.Namespace).Create(resourceCrd); err != nil {
			return nil, errors.Wrapf(err, "creating kube resource %v", resourceCrd.Name)
		}
	}

	// return a read object to update the resource version
	return rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	if err := rc.kube.ResourcesV1().Resources(namespace).Delete(name, nil); err != nil {
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) (resources.ResourceList, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	resourceCrdList, err := rc.kube.ResourcesV1().Resources(namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "listing resources in %v", namespace)
	}
	var resourceList resources.ResourceList
	for _, resourceCrd := range resourceCrdList.Items {
		resource := rc.NewResource()
		if resourceCrd.Spec != nil {
			if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
				return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
			}
		}
		resources.UpdateMetadata(resource, func(meta *core.Metadata) {
			meta.Namespace = resourceCrd.Namespace
			meta.Name = resourceCrd.Name
			meta.ResourceVersion = resourceCrd.ResourceVersion
		})
		resourceList = append(resourceList, resource)
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan resources.ResourceList, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	resourcesChan := make(chan resources.ResourceList)
	errs := make(chan error)
	ctx := opts.Ctx
	if ctx2, err := tag.New(opts.Ctx, tag.Insert(KeyNamespace, namespace), tag.Insert(KeyKind, rc.resourceName)); err == nil {
		ctx = ctx2
	}

	updateResourceList := func() {
		stats.Record(ctx, MLists.M(1))
		list, err := rc.List(namespace, clients.ListOpts{
			Ctx:      ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}
	// watch should open up with an initial read

	go func() {
		defer close(resourcesChan)
		defer close(errs)

		updateResourceList()
		mxdelay := time.Second * 10
		be := contextutils.NewExponentioalBackoff(contextutils.ExponentioalBackoff{MaxDelay: &mxdelay})

		ctx = contextutils.WithLogger(ctx, "watchloop")

		be.Backoff(ctx, func(ctx context.Context) error {
			startTime := time.Now()
			defer func() {
				ms := float64(time.Since(startTime).Nanoseconds()) / 1e6
				stats.Record(ctx, MWatches.M(1), MWatchLen.M(ms))
			}()

			return rc.watch(ctx, namespace, opts.Selector, opts.RefreshRate, updateResourceList, errs)
		})
	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) watch(ctx context.Context, namespace string, selector map[string]string, refresh time.Duration, updateResourceList func(), errs chan<- error) error {

	logger := contextutils.LoggerFrom(ctx)

	watch, err := rc.kube.ResourcesV1().Resources(namespace).Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(selector).String(),
	})
	if err != nil {
		return errors.Wrapf(err, "initiating kube watch in %v", namespace)
	}
	for {
		select {
		case <-time.After(refresh):
			updateResourceList()
		case event, ok := <-watch.ResultChan():
			if !ok {
				logger.Warnf("watch was closed")
				return errors.Errorf("watch closed")
			}
			switch event.Type {
			case kubewatch.Error:
				// TODO(yuval-k): do we want to select on this channel?
				errs <- errors.Errorf("error during watch: %v", event)
			default:
				logger.Debugf("got event - updating %v", event)
				updateResourceList()
			}
		case <-ctx.Done():
			watch.Stop()
			return ctx.Err()
		}
	}
}

func (rc *ResourceClient) exist(namespace, name string) bool {
	_, err := rc.kube.ResourcesV1().Resources(namespace).Get(name, metav1.GetOptions{})
	return err == nil

}
