FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY gloo-function-discovery /
CMD ["/gloo-function-discovery"]
