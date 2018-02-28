package updater

import (
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	lambdaplugin "github.com/solo-io/gloo-plugins/aws"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

const (
	// expected map identifiers for secrets
	awsAccessKey = "access_key"
	awsSecretKey = "secret_key"
)

func getLambdaFuncs(us *v1.Upstream, secrets secretwatcher.SecretMap) ([]*v1.Function, error) {
	lambdaSpec, err := lambdaplugin.DecodeUpstreamSpec(us.Spec)
	if err != nil {
		return nil, errors.Wrap(err, "decoding lambda upstream spec")
	}
	awsSecrets, ok := secrets[lambdaSpec.SecretRef]
	if !ok {
		return nil, errors.Wrapf(err, "secrets not found for secret ref %v", lambdaSpec.SecretRef)
	}

	accessKey, ok := awsSecrets[awsAccessKey]
	if !ok {
		return nil, errors.Errorf("key %v missing from provided secret", awsAccessKey)
	}
	if accessKey != "" && !utf8.Valid([]byte(accessKey)) {
		return nil, errors.Errorf("%s not a valid string", awsAccessKey)
	}
	secretKey, ok := awsSecrets[awsSecretKey]
	if !ok {
		return nil, errors.Errorf("key %v missing from provided secret", awsSecretKey)
	}
	if secretKey != "" && !utf8.Valid([]byte(secretKey)) {
		return nil, errors.Errorf("%s not a valid string", awsSecretKey)
	}

	sess, err := session.NewSession(aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create AWS session")
	}
	svc := lambda.New(sess, &aws.Config{Region: aws.String(lambdaSpec.Region)})
	options := &lambda.ListFunctionsInput{FunctionVersion: aws.String("ALL")}
	results, err := svc.ListFunctions(options)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of functions from AWS")
	}
	return convertResultToFunctionSpec(results), nil
}

func convertResultToFunctionSpec(results *lambda.ListFunctionsOutput) []*v1.Function {
	var lambdaFuncs []*v1.Function
	for _, f := range results.Functions {
		fn := &v1.Function{
			Name: aws.StringValue(f.FunctionName) + ":" + aws.StringValue(f.Version),
			Spec: lambdaplugin.EncodeFunctionSpec(lambdaplugin.FunctionSpec{
				FunctionName: aws.StringValue(f.FunctionName),
				Qualifier:    aws.StringValue(f.Version),
			}),
		}
		lambdaFuncs = append(lambdaFuncs, fn)
	}
	return lambdaFuncs
}
