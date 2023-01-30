package validation

import (
	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	enterprisev1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/enterprise/graphql/v1"
	sticthing "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/v8go"
)

// If settings are configured to reject breaking changes, and the schema diff shows breaking changes, then returns
// an error with details of the breaking changes.
func ValidateSchemaUpdate(
	oldSchemaDef string,
	newSchemaDef string,
	settings *gloov1.Settings) error {

	rejectBreaking := false
	if settings.Spec.GetGraphqlOptions().GetSchemaChangeValidationOptions().GetRejectBreakingChanges() != nil {
		rejectBreaking = settings.Spec.GetGraphqlOptions().GetSchemaChangeValidationOptions().GetRejectBreakingChanges().GetValue()
	}
	// if settings are not configured to reject breaking changes, then we don't have to do any further validation
	if !rejectBreaking {
		return nil
	}

	// construct input to send to the diff function
	diffInput := &enterprisev1.GraphQLInspectorDiffInput{
		OldSchema: oldSchemaDef,
		NewSchema: newSchemaDef,
		Rules:     settings.Spec.GetGraphqlOptions().GetSchemaChangeValidationOptions().GetProcessingRules(),
	}
	runner, err := sticthing.GetStitchingScriptRunner()
	if err != nil {
		return eris.Wrap(err, "could not create the stitching script runner")
	}
	diffOutput, err := runner.RunSchemaDiff(diffInput)
	if err != nil {
		return eris.Wrap(err, "could not get schema diff")
	}

	// if there are any breaking changes, return their corresponding change descriptions
	var multiErr *multierror.Error
	for _, change := range diffOutput.GetChanges() {
		if change.GetCriticality().GetLevel() == enterprisev1.GraphQLInspectorDiffOutput_BREAKING {
			multiErr = multierror.Append(multiErr, eris.New(change.GetMessage()))
		}
	}

	return multiErr.ErrorOrNil()
}
