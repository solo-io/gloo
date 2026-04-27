// Package gomod contains tests that assert invariants on the repository's
// go.mod file.
//
// The go-control-plane module is a complex repository whose proto
// bindings must track the Envoy runtime version that Gloo ships
// with. Newer Envoy releases may break wire and API compatibility
// with older proto definitions, and go-control-plane does not publish
// a compatibility matrix that lets us mechanically determine which
// proto version pairs with which Envoy runtime. In practice we have
// to pin both the replace directive for
// `github.com/envoyproxy/go-control-plane` and the envoy submodule by
// hand after verifying against the Envoy binary we actually deploy.
//
// The pin is enforced through `replace` directives, not `exclude`:
// `replace` redirects every version of a module path to our chosen
// pseudo-version, so transitive bumps (e.g. from a grpc CVE patch) are
// silently rewritten at build time. The `require` lines in go.mod may
// still show the transitively-demanded versions — that is MVS metadata
// and does not affect what code runs — so this test deliberately looks
// at the replace target as the source of truth, falling back to require
// only when no replace is in place.
package gomod_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

const (
	goControlPlaneModule = "github.com/envoyproxy/go-control-plane"
	// maxGoControlPlaneVersion is exclusive: the effective pinned version
	// must compare strictly less than this value.
	//
	// Do not raise this without careful analysis of go-control-plane.
	maxGoControlPlaneVersion = "v0.33.0"
)

func TestGoControlPlaneVersionBelowCap(t *testing.T) {
	mod := loadGoMod(t)

	pinned, ok := effectivePin(mod, goControlPlaneModule)
	if !ok {
		t.Fatalf("could not find %s in go.mod (neither require nor replace matched)", goControlPlaneModule)
	}

	canonical := semver.Canonical(pinned.version)
	if !semver.IsValid(canonical) {
		t.Fatalf("%s %s is %q, which is not a valid semver string", goControlPlaneModule, pinned.origin, pinned.version)
	}

	if semver.Compare(canonical, maxGoControlPlaneVersion) >= 0 {
		t.Fatalf(
			"%s %s is %s, which is >= %s.\n\n"+
				"go-control-plane is a complex repository whose proto bindings must track "+
				"the Envoy runtime version Gloo ships with. Backwards compatibility across "+
				"go-control-plane versions is not documented, so bumping past the cap must "+
				"be done deliberately after verifying against the running Envoy binary. "+
				"If this bump is intentional, update maxGoControlPlaneVersion in %s along "+
				"with the pins in go.mod.",
			goControlPlaneModule, pinned.origin, pinned.version, maxGoControlPlaneVersion, thisFile(),
		)
	}
}

// pinnedMinors captures the major.minor that the effective pinned version
// of each go-control-plane module is allowed to use. The "effective pin"
// is the replace target if a replace directive is in play, otherwise the
// require version. This matches what the build actually resolves.
//
// A mismatch — including a minor bump introduced transitively by an
// unrelated upgrade — fails the test.
//
// Entries are keyed by module path. The value is the expected MAJOR.MINOR
// prefix (e.g. "v1.32"); patch and pseudo-version suffixes are ignored.
//
// An empty value means "no minor has been approved yet": if the module
// ever appears in go.mod, the test fails so the developer must deliberately
// choose and record the minor here.
var pinnedMinors = map[string]string{
	// Root module. Pinned at v0.13 via the replace in go.mod. Transitive
	// require bumps to v0.14+ (historically via grpc CVE upgrades) are
	// rewritten to v0.13.x by that replace.
	"github.com/envoyproxy/go-control-plane": "v0.13",
	// envoy: replace directive pins to v1.32.5-pseudo to match the Envoy
	// v1.32 runtime shipped on the v1.18.x branch. require may drift up
	// (currently v1.36.0) under MVS pressure from grpc; replace wins.
	"github.com/envoyproxy/go-control-plane/envoy": "v1.32",
	// contrib: proto bindings must stay aligned with the envoy submodule
	// above, so they share the v1.32 minor. No replace is needed today
	// because contrib isn't being transitively bumped — the require is
	// the pin. If that changes, add a replace in go.mod.
	"github.com/envoyproxy/go-control-plane/contrib": "v1.32",
	// ratelimit: independent release cadence from envoy/contrib. Pinned
	// at v0.1 so a breaking API reshuffle under v0 semver cannot slip in
	// via a transitive bump.
	"github.com/envoyproxy/go-control-plane/ratelimit": "v0.1",
	// xdsmatcher: not currently a direct dependency. Listed here so that
	// if anything pulls it into go.mod the test fails until a reviewer
	// records the deliberately chosen minor.
	"github.com/envoyproxy/go-control-plane/xdsmatcher": "",
}

