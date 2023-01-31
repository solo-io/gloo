# Envoyinit
The running instance of Envoy. In Gloo Edge, this is commonly referred to as the `gateway-proxy` component.

This is the Enterprise extension of the Open Source Envoyinit. Please refer to the [Open Source README for further details](https://github.com/solo-io/gloo/blob/master/projects/envoyinit/README.md)

## Background
### Source Code
The Gloo Edge service proxies provide all the functionality of the [open source Gateway Proxy](https://github.com/solo-io/envoy-gloo), in addition to some custom extensions. The source code for these proxies is maintained at [envoy-gloo-ee](https://github.com/solo-io/envoy-gloo-ee)

### Versioning
In the [Makefile](https://github.com/solo-io/solo-projects/blob/master/Makefile), the `ENVOY_GLOO_IMAGE` value defines the version of `envoy-gloo-ee` that Gloo Edge depends on.

Envoy publishes new minor releases [each quarter](https://www.envoyproxy.io/docs/envoy/latest/version_history/version_history#). Gloo attempts to follow this cadence, and increment our minor version of `envoy-gloo-ee` as well.

## Build
*All make targets are currently defined in the [Makefile](https://github.com/solo-io/solo-projects/blob/master/Makefile) at the root of the repository.*

The `VERSION` env variable determines the name of the tag for the image.

You may either inject the version yourself:
```bash
VERSION=<version name> make gloo-ee-envoy-wrapper-docker -B
```

Or rely on the auto-generated version:
```shell
make gloo-ee-envoy-wrapper-docker -B
```

## Release
During a Gloo Edge release, the `gloo-ee-envoy-wrapper` image is published to the [Google Cloud Registry](https://console.cloud.google.com/gcr/images/gloo-edge/GLOBAL) and the [Quay repository](https://quay.io/repository/solo-io/gloo-ee-envoy-wrapper).

## Configuration
The Gloo Edge Enterprise Gateway-Proxy manages configuration identically to the Open Source implementation. Please refer to the [Open Source README](https://github.com/solo-io/gloo/blob/master/projects/envoyinit/README.md#configuration) for details