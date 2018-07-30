package file

import (
	"io/ioutil"

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
