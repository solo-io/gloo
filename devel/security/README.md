# Security Policies

## Vulnerabilities in Third-Party Libraries

Gloo is scanned on a weekly basis for security vulnerabilities. We use [trivy](https://github.com/aquasecurity/trivy) to scan our products for security issues and we scan OSS Gloo, Gloo EE, and Portal. Scans are initiated from [a single job in OSS Gloo](https://github.com/solo-io/gloo/tree/main/.github/workflows#trivy-vulnerability-scanning). This job checks out the main branch and uses that branch's `.trivyignore.yaml` while it scans all requested images and repositories. The scans can also [be run locally](https://github.com/solo-io/gloo/tree/main/docs/cmd/securityscanutils).

Occasionally, the best approach to resolve a reported CVE is to ignore it in Trivy (ie. if we determine that the vulnerability does not affect us, or if it is raised as a false positive). To support this, we use a single `.trivyignore.yaml` file [on our main branch](https://github.com/solo-io/gloo/blob/main/.trivyignore.yaml). This file is kept up-to-date with all vulnerabilities that we feel can be safely ignored.

Entries take Trivy's YAML ignore format, which lets an entry be scoped to the artifact it was reported against:

```yaml
vulnerabilities:
  - id: CVE-2026-25681
    paths: ["usr/local/bin/kubectl"]
    statement: Why this is acceptable in this artifact.
```

Prefer `paths` whenever a finding is only acceptable in one image. An entry with no `paths` suppresses that CVE everywhere, including in binaries we build ourselves, so it would also hide the same vulnerability if we later introduced it in our own code. Leave an entry unscoped only when it genuinely applies repo-wide, such as a false positive originating in our own module graph.

Note that Trivy picks its parser from the file extension and only auto-discovers a file named `.trivyignore`, so this file has to be named explicitly. The `Makefile` does that for the scan targets via `TRIVY_IGNOREFILE`, which defaults to `.trivyignore.yaml` and can be overridden. Beware that YAML content in a file named `.trivyignore` is silently not applied.

All LTS branches are scanned from our main branch. This means that the job definitions to perform the scans are not backported to LTS branches, nor is `.trivyignore.yaml`. This is especially relevant if a user is building images for LTS branches/non-`gloo` repositories and attempting to scan them; in these cases, the locally checked out branch will not contain the correct `.trivyignore.yaml` file. Users who wish to perform a scan locally from another branch should download main's `.trivyignore.yaml` to a location that will not be updated by `git` operations and point the scan at it with `make <target> TRIVY_IGNOREFILE=/path/to/.trivyignore.yaml`, or pass Trivy `--ignorefile` directly. 
