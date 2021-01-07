package translator

import "fmt"

func GetGlooInstanceName(cluster, deploymentNamespace string) string {
	return fmt.Sprintf("%v-%v", cluster, deploymentNamespace)
}
