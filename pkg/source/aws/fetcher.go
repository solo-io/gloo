package aws

import (
	"fmt"

	"github.com/solo-io/gloo-function-discovery/pkg/secret"
	"github.com/solo-io/gloo-function-discovery/pkg/source"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

const (
	regionKey     = "region"
	credentialKey = "credential"
	keyIDKey      = "keyid"
	secretKey     = "secretkey"

	functionNameKey = "FunctionName"
	qualifierKey    = "Qualifier"

	awsUpstreamType = "aws"
)

func GetAWSFetcher(s *secret.SecretRepo) *awsFetcher {
	return &awsFetcher{secretRepo: s}
}

type awsFetcher struct {
	secretRepo *secret.SecretRepo
}

func (a *awsFetcher) CanFetch(u *source.Upstream) bool {
	return u.Type == awsUpstreamType
}

func (a *awsFetcher) Fetch(u *source.Upstream) ([]source.Function, error) {
	secretRef := secretRef(u)
	data, exists := a.secretRepo.Get(secretRef)
	if !exists {
		return nil, fmt.Errorf("unable to get credential referenced by %s", secretRef)
	}
	keyID := string(data[keyIDKey])
	secretKey := string(data[secretKey])

	session, err := session.NewSession(aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(keyID, secretKey, "")))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get AWS session")
	}
	svc := lambda.New(session, &aws.Config{Region: aws.String(region(u))})
	options := &lambda.ListFunctionsInput{FunctionVersion: aws.String("ALL")}
	result, err := svc.ListFunctions(options)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of functions from AWS")
	}

	functions := make([]source.Function, len(result.Functions))
	for i, f := range result.Functions {
		functions[i] = source.Function{
			Name: aws.StringValue(f.FunctionName) + ":" + aws.StringValue(f.Version),
			Spec: map[string]interface{}{
				functionNameKey: aws.StringValue(f.FunctionName),
				qualifierKey:    aws.StringValue(f.Version),
			},
		}
	}

	return functions, nil
}

func secretRef(u *source.Upstream) string {
	v, exists := u.Spec[credentialKey]
	if !exists {
		return ""
	}
	return v.(string)
}

func region(u *source.Upstream) string {
	v, exists := u.Spec[regionKey]
	if !exists {
		return ""
	}
	return v.(string)
}
