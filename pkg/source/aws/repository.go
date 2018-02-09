package aws

import "sync"

// TODO(ashish) - convert to generic repo and not just
// for AWS
//memRepo goroutine safe in memory repository for regions
type memRepo struct {
	m    map[string]Region
	lock sync.RWMutex
}

func newRepo() *memRepo {
	return &memRepo{
		m: make(map[string]Region),
	}
}

// Get a region for given key. Returns false if it doesn't exist
func (r *memRepo) get(k string) (Region, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	v, ok := r.m[k]
	return v, ok
}

// Set a region. Returns true if it new and false for update
func (r *memRepo) set(v Region) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.m[v.ID]
	r.m[v.ID] = v
	return !ok
}

// Delete removes the region for given key
func (r *memRepo) delete(k string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.m, k)
}

// Regions returns regions in the repository
func (r *memRepo) regions() []Region {
	r.lock.RLock()
	defer r.lock.RUnlock()

	vs := make([]Region, len(r.m))
	i := 0
	for _, v := range r.m {
		vs[i] = v
		i++
	}
	return vs
}
