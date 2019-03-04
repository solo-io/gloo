package keys

import (
	"context"
	"time"

	"github.com/solo-io/solo-projects/pkg/license/model"
)

type KeyManager interface {
	KeyGenerator() KeyGenerator
	KeyValidator() KeyValidator
}

type KeyGenerator interface {
	GenerateKey(ctx context.Context, expires time.Time) (string, error)
}

type KeyValidator interface {
	ValidateKey(context.Context, string) (*model.KeyInfo, error)
}

func IsExperied(k *model.KeyInfo) bool {
	return time.Now().After(k.Expiration)
}
