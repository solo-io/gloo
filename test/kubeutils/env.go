package kubeutils

import (
	"os"
	"strconv"
)

const (
	TearDown       = "TEAR_DOWN"
	KubeTestType   = "KUBE2E_TESTS"
	GlooLicenseKey = "GLOO_LICENSE_KEY"
)

func ShouldTearDown() bool {
	return IsEnvTruthy(TearDown)
}

func LicenseKey() string {
	return os.Getenv(GlooLicenseKey)
}

func IsKubeTestType(expectedType string) bool {
	return expectedType == os.Getenv(KubeTestType)
}

func IsEnvTruthy(envVarName string) bool {
	envValue, _ := strconv.ParseBool(os.Getenv(envVarName))
	return envValue
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
