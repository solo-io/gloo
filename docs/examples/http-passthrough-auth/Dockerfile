FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETOS
ARG TARGETARCH

RUN apk --no-cache add make
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -o server

FROM --platform=$BUILDPLATFORM alpine
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
