FROM golang:1.22 as build
ENV GOBIN /tmp
ENV GO111MODULE on
ENV CGO_ENABLED 0
RUN go install sigs.k8s.io/cloud-provider-kind@latest

FROM docker:25.0-dind
COPY --from=build --chown=root:root /tmp/cloud-provider-kind /bin/cloud-provider-kind
ENTRYPOINT ["/bin/cloud-provider-kind"]
