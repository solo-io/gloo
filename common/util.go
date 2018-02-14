package common

import (
	"fmt"

	"github.com/solo-io/glue-storage"
)

func UnknownType(item storage.Item) error {
	return fmt.Errorf("Unknown Item Type: %t", item)
}
