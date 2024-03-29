package version

import (
	"bytes"
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/protoutils"
	gloo_version "github.com/solo-io/gloo/pkg/version"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	mock_version "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/version/mocks"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/printers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/version"
)

var _ = Describe("version command", func() {
	var (
		ctrl   *gomock.Controller
		client *mock_version.MockServerVersion
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(T)
		client = mock_version.NewMockServerVersion(ctrl)
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() { cancel() })

	Context("getVersion", func() {
		It("will include an error status if an error occurs while getting the server version", func() {
			fakeErr := eris.New("test")
			client.EXPECT().Get(ctx).Return(nil, fakeErr).Times(1)
			client.EXPECT().GetClusterVersion().Return(nil, nil).Times(1)
			sv := GetClientServerVersions(ctx, client)
			Expect(sv.GetStatus().GetError()).NotTo(BeNil())
			Expect(sv.GetStatus().GetError().GetErrors()).To(ContainElements(UndefinedServer, fakeErr.Error()))
		})
		It("will include an error status if an error occurs while getting the kube cluster version", func() {
			fakeErr := eris.New("test")
			v := make([]*version.ServerVersion, 1)
			client.EXPECT().Get(ctx).Return(v, nil).Times(1)
			client.EXPECT().GetClusterVersion().Return(nil, fakeErr).Times(1)
			sv := GetClientServerVersions(ctx, client)
			Expect(sv.GetStatus().GetError()).NotTo(BeNil())
			Expect(sv.GetStatus().GetError().GetWarnings()).To(ContainElements(UndefinedK8s, fakeErr.Error()))
		})
		It("can get the server version", func() {
			v := make([]*version.ServerVersion, 1)
			client.EXPECT().Get(ctx).Return(v, nil).Times(1)
			client.EXPECT().GetClusterVersion().Return(nil, nil).Times(1)
			vrs := GetClientServerVersions(ctx, client)
			Expect(vrs.GetStatus().GetError()).To(BeNil())
			Expect(v).To(Equal(vrs.Server))
		})
		It("can get kubernetes version", func() {
			expectedVrs := make([]*version.ServerVersion, 1)
			expectedK8s := &version.KubernetesClusterVersion{}
			client.EXPECT().Get(ctx).Return(expectedVrs, nil).Times(1)
			client.EXPECT().GetClusterVersion().Return(expectedK8s, nil).Times(1)
			vrs := GetClientServerVersions(ctx, client)
			Expect(vrs.GetStatus().GetError()).To(BeNil())
			Expect(expectedVrs).To(Equal(vrs.Server))
			Expect(expectedK8s).To(Equal(vrs.KubernetesCluster))
		})
	})

	var (
		namespace = "gloo-system"
	)
	type (
		printTestInput struct {
			serverVersion     []*version.ServerVersion
			k8sClusterVersion *version.KubernetesClusterVersion
			outputType        printers.OutputType
			enterprise        bool
		}
		expectedTestOutput struct {
			err            error
			vrs            *version.Version
			errorStrings   []string // status errors expected to match all of these strings
			warningStrings []string // status warnings expected to match all of these strings
		}
	)
	defaultServerVersion := func(enterprise bool) []*version.ServerVersion {
		return []*version.ServerVersion{{
			Enterprise: enterprise,
			Type:       version.GlooType_Gateway,
			VersionType: &version.ServerVersion_Kubernetes{
				Kubernetes: &version.Kubernetes{
					Containers: []*version.Kubernetes_Container{
						{
							Tag:      "v0.0.1",
							Name:     "gloo",
							Registry: "quay.io/solo-io",
						},
						{
							Tag:      "v0.0.2",
							Name:     "gateway",
							Registry: "quay.io/solo-io",
						},
					},
					Namespace: namespace,
				},
			},
		}}
	}

	defaultKubeClusterVersion := func() *version.KubernetesClusterVersion {
		return &version.KubernetesClusterVersion{
			Major:      "1",
			Minor:      "26",
			GitVersion: "v1.26.14",
			BuildDate:  "2024-02-29T21:46:10Z",
			Platform:   "linux/amd64",
		}
	}

	defaultOutputVersion := func(enterprise bool) *version.Version {
		return &version.Version{
			Client: &version.ClientVersion{
				Version: gloo_version.Version,
			},
			Server:            defaultServerVersion(enterprise),
			KubernetesCluster: defaultKubeClusterVersion(),
			Status: &version.Status{
				Status: &version.Status_Ok{
					Ok: &version.Status_OkStatus{},
				},
			},
		}
	}

	entries := func(enterprise bool, outputType printers.OutputType) []TableEntry {
		repoPrefix := "oss"
		if enterprise {
			repoPrefix = "ee"
		}
		printerSuffix := "json"
		if outputType == printers.YAML {
			printerSuffix = "yaml"
		}
		entryName := func(name string) string {
			return fmt.Sprintf("%s %s %s", repoPrefix, printerSuffix, name)
		}
		return []TableEntry{
			Entry(entryName("complete"), &printTestInput{
				serverVersion:     defaultServerVersion(enterprise),
				k8sClusterVersion: defaultKubeClusterVersion(),
				outputType:        outputType,
				enterprise:        enterprise,
			}, &expectedTestOutput{
				vrs: defaultOutputVersion(enterprise),
			}),
			Entry(entryName("no server version"), &printTestInput{
				serverVersion:     nil,
				k8sClusterVersion: defaultKubeClusterVersion(),
				outputType:        outputType,
				enterprise:        enterprise,
			}, &expectedTestOutput{
				vrs: &version.Version{
					Client: &version.ClientVersion{
						Version: gloo_version.Version,
					},
					Server:            nil,
					KubernetesCluster: defaultKubeClusterVersion(),
					Status: &version.Status{
						Status: &version.Status_Error{
							Error: &version.Status_ErrorStatus{
								Errors: []string{UndefinedServer},
							},
						},
					},
				},
				errorStrings: []string{UndefinedServer},
			}),
			Entry(entryName("no kube cluster"), &printTestInput{
				serverVersion:     defaultServerVersion(enterprise),
				k8sClusterVersion: nil,
				outputType:        outputType,
				enterprise:        enterprise,
			}, &expectedTestOutput{
				vrs: &version.Version{
					Client: &version.ClientVersion{
						Version: gloo_version.Version,
					},
					Server:            defaultServerVersion(enterprise),
					KubernetesCluster: nil,
					Status: &version.Status{
						Status: &version.Status_Error{
							Error: &version.Status_ErrorStatus{
								Warnings: []string{UndefinedK8s},
							},
						},
					},
				},
				warningStrings: []string{UndefinedK8s},
			}),
			Entry(entryName("no server version or kube cluster"), &printTestInput{
				serverVersion:     nil,
				k8sClusterVersion: nil,
				outputType:        outputType,
				enterprise:        enterprise,
			}, &expectedTestOutput{
				vrs: &version.Version{
					Client: &version.ClientVersion{
						Version: gloo_version.Version,
					},
					Server:            nil,
					KubernetesCluster: nil,
					Status: &version.Status{
						Status: &version.Status_Error{
							Error: &version.Status_ErrorStatus{
								Errors:   []string{UndefinedServer},
								Warnings: []string{UndefinedK8s},
							},
						},
					},
				},
				errorStrings:   []string{UndefinedServer},
				warningStrings: []string{UndefinedK8s},
			}),
		}
	}
	allEntries := func() []TableEntry {
		all := make([]TableEntry, 0)
		all = append(all, entries(false, printers.JSON)...) // oss json
		all = append(all, entries(true, printers.JSON)...)  // oss yaml
		all = append(all, entries(false, printers.YAML)...) // ee json
		all = append(all, entries(true, printers.YAML)...)  // ee yaml
		return all
	}()
	FDescribeTable("printing machine readable",
		func(inp *printTestInput, expected *expectedTestOutput) {
			buf := &bytes.Buffer{}
			opts := &options.Options{
				Top: options.Top{
					Output: inp.outputType,
				},
			}

			if inp.serverVersion != nil {
				client.EXPECT().Get(nil).Times(1).Return(inp.serverVersion, nil)
			} else {
				client.EXPECT().Get(nil).Times(1).Return(nil, eris.New("mock err"))
			}

			if inp.k8sClusterVersion != nil {
				client.EXPECT().GetClusterVersion().Times(1).Return(inp.k8sClusterVersion, nil)
			} else {
				client.EXPECT().GetClusterVersion().Times(1).Return(nil, eris.New("mock err"))
			}

			err := printVersion(client, buf, opts)
			Expect(err).NotTo(HaveOccurred())

			var receivedVer version.Version
			switch inp.outputType {
			case printers.YAML:
				err = protoutils.UnmarshalYaml(buf.Bytes(), &receivedVer)
			case printers.JSON:
				err = protoutils.UnmarshalBytes(buf.Bytes(), &receivedVer)
			default:
				Fail("unhandled outputType")
			}
			Expect(err).NotTo(HaveOccurred())

			Expect(receivedVer.GetClient().GetVersion()).To(Equal(gloo_version.Version))
			Expect(receivedVer.GetServer()).To(Equal(expected.vrs.GetServer()))
			Expect(receivedVer.GetKubernetesCluster()).To(Equal(expected.vrs.GetKubernetesCluster()))
			if len(expected.warningStrings)+len(expected.errorStrings) == 0 {
				Expect(receivedVer.GetStatus().GetError()).To(BeNil())
			} else {
				Expect(receivedVer.GetStatus().GetError()).NotTo(BeNil())
				Expect(receivedVer.GetStatus().GetError().GetWarnings()).To(ContainElements(expected.warningStrings))
				Expect(receivedVer.GetStatus().GetError().GetErrors()).To(ContainElements(expected.errorStrings))
			}
		},
		allEntries,
	)
	// TODO(jbohanon): test kube features on tables
	Context("printing human readable", func() {
		var (
			sv  *version.ServerVersion
			buf *bytes.Buffer
		)
		BeforeEach(func() {
			buf = &bytes.Buffer{}

			sv = &version.ServerVersion{
				Type: version.GlooType_Gateway,
				VersionType: &version.ServerVersion_Kubernetes{
					Kubernetes: &version.Kubernetes{
						Containers: []*version.Kubernetes_Container{
							{
								Tag:      "v0.0.1",
								Name:     "gloo",
								Registry: "quay.io/solo-io",
							},
							{
								Tag:      "v0.0.2",
								Name:     "gateway",
								Registry: "quay.io/solo-io",
							},
						},
						Namespace: namespace,
					},
				},
			}
		})

		var osTableOutput = fmt.Sprintf(`Client version: %s
Server version:
+-------------+-----------------+-----------------+
|  NAMESPACE  | DEPLOYMENT-TYPE |   CONTAINERS    |
+-------------+-----------------+-----------------+
| gloo-system | Gateway         | gloo: v0.0.1    |
|             |                 | gateway: v0.0.2 |
+-------------+-----------------+-----------------+
`, gloo_version.Version)

		var eTableOutput = fmt.Sprintf(`Client version: %s
Server version:
+-------------+--------------------+-----------------+
|  NAMESPACE  |  DEPLOYMENT-TYPE   |   CONTAINERS    |
+-------------+--------------------+-----------------+
| gloo-system | Gateway Enterprise | gloo: v0.0.1    |
|             |                    | gateway: v0.0.2 |
+-------------+--------------------+-----------------+
`, gloo_version.Version)

		tests := []struct {
			name       string
			result     string
			outputType printers.OutputType
			enterprise bool
		}{
			{
				name:       "table",
				result:     osTableOutput,
				outputType: printers.TABLE,
				enterprise: false,
			},
			{
				name:       "enterprise table",
				result:     eTableOutput,
				outputType: printers.TABLE,
				enterprise: true,
			},
		}

		for _, test := range tests {
			test := test
			Context(test.name, func() {
				It("can translate with valid server version", func() {
					opts := &options.Options{
						Top: options.Top{
							Output: test.outputType,
						},
					}
					sv.Enterprise = test.enterprise
					client.EXPECT().Get(nil).Times(1).Return([]*version.ServerVersion{sv}, nil)
					client.EXPECT().GetClusterVersion().Times(1).Return(nil, nil)
					err := printVersion(client, buf, opts)
					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(test.result))
				})

				It("can translate with nil server version", func() {
					opts := &options.Options{
						Top: options.Top{
							Output: test.outputType,
						},
					}
					client.EXPECT().Get(nil).Times(1).Return(nil, eris.Errorf("fake rbac error"))
					client.EXPECT().GetClusterVersion().Times(1).Return(nil, nil)
					err := printVersion(client, buf, opts)
					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(ContainSubstring(UndefinedServer))
				})
			})
		}
	})

})
