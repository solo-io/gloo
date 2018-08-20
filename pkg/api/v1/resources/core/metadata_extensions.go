package core

func (m Metadata) Less(m2 Metadata) bool {
	if m.Namespace == m2.Namespace {
		return m.Name < m2.Name
	}
	return m.Namespace < m2.Namespace
}

func (m Metadata) ObjectRef() (string, string) {
	return m.Namespace, m.Name
}
