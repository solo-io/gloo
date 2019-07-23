package generate

type Dependency struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Repository string   `json:"repository"`
	Condition  string   `json:"condition,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type DependencyList struct {
	Dependencies []Dependency `json:"dependencies"`
}
