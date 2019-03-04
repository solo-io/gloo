package defaults

import (
	"fmt"
	"os"

	"github.com/solo-io/solo-projects/pkg/license/db"
	"github.com/solo-io/solo-projects/pkg/license/db/bucket"
	"github.com/solo-io/solo-projects/pkg/license/keys"
	"github.com/solo-io/solo-projects/pkg/license/keys/jwt"
	"github.com/solo-io/solo-projects/pkg/license/notify"
	"github.com/solo-io/solo-projects/pkg/license/notify/hubspot"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
)

var bucketURL = "gs://solo-corp-gloo-pro-backup"

func GetDb() (*db.KeyDb, error) {

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		return nil, fmt.Errorf("please set GOOGLE_APPLICATION_CREDENTIALS")
	}

	b, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return nil, err
	}
	db := bucket.BucketDb{
		Bucket: b,
	}
	return db, nil
}

func GetNotifier() (*notify.Notifier, error) {
	key := os.Getenv("HUBSPOT_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("please set HUBSPOT_API_KEY")
	}

	n := hubspot.HubspotNotifier{
		Hubspot: hubspot.HubspotNotifier(key),
	}
	return n, nil
}

var secret = []byte{0xda, 0x9b, 0x2d, 0xc0, 0x10, 0x73, 0x9f, 0xcb, 0x86, 0x79, 0xc5, 0xc0, 0xb6, 0xaa, 0xc8}

func GetKeyManager() (keys.KeyManager, error) {

	return &keyManager{
		generator: jwt.KeyGenHMAC{
			Secret: []byte(secret),
		},
		verifier: jwt.KeyValidatorHMAC{
			Secret: []byte(secret),
		},
	}, nil
}

type keyManager struct {
	generator keys.KeyGenerator
	verifier  keys.KeyValidator
}

func (k *keyManager) KeyGenerator() KeyGenerator {
	return k.generator
}
func (k *keyManager) KeyValidator() KeyValidator {
	return k.verifier
}
