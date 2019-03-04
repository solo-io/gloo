package notify

import (
	"context"

	"github.com/solo-io/solo-projects/pkg/license/model"
)

type Notifier interface {
	Notify(ctx context.Context, ua model.UserInfo, key string) error
}
