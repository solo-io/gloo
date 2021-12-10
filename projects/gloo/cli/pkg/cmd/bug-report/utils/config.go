package utils

import "time"

// BugReportConfig controls what is captured and Include in the kube-capture tool
// archive.
type BugReportConfig struct {
	// KubeConfigPath is the path to kube config file.
	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
	// Context is the cluster Context in the kube config
	Context string `json:"context,omitempty"`

	// GlooNamespace is the namespace where the istio control plane is installed.
	GlooNamespace string `json:"glooNamespace,omitempty"`

	// DryRun controls whether logs are actually captured and saved.
	DryRun bool `json:"dryRun,omitempty"`

	// FullSecrets controls whether secret contents are included.
	FullSecrets bool `json:"fullSecrets,omitempty"`

	// CommandTimeout is the maximum amount of time running the command
	// before giving up, even if not all logs are captured. Upon timeout,
	// the command creates an archive with only the logs captured so far.
	CommandTimeout time.Duration `json:"commandTimeout,omitempty"`

	//// Include is a list of SelectionSpec entries for resources to include.
	//Include SelectionSpecs `json:"include,omitempty"`
	//// Exclude is a list of SelectionSpec entries for resources t0 exclude.
	//Exclude SelectionSpecs `json:"exclude,omitempty"`

	// StartTime is the start time the the log capture time range.
	// If set, Since must be unset.
	StartTime time.Time `json:"startTime,omitempty"`
	// EndTime is the end time the the log capture time range.
	// Default is now.
	EndTime time.Time `json:"endTime,omitempty"`
	// Since defines the start time the the log capture time range.
	// StartTime is set to EndTime - Since.
	// If set, StartTime must be unset.
	Since time.Duration `json:"since,omitempty"`

	// CriticalErrors is a list of glob pattern matches for errors that,
	// if found in a log, set the highest priority for the log to ensure
	// that it is Include in the capture archive.
	CriticalErrors []string `json:"criticalErrors,omitempty"`
	// IgnoredErrors are glob error patterns which are ignored when
	// calculating the error heuristic for a log.
	IgnoredErrors []string `json:"ignoredErrors,omitempty"`
}

type SelectionSpec struct {
}
