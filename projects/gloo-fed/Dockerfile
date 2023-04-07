#####################################################################################################################
#
#   This Dockerfile describes a container that serves as a reproducible build environment for gloo fed.
#
#####################################################################################################################
ARG GO_BUILD_IMAGE
FROM $GO_BUILD_IMAGE as build-env

ARG VERSION
ARG GCFLAGS
ARG LDFLAGS
ARG GITHUB_TOKEN
ARG DOCKER_GOARCH

# Fail if VERSION is not set
RUN if [[ ! $VERSION ]]; then echo "Required VERSION build argument not set" && exit 1; fi
RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf https://github.com/
RUN go env -w GOPRIVATE=github.com/solo-io

# set up go mod and download go resources, so that we can save this step for caching for performance reasons
ADD ./go.mod /go/src/github.com/solo-io/solo-projects/go.mod
ADD ./go.sum /go/src/github.com/solo-io/solo-projects/go.sum

WORKDIR /go/src/github.com/solo-io/solo-projects

RUN go mod download

WORKDIR /
ADD . /go/src/github.com/solo-io/solo-projects
WORKDIR /go/src/github.com/solo-io/solo-projects

RUN CGO_ENABLED=0 GOARCH=${DOCKER_GOARCH} GOOS=linux \
go build -ldflags="$LDFLAGS" -gcflags="$GCFLAGS" -o gloo-fed-linux-${DOCKER_GOARCH} \
projects/gloo-fed/cmd/main.go


FROM alpine:3.17.3
ARG DOCKER_GOARCH

COPY --from=build-env /go/src/github.com/solo-io/solo-projects/gloo-fed-linux-${DOCKER_GOARCH} /
