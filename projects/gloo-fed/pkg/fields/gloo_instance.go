package fields

import (
	"context"

	"github.com/rotisserie/eris"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func AddGlooInstanceIndexer(ctx context.Context, mgr manager.Manager) error {
	informer, err := mgr.GetCache().GetInformer(ctx, &fedv1.GlooInstance{})
	if err != nil {
		return err
	}

	if err = informer.AddIndexers(map[string]cache.IndexFunc{
		FieldIndexName(ClusterIndex): func(obj interface{}) ([]string, error) {
			instance, ok := obj.(*fedv1.GlooInstance)
			if !ok {
				return nil, eris.New("Wrong type passed into indexer")
			}

			return []string{KeyToNamespacedKey("", instance.Spec.GetCluster())}, nil
		},
	}); err != nil {
		return err
	}

	return nil
}
