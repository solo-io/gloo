package translator

import (
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"context"
)

type Translator interface {
	Translate(ctx context.Context, proxy *v1.Proxy, snap *v1.Snapshot) (envoycache.Snapshot, reporter.ResourceErrors, error)
}

