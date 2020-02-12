package license

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/solo-io/licensing/pkg/defaults"
)

var (
	LicenseEmptyError   = fmt.Errorf("license is empty")
	LicenseExpiredError = fmt.Errorf("license expired")
)

func LicenseStatus(ctx context.Context) error {
	license := os.Getenv("GLOO_LICENSE_KEY")
	license = strings.TrimSpace(license)
	return IsLicenseValid(ctx, license)
}

func IsLicenseValid(ctx context.Context, license string) error {
	if license == "" {
		return LicenseEmptyError
	}
	km, err := defaults.GetKeyManager()
	if err != nil {
		return err
	}

	validator := km.KeyValidator()
	decryptedLicense, err := validator.ValidateKey(ctx, license)
	if err != nil {
		return err
	}
	if decryptedLicense.IsExpired() {
		return LicenseExpiredError
	}

	return nil
}
