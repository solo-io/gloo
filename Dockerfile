FROM alpine:3.7
COPY gloo /
EXPOSE 8081
ENTRYPOINT ["/gloo"]
