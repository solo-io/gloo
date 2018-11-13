package customtypes

type UpstreamQuery struct {
	Namespace string
}

type UpstreamMutation struct {
	Namespace string
}

type VirtualServiceQuery struct {
	Namespace string
}

type VirtualServiceMutation struct {
	Namespace string
}

type ResolverMapQuery struct {
	Namespace string
}

type ResolverMapMutation struct {
	Namespace string
}

type SchemaQuery struct {
	Namespace string
}

type SchemaMutation struct {
	Namespace string
}

type SecretQuery struct {
	Namespace string
}

type SecretMutation struct {
	Namespace string
}

type ArtifactQuery struct {
	Namespace string
}

type ArtifactMutation struct {
	Namespace string
}

type SettingsQuery struct{}

type SettingsMutation struct{}

type VcsMutation struct {
	Username string
}
