package pkgmemoize // import "go.pedge.io/pkg/memoize"

import "sync"

// EMemoFunc is a memoizeable function with errors.
type EMemoFunc func(int) (interface{}, error)

// NewEMemoFunc returns a new EMemoFunc.
//
// values and errors returned from f cannot be mutated, they will be persisted.
func NewEMemoFunc(f func(int, EMemoFunc) (interface{}, error)) EMemoFunc {
	return newEMemoFunc(f).Do
}

// MemoFunc is a memoizeable function without errors.
type MemoFunc func(int) interface{}

// NewMemoFunc returns a new MemoFunc.
//
// values returned from f cannot be mutated, they will be persisted.
func NewMemoFunc(f func(int, MemoFunc) interface{}) MemoFunc {
	return newMemoFunc(f).Do
}

// EMemoFuncProvider is a provider for EMemoFuncs, good for circular dependencies.
type EMemoFuncProvider interface {
	Get() EMemoFunc
}

// EMemoFuncStore is a store for EMemoFuncs, good for circular dependencies.
type EMemoFuncStore interface {
	EMemoFuncProvider
	Set(EMemoFunc)
}

// NewEMemoFuncStore returns a new EMemoFuncStore.
func NewEMemoFuncStore() EMemoFuncStore {
	return newEMemoFuncStore()
}

// MemoFuncProvider is a provider for MemoFuncs, good for circular dependencies.
type MemoFuncProvider interface {
	Get() MemoFunc
}

// MemoFuncStore is a store for MemoFuncs, good for circular dependencies.
type MemoFuncStore interface {
	MemoFuncProvider
	Set(MemoFunc)
}

// NewMemoFuncStore returns a new MemoFuncStore.
func NewMemoFuncStore() MemoFuncStore {
	return newMemoFuncStore()
}

type valueError struct {
	value interface{}
	err   error
}

type eMemoFunc struct {
	f     func(int, EMemoFunc) (interface{}, error)
	cache map[int]*valueError
	lock  *sync.RWMutex
}

func newEMemoFunc(f func(int, EMemoFunc) (interface{}, error)) *eMemoFunc {
	return &eMemoFunc{f, make(map[int]*valueError), &sync.RWMutex{}}
}

func (e *eMemoFunc) Do(i int) (interface{}, error) {
	e.lock.RLock()
	cached, ok := e.cache[i]
	e.lock.RUnlock()
	if ok {
		return cached.value, cached.err
	}
	value, err := e.f(i, e.Do)
	cached = &valueError{value: value, err: err}
	e.lock.Lock()
	e.cache[i] = cached
	e.lock.Unlock()
	return value, err
}

type memoFunc struct {
	f     func(int, MemoFunc) interface{}
	cache map[int]interface{}
	lock  *sync.RWMutex
}

func newMemoFunc(f func(int, MemoFunc) interface{}) *memoFunc {
	return &memoFunc{f, make(map[int]interface{}), &sync.RWMutex{}}
}

func (m *memoFunc) Do(i int) interface{} {
	m.lock.RLock()
	value, ok := m.cache[i]
	m.lock.RUnlock()
	if ok {
		return value
	}
	value = m.f(i, m.Do)
	m.lock.Lock()
	m.cache[i] = value
	m.lock.Unlock()
	return value
}

type eMemoFuncStore struct {
	value EMemoFunc
}

func newEMemoFuncStore() *eMemoFuncStore {
	return &eMemoFuncStore{}
}

func (s *eMemoFuncStore) Get() EMemoFunc {
	return s.value
}

func (s *eMemoFuncStore) Set(value EMemoFunc) {
	s.value = value
}

type memoFuncStore struct {
	value MemoFunc
}

func newMemoFuncStore() *memoFuncStore {
	return &memoFuncStore{}
}

func (s *memoFuncStore) Get() MemoFunc {
	return s.value
}

func (s *memoFuncStore) Set(value MemoFunc) {
	s.value = value
}
