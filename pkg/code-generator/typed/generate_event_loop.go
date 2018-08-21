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
	"context"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type Syncer interface {
	Sync(context.Context, *Snapshot) error
}

type EventLoop interface {
	Run(namespace string, opts clients.WatchOpts) (<-chan error, error)
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

func (el *eventLoop) Run(namespace string, opts clients.WatchOpts) (<-chan error, error) {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "{{ .PackageName }}.event_loop")
	logger := contextutils.LoggerFrom(opts.Ctx)
	logger.Infof("event loop started")

	errs := make(chan error)

	watch, cacheErrs, err := el.cache.Snapshots(namespace, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "starting snapshot watch")
	}
	go errutils.AggregateErrs(opts.Ctx, errs, cacheErrs, "{{ .PackageName }}.cache errors")
	go func() {
		// create a new context for each loop, cancel it before each loop
		var cancel context.CancelFunc
		for {
			select {
			case snapshot := <-watch:
				if cancel != nil {
					cancel()
				}
				ctx, canc := context.WithCancel(opts.Ctx)
				cancel = canc
				err := el.syncer.Sync(ctx, snapshot)
				if err != nil {
					errs <- err
				}
			case <-opts.Ctx.Done():
				return
			}
		}
	}()
	return errs, nil
}

`

const eventLoopTestTemplateContents = `package {{ .PackageName }}

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

var _ = Describe("{{ uppercase .PackageName }}EventLoop", func() {
	var (
		namespace string
		cache     Cache
		err       error
	)

	BeforeEach(func() {
{{- range .ResourceTypes}}

		{{ lowercase . }}ClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		{{ lowercase . }}Client, err := New{{ . }}Client({{ lowercase . }}ClientFactory)
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
		el := NewEventLoop(cache, sync)
		_, err := el.Run(namespace, clients.WatchOpts{})
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
