package templates

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var funcs = template.FuncMap{
	"join":        strings.Join,
	"lowercase":   strings.ToLower,
	"lower_camel": strcase.ToLowerCamel,
	"upper_camel": strcase.ToCamel,
	"snake":       strcase.ToSnake,
}

var ResourceExtensionTemplate = template.Must(template.New("resource_extensions").Funcs(funcs).Parse(`package {{ .Project.PackageName }}

import (
	"sort"

	"github.com/gogo/protobuf/proto"
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

{{- if $.HasStatus }}

func (r *{{ .Name }}) SetStatus(status core.Status) {
	r.Status = status
}
{{- end }}

func (r *{{ .Name }}) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

{{- if .HasData}}

func (r *{{ .Name }}) SetData(data map[string]string) {
	r.Data = data
}
{{- end}}

type {{ .Name }}List []*{{ .Name }}
type {{ .PluralName }}ByNamespace map[string]{{ .Name }}List

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list {{ .Name }}List) Find(namespace, name string) (*{{ .Name }}, error) {
	for _, {{ lower_camel .Name }} := range list {
		if {{ lower_camel .Name }}.Metadata.Name == name {
			if namespace == "" || {{ lower_camel .Name }}.Metadata.Namespace == namespace {
				return {{ lower_camel .Name }}, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find {{ lower_camel .Name }} %v.%v", namespace, name)
}

func (list {{ .Name }}List) AsResources() resources.ResourceList {
	var ress resources.ResourceList 
	for _, {{ lower_camel .Name }} := range list {
		ress = append(ress, {{ lower_camel .Name }})
	}
	return ress
}

{{ if $.HasStatus -}}
func (list {{ .Name }}List) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, {{ lower_camel .Name }} := range list {
		ress = append(ress, {{ lower_camel .Name }})
	}
	return ress
}
{{- end}}

func (list {{ .Name }}List) Names() []string {
	var names []string
	for _, {{ lower_camel .Name }} := range list {
		names = append(names, {{ lower_camel .Name }}.Metadata.Name)
	}
	return names
}

func (list {{ .Name }}List) NamespacesDotNames() []string {
	var names []string
	for _, {{ lower_camel .Name }} := range list {
		names = append(names, {{ lower_camel .Name }}.Metadata.Namespace + "." + {{ lower_camel .Name }}.Metadata.Name)
	}
	return names
}

func (list {{ .Name }}List) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

func (list {{ .Name }}List) Clone() {{ .Name }}List {
	var {{ lower_camel .Name }}List {{ .Name }}List
	for _, {{ lower_camel .Name }} := range list {
		{{ lower_camel .Name }}List = append({{ lower_camel .Name }}List, proto.Clone({{ lower_camel .Name }}).(*{{ .Name }}))
	}
	return {{ lower_camel .Name }}List 
}

func (list {{ .Name }}List) ByNamespace() {{ .PluralName }}ByNamespace {
	byNamespace := make({{ .PluralName }}ByNamespace)
	for _, {{ lower_camel .Name }} := range list {
		byNamespace.Add({{ lower_camel .Name }})
	}
	return byNamespace
}

func (byNamespace {{ .PluralName }}ByNamespace) Add({{ lower_camel .Name }} ... *{{ .Name }}) {
	for _, item := range {{ lower_camel .Name }} {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace {{ .PluralName }}ByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace {{ .PluralName }}ByNamespace) List() {{ .Name }}List {
	var list {{ .Name }}List
	for _, {{ lower_camel .Name }}List := range byNamespace {
		list = append(list, {{ lower_camel .Name }}List...)
	}
	list.Sort()
	return list
}

func (byNamespace {{ .PluralName }}ByNamespace) Clone() {{ .PluralName }}ByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &{{ .Name }}{}

// Kubernetes Adapter for {{ .Name }}

func (o *{{ .Name }}) GetObjectKind() schema.ObjectKind {
	t := {{ .Name }}Crd.TypeMeta()
	return &t
}

func (o *{{ .Name }}) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*{{ .Name }})
}

var {{ .Name }}Crd = crd.NewCrd("{{ .Project.GroupName }}",
	"{{ upper_camel .PluralName }}",
	"{{ .Project.GroupName }}",
	"{{ .Project.Version }}",
	"{{ .Name }}",
	"{{ .ShortName }}",
	&{{ .Name }}{})
`))
var ResourceClientTemplate = template.Must(template.New("resource_client").Funcs(funcs).Parse(`package {{ .Project.PackageName }}

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
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &{{ .Name }}{},
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

const resourceClientTestTemplateContents = `package {{ .Project.PackageName }}

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
	{{- if .HasData }}
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
