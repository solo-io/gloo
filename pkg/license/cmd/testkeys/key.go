package main

import (
	"context"
	"crypto"
	_ "crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha1"
	_ "crypto/sha256"
	"reflect"
	"time"

	"fmt"

	"github.com/solo-io/solo-projects/pkg/license/keys/jwt"
)

/*
utility program that prints keys so that we can see how they look like
*/

func main() {
	ctx := context.Background()
	lengths := []int{512, 1024, 2048}
	for _, hash := range []crypto.Hash{crypto.MD5, crypto.SHA1, crypto.SHA256} {
		fmt.Println("hash", reflect.TypeOf(hash.New()))
		for _, l := range lengths {

			k, err := rsa.GenerateKey(rand.Reader, l)
			if err != nil {
				panic(err)
			}

			kg := jwt.KeyGen{
				Priv: k,
				//Hash: hash,
			}
			kv := jwt.KeyValidator{
				Pub: &k.PublicKey,
			}

			key, err := kg.GenerateKey(ctx, time.Now())
			if err != nil {
				panic(err)
			}
			ki, err := kv.ValidateKey(ctx, key)
			if err != nil {
				panic(err)
			}
			fmt.Printf("length %d\tkeylen %d\tkey %s exp: %v\n", l, len(key), key, ki.Expiration)
		}
		fmt.Println()

	}

	secret := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}
	kg := jwt.KeyGenHMAC{
		Secret: secret,
	}
	kv := jwt.KeyValidatorHMAC{
		Secret: secret,
	}

	key, err := kg.GenerateKey(ctx, time.Now())
	if err != nil {
		panic(err)
	}
	ki, err := kv.ValidateKey(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("keylen %d\tkey %s exp: %v\n", len(key), key, ki.Expiration)
}
