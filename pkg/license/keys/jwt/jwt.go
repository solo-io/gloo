package jwt

import (
	"context"
	_ "crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha1"
	_ "crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/solo-io/solo-projects/pkg/license/keys"
	"github.com/solo-io/solo-projects/pkg/license/model"

	jwt "github.com/dgrijalva/jwt-go"
)

var _ keys.KeyGenerator = new(KeyGen)
var _ keys.KeyValidator = new(KeyValidator)

type KeyGen struct {
	Priv *rsa.PrivateKey
}

type KeyValidator struct {
	Pub *rsa.PublicKey
}

func (k *KeyGen) GenerateKey(ctx context.Context, expires time.Time) (string, error) {
	randbytes := make([]byte, 4)
	_, err := rand.Reader.Read(randbytes)
	if err != nil {
		return "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat":                time.Now().Unix(),
		"exp":                expires.Unix(),
		"http://solo.io/key": base64.RawStdEncoding.EncodeToString(randbytes),
	})

	return t.SignedString(k.Priv)
}

func (k *KeyValidator) ValidateKey(ctx context.Context, key string) (*model.KeyInfo, error) {
	var parser jwt.Parser
	parser.SkipClaimsValidation = false

	token, err := parser.ParseWithClaims(string(key), jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {

		if token.Header["alg"] != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("signing method doesn't match")
		}
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return k.Pub, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}
	var ki model.KeyInfo
	claims := token.Claims
	if mapClaims, ok := claims.(jwt.MapClaims); ok {
		ki.Expiration, err = GetExpiresAt(mapClaims)
		if err != nil {
			return nil, err
		}
	}

	return &ki, nil
}

func GetExpiresAt(m jwt.MapClaims) (time.Time, error) {
	switch exp := m["exp"].(type) {
	case float64:
		return expToTime(int64(exp)), nil
	case json.Number:
		v, _ := exp.Int64()
		return expToTime(v), nil
	}
	return time.Now(), fmt.Errorf("no exp")
}

func expToTime(exp int64) time.Time {
	return time.Unix(exp, 0)
}
