package kubeutils

import (
	"os"

	"github.com/solo-io/gloo/test/testutils"
	enterprisetestutils "github.com/solo-io/solo-projects/test/testutils"
)

func ShouldTearDown() bool {
	return testutils.ShouldTearDown()
}

func LicenseKey() string {
	return os.Getenv(enterprisetestutils.GlooLicenseKey)
}

func IsKubeTestType(expectedType string) bool {
	return expectedType == os.Getenv(testutils.KubeTestType)
}

func IsEnvDefined(envSlice []string) bool {
	for _, e := range envSlice {
		_, defined := os.LookupEnv(e)
		if !defined {
			return false
		}
	}
	return true
}
