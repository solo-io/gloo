# Description
 - This directory contains resources used in the caching kube2e tests, defined in [/test/kube2e/caching](/test/kube2e/caching) the image is also used by caching upgrade tests [/test/kube2e/upgrade](/test/kube2e/upgrade)
 - The files in this directory are used to build the `gcr.io/solo-test-236622/cache_test_service` docker image referenced on the pod resource in the caching tests.
   - This image is a backend service for cache testing, originally developed for the [upstream Envoy cache filter sandbox](https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/cache#cache-filter)
     - The source of this service can be found [here](https://github.com/envoyproxy/envoy/tree/main/examples/cache), in upstream Envoy
   - You can build this image by running [./build.sh](./build.sh), with the `TAG` environment variable set to the value of the tag you want the `gcr.io/solo-test-236622/cache_test_service` image to have
     - `TAG` will default to `0.0.0` if not set
       - the `TAG` does not have to match the tag of any gloo components. The only place this image is referenced is in the resources.yaml file we use to create resources in these caching kube2e tests
   - [./responses.yaml](./responses.yaml) defines the routes and responses that the service will handle
   - The reason we choose to use this service that it is set up to handle cache validation, which other services we use for testing (petstore, postman-echo, httpbin) are not