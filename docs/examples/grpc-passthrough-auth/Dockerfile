FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETOS
ARG TARGETARCH

RUN apk --no-cache add make
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
  go build -a --ldflags '-extldflags "-static"' -tags "netgo" \
  -installsuffix netgo -o server

FROM --platform=$BUILDPLATFORM alpine
RUN ls
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
