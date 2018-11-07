package templates

import (
	"text/template"
)

var ResourceClientTemplate = template.Must(template.New("resource_reconciler").Funcs(funcs).Parse(`package {{ .Project.Version }}

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type {{ .Name }}Client interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*{{ .Name }}, error)
	Write(resource *{{ .Name }}, opts clients.WriteOpts) (*{{ .Name }}, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ({{ .Name }}List, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan {{ .Name }}List, <-chan error, error)
}

type {{ lower_camel .Name }}Client struct {
	rc clients.ResourceClient
}

func New{{ .Name }}Client(rcFactory factory.ResourceClientFactory) ({{ .Name }}Client, error) {
	return New{{ .Name }}ClientWithToken(rcFactory, "")
}

func New{{ .Name }}ClientWithToken(rcFactory factory.ResourceClientFactory, token string) ({{ .Name }}Client, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &{{ .Name }}{},
		Token: token,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base {{ .Name }} resource client")
	}
	return &{{ lower_camel .Name }}Client{
		rc: rc,
	}, nil
}

func (client *{{ lower_camel .Name }}Client) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *{{ lower_camel .Name }}Client) Register() error {
	return client.rc.Register()
}

func (client *{{ lower_camel .Name }}Client) Read(namespace, name string, opts clients.ReadOpts) (*{{ .Name }}, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .Name }}), nil
}

func (client *{{ lower_camel .Name }}Client) Write({{ lower_camel .Name }} *{{ .Name }}, opts clients.WriteOpts) (*{{ .Name }}, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write({{ lower_camel .Name }}, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .Name }}), nil
}

func (client *{{ lower_camel .Name }}Client) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *{{ lower_camel .Name }}Client) List(namespace string, opts clients.ListOpts) ({{ .Name }}List, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertTo{{ .Name }}(resourceList), nil
}

func (client *{{ lower_camel .Name }}Client) Watch(namespace string, opts clients.WatchOpts) (<-chan {{ .Name }}List, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	{{ lower_camel .PluralName }}Chan := make(chan {{ .Name }}List)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				{{ lower_camel .PluralName }}Chan <- convertTo{{ .Name }}(resourceList)
			case <-opts.Ctx.Done():
				close({{ lower_camel .PluralName }}Chan)
				return
			}
		}
	}()
	return {{ lower_camel .PluralName }}Chan, errs, nil
}

func convertTo{{ .Name }}(resources resources.ResourceList) {{ .Name }}List {
	var {{ lower_camel .Name }}List {{ .Name }}List
	for _, resource := range resources {
		{{ lower_camel .Name }}List = append({{ lower_camel .Name }}List, resource.(*{{ .Name }}))
	}
	return {{ lower_camel .Name }}List
}

`))
