package graphql

import (
	"context"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

type settingsQueryResolver struct{ *ApiResolver }

func (r *settingsQueryResolver) Get(ctx context.Context, obj *customtypes.SettingsQuery) (*models.Settings, error) {
	namespace := defaults.GlooSystem
	name := defaults.SettingsName
	settings, err := r.Settings.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSettings(settings), nil
}

type settingsMutationResolver struct{ *ApiResolver }

func (r *settingsMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.SettingsMutation, settings *v1.Settings) (*models.Settings, error) {
	out, err := r.Settings.Write(settings, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSettings(out), nil
}
func (r *settingsMutationResolver) Update(ctx context.Context, obj *customtypes.SettingsMutation, rawUpdates models.InputSettings) (*models.Settings, error) {
	updates, err := NewConverter(r.ApiResolver, ctx).ConvertInputSettings(rawUpdates)
	if err != nil {
		return nil, err
	}

	namespace := updates.Metadata.Namespace
	name := updates.Metadata.Name
	settings, err := r.Settings.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}

	// preserve the given metadata to ensure request was made with latest resourceVersion
	settings.Metadata = updates.Metadata

	// only apply changes to the provided fields
	if updates.RefreshRate != nil {
		settings.RefreshRate = updates.RefreshRate
	}
	if updates.WatchNamespaces != nil {
		settings.WatchNamespaces = updates.WatchNamespaces
	}
	return r.write(true, ctx, obj, settings)
}
