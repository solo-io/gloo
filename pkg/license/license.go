package license

import (
	"context"
	"fmt"
	"os"

	"github.com/solo-io/licensing/pkg/defaults"
	"github.com/solo-io/licensing/pkg/keys"
)

func LicenseStatus(ctx context.Context) error {
	license := os.Getenv("GLOO_LICENSE_KEY")

	km, err := defaults.GetKeyManager()
	if err != nil {
		return err
	}

	validator := km.KeyValidator()
	ki, err := validator.ValidateKey(ctx, license)
	if err != nil {
		return err
	}
	if keys.IsExpired(ki) {
		return fmt.Errorf("license expired")
	}
	return nil
}
