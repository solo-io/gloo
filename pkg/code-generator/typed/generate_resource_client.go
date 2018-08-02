package typed

import (
	"bytes"
	"text/template"
)

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

var typedClientTemplate = template.Must(template.New("typed_client").Funcs(funcs).Parse(typedClientTemplateContents))

var kubeTestTemplate = template.Must(template.New("typed_client_kube_test").Funcs(funcs).Parse(kubeTestTemplateContents))

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

func (r *{{ .ResourceType }}) SetStatus(status core.Status) {
	r.Status = status
}

func (r *{{ .ResourceType }}) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

{{- if .IsDataType}}

func (r *{{ .ResourceType }}) SetData(data map[string]string) {
	r.Data = data
}
{{- end}}

var _ resources.Resource = &{{ .ResourceType }}{}

type {{ .ResourceType }}Client interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*{{ .ResourceType }}, error)
	Write(resource *{{ .ResourceType }}, opts clients.WriteOpts) (*{{ .ResourceType }}, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*{{ .ResourceType }}, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*{{ .ResourceType }}, <-chan error, error)
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

func (client *{{ lowercase .ResourceType }}Client) Register() error {
	return client.rc.Register()
}

func (client *{{ lowercase .ResourceType }}Client) Read(namespace, name string, opts clients.ReadOpts) (*{{ .ResourceType }}, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *{{ lowercase .ResourceType }}Client) Write({{ lowercase .ResourceType }} *{{ .ResourceType }}, opts clients.WriteOpts) (*{{ .ResourceType }}, error) {
	resource, err := client.rc.Write({{ lowercase .ResourceType }}, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*{{ .ResourceType }}), nil
}

func (client *{{ lowercase .ResourceType }}Client) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *{{ lowercase .ResourceType }}Client) List(namespace string, opts clients.ListOpts) ([]*{{ .ResourceType }}, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertTo{{ .ResourceType }}(resourceList), nil
}

func (client *{{ lowercase .ResourceType }}Client) Watch(namespace string, opts clients.WatchOpts) (<-chan []*{{ .ResourceType }}, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	{{ lowercase .ResourceType }}sChan := make(chan []*{{ .ResourceType }})
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				{{ lowercase .ResourceType }}sChan <- convertTo{{ .ResourceType }}(resourceList)
			}
		}
	}()
	return {{ lowercase .ResourceType }}sChan, errs, nil
}

func convertTo{{ .ResourceType }}(resources []resources.Resource) []*{{ .ResourceType }} {
	var {{ lowercase .ResourceType }}List []*{{ .ResourceType }}
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
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
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
		clientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: {{ .ResourceType }}Crd,
			Cfg: cfg,
		})
		client, err = New{{ .ResourceType }}Client(clientFactory)
		Expect(err).NotTo(HaveOccurred())
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

		read, err := client.Read(namespace, name, clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(read).To(Equal(r1))

		_, err = client.Read("doesntexist", name, clients.ReadOpts{})
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
