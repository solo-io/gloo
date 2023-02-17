#####################################################################################################################
#
#   This Dockerfile describes a container that serves as a reproducible build environment for the ext auth server.
#
#   Also, to be able to load Go plugins and to build with goboring, the ext auth server has to be built with CGO_ENABLED=1 (the Go runtime relies
#   on C code and libraries to load plugins), so cross compiling isn't as easy as `GOARCH=amd64 GOOS=linux go build`
#   anymore. Using this container you can build the binary on any architecture.
#
#####################################################################################################################
ARG GO_BUILD_IMAGE
FROM $GO_BUILD_IMAGE as build-env

ARG VERSION
ARG GCFLAGS
ARG LDFLAGS
ARG GITHUB_TOKEN
ARG USE_APK
ARG DOCKER_GOARCH
ARG DOCKER_CGO_ENABLED=0

# Fail if VERSION is not set
RUN if [[ ! $VERSION ]]; then echo "Required VERSION build argument not set" && exit 1; fi

RUN echo $USE_APK
RUN if [[ $USE_APK ]]; then apk add --no-cache gcc musl-dev git; fi

ADD . /go/src/github.com/solo-io/solo-projects
WORKDIR /go/src/github.com/solo-io/solo-projects
RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf https://github.com/

RUN go env -w GOPRIVATE=github.com/solo-io

# Build ext-auth server with CGO enabled
RUN CGO_ENABLED=${DOCKER_CGO_ENABLED} GOARCH=${DOCKER_GOARCH} GOOS=linux go build \
    -gcflags="$GCFLAGS" \
    -ldflags="$LDFLAGS" \
    -o extauth-linux-${DOCKER_GOARCH} \
    projects/extauth/cmd/main.go

# Build script to verify if a plugin can be loaded by the ext-auth server
RUN CGO_ENABLED=${DOCKER_CGO_ENABLED} GOARCH=amd64 GOOS=linux go build \
    -gcflags="$GCFLAGS" \
    -ldflags="$LDFLAGS" \
    -o verify-plugins-linux-amd64 \
    projects/extauth/scripts/verify_plugins.go

FROM alpine:3.17.1
ARG DOCKER_GOARCH
COPY --from=build-env /go/src/github.com/solo-io/solo-projects/extauth-linux-${DOCKER_GOARCH} /
COPY --from=build-env /go/src/github.com/solo-io/solo-projects/verify-plugins-linux-amd64 /
