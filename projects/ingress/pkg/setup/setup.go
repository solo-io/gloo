package setup

import (
	"github.com/solo-io/gloo/pkg/utils/setuputils"
)

func Main() error {
	return setuputils.Main("ingress", Setup)
}
