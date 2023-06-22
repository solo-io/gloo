package dynamodb

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-projects/test/services"

	"k8s.io/utils/pointer"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/gomega"
)

type Instance struct {
	dockerRunArgs []string
	containerName string

	port    uint32
	address string
	region  string
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

func (i *Instance) Port() uint32 {
	return i.port
}

func (i *Instance) Address() string {
	return i.address
}

func (i *Instance) Url() string {
	return fmt.Sprintf("http://%s:%d", i.Address(), i.Port())
}

func (i *Instance) Client() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:   aws.String(i.Url()),
		Region:     aws.String(i.region),
		MaxRetries: pointer.Int(3)},
	))
	return dynamodb.New(sess)
}

func (i *Instance) ListTables() (*dynamodb.ListTablesOutput, error) {
	return i.Client().ListTables(&dynamodb.ListTablesInput{})
}

func (i *Instance) EventuallyIsHealthy() {
	Eventually(func(g Gomega) {
		_, err := i.ListTables()
		g.Expect(err).NotTo(HaveOccurred(), "should be able to list tables")
	}, "5s", ".1s").Should(Succeed())
}
