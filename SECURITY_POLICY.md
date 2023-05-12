# Security policy

## Release model
Under the Gloo Enterprise [release model](https://docs.solo.io/gloo/latest/reference/support/), 
a new stable X.Y.0 release will be published every 3 months. 
New feature development occurs on the main branch, and a new stable version will be released from main
at the end of each quarter. After it has been published, a stable release will be patched only to address 
security vulnerabilities and fix bugs that critically impact existing features or the stability of the product.

### Envoy
Gloo ships with a custom Envoy distribution that includes additional filters while leaving the core Envoy components untouched. 
The Envoy project recently introduced a [stable release model](https://github.com/envoyproxy/envoy/blob/main/RELEASES.md), 
which aims at publishing quarterly stable releases, with 1 year of support for each stable release. 
Gloo Enterprise will always track the latest stable Envoy release. 
In case the expected quarterly stable release of Envoy is not available when a new stable version of Gloo Enterprise 
is released, Gloo Enterprise will temporarily track the main branch of Envoy until said release is available. Since 
Envoy makes strong guarantees of stability and backwards compatibility on the main branch, the solo.io team deems this 
approach to be acceptable.

When starting the release process for a new stable version of Gloo Enterprise, one of two scenarios 
can occur with respect to the upstream Envoy project:
1. The new quarterly stable Envoy version **has already been released**: in this case, the new stable 
Gloo Enterprise release will track this release.
2. The new quarterly stable Envoy version **has not yet been released** due to delays (up to 3 weeks): 
in this case, the new stable Gloo Enterprise release will track the upstream Envoy main branch 
until the stable Envoy release has been published; at that point, the stable Gloo Enterprise 
release will be patched to track it.

## Security Releases
Solo.io will support the 3 most recent stable releases of Gloo Enterprise and backport fixes to critical 
bugs and vulnerabilities in the Gloo components, as well as in Envoy and in other open source dependencies.

### Gloo components
When there is a critical bug or vulnerability in a Gloo component (`gloo`, `gateway`, `discovery`, `extauth`, 
`rate-limit`, or `observability`), the Solo.io team will prepare a fix for the main branch, and backport 
the fix to all supported stable release branches. 

Once all of these releases have been prepared and made available, then the Solo team will send notifications 
to all customers through Slack and email about the security update. 

### Envoy
When there is a critical bug or vulnerability in Envoy, Solo.io will find out either when the fix has been merged 
into Envoy, or with advanced notice behind a communication embargo. As soon as the fix has been released, Solo.io will:
1. Update each of the relevant versions of the custom Envoy distribution shipped with Gloo to 
track the corresponding patches stable Envoy versions;
2. Update each relevant stable version of Gloo to use the corresponding updated custom Envoy distribution.

Once all the supported stable versions of Gloo have been updated, the releases will be prepared and communicated as above. 

The [12 months support window](https://github.com/envoyproxy/envoy/blob/main/RELEASES.md#stable-releases) for Envoy 
means that security vulnerabilities and critical bug fixes will be ported to the 4 last quarterly stable releases on Envoy. 
In turn, given the release model outlined earlier in this document and the Gloo Enterprise 
[support policy](https://docs.solo.io/gloo/latest/reference/support/#support-will-be-release-n-through-n-2), 
this means that Solo.io should always be able to patch each of the supported Gloo Enterprise versions by just switching 
to the corresponding patched Envoy release, without the risk of running into compatibility issues.

If a customer is paying for extended support beyond the 12 months Envoy support window, it may not be possible to 
release a security fix by simply bumping the dependency on upstream Envoy, as the fix will not have been ported to the 
Envoy version tracked by the Gloo Enterprise version the customer is running. 
How to move forward in these cases will be discussed between Solo.io and affected customers on a case-by-case basis. 

### Other open source dependencies (Prometheus, Grafana, Redis)
In the unlikely event of a vulnerability in Prometheus, Grafana, Redis, or other 3rd party open source dependencies, 
Solo.io will work with those communities to identify safe versions to use and then upgrade to those versions on all 
the stable branches of Gloo. Like above, Solo will prepare the updated releases and communicate them to customers.

## Recent fixes

### Envoy
The following table shows recent CVEs that affected the custom Envoy distribution shipped with Gloo. The last three columns indicate:
- when the vulnerability has been fixed in upstream Envoy,
- when the fix has been ported to the open source version of Gloo,
- when the fix has been ported to Gloo Enterprise.

| CVE Reference  | Base score     |       Envoy Release | Gloo Release         | GlooE Release                              |
|----------------|----------------|--------------------:|----------------------|--------------------------------------------|
| CVE-2020-8659  | 7.5 (High)     | 1.12.3 (03/03/2020) | V1.2.23 (03/04/2020) | V1.2.10 (03/04/2020)                       |
| CVE-2020-8661  | 7.5 (High)     | 1.12.3 (03/03/2020) | V1.2.23 (03/04/2020) | V1.2.10 (03/04/2020)                       |
| CVE-2020-8664  | 5.3 (Medium)   | 1.12.3 (03/03/2020) | V1.2.23 (03/04/2020) | V1.2.10 (03/04/2020)                       |
| CVE-2020-8660  | 5.3 (Medium)   | 1.12.3 (03/03/2020) | V1.2.23 (03/04/2020) | V1.2.10 (03/04/2020)                       |
| CVE-2019-18801 | 9.8 (Critical) | 1.12.2 (12/10/2019) | 1.2.5 (12/10/2019)   | 0.21.2 (12/18/2019) 1.0.0-rc5 (12/11/2019) |
| CVE-2019-18802 | 9.8 (Critical) | 1.12.2 (12/10/2019) | 1.2.5 (12/10/2019)   | 0.21.2 (12/18/2019) 1.0.0-rc5 (12/11/2019) |
| CVE-2019-18838 | 7.5 (High)     | 1.12.2 (12/10/2019) | 1.2.5 (12/10/2019)   | 0.21.2 (12/18/2019) 1.0.0-rc5 (12/11/2019) |
| CVE-2019-15226 | 7.5 (High)     | 1.11.2 (10/08/2019) | 0.20.4 (10/08/2019)  | 0.20.2 (10/09/2019)                        |
| CVE-2019-15225 | 7.5 (High)     | 1.11.2 (10/08/2019) | 0.20.4 (10/08/2019)  | 0.20.2 (10/09/2019)                        |

### Gloo
The following table shows turnaround times for fixing recently reported critical bugs to Gloo Enterprise.

| Reported date | Description                                                      | Most recent affected version | GlooE Release      |
|---------------|------------------------------------------------------------------|------------------------------|--------------------|
| 02/06/2020    | RBAC route policies break when multiple HTTP methods are defined | 1.2.3                        | 1.2.4 (02/08/2020) |
| 02/04/2020    | Intermittent 503 responses while upstreams are being updated     | 1.2.3                        | 1.2.4 (02/08/2020) |
| 01/21/2020    | High CPU usage when processing config maps                       | 1.2.2                        | 1.2.3 (01/30/2020) |
| 01/16/2020    | Upgrade via Helm 2 broken                                        | 1.2.2                        | 1.2.3 (01/30/2020) |
| 01/14/2020    | Grafana "Average Request Time" panel shows zero values           | 1.2.1                        | 1.2.6 (02/19/2020) |

## Additional resources
- [Original Google Doc](https://docs.google.com/document/d/1M5Sto821S5fYLRdP-XVYdz4hd397HWbAqLM731UJIdQ)
