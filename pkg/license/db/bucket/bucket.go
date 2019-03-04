package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"gocloud.dev/blob"

	"github.com/solo-io/solo-projects/pkg/license/db"
	"github.com/solo-io/solo-projects/pkg/license/model"
)

var _ db.KeyDb = new(BucketDb)

type BucketDb struct {
	Bucket *blob.Bucket
}

func (d *BucketDb) Save(ctx context.Context, r *model.Request, ua model.UserInfo, key string) error {

	data := model.RequestAndKey{
		Request: r,
		Key:     key,
	}

	jsondata, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	randint := rand.Uint32()
	return d.Bucket.WriteAll(ctx, fmt.Sprintf("%v-%v", randint, ua.Email), []byte(jsondata), nil)
}

func init() {
	rand.Seed(time.Now().Unix())
}
