package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

func AWSFetcher(region string, t AccessToken) ([]Lambda, error) {
	session, err := session.NewSession(aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(t.ID, t.Secret, "")))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get AWS session")
	}
	svc := lambda.New(session, &aws.Config{Region: aws.String(region)})
	options := &lambda.ListFunctionsInput{FunctionVersion: aws.String("ALL")}
	result, err := svc.ListFunctions(options)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get list of functions from AWS")
	}

	lambdas := make([]Lambda, len(result.Functions))
	for i, f := range result.Functions {
		lambdas[i] = Lambda{
			Name:      aws.StringValue(f.FunctionName),
			Qualifier: aws.StringValue(f.Version),
		}
	}

	return lambdas, nil
}
