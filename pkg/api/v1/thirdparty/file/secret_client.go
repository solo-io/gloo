package file

import (
	"sort"

	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"github.com/radovskyb/watcher"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/thirdparty"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
)

type secretClient struct {
	dir string
}

func NewSecretClient(dir string) thirdparty.ThirdPartyResourceClient {
	return &secretClient{
		dir: dir,
	}
}

func readFileIntoSecret(file string) (*thirdparty.Secret, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var secret thirdparty.Secret
	err = yaml.Unmarshal(b, &secret)
	return &secret, err
}

func (rc *secretClient) Read(namespace, name string, opts clients.ReadOpts) (thirdparty.ThirdPartyResource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	path := rc.filename(namespace, name)
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil, errors.NewNotExistErr(namespace, name, err)
	}

	return readFileIntoSecret(path)
}

func (rc *secretClient) Write(resource thirdparty.ThirdPartyResource, opts clients.WriteOpts) (thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	if err := resources.ValidateName(resource.GetMetadata().Name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	meta.Namespace = clients.DefaultNamespaceIfEmpty(meta.Namespace)

	original, err := rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{})
	if original != nil && err == nil {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(meta)
		}
		if meta.ResourceVersion != original.GetMetadata().ResourceVersion {
			return nil, errors.Errorf("resource version error. must update new resource version to match current")
		}
	}

	if rc.exist(meta.Namespace, meta.Name) && !opts.OverwriteExisting {
		return nil, errors.NewExistErr(meta)
	}

	meta.ResourceVersion = newOrIncrementResourceVer(meta.ResourceVersion)
	clone := &thirdparty.Secret{
		Data: thirdparty.Data{
			Metadata: meta,
			Values:   resource.GetData(),
		},
	}

	path := rc.filename(meta.Namespace, meta.Name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil && !os.IsExist(err) {
		return nil, errors.Wrapf(err, "creating directory")
	}
	if err := writeThirdPartyResource(path, clone); err != nil {
		return nil, errors.Wrapf(err, "writing resource to disk")
	}

	// return a read object to update the resource version
	return rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *secretClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	path := rc.filename(namespace, name)
	return os.Remove(path)
}

func (rc *secretClient) List(namespace string, opts clients.ListOpts) ([]thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	namespaceDir := filepath.Join(rc.dir, namespace)
	files, err := ioutil.ReadDir(namespaceDir)
	if err != nil {
		return nil, errors.Wrapf(err, "reading namespace dir")
	}

	var resourceList []thirdparty.ThirdPartyResource
	for _, file := range files {
		path := filepath.Join(namespaceDir, file.Name())
		resource, err := readFileIntoSecret(path)
		if err != nil {
			return nil, errors.Wrapf(err, "reading file as secret")
		}
		if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(resource.GetMetadata().Labels)) {
			resourceList = append(resourceList, resource)
		}
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *secretClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []thirdparty.ThirdPartyResource, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	dir := filepath.Join(rc.dir, namespace)
	events, errs, err := rc.events(opts.Ctx, dir, opts.RefreshRate)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting watch on namespace dir")
	}
	resourcesChan := make(chan []thirdparty.ThirdPartyResource)
	go func() {
		// watch should open up with an initial read
		list, err := rc.List(namespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}()
	go func() {
		for {
			select {
			case <-events:
				list, err := rc.List(namespace, clients.ListOpts{
					Ctx:      opts.Ctx,
					Selector: opts.Selector,
				})
				if err != nil {
					errs <- err
					continue
				}
				resourcesChan <- list
			}
		}
	}()

	return resourcesChan, errs, nil
}

func (rc *secretClient) exist(namespace, name string) bool {
	_, err := os.Stat(rc.filename(namespace, name))
	return err == nil
}

func (rc *secretClient) filename(namespace, name string) string {
	return filepath.Join(rc.dir, namespace, name) + ".yaml"
}

func (rc *secretClient) events(ctx context.Context, dir string, refreshRate time.Duration) (<-chan struct{}, chan error, error) {
	events := make(chan struct{})
	errs := make(chan error)
	w := watcher.New()
	w.SetMaxEvents(0)
	w.FilterOps(watcher.Create, watcher.Write, watcher.Remove, watcher.Rename, watcher.Move)
	if err := w.AddRecursive(dir); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to watch directory %v", dir)
	}
	go func() {
		if err := w.Start(refreshRate); err != nil {
			errs <- err
		}
	}()
	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.IsDir() {
					continue
				}
				events <- struct{}{}
			case err := <-w.Error:
				errs <- errors.Wrapf(err, "file watcher error")
			case <-ctx.Done():
				w.Close()
				return
			}
		}
	}()
	return events, errs, nil
}
