package core

func (r ResourceRef) Strings() (string, string) {
	return r.Namespace, r.Name
}
