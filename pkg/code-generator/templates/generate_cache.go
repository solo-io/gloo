package templates

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
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type Snapshot struct {
{{- range .ResourceTypes}}
	{{ (resource . $).PluralName }} {{ (resource . $).PluralName }}ByNamespace
{{- end}}
}

func (s Snapshot) Clone() Snapshot {
	return Snapshot{
{{- range .ResourceTypes}}
		{{ (resource . $).PluralName }}: s.{{ (resource . $).PluralName }}.Clone(),
{{- end}}
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
{{- range .ResourceTypes}}
	for _, {{ lower_camel . }} := range snapshotForHashing.{{ (resource . $).PluralName }}.List() {
		resources.UpdateMetadata({{ lower_camel . }}, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
{{- if (index $.ResourceLevelParams .).IsInputType }}
		{{ lower_camel . }}.SetStatus(core.Status{})
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
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache({{ clients . true }}) Cache {
	return &cache{
{{- range .ResourceTypes}}
		{{ lower_camel . }}: {{ lower_camel . }}Client,
{{- end}}
	}
}

type cache struct {
{{- range .ResourceTypes}}
	{{ lower_camel . }} {{ . }}Client
{{- end}}
}

func (c *cache) Register() error {
{{- range .ResourceTypes}}
	if err := c.{{ lower_camel . }}.Register(); err != nil {
		return err
	}
{{- end}}
	return nil
}

{{- range .ResourceTypes}}

func (c *cache) {{ . }}() {{ . }}Client {
	return c.{{ lower_camel . }}
}
{{- end}}

func (c *cache) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
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

	for _, namespace := range watchNamespaces {
{{- range .ResourceTypes}}
		{{ lower_camel . }}Chan, {{ lower_camel . }}Errs, err := c.{{ lower_camel . }}.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting {{ . }} watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, {{ lower_camel . }}Errs, namespace+"-{{ lower_camel (resource . $).PluralName }}")
		go func() {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case {{ lower_camel . }}List := <-{{ lower_camel . }}Chan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.{{ (resource . $).PluralName }}.Clear(namespace)
					newSnapshot.{{ (resource . $).PluralName }}.Add({{ lower_camel . }}List...)
					sync(newSnapshot)
				}
			}
		}()
{{- end}}
	}


	go func() {
		select {
		case <-opts.Ctx.Done():
			close(snapshots)
			close(errs)
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
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
{{- if (need_memory_client .) }}
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
{{- end}}
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/rest"
{{- if (need_clientset .) }}
	"k8s.io/client-go/kubernetes"
{{- end}}
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
		{{ lower_camel . }}Client {{ . }}Client
{{- end}}
{{- if (need_clientset .) }}
		kube               kubernetes.Interface
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

		// {{ . }} Constructor
{{- if (index $.ResourceLevelParams .).IsInputType }}
		{{ lower_camel . }}ClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: {{ . }}Crd,
			Cfg: cfg,
		})
{{- else if (index $.ResourceLevelParams .).IsDataType }}
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
{{/* TODO(ilackarms): Come with a more generic way to specify that a resource is "Secret"*/}}
{{- if (eq . "Secret") }}
		{{ lower_camel . }}ClientFactory := factory.NewResourceClientFactory(&factory.KubeSecretClientOpts{
			Clientset: kube,
		})
{{- else }}
		{{ lower_camel . }}ClientFactory := factory.NewResourceClientFactory(&factory.KubeConfigMapClientOpts{
			Clientset: kube,
		})
{{- end }}
{{- else }}
		{{ lower_camel . }}ClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
{{- end }}
		{{ lower_camel . }}Client, err = New{{ . }}Client({{ lower_camel . }}ClientFactory)
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

		snapshots, errs, err := cache.Snapshots([]string{namespace}, clients.WatchOpts{
			RefreshRate: time.Minute,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *Snapshot

{{- range .ResourceTypes }}
		{{ lower_camel . }}1, err := {{ lower_camel . }}Client.Write(New{{ . }}(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drain{{ lower_camel . }}:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drain{{ lower_camel . }}
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.{{ (resource . $).PluralName }}).To(ContainElement({{ lower_camel . }}1))

		{{ lower_camel . }}2, err := {{ lower_camel . }}Client.Write(New{{ . }}(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ (resource . $).PluralName }}).To(ContainElement({{ lower_camel . }}1))
			Expect(snap.{{ (resource . $).PluralName }}).To(ContainElement({{ lower_camel . }}2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
{{- end}}

{{- range .ResourceTypes }}
		err = {{ lower_camel . }}Client.Delete({{ lower_camel . }}2.Metadata.Namespace, {{ lower_camel . }}2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ (resource . $).PluralName }}).To(ContainElement({{ lower_camel . }}1))
			Expect(snap.{{ (resource . $).PluralName }}).NotTo(ContainElement({{ lower_camel . }}2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = {{ lower_camel . }}Client.Delete({{ lower_camel . }}1.Metadata.Namespace, {{ lower_camel . }}1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ (resource . $).PluralName }}).NotTo(ContainElement({{ lower_camel . }}1))
			Expect(snap.{{ (resource . $).PluralName }}).NotTo(ContainElement({{ lower_camel . }}2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
{{- end}}
	})
})
`
