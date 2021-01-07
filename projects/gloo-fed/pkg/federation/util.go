package federation

import (
	"github.com/solo-io/skv2/pkg/resource"
)

func GetOwnerLabel(r resource.Resource) map[string]string {
	return map[string]string{HubOwner: r.GetNamespace() + "." + r.GetName()}
}

// Merge merges any number of map[string]string into a new map.
// Key value pairs provided in later maps will overwrite same-key pairs in earlier maps.
func Merge(maps ...map[string]string) map[string]string {
	output := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			output[k] = v
		}
	}

	return output
}
