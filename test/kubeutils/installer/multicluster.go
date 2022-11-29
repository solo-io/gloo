package installer

import (
	"context"

	"github.com/solo-io/solo-projects/test/kubeutils"

	"github.com/onsi/ginkgo"

	"golang.org/x/sync/errgroup"
)

type MultiClusterInstaller struct {
	installersByCluster map[string][]Installer

	orchestrator kubeutils.Orchestrator
}

func NewMultiCluster(orchestrator kubeutils.Orchestrator) *MultiClusterInstaller {
	return &MultiClusterInstaller{
		installersByCluster: map[string][]Installer{},
		orchestrator:        orchestrator,
	}
}

func (m *MultiClusterInstaller) GetContext() (string, string) {
	return "", ""
}

func (m *MultiClusterInstaller) Install(ctx context.Context) error {
	for clusterName, installers := range m.GetInstallers() {
		if err := m.doInstallBatch(ctx, clusterName, installers); err != nil {
			return err
		}
	}

	return nil
}

func (m *MultiClusterInstaller) doInstallBatch(ctx context.Context, clusterName string, clusterInstallers []Installer) error {
	errorGroup, ctx := errgroup.WithContext(ctx)

	if err := m.orchestrator.SetClusterContext(ctx, clusterName); err != nil {
		return err
	}

	for _, clusterInstaller := range clusterInstallers {
		clusterInstallerCpy := clusterInstaller
		errorGroup.Go(func() error {
			defer ginkgo.GinkgoRecover()
			return clusterInstallerCpy.Install(ctx)

		})
	}

	// This will throw an error if any of the underlying installers failed
	return errorGroup.Wait()
}

func (m *MultiClusterInstaller) Uninstall(ctx context.Context) error {
	for clusterName, installers := range m.GetInstallers() {
		if err := m.doUninstallBatch(ctx, clusterName, installers); err != nil {
			return err
		}
	}

	return nil
}

func (m *MultiClusterInstaller) doUninstallBatch(ctx context.Context, clusterName string, clusterInstallers []Installer) error {
	errorGroup, ctx := errgroup.WithContext(ctx)

	if err := m.orchestrator.SetClusterContext(ctx, clusterName); err != nil {
		return err
	}

	for _, clusterInstaller := range clusterInstallers {
		clusterInstallerCpy := clusterInstaller
		errorGroup.Go(func() error {
			defer ginkgo.GinkgoRecover()
			return clusterInstallerCpy.Uninstall(ctx)

		})
	}

	// This will throw an error if any of the underlying installers failed
	return errorGroup.Wait()
}

func (m *MultiClusterInstaller) RegisterInstaller(installer Installer) {
	installCluster, _ := installer.GetContext()
	m.installersByCluster[installCluster] = append(m.installersByCluster[installCluster], installer)
}

func (m *MultiClusterInstaller) GetInstallers() map[string][]Installer {
	return m.installersByCluster
}
