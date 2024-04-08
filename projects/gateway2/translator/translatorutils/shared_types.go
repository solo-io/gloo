package translatorutils

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

type TranslationReports struct {
	ProxyReport     *validation.ProxyReport
	ResourceReports reporter.ResourceReports
}
