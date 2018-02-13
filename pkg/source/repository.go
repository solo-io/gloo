package source

import "sync"

//memRepo goroutine safe in memory repository for regions
type memRepo struct {
	m    map[string]Upstream
	lock sync.RWMutex
}

func newRepo() *memRepo {
	return &memRepo{
		m: make(map[string]Upstream),
	}
}

// Get an upstream for given id. Returns false if it doesn't exist
func (r *memRepo) get(k string) (Upstream, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	v, ok := r.m[k]
	return v, ok
}

// Set an upstream. Returns true if it new and false for update
func (r *memRepo) set(v Upstream) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.m[v.ID]
	r.m[v.ID] = v
	return !ok
}

// Delete removes the upstream for given id
func (r *memRepo) delete(k string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.m, k)
}

// Upstream returns all the upstreams in the repo
func (r *memRepo) upstreams() []Upstream {
	r.lock.RLock()
	defer r.lock.RUnlock()

	vs := make([]Upstream, len(r.m))
	i := 0
	for _, v := range r.m {
		vs[i] = v
		i++
	}
	return vs
}
