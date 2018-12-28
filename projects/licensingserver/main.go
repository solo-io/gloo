package licensingserver

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	pb "github.com/solo-io/solo-projects/projects/licensingserver/pkg/api/v1"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type LicensingServerClient struct {
	conn             *grpc.ClientConn
	healthClient     healthpb.HealthClient
	validationClient pb.LicenseValidationClient
	host             string
	grpcPort         string
	healthPort       string
}

type LicensingServerClientOptions func(*LicensingServerClient) error

func NewLicensingServerClient(hostname, grpcPort, healthPort string, options ...LicensingServerClientOptions) (*LicensingServerClient, error) {
	if hostname == "" {
		hostname = "localhost"
	}
	if grpcPort == "" || healthPort == "" {
		return nil, fmt.Errorf("ports cannot be empty")
	}

	lsc := &LicensingServerClient{
		host:       hostname,
		grpcPort:   grpcPort,
		healthPort: healthPort,
	}

	err := lsc.Options(options...)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(lsc.grpcUrl(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	lsc.conn = conn
	lsc.healthClient = healthpb.NewHealthClient(conn)
	lsc.validationClient = pb.NewLicenseValidationClient(conn)

	return lsc, nil
}

func (lsc *LicensingServerClient) grpcUrl() string {
	return fmt.Sprintf("%s:%s", lsc.host, lsc.grpcPort)
}
func (lsc *LicensingServerClient) httpUrl() string {
	return fmt.Sprintf("%s:%s", lsc.host, lsc.healthPort)
}

func (lsc *LicensingServerClient) Validate(key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key cannot be empty")
	}

	req := pb.ValidateKeyRequest{
		Key: &pb.LicenseKey{
			Key: key,
		},
	}
	resp, err := lsc.validationClient.ValidateKey(context.TODO(), &req)
	if err != nil {
		return false, err
	}
	return resp.Valid, nil

}

func (lsc *LicensingServerClient) HealthCheckGRPC() (healthpb.HealthCheckResponse_ServingStatus, error) {
	req := healthpb.HealthCheckRequest{Service: ""}
	resp, err := lsc.healthClient.Check(context.TODO(), &req)
	if err != nil {
		return 0, err
	}
	return resp.Status, nil
}

func (lsc *LicensingServerClient) HealthCheckHTTP() (healthpb.HealthCheckResponse_ServingStatus, error) {
	client := http.DefaultClient
	resp, err := client.Get(fmt.Sprintf("http://%s/healthcheck", lsc.httpUrl()))
	if err != nil {
		return 0, err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if string(dat) == "OK" {
		return 1, nil
	} else {
		return 2, nil
	}
}

func (lsc *LicensingServerClient) Close() error {
	return lsc.conn.Close()
}

func (lsc *LicensingServerClient) Options(options ...LicensingServerClientOptions) error {
	for _, v := range options {
		err := v(lsc)
		if err != nil {
			return err
		}
	}
	return nil
}
