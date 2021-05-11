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
ARG GITHUB_TOKEN

# Fail if VERSION is not set
RUN if [[ ! $VERSION ]]; then echo "Required VERSION build argument not set" && exit 1; fi

RUN apk add --no-cache gcc musl-dev git

ADD . /go/src/github.com/solo-io/solo-projects
WORKDIR /go/src/github.com/solo-io/solo-projects

RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf https://github.com/

RUN go env -w GOPRIVATE=github.com/solo-io/ext-auth-service

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

FROM alpine:3.13.5
COPY --from=build-env /go/src/github.com/solo-io/solo-projects/extauth-linux-amd64 /
COPY --from=build-env /go/src/github.com/solo-io/solo-projects/verify-plugins-linux-amd64 /