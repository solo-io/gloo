package federation

import (
	"crypto/md5"
	"fmt"

	"github.com/solo-io/skv2/pkg/ezkube"
)

// GetIdentifier returns a string that uniquely identifies the given federated resource on this cluster.
func GetIdentifier(fedResource ezkube.ResourceId) string {
	return fedResource.GetNamespace() + "." + fedResource.GetName()
}

// GetShortenedIdentifier returns a string, meant to be used as a label value, identifying the given federated resource on this cluster.
// Since label values have a max length of 63 characters, values that are too long are hashed/truncated to fit within the limit.
func GetShortenedIdentifier(fedResource ezkube.ResourceId) string {
	return getShortenedLabelValue(GetIdentifier(fedResource))
}

// GetOwnerAnnotation returns an annotation key/value pair containing the unique identifier for the given federated resource.
// If a resource contains this annotation, it indicates that the given federated resource owns/manages that resource.
func GetOwnerAnnotation(fedResource ezkube.ResourceId) map[string]string {
	return map[string]string{HubOwner: GetIdentifier(fedResource)}
}

// GetOwnerLabel returns a label key/value pair containing an identifier for the given federated resource.
// It can be used for querying for resources that are owned/managed by this federated resource.
// Since label values have a max length of 63 characters, values that are too long are hashed/truncated to fit within the limit.
// In the case of truncation, the owner can be verified using the owner annotation, which contains the full owner identifier.
func GetOwnerLabel(fedResource ezkube.ResourceId) map[string]string {
	return map[string]string{HubOwner: GetShortenedIdentifier(fedResource)}
}

// Copied from https://github.com/solo-io/gloo-mesh-enterprise/blob/23dbc486914629f822f4f8def682f778bad27219/pkg/labels/labels.go#L317
// getShortenedLabelValue returns a shortened version of the input string.
// It is based on the `kubeutils.SanitizeNameV2` function, but it just does the shortening part.
func getShortenedLabelValue(val string) string {
	if len(val) > 63 {
		hash := md5.Sum([]byte(val))
		val = fmt.Sprintf("%s-%x", val[:31], hash)
		val = val[:63]
	}
	return val
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
