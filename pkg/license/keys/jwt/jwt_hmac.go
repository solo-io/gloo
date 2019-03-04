package jwt

import (
	"context"
	_ "crypto/md5"
	"crypto/rand"
	_ "crypto/sha1"
	_ "crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/solo-io/solo-projects/pkg/license/keys"
	"github.com/solo-io/solo-projects/pkg/license/model"

	jwt "github.com/dgrijalva/jwt-go"
)

// this is base 64 for: '{"alg":"HS256","typ":"JWT"}'
const header = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."

var _ keys.KeyGenerator = new(KeyGenHMAC)
var _ keys.KeyValidator = new(KeyValidatorHMAC)

type KeyGenHMAC struct {
	Secret []byte
}

type KeyValidatorHMAC struct {
	Secret []byte
}

func (k *KeyGenHMAC) GenerateKey(ctx context.Context, expires time.Time) (string, error) {
	randbytes := make([]byte, 4)
	_, err := rand.Reader.Read(randbytes)
	if err != nil {
		return "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": expires.Unix(),
		"k":   base64.RawStdEncoding.EncodeToString(randbytes),
	})

	key, err := t.SignedString(k.Secret)
	// remove header as in our usecase it is always the same
	return strings.TrimPrefix(key, header), err
}

func (k *KeyValidatorHMAC) ValidateKey(ctx context.Context, key string) (*model.KeyInfo, error) {
	if strings.Count(key, ".") == 1 {
		// add header if it was removed so we can parse it with ease
		key = header + key
	}
	var parser jwt.Parser
	parser.SkipClaimsValidation = false

	token, err := parser.ParseWithClaims(string(key), jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {

		if token.Header["alg"] != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("signing method doesn't match")
		}
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return k.Secret, nil
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
