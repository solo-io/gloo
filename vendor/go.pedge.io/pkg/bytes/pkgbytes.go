package pkgbytes // import "go.pedge.io/pkg/bytes"

import (
	"bytes"

	"github.com/oxtoacart/bpool"
)

// BufferPool is a pool of buffers.
type BufferPool interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

// NewSizedBufferPool creates a new sized buffer pool.
func NewSizedBufferPool(size int, alloc int) BufferPool {
	return bpool.NewSizedBufferPool(size, alloc)
}
