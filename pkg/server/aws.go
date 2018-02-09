package server

import (
	"fmt"
	"log"
	"time"

	"github.com/solo-io/glue-discovery/pkg/secret"

	"github.com/pkg/errors"
	"github.com/solo-io/glue-discovery/pkg/source/aws"
	apiv1 "github.com/solo-io/glue/pkg/api/types/v1"
	solov1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
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

// adapter between aws poller and what controller expects
// can be removed if we make aws poller implement necessary
// methods; not doing for now since aws poller doesn't
// need to be aware of any data type outside the package
type awsHandler struct {
	controller *controller
	secretRepo *secret.SecretRepo
	poller     *aws.AWSPoller
}

func newAWSHandler(c *controller, s *secret.SecretRepo) awsHandler {
	updater := func(r aws.Region) error {
		upstream, exists, err := c.get(r.ID)
		if err != nil {
			return errors.Wrapf(err, "Unable to update upstream %s", r.ID)
		}
		if !exists {
			log.Printf("upstream %s not found, will not update", r.ID)
			return nil
		}
		upstream.Spec.Functions = toFunctions(r.Lambdas)
		log.Println("updating upstream ", r.ID)
		return c.set(upstream)
	}
	fetcher := func(region, tokenRef string) ([]aws.Lambda, error) {
		data, exists := s.Get(tokenRef)
		if !exists {
			return nil, fmt.Errorf("Unable to get credential referenced by %s", tokenRef)
		}
		token := aws.AccessToken{ID: string(data[keyIDKey]), Secret: string(data[secretKey])}
		return aws.AWSFetcher(region, token)
	}
	poller := aws.NewAWSPoller(fetcher, updater)
	return awsHandler{controller: c, poller: poller}
}

func (a awsHandler) Update(u *solov1.Upstream) {
	if u.Spec.Type == awsUpstreamType {
		a.poller.AddUpdateRegion(toRegion(u))
	}
}

func (a awsHandler) Remove(u *solov1.Upstream) {
	a.poller.RemoveRegion(toID(u))
}

func (a awsHandler) Start(stop <-chan struct{}) {
	a.poller.Start(1*time.Minute, stop)
}

func toRegion(u *solov1.Upstream) aws.Region {
	r := aws.Region{
		ID:       toID(u),
		Name:     u.Spec.Spec[regionKey].(string),
		TokenRef: toTokenRef(u),
		Lambdas:  toLambdas(u.Spec.Functions),
	}
	return r
}

func toID(u *solov1.Upstream) string {
	return fmt.Sprintf("%s/%s", u.Namespace, u.Name)
}

func toTokenRef(u *solov1.Upstream) string {
	return fmt.Sprintf("%s/%s", u.Namespace, u.Spec.Spec[credentialKey].(string))
}

func toLambdas(functions []apiv1.Function) []aws.Lambda {
	lambdas := make([]aws.Lambda, len(functions))
	for i, f := range functions {
		lambdas[i] = aws.Lambda{
			Name:      f.Spec[functionNameKey].(string),
			Qualifier: f.Spec[qualifierKey].(string),
		}
	}
	return lambdas
}

func toFunctions(lambdas []aws.Lambda) []apiv1.Function {
	functions := make([]apiv1.Function, len(lambdas))
	for i, l := range lambdas {
		functions[i] = apiv1.Function{
			Name: l.Name + ":" + l.Qualifier,
			Spec: map[string]interface{}{
				functionNameKey: l.Name,
				qualifierKey:    l.Qualifier,
			},
		}
	}
	return functions
}
