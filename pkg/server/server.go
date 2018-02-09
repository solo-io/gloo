package server

import (
	"github.com/solo-io/glue-discovery/pkg/secret"
	"github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned/typed/solo.io/v1"
)

// Server represents the service discovery service
type Server struct {
	UpstreamRepo v1.UpstreamInterface
	SecretRepo   *secret.SecretRepo
	Port         int
}

// Start starts the discovery service and its components
func (s *Server) Start(stop <-chan struct{}) {
	ctrlr := newController(s.UpstreamRepo)
	ctrlr.Run(stop)
	s.SecretRepo.Run(stop)
	aws := newAWSHandler(ctrlr, s.SecretRepo)
	ctrlr.AddHandler(aws)
	aws.Start(stop)
}
