package setuputils_test

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	"github.com/solo-io/gloo/pkg/utils/settingsutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	external_kubernetes_namespace "github.com/solo-io/solo-kit/api/external/kubernetes/namespace"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/solo-io/gloo/pkg/utils/setuputils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("SetupSyncer", func() {
	It("calls the setup function with the referenced settings crd", func() {
		var actualSettings *v1.Settings
		expectedSettings := &v1.Settings{
			Metadata: &core.Metadata{Name: "hello", Namespace: "goodbye"},
		}
		setupSyncer := NewSetupSyncer(
			expectedSettings.Metadata.Ref(),
			func(
				ctx context.Context,
				kubeCache kube.SharedCache,
				inMemoryCache memory.InMemoryResourceCache,
				settings *v1.Settings,
				identity leaderelector.Identity) error {
				actualSettings = expectedSettings
				return nil
			},
			singlereplica.Identity())
		err := setupSyncer.Sync(context.TODO(), &v1.SetupSnapshot{
			Settings: v1.SettingsList{expectedSettings},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(actualSettings).To(Equal(expectedSettings))
	})

	When("calling ShouldSync", func() {
		var snapshot *v1.SetupSnapshot
		var settings *v1.Settings
		var namespaces kubernetes.KubeNamespaceList
		var setupSyncer *SetupSyncer

		BeforeEach(func() {
			labelSelectors := []*v1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"match": "me",
					},
					MatchExpressions: []*v1.LabelSelectorRequirement{
						{
							Key:      "common",
							Operator: "In",
							Values:   []string{"label"},
						},
					},
				},
			}
			settings = &v1.Settings{
				Metadata:                &core.Metadata{Name: "hello", Namespace: "goodbye"},
				WatchNamespaceSelectors: labelSelectors,
				DiscoveryNamespace:      "old",
			}

			matchedNamespace := &kubernetes.KubeNamespace{
				KubeNamespace: external_kubernetes_namespace.KubeNamespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "match-me",
						Labels: map[string]string{
							"match":  "me",
							"common": "label",
						},
					},
				},
			}
			skippedNamespace := &kubernetes.KubeNamespace{
				KubeNamespace: external_kubernetes_namespace.KubeNamespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "skip-me",
						Labels: map[string]string{
							"skip":   "me",
							"common": "label",
						},
					},
				},
			}
			namespaces = kubernetes.KubeNamespaceList{
				matchedNamespace,
				skippedNamespace,
			}

			snapshot = &v1.SetupSnapshot{
				Settings:       v1.SettingsList{settings},
				Kubenamespaces: namespaces,
			}

			setupSyncer = NewSetupSyncer(
				settings.Metadata.Ref(),
				func(
					ctx context.Context,
					kubeCache kube.SharedCache,
					inMemoryCache memory.InMemoryResourceCache,
					settings *v1.Settings,
					identity leaderelector.Identity) error {
					return nil
				},
				singlereplica.Identity())
		})

		BeforeEach(func() {
			// Handle the initial check where the previous settings hash is nil
			setupSyncer.Sync(context.TODO(), snapshot)
		})

		It("return false if nothing has changed", func() {
			newSnapshot := snapshot.Clone()
			Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeFalse())
		})

		It("return true if the settings object has changed", func() {
			// Ensure only the settings has changed
			newSnapshot := snapshot.Clone()
			newSnapshot.Settings[0].DiscoveryNamespace = "new"
			Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeTrue())
		})

		When("The kubenamespaces has changed", func() {
			It("A new namespace to watch has been added", func() {
				newSnapshot := snapshot.Clone()
				newSnapshot.Kubenamespaces = append(newSnapshot.Kubenamespaces, &kubernetes.KubeNamespace{
					KubeNamespace: external_kubernetes_namespace.KubeNamespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: "match-me-again",
							Labels: map[string]string{
								"match":   "me",
								"common":  "label",
								"matched": "again",
							},
						},
					},
				})
				Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeTrue())
			})

			It("A namespace to watch has been deleted", func() {
				newSnapshot := snapshot.Clone()
				newSnapshot.Kubenamespaces = kubernetes.KubeNamespaceList{newSnapshot.Kubenamespaces[1]}
				Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeTrue())
			})

			It("A namespace to watch has been modified to unwatch it", func() {
				newSnapshot := snapshot.Clone()
				newSnapshot.Kubenamespaces[0].Labels = nil
				Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeTrue())
			})

			It("A namespace to watch has been modified that doesn't affect how we watch it", func() {
				newSnapshot := snapshot.Clone()
				newSnapshot.Kubenamespaces[0].Labels["random"] = "label"
				Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeFalse())
			})

			It("A namespace to has been added but we don't need to watch it", func() {
				newSnapshot := snapshot.Clone()
				newSnapshot.Kubenamespaces = append(newSnapshot.Kubenamespaces, &kubernetes.KubeNamespace{
					KubeNamespace: external_kubernetes_namespace.KubeNamespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: "skip-me-again",
							Labels: map[string]string{
								"skip":    "me",
								"common":  "label",
								"skipped": "again",
							},
						},
					},
				})
				Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeFalse())
			})

			It("A namespace to has been deleted but we don't watch it", func() {
				newSnapshot := snapshot.Clone()
				newSnapshot.Kubenamespaces = kubernetes.KubeNamespaceList{newSnapshot.Kubenamespaces[0]}
				Expect(setupSyncer.ShouldSync(context.TODO(), snapshot, &newSnapshot)).To(BeFalse())
			})
		})
	})

	When("a component in the setup run holds restarts", func() {
		var (
			settings    *v1.Settings
			snapshot    *v1.SetupSnapshot
			hold        bool
			setupErr    error
			releases    chan [2]func()
			setupSyncer *SetupSyncer
		)

		BeforeEach(func() {
			settings = &v1.Settings{Metadata: &core.Metadata{Name: "hello", Namespace: "goodbye"}}
			snapshot = &v1.SetupSnapshot{Settings: v1.SettingsList{settings}}
			hold = true
			setupErr = nil
			releases = make(chan [2]func(), 1)
			setupSyncer = NewSetupSyncer(
				settings.Metadata.Ref(),
				func(
					ctx context.Context,
					kubeCache kube.SharedCache,
					inMemoryCache memory.InMemoryResourceCache,
					settings *v1.Settings,
					identity leaderelector.Identity) error {
					if hold {
						// Two components in the same run each hold restarts.
						releases <- [2]func(){HoldRestartsUntilFirstSync(ctx), HoldRestartsUntilFirstSync(ctx)}
					}
					return setupErr
				},
				singlereplica.Identity())
		})

		// startSync runs Sync in the background and returns a channel that receives its error.
		startSync := func(ctx context.Context) <-chan error {
			done := make(chan error, 1)
			go func() { done <- setupSyncer.Sync(ctx, snapshot) }()
			return done
		}
		// release receives the setup run's release funcs and calls the first.
		release := func() [2]func() {
			var r [2]func()
			Eventually(releases, "2s").Should(Receive(&r))
			r[0]()
			return r
		}

		It("returns from Sync once the hold is released", func() {
			done := startSync(context.Background())
			Consistently(done, "500ms").ShouldNot(Receive())

			release()
			Eventually(done, "2s").Should(Receive(BeNil()))
		})

		It("returns from Sync after the configured timeout if the hold is never released", func() {
			GinkgoT().Setenv(FirstSyncTimeoutEnv, "200ms")
			done := startSync(context.Background())
			Eventually(done, "2s").Should(Receive(BeNil()))
		})

		It("returns from Sync when its context ends", func() {
			ctx, cancel := context.WithCancel(context.Background())
			done := startSync(ctx)
			Consistently(done, "200ms").ShouldNot(Receive())

			cancel()
			Eventually(done, "2s").Should(Receive(BeNil()))
		})

		It("shares one release between callers in the same run", func() {
			done := startSync(context.Background())
			Consistently(done, "200ms").ShouldNot(Receive())

			r := release()
			Expect(r[1]).NotTo(Panic())
			Eventually(done, "2s").Should(Receive(BeNil()))
		})

		When("the setup fails", func() {
			BeforeEach(func() { setupErr = errors.New("setup failed") })

			It("returns the error without waiting", func() {
				Eventually(startSync(context.Background()), "2s").Should(Receive(MatchError(setupErr)))
			})
		})

		When("the setup does not hold restarts", func() {
			BeforeEach(func() { hold = false })

			It("returns from Sync without waiting", func() {
				Eventually(startSync(context.Background()), "2s").Should(Receive(BeNil()))
			})
		})
	})

	// Exercises the generated setup emitter and event loop with in-memory clients.
	When("running in the setup event loop with watchNamespaceSelectors", func() {
		type setupCall struct {
			ctx           context.Context
			namespaces    []string
			firstSyncDone func()
		}

		It("does not restart setup for a namespace change until the running setup has run its first sync", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			rcFactory := &factory.MemoryResourceClientFactory{Cache: memory.NewInMemoryResourceCache()}
			settingsClient, err := v1.NewSettingsClient(ctx, rcFactory)
			Expect(err).NotTo(HaveOccurred())
			namespaceClient, err := kubernetes.NewKubeNamespaceClient(ctx, rcFactory)
			Expect(err).NotTo(HaveOccurred())

			settings := &v1.Settings{
				Metadata: &core.Metadata{Name: "default", Namespace: "gloo-system"},
				WatchNamespaceSelectors: []*v1.LabelSelector{{
					MatchLabels: map[string]string{"gloo-watch": "enabled"},
				}},
			}
			_, err = settingsClient.Write(settings, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())
			writeWatchedNamespace := func(name string) {
				_, err := namespaceClient.Write(&kubernetes.KubeNamespace{
					KubeNamespace: external_kubernetes_namespace.KubeNamespace{
						ObjectMeta: metav1.ObjectMeta{Name: name, Labels: map[string]string{"gloo-watch": "enabled"}},
					},
				}, clients.WriteOpts{Ctx: ctx})
				Expect(err).NotTo(HaveOccurred())
			}
			writeWatchedNamespace("ns-1")

			var (
				mu    sync.Mutex
				calls []setupCall
			)
			getCalls := func() []setupCall {
				mu.Lock()
				defer mu.Unlock()
				return append([]setupCall(nil), calls...)
			}
			emitter := v1.NewSetupEmitter(settingsClient, namespaceClient)
			setupSyncer := NewSetupSyncer(
				settings.GetMetadata().Ref(),
				func(
					ctx context.Context,
					kubeCache kube.SharedCache,
					inMemoryCache memory.InMemoryResourceCache,
					settings *v1.Settings,
					identity leaderelector.Identity) error {
					mu.Lock()
					defer mu.Unlock()
					calls = append(calls, setupCall{
						ctx:           ctx,
						namespaces:    settingsutil.GetNamespacesToWatch(settings),
						firstSyncDone: HoldRestartsUntilFirstSync(ctx),
					})
					return nil
				},
				singlereplica.Identity())
			errs, err := v1.NewSetupEventLoop(emitter, setupSyncer).Run([]string{"gloo-system"}, clients.WatchOpts{
				Ctx:         ctx,
				RefreshRate: 100 * time.Millisecond,
			})
			Expect(err).NotTo(HaveOccurred())
			go func() {
				for range errs {
				}
			}()

			Eventually(getCalls, "5s").Should(HaveLen(1))
			first := getCalls()[0]
			Expect(first.namespaces).To(ConsistOf("gloo-system", "ns-1"))

			writeWatchedNamespace("ns-2")
			Consistently(func(g Gomega) {
				g.Expect(first.ctx.Err()).NotTo(HaveOccurred(), "setup ctx cancelled before its first sync")
				g.Expect(getCalls()).To(HaveLen(1))
			}, "2s").Should(Succeed())
			writeWatchedNamespace("ns-3")
			Consistently(getCalls, "1s").Should(HaveLen(1))

			// Both changes made during the hold restart setup once, with the latest namespaces.
			first.firstSyncDone()
			Eventually(first.ctx.Done(), "5s").Should(BeClosed())
			Eventually(getCalls, "5s").Should(HaveLen(2))
			Expect(getCalls()[1].namespaces).To(ConsistOf("gloo-system", "ns-1", "ns-2", "ns-3"))
			getCalls()[1].firstSyncDone()
			Consistently(getCalls, "1s").Should(HaveLen(2))
		})
	})
})
