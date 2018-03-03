FROM scratch
COPY gloo /
EXPOSE 8081
ENTRYPOINT ["/gloo"]
