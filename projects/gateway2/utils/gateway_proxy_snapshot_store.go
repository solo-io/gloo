package utils

import (
	"context"
	"sync"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

// GatewayProxySnapshotStore retains the Gateway API controller's complete
// desired Proxy list. Setup runs read it independently, so a departing run
// cannot consume an update that its replacement needs.
type GatewayProxySnapshotStore struct {
	mu      sync.Mutex
	proxies v1.ProxyList
	version uint64
	changed chan struct{}
}

func NewGatewayProxySnapshotStore() *GatewayProxySnapshotStore {
	return &GatewayProxySnapshotStore{changed: make(chan struct{})}
}

// Publish accepts an empty list as a complete desired state. The store owns a
// copy because reconciliation mutates Proxy metadata.
func (s *GatewayProxySnapshotStore) Publish(proxies v1.ProxyList) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.proxies = proxies.Clone()
	s.version++
	close(s.changed)
	s.changed = make(chan struct{})
}

// Current returns the most recently published list and its version. The third
// result distinguishes an empty desired list from no publication yet.
func (s *GatewayProxySnapshotStore) Current() (v1.ProxyList, uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.version == 0 {
		return nil, 0, false
	}
	return s.proxies.Clone(), s.version, true
}

// WaitForNewer returns the newest list published after version. It never
// removes a list, and checks cancellation before returning to an old run.
func (s *GatewayProxySnapshotStore) WaitForNewer(ctx context.Context, version uint64) (v1.ProxyList, uint64, error) {
	for {
		s.mu.Lock()
		if err := ctx.Err(); err != nil {
			s.mu.Unlock()
			return nil, 0, err
		}
		if s.version > version {
			proxies, currentVersion := s.proxies.Clone(), s.version
			s.mu.Unlock()
			return proxies, currentVersion, nil
		}
		changed := s.changed
		s.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		case <-changed:
		}
	}
}
