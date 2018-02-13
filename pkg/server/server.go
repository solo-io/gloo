package server

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/glue-discovery/pkg/secret"
	"github.com/solo-io/glue-discovery/pkg/source"
	"github.com/solo-io/glue-discovery/pkg/source/aws"
	apiv1 "github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned/typed/solo.io/v1"
)

// Server represents the service discovery service
type Server struct {
	UpstreamRepo v1.UpstreamInterface
	SecretRepo   *secret.SecretRepo
	Port         int
}

// Start starts the discovery service and its components
func (s *Server) Start(resyncPeriod time.Duration, stop <-chan struct{}) {
	ctrlr := newController(resyncPeriod, s.UpstreamRepo)

	source.FetcherRegistry.Add(aws.GetAWSFetcher(s.SecretRepo))

	updater := func(u source.Upstream) error {
		crd, exists, err := ctrlr.get(u.ID)
		if err != nil {
			return errors.Wrapf(err, "unable to update stream %s", u.ID)
		}
		if !exists {
			log.Printf("upstream %s not found, will not update", u.ID)
			return nil
		}
		crd.Spec.Functions = toCRDFunctions(u.Functions)
		log.Println("updating upstream ", u.ID)
		return ctrlr.set(crd)
	}
	poller := source.NewPoller(updater)
	poller.Start(resyncPeriod, stop)
	ctrlr.AddHandler(&handler{poller: poller})
	s.SecretRepo.Run(stop)
	ctrlr.Run(stop)
}

func toCRDFunctions(functions []source.Function) []apiv1.Function {
	crds := make([]apiv1.Function, len(functions))
	for i, f := range functions {
		crds[i] = apiv1.Function{
			Name: f.Name,
			Spec: f.Spec,
		}
	}
	return crds
}
