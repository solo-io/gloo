package discoveries

import (
	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	printer2 "github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/printer"
)

// GraphQLApiOptions should only have one Local or Remote executor.
type GraphQLApiOptions struct {
	Local  *LocalExecutor
	Remote *RemoteExecutor
	Schema string
}

// LocalExecutor is used for options for constructing the local executor.
type LocalExecutor struct {
	// Resolutions are the local resolutions used to resolve the GraphQL queries.
	Resolutions map[string]*v1beta1.Resolution
}

// RemoteExecutor is used for options for constructing the remote exeuctor.
type RemoteExecutor struct {
	// Path is the path for quering the introspection query, or any query for that matter.
	Path string
}

// NewGraphQLApi creates a GraphQLApi resource.
func NewGraphQLApi(upstream *v1.Upstream, options GraphQLApiOptions) (*v1beta1.GraphQLApi, error) {
	var executor *v1beta1.Executor
	if options.Local != nil && options.Remote != nil {
		return nil, eris.New("cannot create a graphQLApi resource with both a local and remote executor")
	}
	if options.Local == nil && options.Remote == nil {
		return nil, eris.New("cannot create a graphQLApi without a remote or local executor")
	}
	if options.Local != nil {
		executor = &v1beta1.Executor{
			Executor: &v1beta1.Executor_Local_{
				Local: &v1beta1.Executor_Local{
					Resolutions:         options.Local.Resolutions,
					EnableIntrospection: true,
				},
			},
		}
	} else {
		executor = &v1beta1.Executor{
			Executor: &v1beta1.Executor_Remote_{
				Remote: &v1beta1.Executor_Remote{
					UpstreamRef: upstream.GetMetadata().Ref(),
					Headers: map[string]string{
						":path": options.Remote.Path,
					},
				},
			},
		}
	}
	return &v1beta1.GraphQLApi{
		Metadata: &core.Metadata{
			Name:      upstream.GetMetadata().GetName(),
			Namespace: upstream.GetMetadata().GetNamespace(),
		},
		Schema: &v1beta1.GraphQLApi_ExecutableSchema{
			ExecutableSchema: &v1beta1.ExecutableSchema{
				Executor:         executor,
				SchemaDefinition: printer2.PrettyPrintKubeString(options.Schema),
			},
		},
	}, nil
}
