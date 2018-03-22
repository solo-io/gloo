FROM alpine:3.7
COPY gloo-static /gloo
EXPOSE 8081
ENTRYPOINT ["/gloo"]
