# Snapshot End-to-End Tests

## Background

Snapshot e2e tests are run in a kind cluster environment. 

## Running Snapshot E2E Tests

To run the snapshot e2e tests, you can use the following make targets:

```bash
 make run-classic-snapshot-e2e-tests
```

```bash
 make run-gateway-snapshot-e2e-tests
```

The above commands will run the declarative env setup and the snapshot e2e tests for the classic and gateway snapshots respectively.

## Creating a new Snapshot E2E Test

The testcases for the snapshot e2e tests are located in the `test/snapshot/testcases` directory. The Gateway and Edge 
versions of the setup are defined in `test/snapshot/edge_classic_e2e` or `test/snapshot/gloo_gateway_e2e` call the testcase. 

For example, the `testcases.TestGatewayIngress` function is defined in the `test/snapshot/testcases/ingress.go` and called from
`test/snapshot/gloo_gateway_e2e/gateway_test.go` as follows:

``` 
testcases.TestGatewayIngress(
				ctx,
				runner,
				&snapshot.TestEnv{
					GatewayName:      gatewayDeploymentName,
					GatewayNamespace: gloodefaults.GlooSystem,
					GatewayPort:      gatewayPort,
					ClusterContext:   kubeCtx,
					ClusterName:      clusterName,
				},
				customAssertions(),
			)
```