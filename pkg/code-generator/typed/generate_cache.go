package typed

import (
	"bytes"
	"text/template"
)

func GenerateCacheCode(params PackageLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := cacheTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateCacheTestCode(params PackageLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := cacheTestTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var cacheTemplate = template.Must(template.New("cache").Funcs(funcs).Parse(cacheTemplateContents))

var cacheTestTemplate = template.Must(template.New("cache_test").Funcs(funcs).Parse(cacheTestTemplateContents))

const cacheTemplateContents = `package {{ .PackageName }}

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
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
	for _, {{ lowercase . }} := range snapshotForHashing.{{ . }}List {
		resources.UpdateMetadata({{ lowercase . }}, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
{{- if (index $.ResourceLevelParams .).IsInputType }}
		{{ lowercase . }}.SetStatus(core.Status{})
{{- end }}
	}
{{- end}}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}

type Cache interface {
	Register() error
{{- range .ResourceTypes}}
	{{ . }}() {{ . }}Client
{{- end}}
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache({{ clients . true }}) Cache {
	return &cache{
{{- range .ResourceTypes}}
		{{ lowercase . }}: {{ lowercase . }}Client,
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

func (c *cache) Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
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
	{{ lowercase . }}Chan, {{ lowercase . }}Errs, err := c.{{ lowercase . }}.Watch(namespace, opts)
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

const cacheTestTemplateContents = `package {{ .PackageName }}

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("{{ uppercase .PackageName }}Cache", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace          string
		cfg                *rest.Config
		cache              Cache
{{- range .ResourceTypes }}
		{{ lowercase . }}Client {{ . }}Client
{{- end}}
	)

	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

{{- range .ResourceTypes }}

		{{ lowercase . }}ClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: {{ . }}Crd,
			Cfg: cfg,
		})
		{{ lowercase . }}Client, err = New{{ . }}Client({{ lowercase . }}ClientFactory)
		Expect(err).NotTo(HaveOccurred())
{{- end}}
		cache = NewCache({{ clients . false }})
	})
	AfterEach(func() {
		services.TeardownKube(namespace)
	})
	It("tracks snapshots on changes to any resource", func() {
		err := cache.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := cache.Snapshots(namespace, clients.WatchOpts{
			RefreshRate: time.Minute,
		})
		Expect(err).NotTo(HaveOccurred())

{{- range .ResourceTypes }}
		{{ lowercase . }}1, err := {{ lowercase . }}Client.Write(New{{ . }}(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ . }}List).To(HaveLen(1))
			Expect(snap.{{ . }}List).To(ContainElement({{ lowercase . }}1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		{{ lowercase . }}2, err := {{ lowercase . }}Client.Write(New{{ . }}(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ . }}List).To(HaveLen(2))
			Expect(snap.{{ . }}List).To(ContainElement({{ lowercase . }}1))
			Expect(snap.{{ . }}List).To(ContainElement({{ lowercase . }}2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
{{- end}}

{{- range .ResourceTypes }}
		err = {{ lowercase . }}Client.Delete({{ lowercase . }}2.Metadata.Namespace, {{ lowercase . }}2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ . }}List).To(HaveLen(1))
			Expect(snap.{{ . }}List).To(ContainElement({{ lowercase . }}1))
			Expect(snap.{{ . }}List).NotTo(ContainElement({{ lowercase . }}2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = {{ lowercase . }}Client.Delete({{ lowercase . }}1.Metadata.Namespace, {{ lowercase . }}1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ . }}List).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
{{- end}}
	})
})
`
