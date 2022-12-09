#!/bin/bash
# Step 1: build intermediate docker image from envoy repo
DOCKER_BUILDKIT=0 docker build https://github.com/envoyproxy/envoy.git\#main:examples/shared/flask -t flask_service:python-3.10-slim-bullseye

# Step 2: build final docker image from intermediate image

# if you build a version of this image with a new tag, update ./resources/pod.yaml
# to use the new image in the caching tests
# image tag defaults to 0.0.0 if not set

# please note that this image is not published on release
# if you make a new build of this image, you will need to manually push it to GCR
# in order for it to be accessible outside of your local environment

TAG=${TAG:-0.0.0}
docker build -t gcr.io/solo-test-236622/cache_test_service:$TAG .
