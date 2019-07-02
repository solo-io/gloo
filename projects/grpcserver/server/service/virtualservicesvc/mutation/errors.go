package mutation

import "github.com/pkg/errors"

var (
	NoRouteProvidedError = errors.Errorf("No route provided.")

	IndexOutOfBoundsError = errors.Errorf("Index out of bounds.")

	AlreadyConfiguredSslWithFiles = errors.Errorf("SSL is already configured with SSL files")

	AlreadyConfiguredSslWithSds = errors.Errorf("SSL is already configured with SDS")
)
