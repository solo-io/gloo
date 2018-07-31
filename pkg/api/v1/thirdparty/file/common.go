package file

import (
	"io/ioutil"

	"fmt"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/solo-io/solo-kit/pkg/api/v1/thirdparty"
)

func writeThirdPartyResource(file string, resource thirdparty.ThirdPartyResource) error {
	b, err := yaml.Marshal(resource)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, b, 0644)
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr+1)
}
