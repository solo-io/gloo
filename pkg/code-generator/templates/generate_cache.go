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
		go func(namespace string, {{ lower_camel . }}Chan  <- chan {{ . }}List) {
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
		}(namespace, {{ lower_camel . }}Chan)
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
	"context"
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
		namespace1          string
		namespace2          string
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
		namespace1 = helpers.RandString(8)
		namespace2 = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace1)
		Expect(err).NotTo(HaveOccurred())
		err = services.SetupKubeForTest(namespace2)
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
		services.TeardownKube(namespace1)
		services.TeardownKube(namespace2)
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := cache.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := cache.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx: ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *Snapshot

{{- range .ResourceTypes }}
		{{ lower_camel . }}1a, err := {{ lower_camel . }}Client.Write(New{{ . }}(namespace1, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		{{ lower_camel . }}1b, err := {{ lower_camel . }}Client.Write(New{{ . }}(namespace2, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drain{{ lower_camel . }}:
		for {
			select {
			case snap = <-snapshots:
				// Expect(snap.{{ (resource . $).PluralName }}.List()).To(ContainElement({{ lower_camel . }}1a))
				// Expect(snap.{{ (resource . $).PluralName }}.List()).To(ContainElement({{ lower_camel . }}1b))
				_, err1 := snap.{{ (resource . $).PluralName }}.List().Find({{ lower_camel . }}1a.Metadata.ObjectRef())
				_, err2 := snap.{{ (resource . $).PluralName }}.List().Find({{ lower_camel . }}1b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil {
					break drain{{ lower_camel . }}
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := {{ lower_camel . }}Client.List(namespace1, clients.ListOpts{})
				nsList2, _ := {{ lower_camel . }}Client.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				msg := log.Sprintf("expected final snapshot before 10 seconds.\nexpected %v\nreceived", combined.List(), snap.{{ (resource . $).PluralName }}.List())
				Fail(msg)
			}
		}

		{{ lower_camel . }}2a, err := {{ lower_camel . }}Client.Write(New{{ . }}(namespace1, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		{{ lower_camel . }}2b, err := {{ lower_camel . }}Client.Write(New{{ . }}(namespace2, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drain{{ lower_camel . }}2:
		for {
			select {
			case snap = <-snapshots:
				_, err1 := snap.{{ (resource . $).PluralName }}.List().Find({{ lower_camel . }}1a.Metadata.ObjectRef())
				_, err2 := snap.{{ (resource . $).PluralName }}.List().Find({{ lower_camel . }}1b.Metadata.ObjectRef())
				_, err3 := snap.{{ (resource . $).PluralName }}.List().Find({{ lower_camel . }}2a.Metadata.ObjectRef())
				_, err4 := snap.{{ (resource . $).PluralName }}.List().Find({{ lower_camel . }}2b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					break drain{{ lower_camel . }}2
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := {{ lower_camel . }}Client.List(namespace1, clients.ListOpts{})
				nsList2, _ := {{ lower_camel . }}Client.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				Fail("expected final snapshot before 10 seconds. expected "+log.Sprintf("%v", combined))
			}
		}
{{- end}}

{{- range .ResourceTypes }}
		err = {{ lower_camel . }}Client.Delete({{ lower_camel . }}2a.Metadata.Namespace, {{ lower_camel . }}2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = {{ lower_camel . }}Client.Delete({{ lower_camel . }}2b.Metadata.Namespace, {{ lower_camel . }}2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ (resource . $).PluralName }}.List()).To(ContainElement({{ lower_camel . }}1a))
			Expect(snap.{{ (resource . $).PluralName }}.List()).To(ContainElement({{ lower_camel . }}1b))
			Expect(snap.{{ (resource . $).PluralName }}.List()).NotTo(ContainElement({{ lower_camel . }}2a))
			Expect(snap.{{ (resource . $).PluralName }}.List()).NotTo(ContainElement({{ lower_camel . }}2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = {{ lower_camel . }}Client.Delete({{ lower_camel . }}1a.Metadata.Namespace, {{ lower_camel . }}1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = {{ lower_camel . }}Client.Delete({{ lower_camel . }}1b.Metadata.Namespace, {{ lower_camel . }}1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.{{ (resource . $).PluralName }}.List()).NotTo(ContainElement({{ lower_camel . }}1a))
			Expect(snap.{{ (resource . $).PluralName }}.List()).NotTo(ContainElement({{ lower_camel . }}1b))
			Expect(snap.{{ (resource . $).PluralName }}.List()).NotTo(ContainElement({{ lower_camel . }}2a))
			Expect(snap.{{ (resource . $).PluralName }}.List()).NotTo(ContainElement({{ lower_camel . }}2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
{{- end}}
	})
})
`
