package testutils

// This file is an extension of the Gloo OSS testutils env.go file
// https://github.com/solo-io/gloo/blob/main/test/testutils/env.go
// As a pattern, we should inherit the OSS testutils wherever possible,
// but if there are Enterprise-only values, we can define them here

const (
	// RedisBinary is used in e2e tests to specify the path to the consul binary to use for the tests
	// See redis/factory.go for more details
	RedisBinary = "REDIS_BINARY"
)
