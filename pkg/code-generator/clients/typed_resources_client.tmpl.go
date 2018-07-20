package clients

const templatecontents = `import (
	"time"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
)

func (r *{{ .ResourceType }}) SetStatus(status core.Status) {
	r.Status = status
}

func (r *{{ .ResourceType }}) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &{{ .ResourceType }}{}

type {{ .ResourceType }}Client interface {
	Register() error
	Read(name string, opts clients.ReadOpts) (*{{ .ResourceType }}, error)
	Write(resource *{{ .ResourceType }}, opts clients.WriteOpts) (*{{ .ResourceType }}, error)
	Delete(name string, opts clients.DeleteOpts) error
	List(opts clients.ListOpts) ([]*{{ .ResourceType }}, error)
	Watch(opts clients.WatchOpts) (<-chan []*{{ .ResourceType }}, <-chan error, error)
}

type typedResourceClient struct {
	rc clients.ResourceClient
}

func New{{ .ResourceType }}Client(factory *factory.ResourceClientFactory) {{ .ResourceType }}Client {
	return &typedResourceClient{
		rc: factory.NewResourceClient(&{{ .ResourceType }}{}),
	}
}

func (client *typedResourceClient) Register() error {
	return client.rc.Register()
}

func (client *typedResourceClient) Read(name string, opts clients.ReadOpts) (*{{ .ResourceType }}, error) {
	resource, err := client.rc.Read(name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *typedResourceClient) Write(typedResource *{{ .ResourceType }}, opts clients.WriteOpts) (*{{ .ResourceType }}, error) {
	resource, err := client.rc.Write(typedResource, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *typedResourceClient) Delete(name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(name, opts)
}

func (client *typedResourceClient) List(opts clients.ListOpts) ([]*{{ .ResourceType }}, error) {
	resourceList, err := client.rc.List(opts)
	if err != nil {
		return nil, err
	}
	return convertResources(resourceList), nil
}

func (client *typedResourceClient) Watch(opts clients.WatchOpts) (<-chan []*{{ .ResourceType }}, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	typedResourcesChan := make(chan []*{{ .ResourceType }})
	go func(){
		for {
			select {
			case resourceList := <- resourcesChan:
				typedResourcesChan <- convertResources(resourceList)
			}
		}
	}()
	return typedResourcesChan, errs, nil
}

func convertResources(resources []resources.Resource) []*{{ .ResourceType }} {
	var typedResourceList []*{{ .ResourceType }}
	for _, resource := range resources {
		typedResourceList = append(typedResourceList, resource.(*{{ .ResourceType }}))
	}
	return typedResourceList
}
`
