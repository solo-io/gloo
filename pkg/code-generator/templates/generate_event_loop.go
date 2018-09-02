package templates

const eventLoopTemplateContents = `package {{ .PackageName }}


`

const eventLoopTestTemplateContents = `package {{ .Project.PackageName }}

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

var _ = Describe("{{ .GoName }}EventLoop", func() {
	var (
		namespace string
		cache     Cache
		err       error
	)

	BeforeEach(func() {
{{- range .ResourceTypes}}

		{{ lower_camel . }}ClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		{{ lower_camel . }}Client, err := New{{ . }}Client({{ lower_camel . }}ClientFactory)
		Expect(err).NotTo(HaveOccurred())
{{- end}}

		cache = NewCache({{ clients . false }})
	})
	It("runs sync function on a new snapshot", func() {
{{- range .ResourceTypes}}
		_, err = cache.{{ . }}().Write(New{{ . }}(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
{{- end}}
		sync := &mockSyncer{}
		el := New{{ .GoName }}EventLoop(cache, sync)
		_, err := el.Run([]string{namespace}, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockSyncer struct {
	synced bool
}

func (s *mockSyncer) Sync(ctx context.Context, snap *Snapshot) error {
	s.synced = true
	return nil
}
`
