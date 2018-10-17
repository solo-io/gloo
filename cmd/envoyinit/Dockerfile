FROM soloio/envoy-gloo@sha256:93ffdf483931204a1c42b91f21c5afae07a218b956d227e29a87b2c58d99fe56
# this is commit 1ba560f243faec138b07b0252d6ea3081b32fa13

COPY envoyinit-linux-amd64 /usr/local/bin/envoyinit

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/usr/local/bin/envoyinit"]
CMD []