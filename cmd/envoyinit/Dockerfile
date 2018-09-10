FROM soloio/envoy-gloo@sha256:f4ef533543644e4d8ca282dcd331c0b5febffaaf9984c17c033de02802e43f2c

COPY envoyinit-linux-amd64 /usr/local/bin/envoyinit

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/usr/local/bin/envoyinit"]
CMD []