func TestGoControlPlaneMinorVersionsPinned(t *testing.T) {
	mod := loadGoMod(t)

	// Sort keys for deterministic subtest order and output.
	paths := make([]string, 0, len(pinnedMinors))
	for p := range pinnedMinors {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, path := range paths {
		expectedMinor := pinnedMinors[path]
		t.Run(path, func(t *testing.T) {
			pinned, ok := effectivePin(mod, path)
			if !ok {
				if expectedMinor == "" {
					// Module not in go.mod and no baseline recorded — nothing to enforce yet.
					t.Skipf("%s is not a dependency and has no pinned minor; nothing to check", path)
				}
				// Having a pinned minor but no entry in go.mod means the module
				// was removed. That is fine, but the reviewer should clear the
				// pin so the intent stays accurate.
				t.Skipf("%s is no longer in go.mod; consider removing the pin from pinnedMinors", path)
			}

			if expectedMinor == "" {
				t.Fatalf(
					"%s is now present in go.mod (%s %s), but pinnedMinors has no approved minor for it.\n\n"+
						"Bumping, adding, or removing go-control-plane modules must be a deliberate decision "+
						"because their proto bindings couple tightly to the Envoy runtime and their "+
						"cross-minor compatibility is not documented. Update %s to record the intended minor.",
					path, pinned.origin, pinned.version, thisFile(),
				)
			}

			canonical := semver.Canonical(pinned.version)
			if !semver.IsValid(canonical) {
				t.Fatalf("%s %s %q is not a valid semver string", path, pinned.origin, pinned.version)
			}

			actualMinor := semver.MajorMinor(canonical)
			if actualMinor != expectedMinor {
				t.Fatalf(
					"%s %s is %s (major.minor %s), but pinnedMinors requires %s.\n\n"+
						"A minor-version change on a go-control-plane module is never safe to accept "+
						"implicitly: the proto bindings must line up with the Envoy runtime in use, and "+
						"go-control-plane does not publish a compatibility matrix we can check against. "+
						"Hold the minor by updating the replace directive in go.mod (preferred — it "+
						"redirects any transitive bump cleanly), or by lowering the require line if "+
						"there is no replace. If this bump is intentional, update %s and verify xDS "+
						"against the running Envoy binary before merging.",
					path, pinned.origin, pinned.version, actualMinor, expectedMinor, thisFile(),
				)
			}
		})
	}
}

// modulePin describes the version that the build will actually use for a
// module path, along with a human-readable origin ("replace target" or
// "require") for error messages.
type modulePin struct {
	origin  string
	version string
}

// effectivePin returns the version the build will resolve for path. A
// replace directive that targets path wins over a require — which matches
// Go's own resolution and means this test treats `replace` as the
// authoritative pinning mechanism. Only when no replace is in play does
// the require line become the pin.
func effectivePin(mod *modfile.File, path string) (modulePin, bool) {
	for _, r := range mod.Replace {
		if r.New.Path == path {
			return modulePin{
				origin:  fmt.Sprintf("replace target (%s => %s)", r.Old.Path, r.New.Path),
				version: r.New.Version,
			}, true
		}
	}
	for _, r := range mod.Require {
		if r.Mod.Path == path {
			return modulePin{origin: "require", version: r.Mod.Version}, true
		}
	}
	return modulePin{}, false
}

// loadGoMod parses the repository's go.mod file.
func loadGoMod(t *testing.T) *modfile.File {
	t.Helper()

	path := repoGoModPath(t)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	mod, err := modfile.Parse(path, data, nil)
	if err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}
	return mod
}

// repoGoModPath walks up from this test file to the repository root and
// returns the path to go.mod. Using runtime.Caller keeps the test
// independent of the working directory it is invoked from.
func repoGoModPath(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed; cannot locate go.mod")
	}
	// test/gomod/<this file> -> repo root is two directories up.
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "go.mod")
}

func thisFile() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		return "test/gomod/gocontrolplane_version_test.go"
	}
	return f
}
