# Performance Tests

## Label
Performance tests are labeled with `performance`.

## When are these tests run?
- These tests are skipped on CI for PRs.
- They were run on a schedule as part of our [nightly tests](nightly-tests.md), but have since been disabled (https://github.com/solo-io/gloo/issues/10534)


## How can you run these tests?
You can run the `performance_test.go` locally with the make target:
```bash
make install-test-tools run-performance-tests
```

Or through GitHub actions via the `performance-test` workflow.