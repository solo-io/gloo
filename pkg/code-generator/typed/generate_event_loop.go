package typed

import (
	"bytes"
	"text/template"
)

func GenerateEventLoopCode(params PackageLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := eventLoopTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateEventLoopTestCode(params PackageLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := eventLoopTestTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var eventLoopTemplate = template.Must(template.New("event_loop").Funcs(funcs).Parse(eventLoopTemplateContents))

var eventLoopTestTemplate = template.Must(template.New("event_loop_test").Funcs(funcs).Parse(eventLoopTestTemplateContents))

const eventLoopTemplateContents = `package {{ .PackageName }}

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

type Syncer interface {
	Sync(*Snapshot) error
}

type EventLoop interface {
	Run(namespace string, opts clients.WatchOpts) error
}

type eventLoop struct {
	cache  Cache
	syncer Syncer
}

func NewEventLoop(cache Cache, syncer Syncer) EventLoop {
	return &eventLoop{
		cache:  cache,
		syncer: syncer,
	}
}

func (el *eventLoop) Run(namespace string, opts clients.WatchOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "{{ .PackageName }}.event_loop")
	logger := contextutils.LoggerFrom(opts.Ctx)
	logger.Infof("event loop started")
	errorHandler := contextutils.ErrorHandlerFrom(opts.Ctx)
	watch, errs, err := el.cache.Snapshots(namespace, opts)
	if err != nil {
		return errors.Wrapf(err, "starting snapshot watch")
	}
	for {
		select {
		case snapshot := <-watch:
			err := el.syncer.Sync(snapshot)
			if err != nil {
				errorHandler.HandleErr(err)
			}
		case err := <-errs:
			errorHandler.HandleErr(err)
		}
	}
}

`

const eventLoopTestTemplateContents = `package {{ .PackageName }}

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
)

var _ = Describe("{{ uppercase .PackageName }}EventLoop", func() {
	var (
		namespace string
		cache     Cache
		err       error
	)

	BeforeEach(func() {
{{- range .ResourceTypes}}
		{{ lowercase . }}ClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{})
		{{ lowercase . }}Client := New{{ . }}Client({{ lowercase . }}ClientFactory)
{{- end}}
		cache = NewCache({{ clients . false }})
	})
	It("runs sync function on a new snapshot", func() {
{{- range .ResourceTypes}}
		_, err = cache.{{ . }}().Write(New{{ . }}(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
{{- end}}
		sync := &mockSyncer{}
		el := NewEventLoop(cache, sync)
		go func() {
			defer GinkgoRecover()
			err := el.Run(clients.WatchOpts{Namespace: namespace})
			Expect(err).NotTo(HaveOccurred())
		}()
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockSyncer struct {
	synced bool
}

func (s *mockSyncer) Sync(snap *Snapshot) error {
	s.synced = true
	return nil
}
`
