package server

import (
	"log"
	"time"

	"github.com/solo-io/gloo/pkg/protoutil"

	google_protobuf "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	apiv1 "github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/secret"
	"github.com/solo-io/gloo-function-discovery/pkg/source"
	"github.com/solo-io/gloo-function-discovery/pkg/source/aws"
	"github.com/solo-io/gloo-function-discovery/pkg/source/gcf"
	"github.com/solo-io/gloo-function-discovery/pkg/source/openapi"
	storage "github.com/solo-io/gloo-storage"
)

// Server represents the service discovery service
type Server struct {
	Upstreams  storage.Upstreams
	SecretRepo *secret.SecretRepo
	Port       int
}

// Start starts the discovery service and its components
func (s *Server) Start(resyncPeriod time.Duration, stop <-chan struct{}) {
	s.SecretRepo.Run(stop)
	ctrlr := newController(resyncPeriod, s.Upstreams)

	source.FetcherRegistry.Add(aws.GetAWSFetcher(s.SecretRepo))
	source.FetcherRegistry.Add(gcf.GetGCFFetcher(s.SecretRepo))
	source.FetcherRegistry.Add(openapi.GetOpenAPIFetcher())

	updater := func(u source.Upstream) error {
		crd, err := s.Upstreams.Get(u.Name)
		if err != nil {
			return errors.Wrapf(err, "unable to update stream %s", u.Name)
		}
		crd.Functions = toAPIFunctions(u.Functions)
		log.Println("updating upstream ", u.Name)
		_, err = s.Upstreams.Update(crd)
		return err
	}
	poller := source.NewPoller(updater)
	poller.Start(resyncPeriod, stop)
	ctrlr.AddHandler(&handler{poller: poller})
	ctrlr.Run(stop)
}

func toAPIFunctions(functions []source.Function) []*apiv1.Function {
	crds := make([]*apiv1.Function, len(functions))
	for i, f := range functions {
		s := &google_protobuf.Struct{}
		err := protoutil.UnmarshalMap(f.Spec, s)
		if err != nil {
			log.Println("unexpected error unmarshalling function: ", f.Name, err)
			return []*apiv1.Function{}
		}
		crds[i] = &apiv1.Function{
			Name: f.Name,
			Spec: s,
		}
	}
	return crds
}
