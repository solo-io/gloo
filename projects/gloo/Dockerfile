#####################################################################################################################
#
#   This Dockerfile describes a container that serves as a reproducible build environment for gloo with race detection enabled.
#
#   Also, to be able to build with -race, the gloo image has to be built with CGO_ENABLED=1 (the Go runtime relies
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
# define race as "-race" to include race
ARG RACE
ARG BUILDARCH

# Fail if VERSION is not set
RUN if [[ ! $VERSION ]]; then echo "Required VERSION build argument not set" && exit 1; fi
RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf https://github.com/
RUN go env -w GOPRIVATE=github.com/solo-io

# get the cross compilers for both gcc and g++
# symlink because the files are misplaced on arm64
RUN if [ "${BUILDARCH}" = "arm64" ] && [ "${DOCKER_GOARCH}" = "amd64" ]; then apt-get update && apt-get install -y gcc-x86-64-linux-gnu g++-x86-64-linux-gnu; ln -s /usr/x86_64-linux-gnu/lib/* /usr/lib/x86_64-linux-gnu/; fi && \
    if [ "${BUILDARCH}" = "amd64" ] && [ "${DOCKER_GOARCH}" = "arm64" ]; then apt-get update && apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu; fi

# set up go mod and download go resources, so that we can save this step for caching for performance reasons
ADD ./go.mod /go/src/github.com/solo-io/solo-projects/go.mod
ADD ./go.sum /go/src/github.com/solo-io/solo-projects/go.sum

WORKDIR /go/src/github.com/solo-io/solo-projects

RUN go mod download

WORKDIR /
ADD . /go/src/github.com/solo-io/solo-projects
WORKDIR /go/src/github.com/solo-io/solo-projects

# make sure that if we need cross compilers that we set the appropriate variables for CC and CCX
RUN \
    if [ "${BUILDARCH}" = "arm64" ] && [ "${DOCKER_GOARCH}" = "amd64" ]; then export CC=x86_64-linux-gnu-gcc && export CXX=x86_64-linux-gnu-g++; fi && \
    if [ "${BUILDARCH}" = "amd64" ] && [ "${DOCKER_GOARCH}" = "arm64" ]; then export CC=aarch64-linux-gnu-gcc && export CXX=aarch64-linux-gnu-g++; fi && \
    CGO_ENABLED=1 GOARCH=${DOCKER_GOARCH} GOOS=linux go build \
    ${RACE} \
    -gcflags="$GCFLAGS" \
    -ldflags="$LDFLAGS" \
    -o gloo-linux-${DOCKER_GOARCH} \
    projects/gloo/cmd/main.go

FROM alpine:3.17.3
ARG DOCKER_GOARCH

COPY --from=build-env /go/src/github.com/solo-io/solo-projects/gloo-linux-${DOCKER_GOARCH} /
