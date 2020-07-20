package collectors

import (
	"github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
)

var UnknownCollectorTypeErr = func(typ CollectorType) error {
	return eris.Errorf("unknown rate limit config collector type [%v]", typ)
}

type collectorFactory struct {
	settings                *ratelimit.ServiceSettings
	crdTranslator           shims.RateLimitConfigTranslator
	ingressConfigTranslator translation.BasicRateLimitTranslator
}

func NewCollectorFactory(
	settings *ratelimit.ServiceSettings,
	crdTranslator shims.RateLimitConfigTranslator,
	ingressConfigTranslator translation.BasicRateLimitTranslator,
) ConfigCollectorFactory {
	return collectorFactory{
		settings:                settings,
		crdTranslator:           crdTranslator,
		ingressConfigTranslator: ingressConfigTranslator,
	}
}

func (f collectorFactory) MakeInstance(typ CollectorType, snapshot *gloov1.ApiSnapshot, reports reporter.ResourceReports) (ConfigCollector, error) {
	switch typ {
	case Global:
		return NewGlobalConfigCollector(f.settings), nil
	case Basic:
		return NewBasicConfigCollector(reports, f.ingressConfigTranslator), nil
	case Crd:
		return NewCrdConfigCollector(snapshot, reports, f.crdTranslator), nil
	default:
		return nil, UnknownCollectorTypeErr(typ)
	}
}
