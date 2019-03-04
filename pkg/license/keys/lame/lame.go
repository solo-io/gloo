package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/solo-io/solo-projects/pkg/license/keys"
	"github.com/solo-io/solo-projects/pkg/license/model"
)

var _ keys.KeyGenerator = new(KeyGen)
var _ keys.KeyValidator = new(KeyValidator)

type KeyGen struct {
	Secret []byte
}

func (k *KeyGen) GenerateKey(ctx context.Context, expires time.Time) (string, error) {
	randbytes := make([]byte, 4)
	_, err := rand.Reader.Read(randbytes)
	if err != nil {
		return "", err
	}

	panic("test")

}

type KeyValidator struct {
	Pub *rsa.PublicKey
}

func (k *KeyValidator) ValidateKey(ctx context.Context, key string) (*model.KeyInfo, error) {
	panic("test")
}
