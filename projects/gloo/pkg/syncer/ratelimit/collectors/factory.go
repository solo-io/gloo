package collectors

import (
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
	"go.uber.org/zap"
)

var UnknownCollectorTypeErr = func(typ CollectorType) error {
	return eris.Errorf("unknown rate limit config collector type [%v]", typ)
}

type collectorFactory struct {
	settings                *ratelimit.ServiceSettings
	globalTranslator        shims.GlobalRateLimitTranslator
	crdTranslator           shims.RateLimitConfigTranslator
	ingressConfigTranslator translation.BasicRateLimitTranslator
}

func NewCollectorFactory(
	settings *ratelimit.ServiceSettings,
	globalTranslator shims.GlobalRateLimitTranslator,
	crdTranslator shims.RateLimitConfigTranslator,
	ingressConfigTranslator translation.BasicRateLimitTranslator,
) ConfigCollectorFactory {
	return collectorFactory{
		settings:                settings,
		globalTranslator:        globalTranslator,
		crdTranslator:           crdTranslator,
		ingressConfigTranslator: ingressConfigTranslator,
	}
}

func (f collectorFactory) MakeInstance(typ CollectorType, snapshot *gloov1snap.ApiSnapshot, logger *zap.SugaredLogger) (ConfigCollector, error) {
	switch typ {
	case Global:
		return NewGlobalConfigCollector(f.settings, logger, f.globalTranslator), nil
	case Basic:
		return NewBasicConfigCollector(f.ingressConfigTranslator), nil
	case Crd:
		return NewCrdConfigCollector(snapshot, f.crdTranslator), nil
	default:
		return nil, UnknownCollectorTypeErr(typ)
	}
}
