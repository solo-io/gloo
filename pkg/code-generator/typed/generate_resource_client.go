package typed

import (
	"bytes"
	"text/template"
)

func GenerateResourceClientCode(params ResourceLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := typedClientTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateResourceClientTestCode(params ResourceLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := resourceClientTestTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var typedClientTemplate = template.Must(template.New("typed_client").Funcs(funcs).Parse(typedClientTemplateContents))

var resourceClientTestTemplate = template.Must(template.New("typed_client_kube_test").Funcs(funcs).Parse(resourceClientTestTemplateContents))

const typedClientTemplateContents = `package {{ .PackageName }}

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func New{{ .ResourceType }}(namespace, name string) *{{ .ResourceType }} {
	return &{{ .ResourceType }}{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

{{- if $.IsInputType }}

func (r *{{ .ResourceType }}) SetStatus(status core.Status) {
	r.Status = status
}
{{- end }}

func (r *{{ .ResourceType }}) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

{{- if .IsDataType}}

func (r *{{ .ResourceType }}) SetData(data map[string]string) {
	r.Data = data
}
{{- end}}

type {{ .ResourceType }}List []*{{ .ResourceType }}

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list {{ .ResourceType }}List) Find(namespace, name string) (*{{ .ResourceType }}, error) {
	for _, {{ lowercase .ResourceType }} := range list {
		if {{ lowercase .ResourceType }}.Metadata.Name == name {
			if namespace == "" || {{ lowercase .ResourceType }}.Metadata.Namespace == namespace {
				return {{ lowercase .ResourceType }}, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find {{ lowercase .ResourceType }} %v.%v", namespace, name)
}

func (list {{ .ResourceType }}List) AsResources() []resources.Resource {
	var ress []resources.Resource 
	for _, {{ lowercase .ResourceType }} := range list {
		ress = append(ress, {{ lowercase .ResourceType }})
	}
	return ress
}

{{ if $.IsInputType -}}
func (list {{ .ResourceType }}List) AsInputResources() []resources.InputResource {
	var ress []resources.InputResource
	for _, {{ lowercase .ResourceType }} := range list {
		ress = append(ress, {{ lowercase .ResourceType }})
	}
	return ress
}
{{- end}}

func (list {{ .ResourceType }}List) Names() []resources.Resource {
	var names []string
	for _, {{ lowercase .ResourceType }} := range list {
		names = append(names, {{ lowercase .ResourceType }}.Metadata.Name)
	}
	return names
}

func (list {{ .ResourceType }}List) NamespacesDotNames() []resources.Resource {
	var names []string
	for _, {{ lowercase .ResourceType }} := range list {
		names = append(names, {{ lowercase .ResourceType }}.Metadata.Namespace + "." + {{ lowercase .ResourceType }}.Metadata.Name)
	}
	return names
}

func (list {{ .ResourceType }}List) () []resources.Resource {
	var names []string
	for _, {{ lowercase .ResourceType }} := range list {
		names = append(names, {{ lowercase .ResourceType }}.Metadata.Namespace + "." + {{ lowercase .ResourceType }}.Metadata.Name)
	}
	return names
}

var _ resources.Resource = &{{ .ResourceType }}{}

type {{ .ResourceType }}Client interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*{{ .ResourceType }}, error)
	Write(resource *{{ .ResourceType }}, opts clients.WriteOpts) (*{{ .ResourceType }}, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ({{ .ResourceType }}List, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan {{ .ResourceType }}List, <-chan error, error)
}

type {{ lowercase .ResourceType }}Client struct {
	rc clients.ResourceClient
}

func New{{ .ResourceType }}Client(rcFactory factory.ResourceClientFactory) ({{ .ResourceType }}Client, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &{{ .ResourceType }}{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base {{ .ResourceType }} resource client")
	}
	return &{{ lowercase .ResourceType }}Client{
		rc: rc,
	}, nil
}

func (client *{{ lowercase .ResourceType }}Client) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *{{ lowercase .ResourceType }}Client) Register() error {
	return client.rc.Register()
}

func (client *{{ lowercase .ResourceType }}Client) Read(namespace, name string, opts clients.ReadOpts) (*{{ .ResourceType }}, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *{{ lowercase .ResourceType }}Client) Write({{ lowercase .ResourceType }} *{{ .ResourceType }}, opts clients.WriteOpts) (*{{ .ResourceType }}, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write({{ lowercase .ResourceType }}, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *{{ lowercase .ResourceType }}Client) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *{{ lowercase .ResourceType }}Client) List(namespace string, opts clients.ListOpts) ({{ .ResourceType }}List, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertTo{{ .ResourceType }}(resourceList), nil
}

func (client *{{ lowercase .ResourceType }}Client) Watch(namespace string, opts clients.WatchOpts) (<-chan {{ .ResourceType }}List, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	{{ lowercase .ResourceType }}sChan := make(chan {{ .ResourceType }}List)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				{{ lowercase .ResourceType }}sChan <- convertTo{{ .ResourceType }}(resourceList)
			case <-opts.Ctx.Done():
				close({{ lowercase .ResourceType }}sChan)
				return
			}
		}
	}()
	return {{ lowercase .ResourceType }}sChan, errs, nil
}

func convertTo{{ .ResourceType }}(resources []resources.Resource) {{ .ResourceType }}List {
	var {{ lowercase .ResourceType }}List {{ .ResourceType }}List
	for _, resource := range resources {
		{{ lowercase .ResourceType }}List = append({{ lowercase .ResourceType }}List, resource.(*{{ .ResourceType }}))
	}
	return {{ lowercase .ResourceType }}List
}

// Kubernetes Adapter for {{ .ResourceType }}

func (o *{{ .ResourceType }}) GetObjectKind() schema.ObjectKind {
	t := {{ .ResourceType }}Crd.TypeMeta()
	return &t
}

func (o *{{ .ResourceType }}) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*{{ .ResourceType }})
}

var {{ .ResourceType }}Crd = crd.NewCrd("{{ .GroupName }}",
	"{{ .PluralName }}",
	"{{ .GroupName }}",
	"{{ .Version }}",
	"{{ .ResourceType }}",
	"{{ .ShortName }}",
	&{{ .ResourceType }}{})
`

const resourceClientTestTemplateContents = `package {{ .PackageName }}

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/bxcodec/faker"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/tests/typed"
)

var _ = Describe("{{ .ResourceType }}Client", func() {
	var (
		namespace string
	)
	for _, test := range []typed.ResourceClientTester{
		&typed.KubeRcTester{Crd: {{ .ResourceType }}Crd},
		&typed.ConsulRcTester{},
		&typed.FileRcTester{},
		&typed.MemoryRcTester{},
	{{- if .IsDataType }}
		&typed.VaultRcTester{},
		&typed.KubeSecretRcTester{},
		&typed.KubeConfigMapRcTester{},
	{{- end}}
	} {
		Context("resource client backed by "+test.Description(), func() {
			var (
				client {{ .ResourceType }}Client
				err    error
			)
			BeforeEach(func() {
				namespace = helpers.RandString(6)
				factoryOpts := test.Setup(namespace)
				client, err = New{{ .ResourceType }}Client(factory.NewResourceClientFactory(factoryOpts))
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				test.Teardown(namespace)
			})
			It("CRUDs {{ .ResourceType }}s", func() {
				{{ .ResourceType }}ClientTest(namespace, client)
			})
		})
	}
})

func {{ .ResourceType }}ClientTest(namespace string, client {{ .ResourceType }}Client) {
	err := client.Register()
	Expect(err).NotTo(HaveOccurred())

	name := "foo"
	input := New{{ .ResourceType }}(namespace, name)
	input.Metadata.Namespace = namespace
	r1, err := client.Write(input, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	_, err = client.Write(input, clients.WriteOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsExist(err)).To(BeTrue())

	Expect(r1).To(BeAssignableToTypeOf(&{{ .ResourceType }}{}))
	Expect(r1.GetMetadata().Name).To(Equal(name))
	Expect(r1.GetMetadata().Namespace).To(Equal(namespace))
	{{- range .Fields }}
	Expect(r1.{{ . }}).To(Equal(input.{{ . }}))
	{{- end }}

	_, err = client.Write(input, clients.WriteOpts{
		OverwriteExisting: true,
	})
	Expect(err).To(HaveOccurred())

	input.Metadata.ResourceVersion = r1.GetMetadata().ResourceVersion
	r1, err = client.Write(input, clients.WriteOpts{
		OverwriteExisting: true,
	})
	Expect(err).NotTo(HaveOccurred())

	read, err := client.Read(namespace, name, clients.ReadOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(read).To(Equal(r1))

	_, err = client.Read("doesntexist", name, clients.ReadOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsNotExist(err)).To(BeTrue())

	name = "boo"
	input = &{{ .ResourceType }}{}

	// ignore return error because interfaces / oneofs mess it up
	faker.FakeData(input)

	input.Metadata = core.Metadata{
		Name:      name,
		Namespace: namespace,
	}

	r2, err := client.Write(input, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	list, err := client.List(namespace, clients.ListOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(list).To(ContainElement(r1))
	Expect(list).To(ContainElement(r2))

	err = client.Delete(namespace, "adsfw", clients.DeleteOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsNotExist(err)).To(BeTrue())

	err = client.Delete(namespace, "adsfw", clients.DeleteOpts{
		IgnoreNotExist: true,
	})
	Expect(err).NotTo(HaveOccurred())

	err = client.Delete(namespace, r2.GetMetadata().Name, clients.DeleteOpts{})
	Expect(err).NotTo(HaveOccurred())
	list, err = client.List(namespace, clients.ListOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(list).To(ContainElement(r1))
	Expect(list).NotTo(ContainElement(r2))

	w, errs, err := client.Watch(namespace, clients.WatchOpts{
		RefreshRate: time.Hour,
	})
	Expect(err).NotTo(HaveOccurred())

	var r3 resources.Resource
	wait := make(chan struct{})
	go func() {
		defer close(wait)
		defer GinkgoRecover()

		resources.UpdateMetadata(r2, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		r2, err = client.Write(r2, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		name = "goo"
		input = &{{ .ResourceType }}{}
		// ignore return error because interfaces / oneofs mess it up
		faker.FakeData(input)
		Expect(err).NotTo(HaveOccurred())
		input.Metadata = core.Metadata{
			Name:      name,
			Namespace: namespace,
		}

		r3, err = client.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
	}()
	<-wait

	select {
	case err := <-errs:
		Expect(err).NotTo(HaveOccurred())
	case list = <-w:
	case <-time.After(time.Millisecond * 5):
		Fail("expected a message in channel")
	}

drain:
	for {
		select {
		case list = <-w:
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Millisecond * 500):
			break drain
		}
	}

	Expect(list).To(ContainElement(r1))
	Expect(list).To(ContainElement(r2))
	Expect(list).To(ContainElement(r3))
}
`
