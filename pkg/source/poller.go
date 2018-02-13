package source

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

type UpdaterFunc func(Upstream) error

type Poller struct {
	repo    *memRepo
	updater UpdaterFunc
}

func NewPoller(u UpdaterFunc) *Poller {
	return &Poller{repo: newRepo(), updater: u}
}

func (p *Poller) Update(u Upstream) {
	log.Printf("Poller updating upstream %s", u.ID)
	existing, exists := p.repo.get(u.ID)
	if exists {
		u.Functions = existing.Functions
	}
	p.repo.set(u)
}

func (p *Poller) Remove(upstreamID string) {
	log.Printf("Poller removing upstream %s", upstreamID)
	p.repo.delete(upstreamID)
}

func (p *Poller) Start(pollPeriod time.Duration, stop <-chan struct{}) {
	go wait.Until(p.run, pollPeriod, stop)
	log.Println("Poller started")
}

func (p *Poller) run() {
	upstreams := p.repo.upstreams()
	log.Printf("polling %d upstreams...", len(upstreams))
	for _, u := range upstreams {
		// TODO (ashish) instead of doing this everytime maintain a map of upstream to fetcher
		fetcher := FetcherRegistry.Fetcher(&u)
		if fetcher == nil {
			log.Printf("Unable to find fetcher for %s", u.ID)
			continue
		}
		// Runs fetchers in current Go routine; TODO(ashish) - use a worker pool
		newFunctions, err := fetcher.Fetch(&u)
		if err != nil {
			log.Printf("Unable to get functions for %s, %q\n", u.ID, err)
			continue
		}

		if diff(newFunctions, u.Functions) {
			updated := Upstream{
				ID:        u.ID,
				Name:      u.Name,
				Type:      u.Type,
				Spec:      u.Spec,
				Functions: newFunctions,
			}

			if err := p.updater(updated); err != nil {
				log.Printf("unable to update change in functions for %s: %q\n", u.ID, err)
				continue
			}
			p.repo.set(updated)
		}
	}
}

func diff(l, r []Function) bool {
	if len(l) != len(r) {
		return true
	}

	m := make(map[string]bool, len(l))
	for _, li := range l {
		m[li.Name] = true
	}
	for _, ri := range r {
		_, exists := m[ri.Name]
		if !exists {
			return true
		}
	}
	return false
}
