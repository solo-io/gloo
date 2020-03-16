package version

var (
	UndefinedVersion = "undefined"
	// Will be set by the linker during build. Does not include "v" prefix.
	Version string
	// default version set if running "make glooctl"
	DevVersion = "dev"
)

func init() {
	if Version == "" {
		Version = UndefinedVersion
	}
}

func IsReleaseVersion() bool {
	return Version != UndefinedVersion && Version != DevVersion
}
