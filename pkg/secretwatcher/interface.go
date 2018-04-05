package secretwatcher

type SecretMap map[string]map[string]string

// Interface is responsible for watching secrets referenced by a config
type Interface interface {
	TrackSecrets(secretRefs []string)

	// secrets are pushed here whenever they are read
	Secrets() <-chan SecretMap

	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
