FROM alpine:3.22.2

COPY cert.pem /cert.pem
COPY key.pem /key.pem
COPY example-tls-server /usr/local/bin/example-tls-server

ENTRYPOINT ["/usr/local/bin/example-tls-server"]
