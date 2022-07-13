package apiserverutils

import (
	"context"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// These are the defaults we use if the Settings resource does not have the options specified
	readOnlyDefault           = true
	apiExplorerEnabledDefault = true
)

// Returns the Settings from the current Gloo Edge installation.
func GetSettings(ctx context.Context, settingsClient gloo_v1.SettingsClient) (*gloo_v1.Settings, error) {
	return settingsClient.GetSettings(ctx, client.ObjectKey{
		Namespace: GetInstallNamespace(),
		Name:      defaults.SettingsName,
	})
}

// Returns the console (UI) options from the current Gloo Edge installation's settings.
func GetConsoleOptions(ctx context.Context, settingsClient gloo_v1.SettingsClient) (*rpc_edge_v1.ConsoleOptions, error) {
	settings, err := GetSettings(ctx, settingsClient)
	if err != nil {
		return nil, err
	}

	var readOnly = readOnlyDefault
	if settings.Spec.GetConsoleOptions().GetReadOnly() != nil {
		readOnly = settings.Spec.GetConsoleOptions().GetReadOnly().GetValue()
	}
	var apiExplorerEnabled = apiExplorerEnabledDefault
	if settings.Spec.GetConsoleOptions().GetApiExplorerEnabled() != nil {
		apiExplorerEnabled = settings.Spec.GetConsoleOptions().GetApiExplorerEnabled().GetValue()
	}

	return &rpc_edge_v1.ConsoleOptions{
		ReadOnly:           readOnly,
		ApiExplorerEnabled: apiExplorerEnabled,
	}, nil
}

// Throws an error if this Gloo Edge UI instance is read-only
func CheckUpdatesAllowed(ctx context.Context, settingsClient gloo_v1.SettingsClient) error {
	consoleOptions, err := GetConsoleOptions(ctx, settingsClient)
	if err != nil {
		return err
	}

	if consoleOptions.GetReadOnly() {
		return eris.New("Cannot perform update: UI is read-only.")
	}
	return nil
}
