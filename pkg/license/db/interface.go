package db

import (
	"context"

	"github.com/solo-io/solo-projects/pkg/license/model"
)

type KeyDb interface {
	Save(ctx context.Context, r *model.Request, ua model.UserInfo, key string) error
}
