package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	DynamoDbContainerName = "dynamodb"
	DynamoDbPort          = "8000"
)

func RunDynamoDbContainer() error {
	image := "amazon/dynamodb-local:latest"
	args := []string{"-d", "--rm",
		// we need to port-forward to docker host for locally running tests
		// (i.e., e2e tests not in docker on the same network)
		"-p", DynamoDbPort + ":" + DynamoDbPort,
		"--net", GetContainerNetwork(),
		image,
	}
	return RunContainer(DynamoDbContainerName, args)
}

func GetDynamoDbHost() string {
	return GetDockerHost(DynamoDbContainerName)
}

type HealthCheck struct {
	IsHealthy bool
	Error     error
}

func DynamoDbHealthCheck(awsEndpoint string) HealthCheck {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(awsEndpoint),
		Region:   aws.String("us-east-2")},
	))
	client := dynamodb.New(sess)
	input := dynamodb.ListTablesInput{}
	_, err := client.ListTables(&input)
	if err != nil {
		return HealthCheck{
			IsHealthy: false,
			Error:     err,
		}
	} else {
		// can't return nil because gomega matcher dereferences the returned value
		return HealthCheck{IsHealthy: true}
	}
}
