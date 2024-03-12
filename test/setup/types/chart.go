package types

type (
	Chart struct {
		Name       string         `yaml:"name"`
		Namespace  string         `yaml:"namespace"`
		Version    string         `yaml:"version"`
		Values     map[string]any `yaml:"values"`
		Local      string         `yaml:"local"`      // Local path to the chart
		Remote     string         `yaml:"remote"`     // Remote repository or path to the chart
		Prioritize bool           `yaml:"prioritize"` // Prioritized charts will be installed before any other component.
	}
)
