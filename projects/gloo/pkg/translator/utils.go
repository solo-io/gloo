package translator

import core_solo_io1 "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

func UpstreamToClusterName(upstream core_solo_io1.ResourceRef) string {
	return upstream.Key()
}
