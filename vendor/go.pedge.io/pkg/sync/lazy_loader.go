package pkgsync

import (
	"sync"
	"sync/atomic"
)

type lazyLoader struct {
	once  *sync.Once
	f     func() (interface{}, error)
	value *atomic.Value
	err   *atomic.Value
}

func newLazyLoader(f func() (interface{}, error)) *lazyLoader {
	return &lazyLoader{&sync.Once{}, f, &atomic.Value{}, &atomic.Value{}}
}

func (l *lazyLoader) Load() (interface{}, error) {
	l.once.Do(func() {
		value, err := l.f()
		if value != nil {
			l.value.Store(value)
		}
		if err != nil {
			l.err.Store(err)
		}
	})
	value := l.value.Load()
	err := l.err.Load()
	if err != nil {
		return value, err.(error)
	}
	return value, nil
}
