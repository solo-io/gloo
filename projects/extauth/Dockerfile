#####################################################################################################################
#
#   This Dockerfile describes a container that serves as a reproducible build environment for the ext auth server.
#
#   Also, to be able to load Go plugins, the ext auth server has to be built with CGO_ENABLED=1 (the Go runtime relies
#   on C code and libraries to load plugins), so cross compiling isn't as easy as `GOARCH=amd64 GOOS=linux go build`
#   anymore. Using this container you can build the binary on any architecture.
#
#####################################################################################################################
ARG GO_BUILD_IMAGE
FROM $GO_BUILD_IMAGE as build-env

ARG VERSION
ARG GCFLAGS

# Fail if VERSION is not set
RUN if [[ ! $VERSION ]]; then echo "Required VERSION build argument not set" && exit 1; fi

RUN apk add --no-cache gcc musl-dev git

ADD . /go/src/github.com/solo-io/solo-projects
WORKDIR /go/src/github.com/solo-io/solo-projects

# De-vendor all the dependencies and move them to the GOPATH, so they will be loaded from there.
# We need this so that the import path for any library shared between the plugins and Gloo is the same.
#
# For example, if we were to vendor the ext-auth-plugin dependency, the ext-auth-server would load the plugin interface
# as `github.com/solo-io/solo-projects/vendor/github.com/solo-io/ext-auth-plugins/api.ExtAuthPlugin`, while the plugin
# would instead implement `github.com/solo-io/ext-auth-plugins/api.ExtAuthPlugin`. Consequently, the type assertion that
# we perform to get a typed plugin instance would fail.
# Also, some packages cause problems if loaded more than once. For example, loading `golang.org/x/net/trace` twice
# causes a panic (see here: https://github.com/golang/go/issues/24137). By flattening the dependencies this way we
# prevent these sorts of problems.
RUN cp -a vendor/. /go/src/ && rm -rf vendor

# Build ext-auth server with CGO enabled
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build \
    -gcflags="$GCFLAGS" \
    -ldflags="-X github.com/solo-io/solo-projects/pkg/version.Version=$VERSION" \
    -o extauth-linux-amd64 \
    projects/extauth/cmd/main.go

# Build script to verify if a plugin can be loaded by the ext-auth server
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build \
    -gcflags="$GCFLAGS" \
    -o verify-plugins-linux-amd64 \
    projects/extauth/scripts/verify_plugins.go

FROM alpine:3.10.1
COPY --from=build-env /go/src/github.com/solo-io/solo-projects/extauth-linux-amd64 /
COPY --from=build-env /go/src/github.com/solo-io/solo-projects/verify-plugins-linux-amd64 /
CMD echo "Success!\n" && ls -lha