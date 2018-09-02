package templates

import (
	"text/template"
)

var ResourceGroupEmitterTestTemplate = template.Must(template.New("resource_group_emitter_test").Funcs(funcs).Parse(`package {{ .Project.PackageName }}

{{- $need_kube_clientset := new_bool }}
{{- range .Resources }}
	{{- if (.HasData) and (not .HasStatus) }}
		{{- $need_memory_client := (set_bool $need_memory_client true) }}
	{{- end}}
{{- end}}

{{- $need_memory_client := new_bool }}
{{- range .Resources }}
	{{- if (not .HasData) and (not .HasStatus) }}
		{{- $need_memory_client := (set_bool $need_memory_client true) }}
	{{- end}}
{{- end}}

{{- $clients := new_str_slice }}
{{- range .Resources}}
{{- $clients := (append_str_slice $clients (printf "%vClient"  (lower_camel .Name))) }}
{{- end}}
{{- $clients := (join_str_slice $clients ", ") }}

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
{{- if $need_memory_client }}
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
{{- end}}
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/rest"
{{- if $need_kube_clientset }}
	"k8s.io/client-go/kubernetes"
{{- end}}
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("{{ upper_camel .Project.PackageName }}Emitter", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace1          string
		namespace2          string
		cfg                *rest.Config
		emitter            {{ .GoName }}Emitter
{{- range .Resources }}
		{{ lower_camel .Name }}Client {{ .Name }}Client
{{- end}}
{{- if $need_kube_clientset }}
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

{{- range .Resources }}

		// {{ .Name }} Constructor
{{- if .HasStatus }}
		{{ lower_camel .Name }}ClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: {{ .Name }}Crd,
			Cfg: cfg,
		})
{{- else if .HasData }}
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
{{/* TODO(ilackarms): Come with a more generic way to specify that a resource is "Secret"*/}}
{{- if (eq . "Secret") }}
		{{ lower_camel .Name }}ClientFactory := factory.NewResourceClientFactory(&factory.KubeSecretClientOpts{
			Clientset: kube,
		})
{{- else }}
		{{ lower_camel .Name }}ClientFactory := factory.NewResourceClientFactory(&factory.KubeConfigMapClientOpts{
			Clientset: kube,
		})
{{- end }}
{{- else }}
		{{ lower_camel .Name }}ClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			{{ .GoName }}Emitter: memory.NewInMemoryResourceCache(),
		})
{{- end }}
		{{ lower_camel .Name }}Client, err = New{{ .Name }}Client({{ lower_camel .Name }}ClientFactory)
		Expect(err).NotTo(HaveOccurred())
{{- end}}
		emitter = New{{ .GoName }}Emitter({{ $clients }})
	})
	AfterEach(func() {
		services.TeardownKube(namespace1)
		services.TeardownKube(namespace2)
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx: ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *Snapshot
{{- range .Resources }}

		/*
			{{ .Name }}
		*/
		
		assertSnapshot{{ .PluralName }} := func(expect{{ .PluralName }} {{ .Name }}List, unexpect{{ .PluralName }} {{ .Name }}List) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expect{{ .PluralName }} {
						if _, err := snap.{{ .PluralName }}.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpect{{ .PluralName }} {
						if _, err := snap.{{ .PluralName }}.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := {{ lower_camel .Name }}Client.List(namespace1, clients.ListOpts{})
					nsList2, _ := {{ lower_camel .Name }}Client.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}	


		{{ lower_camel .Name }}1a, err := {{ lower_camel .Name }}Client.Write(New{{ .Name }}(namespace1, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		{{ lower_camel .Name }}1b, err := {{ lower_camel .Name }}Client.Write(New{{ .Name }}(namespace2, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshot{{ .PluralName }}({{ .Name }}List{ {{ lower_camel .Name }}1a, {{ lower_camel .Name }}1b }, nil)

		{{ lower_camel .Name }}2a, err := {{ lower_camel .Name }}Client.Write(New{{ .Name }}(namespace1, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		{{ lower_camel .Name }}2b, err := {{ lower_camel .Name }}Client.Write(New{{ .Name }}(namespace2, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshot{{ .PluralName }}({{ .Name }}List{ {{ lower_camel .Name }}1a, {{ lower_camel .Name }}1b,  {{ lower_camel .Name }}2a, {{ lower_camel .Name }}2b  }, nil)

		err = {{ lower_camel .Name }}Client.Delete({{ lower_camel .Name }}2a.Metadata.Namespace, {{ lower_camel .Name }}2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = {{ lower_camel .Name }}Client.Delete({{ lower_camel .Name }}2b.Metadata.Namespace, {{ lower_camel .Name }}2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshot{{ .PluralName }}({{ .Name }}List{ {{ lower_camel .Name }}1a, {{ lower_camel .Name }}1b }, {{ .Name }}List{ {{ lower_camel .Name }}2a, {{ lower_camel .Name }}2b })

		err = {{ lower_camel .Name }}Client.Delete({{ lower_camel .Name }}1a.Metadata.Namespace, {{ lower_camel .Name }}1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = {{ lower_camel .Name }}Client.Delete({{ lower_camel .Name }}1b.Metadata.Namespace, {{ lower_camel .Name }}1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshot{{ .PluralName }}(nil, {{ .Name }}List{ {{ lower_camel .Name }}1a, {{ lower_camel .Name }}1b, {{ lower_camel .Name }}2a, {{ lower_camel .Name }}2b })
{{- end}}
	})
})

`))
