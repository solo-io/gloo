package plugin

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/translator"
)

type State struct {
	Dependencies translator.Dependencies
	Config       *v1.Config
}

type PluginInputs struct {
	State          *State
	NameTranslator translator.NameTranslator
}

func (s *State) GetFunction(fd *v1.FunctionDestination) *v1.Function {
	us := s.GetUpstream(fd.UpstreamName)
	if us != nil {

		for _, function := range us.Functions {
			if function.Name == fd.FunctionName {
				return &function
			}
		}
	}
	return nil
}

func (s *State) GetUpstream(name string) *v1.Upstream {
	for _, us := range s.Config.Upstreams {
		if us.Name == name {
			return &us
		}
	}
	return nil
}
