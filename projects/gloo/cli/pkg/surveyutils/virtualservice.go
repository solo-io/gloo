package surveyutils

import (
	"fmt"

	"github.com/solo-io/solo-projects/pkg/cliutil"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
)

func AddVirtualServiceFlagsInteractive(rl *options.RateLimit) error {

	if err := rateLimitingSurvey(rl); err != nil {
		return err
	}

	return nil
}

// TODO: move this to input virtual  host
func rateLimitingSurvey(input *options.RateLimit) error {
	yes, err := cliutil.GetYesInput("do you wish to add rate limiting to the virtual service [y/n]?")
	if err != nil {
		return err
	}

	if !yes {
		return nil
	}

	input.Enable = true

	if err := cliutil.ChooseFromList(
		"Limit requests over this unit of time for clients: ",
		&input.TimeUnit,
		options.RateLimit_TimeUnits,
	); err != nil {
		return err
	}

	if err := cliutil.GetUint32InputDefault(
		fmt.Sprintf("how many requests per %v?", input.TimeUnit),
		&input.RequestsPerTimeUnit,
		100,
	); err != nil {
		return err
	}
	return nil
}
