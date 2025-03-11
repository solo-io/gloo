# GRPC Passthrough Auth Server

This server's `/auth` endpoint will return a 200 if the `authorization: authorize me` header is set, and a 403 otherwise

If the `REQUEST_LOGGING` environment variable is set the full request will be written to the logs

Build the `quay.io/solo-io/passthrough-grpc-service-example` image by running `make docker-local` from this directory.
