package common

import (
	"fmt"

	"github.com/solo-io/glue-storage/pkg/storage"
)

func UnknownType(item storage.Item) error {
	return fmt.Errorf("Unknown Item Type: %t", item)
}
