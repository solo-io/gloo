package typed

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

type ResourceLevelTemplateParams struct {
	PackageName           string
	ResourceType          string
	ResourceTypeLowerCase string
	PluralName            string
	GroupName             string
	Version               string
	ShortName             string
	Fields                []string
}

type Field struct {
	Name   string
	Type   string
	Fields []Field
}

type PackageLevelTemplateParams struct {
	PackageName   string
	ResourceTypes []string
}

var funcs = template.FuncMap{
	"join":      strings.Join,
	"lowercase": strcase.ToLowerCamel,
}

func GenerateTypedClientCode(params ResourceLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := typedClientTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateTypedClientKubeTestCode(params ResourceLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := kubeTestTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateTestSuiteCode(params PackageLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := testSuiteTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateInventoryCode(params PackageLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := inventoryTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var typedClientTemplate = template.Must(template.New("typed_client").Funcs(funcs).Parse(typedClientTemplateContents))
var kubeTestTemplate = template.Must(template.New("typed_client_kube_test").Funcs(funcs).Parse(kubeTestTemplateContents))
var testSuiteTemplate = template.Must(template.New("typed_client_test_suite").Funcs(funcs).Parse(testSuiteTemplateContents))
var inventoryTemplate = template.Must(template.New("inventory").Funcs(funcs).Parse(inventoryTemplateContents))

const typedClientTemplateContents = `package {{ .PackageName }}

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime"
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

type {{ .ResourceTypeLowerCase }}Client struct {
	rc clients.ResourceClient
}

func New{{ .ResourceType }}Client(factory *factory.ResourceClientFactory) {{ .ResourceType }}Client {
	return &{{ .ResourceTypeLowerCase }}Client{
		rc: factory.NewResourceClient(&{{ .ResourceType }}{}),
	}
}

func (client *{{ .ResourceTypeLowerCase }}Client) Register() error {
	return client.rc.Register()
}

func (client *{{ .ResourceTypeLowerCase }}Client) Read(name string, opts clients.ReadOpts) (*{{ .ResourceType }}, error) {
	resource, err := client.rc.Read(name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *{{ .ResourceTypeLowerCase }}Client) Write({{ .ResourceTypeLowerCase }} *{{ .ResourceType }}, opts clients.WriteOpts) (*{{ .ResourceType }}, error) {
	resource, err := client.rc.Write({{ .ResourceTypeLowerCase }}, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *{{ .ResourceTypeLowerCase }}Client) Delete(name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(name, opts)
}

func (client *{{ .ResourceTypeLowerCase }}Client) List(opts clients.ListOpts) ([]*{{ .ResourceType }}, error) {
	resourceList, err := client.rc.List(opts)
	if err != nil {
		return nil, err
	}
	return convertTo{{ .ResourceType }}(resourceList), nil
}

func (client *{{ .ResourceTypeLowerCase }}Client) Watch(opts clients.WatchOpts) (<-chan []*{{ .ResourceType }}, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	{{ .ResourceTypeLowerCase }}sChan := make(chan []*{{ .ResourceType }})
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				{{ .ResourceTypeLowerCase }}sChan <- convertTo{{ .ResourceType }}(resourceList)
			}
		}
	}()
	return {{ .ResourceTypeLowerCase }}sChan, errs, nil
}

func convertTo{{ .ResourceType }}(resources []resources.Resource) []*{{ .ResourceType }} {
	var {{ .ResourceTypeLowerCase }}List []*{{ .ResourceType }}
	for _, resource := range resources {
		{{ .ResourceTypeLowerCase }}List = append({{ .ResourceTypeLowerCase }}List, resource.(*{{ .ResourceType }}))
	}
	return {{ .ResourceTypeLowerCase }}List
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

const testSuiteTemplateContents = `package {{ .PackageName }}

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test{{ join .ResourceTypes "" }}(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "{{ join .ResourceTypes "" }} Suite")
}
`

const kubeTestTemplateContents = `package {{ .PackageName }}

import (
	"os"
	"path/filepath"
	"time"

	"github.com/bxcodec/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("{{ .ResourceType }}Client", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace string
		cfg       *rest.Config
		client    {{ .ResourceType }}Client
	)
	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		apiextsClient, err := apiexts.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		resourceClient, err := versioned.NewForConfig(cfg, {{ .ResourceType }}Crd)
		Expect(err).NotTo(HaveOccurred())
		clientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd:     {{ .ResourceType }}Crd,
			Kube:    resourceClient,
			ApiExts: apiextsClient,
		})
		client = New{{ .ResourceType }}Client(clientFactory)
	})
	AfterEach(func() {
		services.TeardownKube(namespace)
	})
	It("CRUDs resources", func() {
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
		Expect(r1.GetMetadata().ResourceVersion).NotTo(Equal("7"))

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

		read, err := client.Read(name, clients.ReadOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(read).To(Equal(r1))

		_, err = client.Read(name, clients.ReadOpts{Namespace: "doesntexist"})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		name = "boo"
		input = &{{ .ResourceType }}{}
		err = faker.FakeData(input)
		Expect(err).NotTo(HaveOccurred())
		input.Metadata = core.Metadata{
			Name:      name,
			Namespace: namespace,
		}

		r2, err := client.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		list, err := client.List(clients.ListOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).To(ContainElement(r2))

		err = client.Delete("adsfw", clients.DeleteOpts{
			Namespace: namespace,
		})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		err = client.Delete("adsfw", clients.DeleteOpts{
			IgnoreNotExist: true,
			Namespace:      namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		err = client.Delete(r2.GetMetadata().Name, clients.DeleteOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		list, err = client.List(clients.ListOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).NotTo(ContainElement(r2))

		w, errs, err := client.Watch(clients.WatchOpts{
			Namespace:   namespace,
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
			err = faker.FakeData(input)
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
	})
})
`
const inventoryTemplateContents = `package {{ .PackageName }}

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/mitchellh/hashstructure"
	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

type Snapshot struct {
{{- range .ResourceTypes}}
	{{ . }}List []*{{.}}
{{- end}}
}

func (s Snapshot) Clone() Snapshot {
{{- range .ResourceTypes}}
	var {{ lowercase . }}List []*{{ . }}
	for _, {{ lowercase . }} := range s.{{ . }}List {
		{{ lowercase . }}List = append({{ lowercase . }}List, proto.Clone({{ lowercase . }}).(*{{ . }}))
	}
{{- end}}
	return Snapshot{
{{- range .ResourceTypes}}
		{{ . }}List: {{ lowercase . }}List,
{{- end}}
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
{{- range .ResourceTypes}}
	for _, {{ lowercase . }} := range s.{{ . }}List {
		resources.UpdateMetadata({{ lowercase . }}, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		{{ lowercase . }}.SetStatus(core.Status{})
	}
{{- end}}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}

type Cache interface {
{{- range .ResourceTypes}}
	{{ . }}() {{ . }}Client
{{- end}}
	Snapshots(opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(factory *factory.ResourceClientFactory) Cache {
	return &cache{
{{- range .ResourceTypes}}
		{{ lowercase . }}: New{{ . }}Client(factory),
{{- end}}
	}
}

type cache struct {
{{- range .ResourceTypes}}
	{{ lowercase . }} {{ . }}Client
{{- end}}
}

func (c *cache) Register() error {
{{- range .ResourceTypes}}
	if err := c.{{ lowercase . }}.Register(); err != nil {
		return err
	}
{{- end}}
	return nil
}

{{- range .ResourceTypes}}

func (c *cache) {{ . }}() {{ . }}Client {
	return c.{{ lowercase . }}
}
{{- end}}

func (c *cache) Snapshots(opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
	snapshots := make(chan *Snapshot)
	errs := make(chan error)

	currentSnapshot := Snapshot{}

	sync := func(newSnapshot Snapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}

{{- range .ResourceTypes}}
	{{ lowercase . }}Chan, {{ lowercase . }}Errs, err := c.{{ lowercase . }}.Watch(opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting {{ . }} watch")
	}
{{- end}}

	go func() {
		for {
			select {
{{- range .ResourceTypes}}
			case {{ lowercase . }}List := <-{{ lowercase . }}Chan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.{{ . }}List = {{ lowercase . }}List
				sync(newSnapshot)
{{- end}}
{{- range .ResourceTypes}}
			case err := <-{{ lowercase . }}Errs:
				errs <- err
{{- end}}
			}
		}
	}()
	return snapshots, errs, nil
}
`
