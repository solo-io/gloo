package debug

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-projects/pkg/license/model"
	"github.com/solo-io/solo-projects/pkg/license/notify"
)

var _ notify.Notifier = new(DebugNotifier)

type DebugNotifier struct{}

func (d *DebugNotifier) Notify(ctx context.Context, ua model.UserInfo, key string) error {
	fmt.Println(ua.Email, key)
	return nil
}
