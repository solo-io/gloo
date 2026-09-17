package reporting

import (
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

// returns true if all the source objects for the config object were accepted in the resource reports
// errors if parsing the config obj metadata fails (should never happen)
func AllSourcesAccepted(reports reporter.ResourceReports, configObj translator.ObjectWithMetadata) (bool, error) {
	allSourcesAccepted := true

	if err := translator.ForEachSource(configObj, func(src translator.SourceRef) error {
		_, report := reports.Find(src.ResourceKind, &core.ResourceRef{Name: src.Name, Namespace: src.Namespace})

		if report.Errors != nil {
			allSourcesAccepted = false
		}

		return nil
	}); err != nil {
		return false, err
	}

	return allSourcesAccepted, nil
}

// SourceError is the error recorded against one rejected source object of a config object.
type SourceError struct {
	// ResourceKind is the kind as recorded in the config object's source metadata, which is a Go
	// type name such as "*v1.Gateway" rather than a Kubernetes kind.
	ResourceKind string
	Ref          *core.ResourceRef
	Err          error
}

// SourceErrors returns an entry for each rejected source object of the config object, in source
// metadata order. These are the reports AllSourcesAccepted decides on, so a log line recording that
// decision should carry them. Entries are unformatted; callers render the kind and error.
func SourceErrors(reports reporter.ResourceReports, configObj translator.ObjectWithMetadata) ([]SourceError, error) {
	var sourceErrors []SourceError

	if err := translator.ForEachSource(configObj, func(src translator.SourceRef) error {
		ref := &core.ResourceRef{Name: src.Name, Namespace: src.Namespace}
		_, report := reports.Find(src.ResourceKind, ref)

		if report.Errors != nil {
			sourceErrors = append(sourceErrors, SourceError{
				ResourceKind: src.ResourceKind,
				Ref:          ref,
				Err:          report.Errors,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return sourceErrors, nil
}
