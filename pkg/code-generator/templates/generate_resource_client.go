package templates

import (
	"text/template"
)

var ResourceClientTemplate = template.Must(template.New("resource_client").Parse(`package {{ .PackageName }}

import (
	"sort"

	"github.com/gogo/protobuf/proto"
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
func New{{ .Name }}(namespace, name string) *{{ .Name }} {
	return &{{ .Name }}{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

{{- if $.IsInputType }}

func (r *{{ .Name }}) SetStatus(status core.Status) {
	r.Status = status
}
{{- end }}

func (r *{{ .Name }}) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

{{- if .IsDataType}}

func (r *{{ .Name }}) SetData(data map[string]string) {
	r.Data = data
}
{{- end}}

type {{ .Name }}List []*{{ .Name }}
type {{ .PluralName }}ByNamespace map[string]{{ .Name }}List

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list {{ .Name }}List) Find(namespace, name string) (*{{ .Name }}, error) {
	for _, {{ lower_camel .ResourceType }} := range list {
		if {{ lower_camel .ResourceType }}.Metadata.Name == name {
			if namespace == "" || {{ lower_camel .ResourceType }}.Metadata.Namespace == namespace {
				return {{ lower_camel .ResourceType }}, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find {{ lower_camel .ResourceType }} %v.%v", namespace, name)
}

func (list {{ .Name }}List) AsResources() resources.ResourceList {
	var ress resources.ResourceList 
	for _, {{ lower_camel .ResourceType }} := range list {
		ress = append(ress, {{ lower_camel .ResourceType }})
	}
	return ress
}

{{ if $.IsInputType -}}
func (list {{ .Name }}List) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, {{ lower_camel .ResourceType }} := range list {
		ress = append(ress, {{ lower_camel .ResourceType }})
	}
	return ress
}
{{- end}}

func (list {{ .Name }}List) Names() []string {
	var names []string
	for _, {{ lower_camel .ResourceType }} := range list {
		names = append(names, {{ lower_camel .ResourceType }}.Metadata.Name)
	}
	return names
}

func (list {{ .Name }}List) NamespacesDotNames() []string {
	var names []string
	for _, {{ lower_camel .ResourceType }} := range list {
		names = append(names, {{ lower_camel .ResourceType }}.Metadata.Namespace + "." + {{ lower_camel .ResourceType }}.Metadata.Name)
	}
	return names
}

func (list {{ .Name }}List) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

func (list {{ .Name }}List) Clone() {{ .Name }}List {
	var {{ lower_camel .ResourceType }}List {{ .Name }}List
	for _, {{ lower_camel .ResourceType }} := range list {
		{{ lower_camel .ResourceType }}List = append({{ lower_camel .ResourceType }}List, proto.Clone({{ lower_camel .ResourceType }}).(*{{ .Name }}))
	}
	return {{ lower_camel .ResourceType }}List 
}

func (list {{ .Name }}List) ByNamespace() {{ .PluralName }}ByNamespace {
	byNamespace := make({{ .PluralName }}ByNamespace)
	for _, {{ lower_camel .ResourceType }} := range list {
		byNamespace.Add({{ lower_camel .ResourceType }})
	}
	return byNamespace
}

func (byNamespace {{ .PluralName }}ByNamespace) Add({{ lower_camel .ResourceType }} ... *{{ .Name }}) {
	for _, item := range {{ lower_camel .ResourceType }} {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace {{ .PluralName }}ByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace {{ .PluralName }}ByNamespace) List() {{ .Name }}List {
	var list {{ .Name }}List
	for _, {{ lower_camel .ResourceType }}List := range byNamespace {
		list = append(list, {{ lower_camel .ResourceType }}List...)
	}
	list.Sort()
	return list
}

func (byNamespace {{ .PluralName }}ByNamespace) Clone() {{ .PluralName }}ByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &{{ .Name }}{}

type {{ .Name }}Client interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*{{ .Name }}, error)
	Write(resource *{{ .Name }}, opts clients.WriteOpts) (*{{ .Name }}, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ({{ .Name }}List, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan {{ .Name }}List, <-chan error, error)
}

type {{ lower_camel .ResourceType }}Client struct {
	rc clients.ResourceClient
}

func New{{ .Name }}Client(rcFactory factory.ResourceClientFactory) ({{ .Name }}Client, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &{{ .Name }}{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base {{ .Name }} resource client")
	}
	return &{{ lower_camel .ResourceType }}Client{
		rc: rc,
	}, nil
}

func (client *{{ lower_camel .ResourceType }}Client) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *{{ lower_camel .ResourceType }}Client) Register() error {
	return client.rc.Register()
}

func (client *{{ lower_camel .ResourceType }}Client) Read(namespace, name string, opts clients.ReadOpts) (*{{ .Name }}, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .Name }}), nil
}

func (client *{{ lower_camel .ResourceType }}Client) Write({{ lower_camel .ResourceType }} *{{ .Name }}, opts clients.WriteOpts) (*{{ .Name }}, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write({{ lower_camel .ResourceType }}, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .Name }}), nil
}

func (client *{{ lower_camel .ResourceType }}Client) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *{{ lower_camel .ResourceType }}Client) List(namespace string, opts clients.ListOpts) ({{ .Name }}List, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertTo{{ .Name }}(resourceList), nil
}

func (client *{{ lower_camel .ResourceType }}Client) Watch(namespace string, opts clients.WatchOpts) (<-chan {{ .Name }}List, <-chan error, error) {
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
	var {{ lower_camel .ResourceType }}List {{ .Name }}List
	for _, resource := range resources {
		{{ lower_camel .ResourceType }}List = append({{ lower_camel .ResourceType }}List, resource.(*{{ .Name }}))
	}
	return {{ lower_camel .ResourceType }}List
}

// Kubernetes Adapter for {{ .Name }}

func (o *{{ .Name }}) GetObjectKind() schema.ObjectKind {
	t := {{ .Name }}Crd.TypeMeta()
	return &t
}

func (o *{{ .Name }}) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*{{ .Name }})
}

var {{ .Name }}Crd = crd.NewCrd("{{ .GroupName }}",
	"{{ lowercase .PluralName }}",
	"{{ .GroupName }}",
	"{{ .Version }}",
	"{{ .Name }}",
	"{{ .ShortName }}",
	&{{ .Name }}{})
`))

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

var _ = Describe("{{ .Name }}Client", func() {
	var (
		namespace string
	)
	for _, test := range []typed.ResourceClientTester{
		&typed.KubeRcTester{Crd: {{ .Name }}Crd},
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
				client {{ .Name }}Client
				err    error
			)
			BeforeEach(func() {
				namespace = helpers.RandString(6)
				factoryOpts := test.Setup(namespace)
				client, err = New{{ .Name }}Client(factory.NewResourceClientFactory(factoryOpts))
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				test.Teardown(namespace)
			})
			It("CRUDs {{ .Name }}s", func() {
				{{ .Name }}ClientTest(namespace, client)
			})
		})
	}
})

func {{ .Name }}ClientTest(namespace string, client {{ .Name }}Client) {
	err := client.Register()
	Expect(err).NotTo(HaveOccurred())

	name := "foo"
	input := New{{ .Name }}(namespace, name)
	input.Metadata.Namespace = namespace
	r1, err := client.Write(input, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	_, err = client.Write(input, clients.WriteOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsExist(err)).To(BeTrue())

	Expect(r1).To(BeAssignableToTypeOf(&{{ .Name }}{}))
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
	input = &{{ .Name }}{}

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
		input = &{{ .Name }}{}
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
