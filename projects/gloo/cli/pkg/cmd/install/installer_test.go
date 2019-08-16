package install_test

import (
	"context"

	install3 "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/install"
	options2 "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"

	"github.com/solo-io/gloo/pkg/cliutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	install2 "github.com/solo-io/gloo/pkg/cliutil/install"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/install"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
)

type MockInstallClient struct {
	expectedCrds     []string
	applied          bool
	waited           bool
	resources        []install2.ResourceType
	knativeInstalled bool
	knativeOurs      bool
}

func (i *MockInstallClient) KubectlApply(manifest []byte) error {
	Expect(i.applied).To(BeFalse())
	i.applied = true
	resources, err := install2.GetResources(string(manifest))
	Expect(err).NotTo(HaveOccurred())
	i.resources = resources
	return nil
}

func (i *MockInstallClient) WaitForCrdsToBeRegistered(ctx context.Context, crds []string) error {
	Expect(i.waited).To(BeFalse())
	i.waited = true
	Expect(crds).To(ConsistOf(i.expectedCrds))
	return nil
}

func (i *MockInstallClient) CheckKnativeInstallation() (bool, bool, error) {
	return i.knativeInstalled, i.knativeOurs, nil
}

var _ = Describe("Install", func() {

	var (
		installer install.GlooStagedInstaller
		opts      options.Options
		extOpts   options2.ExtraOptions
		validator MockInstallClient
	)

	BeforeEach(func() {
		opts.Install.Namespace = "gloo-system"
		opts.Install.HelmChartOverride = file
		extOpts.Install.LicenseKey = "eyJleHAiOjE1NTQ1MTYyNTEsImlhdCI6MTU1MTgzNzg1MSwiayI6ImVqMVYyUSJ9.5lDPOuRWo4_qr3r9PXBv6lYIut3DbBrqqRauwSQZm4E"
	})

	expectKinds := func(resources []install2.ResourceType, kinds []string) {
		for _, resource := range resources {
			ExpectWithOffset(1, kinds).To(ContainElement(resource.Kind))
		}
	}

	expectNames := func(resources []install2.ResourceType, names []string) {
		for _, resource := range resources {
			ExpectWithOffset(1, names).To(ContainElement(resource.Metadata.Name))
		}
	}

	expectLabels := func(resources []install2.ResourceType, labels map[string]string) {
		for _, resource := range resources {
			actualLabels := resource.Metadata.Labels
			for k, v := range labels {
				val, ok := actualLabels[k]
				ExpectWithOffset(1, ok).To(BeTrue())
				// Currently all the observability stuff has custom app labels
				if val == "grafana" || val == "prometheus" || val == "gloo-ee" {
					continue
				}
				ExpectWithOffset(1, v).To(BeEquivalentTo(val))
			}
		}
	}

	expectNamespace := func(resources []install2.ResourceType, namespace string) {
		globalKinds := []string{
			"Namespace",
			"ClusterRole",
			"ClusterRoleBinding",
		}
		for _, resource := range resources {
			if cliutil.Contains(globalKinds, resource.TypeMeta.Kind) {
				continue
			}
			ExpectWithOffset(1, resource.Metadata.Namespace).To(BeEquivalentTo(namespace))
		}
	}

	Context("Gateway with default values", func() {
		BeforeEach(func() {
			spec, err := install3.GetInstallSpec(&opts, &extOpts)
			Expect(err).NotTo(HaveOccurred())
			validator = MockInstallClient{
				expectedCrds: install.GlooCrdNames,
			}
			installer, err = install.NewGlooStagedInstaller(&opts, *spec, &validator)
			Expect(err).NotTo(HaveOccurred())
		})

		It("installs expected crds for gloo", func() {
			err := installer.DoCrdInstall()
			Expect(err).NotTo(HaveOccurred())
			Expect(validator.applied).To(BeTrue())
			Expect(validator.waited).To(BeTrue())
			expectKinds(validator.resources, []string{"CustomResourceDefinition"})
			expectNames(validator.resources, install.GlooCrdNames)
		})

		It("only prepares the namespace CRD preinstall", func() {
			err := installer.DoPreInstall()
			Expect(err).NotTo(HaveOccurred())
			// the namespace CRD triggers the "applied" property
			Expect(validator.applied).To(BeTrue())
			Expect(validator.waited).To(BeFalse())
			expectKinds(validator.resources, install.GlooPreInstallKinds)
			expectLabels(validator.resources, install.ExpectedLabels)
			expectNamespace(validator.resources, "gloo-system")
		})

		It("installs expected kinds for gloo", func() {
			err := installer.DoInstall()
			Expect(err).NotTo(HaveOccurred())
			Expect(validator.applied).To(BeTrue())
			Expect(validator.waited).To(BeFalse())
			glooeInstallKinds := append(install.GlooInstallKinds,
				"PodSecurityPolicy", "Secret", "ServiceAccount", "ClusterRole", "ClusterRoleBinding", "Role", "RoleBinding", "Upstream", "PersistentVolumeClaim", "Settings")
			expectKinds(validator.resources, glooeInstallKinds)
			expectLabels(validator.resources, install.ExpectedLabels)
		})

	})
})
