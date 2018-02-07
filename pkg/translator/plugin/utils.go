package plugin

type SimpleDependenciesDescription struct {
	secretRefs []string
}

func (s *SimpleDependenciesDescription) AddSecretRef(sr string) {
	s.secretRefs = append(s.secretRefs, sr)
}

func (s *SimpleDependenciesDescription) SecretRefs() []string {
	return s.secretRefs
}
