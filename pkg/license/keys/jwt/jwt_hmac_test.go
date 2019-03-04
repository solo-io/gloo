package jwt_test

import (
	"context"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/pkg/license/keys/jwt"
)

var _ = Describe("JwtHmac", func() {

	It("should create a key with no header", func() {

		secret := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}
		kg := KeyGenHMAC{
			Secret: secret,
		}
		ctx := context.Background()
		key, err := kg.GenerateKey(ctx, time.Now())
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Count(key, ".")).To(Equal(1))
	})

	It("should work", func() {

		secret := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}
		kg := KeyGenHMAC{
			Secret: secret,
		}
		kv := KeyValidatorHMAC{
			Secret: secret,
		}
		ctx := context.Background()
		exp := time.Now()
		key, err := kg.GenerateKey(ctx, exp)
		Expect(err).NotTo(HaveOccurred())
		ki, err := kv.ValidateKey(ctx, key)
		Expect(err).NotTo(HaveOccurred())

		// remove nano sec component
		expected := time.Unix(exp.Unix(), 0)
		Expect(ki.Expiration).To(Equal(expected))
	})
})
