
When creating a new resource Group extend the Validator.go method`ValidationIsSupported()`

When supporting a new resource we can extend based off the interfaces of the Group.
Currently we support Gloo and Gateway resources, via the `gateway_resource_validation.go` and `gloo_resource_validation.go` files.

