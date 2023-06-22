package aerospike

import (
	"context"
	"fmt"

	ae "github.com/aerospike/aerospike-client-go"

	"github.com/solo-io/solo-projects/test/services"

	. "github.com/onsi/gomega"
)

type Instance struct {
	dockerRunArgs []string
	containerName string

	port    uint32
	address string

	namespace string
}

func (i *Instance) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		i.Clean()
	}()

	err := services.RunContainer(i.containerName, i.dockerRunArgs)
	Expect(err).NotTo(HaveOccurred(), "should be able to run container")

	i.EventuallyIsHealthy()
}

func (i *Instance) Clean() {
	services.MustKillAndRemoveContainer(i.containerName)
}

func (i *Instance) Port() int {
	return int(i.port)
}

func (i *Instance) Address() string {
	return i.address
}

func (i *Instance) Namespace() string {
	return i.namespace
}

func (i *Instance) Url() string {
	return fmt.Sprintf("http://%s:%d", i.Address(), i.Port())
}

func (i *Instance) ClientOrError() (*ae.Client, error) {
	return ae.NewClient(i.Address(), i.Port())
}

func (i *Instance) ConfigureSettingsForRateLimiter() {
	args := []string{
		// Have to set the nsup-period so that the container will work with the rate limiter
		// if the nsup-period is not set, then a request from the rate limiter will be blocked
		// because TTL is used in the rate limiter
		"asinfo", "-v", fmt.Sprintf("set-config:context=namespace;id=%s;nsup-period=1", i.namespace),
	}
	_, err := services.ExecOnContainer(i.containerName, args)
	Expect(err).NotTo(HaveOccurred(), "should be able to configure rate limiter settings")
}

func (i *Instance) ErrorIfUnhealthy() error {
	args := []string{
		"asinfo", "-v", "status",
	}
	out, err := services.ExecOnContainer(i.containerName, args)
	if err != nil {
		return err
	}
	if len(out) > 0 && string(out) == "ok\n" {
		// if we can connect to the client, then we are good
		_, connectErr := i.ClientOrError()
		return connectErr
	}

	return nil
}

func (i *Instance) EventuallyIsHealthy() {
	Eventually(func(g Gomega) {
		err := i.ErrorIfUnhealthy()
		g.Expect(err).NotTo(HaveOccurred(), "should be able to connect to aerospike")
	}, "5s", ".1s").Should(Succeed())
}
