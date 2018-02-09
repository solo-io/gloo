package aws

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// TODO(ashish) - when adding next source abstract this
// to generic poller that takes fetchers and updaters

// AccessToken used to talk to AWS
type AccessToken struct {
	ID     string
	Secret string
}

// Region represents the AWS region and Lambdas in the region
type Region struct {
	ID       string
	Name     string
	TokenRef string
	Lambdas  []Lambda
}

// Lambda represents AWS Lambda, each qualifier is treated as separate Lambda
type Lambda struct {
	Name      string
	Qualifier string
}

// Fetcher gets collection of Lambdas for given region
type FetcherFunc func(string, string) ([]Lambda, error)

// Updater updates changes in Lambdas; for example saves
// it in CRDs
type UpdaterFunc func(Region) error

type AWSPoller struct {
	repo    *memRepo
	fetcher FetcherFunc
	updater UpdaterFunc
}

func NewAWSPoller(f FetcherFunc, u UpdaterFunc) *AWSPoller {
	return &AWSPoller{
		repo:    newRepo(),
		fetcher: f,
		updater: u,
	}
}

func (a *AWSPoller) AddUpdateRegion(region Region) {
	existing, exists := a.repo.get(region.ID)
	if exists {
		region.Lambdas = existing.Lambdas
	}
	a.repo.set(region)
}

func (a *AWSPoller) RemoveRegion(regionID string) {
	a.repo.delete(regionID)
}

func (a *AWSPoller) Start(pollPeriod time.Duration, stop <-chan struct{}) {
	go wait.Until(a.run, pollPeriod, stop)
	log.Println("AWS Poller started")
}

func (a *AWSPoller) run() {
	regions := a.repo.regions()
	log.Printf("polling aws for %d regions\n", len(regions))
	for _, r := range regions {
		newLambdas, err := a.fetcher(r.Name, r.TokenRef)
		if err != nil {
			log.Printf("Unable to get lambdas for %s, %q\n", r.Name, err)
			continue
		}

		if diff(newLambdas, r.Lambdas) {
			updated := Region{
				ID:       r.ID,
				Name:     r.Name,
				TokenRef: r.TokenRef,
				Lambdas:  newLambdas}
			if err := a.updater(updated); err != nil {
				log.Printf("unable to update change in Lambdas for %s: %q\n", r.Name, err)
				continue
			}
			a.repo.set(updated)
		}
	}
}

func diff(l, r []Lambda) bool {
	if len(l) != len(r) {
		return true
	}

	m := make(map[string]bool, len(l))
	for _, li := range l {
		m[li.Name+":"+li.Qualifier] = true
	}

	for _, ri := range r {
		_, exists := m[ri.Name+":"+ri.Qualifier]
		if !exists {
			return true
		}
	}

	return false
}
