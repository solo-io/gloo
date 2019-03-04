package debug

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-projects/pkg/license/db"
	"github.com/solo-io/solo-projects/pkg/license/model"
)

var _ db.KeyDb = new(DebugKeyDb)

type DebugKeyDb struct{}

func (d *DebugKeyDb) Save(ctx context.Context, r *model.Request, ua model.UserInfo, key string) error {
	fmt.Println(ua.Email, key)
	return nil
}
