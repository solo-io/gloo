#####################################################################################################################
#
#   This Dockerfile describes a container that serves as a reproducible build environment for the caching server.
#
#   Also, to be able to build with goboring, the caching server has to be built with CGO_ENABLED=1 (the Go runtime relies
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

# Fail if VERSION is not set
RUN if [[ ! $VERSION ]]; then echo "Required VERSION build argument not set" && exit 1; fi

RUN echo $USE_APK
RUN if [[ $USE_APK ]]; then apk add --no-cache gcc musl-dev git; fi

ADD . /go/src/github.com/solo-io/solo-projects
WORKDIR /go/src/github.com/solo-io/solo-projects

RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf https://github.com/

RUN go env -w GOPRIVATE=github.com/solo-io

# Build caching server with CGO enabled
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build \
    -gcflags="$GCFLAGS" \
    -ldflags="-X github.com/solo-io/solo-projects/pkg/version.Version=$VERSION $LDFLAGS" \
    -o caching-linux-amd64 \
    projects/caching/cmd/main.go


FROM alpine:3.15.6
COPY --from=build-env /go/src/github.com/solo-io/solo-projects/caching-linux-amd64 /
