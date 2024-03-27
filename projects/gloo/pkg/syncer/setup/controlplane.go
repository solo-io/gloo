package setup

import (
	"context"
	"fmt"
	"net"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	xdsserver "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	NoXdsPortFoundError = func(portName string, svcNamespace string, svcName string) error {
		return eris.Errorf("no port with the name %s found in service %s.%s", portName, svcNamespace, svcName)
	}
)

func NewControlPlane(ctx context.Context, grpcServer *grpc.Server, bindAddr net.Addr, callbacks xdsserver.Callbacks,
	xdsHost string, xdsPort int32, start bool) bootstrap.ControlPlane {
	fmt.Printf("xxxxxxxxxx NewControlPlane host=%s, port=%v\n", xdsHost, xdsPort)
	snapshotCache := xds.NewAdsSnapshotCache(ctx)
	xdsServer := server.NewServer(ctx, snapshotCache, callbacks)
	reflection.Register(grpcServer)

	return bootstrap.ControlPlane{
		GrpcService: &bootstrap.GrpcService{
			GrpcServer:      grpcServer,
			StartGrpcServer: start,
			BindAddr:        bindAddr,
			Ctx:             ctx,
		},
		SnapshotCache: snapshotCache,
		XDSServer:     xdsServer,
		XdsHost:       xdsHost,
		XdsPort:       xdsPort,
	}
}

// GetControlPlaneXdsPort gets the xDS port from the gloo Service.
func GetControlPlaneXdsPort(ctx context.Context, svcClient skkube.ServiceClient) (int32, error) {
	// this is the namespace where gloo is running
	svcNamespace := utils.GetPodNamespace()
	glooSvc, err := svcClient.Read(svcNamespace, kubeutils.GlooServiceName, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return 0, err
	}

	// find the xds port on the gloo service
	for _, port := range glooSvc.Spec.Ports {
		if port.Name == kubeutils.GlooXdsPortName {
			return port.Port, nil
		}
	}
	return 0, NoXdsPortFoundError(kubeutils.GlooXdsPortName, svcNamespace, kubeutils.GlooServiceName)
}

// GetControlPlaneXdsHost gets the xDS address from the gloo Service.
func GetControlPlaneXdsHost() string {
	return fmt.Sprintf("%s.%s.svc.%s", kubeutils.GlooServiceName, utils.GetPodNamespace(), "cluster.local")
}
