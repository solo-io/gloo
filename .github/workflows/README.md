# Github Actions
The following is a set of instruction on how to use the Github Actions to perform manual actions if necessary.


## Create Test Release Action (Release-branch.yaml)
Instructions on how to use [the Create Test Release action](https://github.com/solo-io/solo-projects/actions/workflows/release-branch.yaml)
- visit [Create Test Release workflow](https://github.com/solo-io/solo-projects/actions/workflows/release-branch.yaml)
- click “run workflow” (top right of central table)
- specify job parameters
    - branch - which solo-projects branch to run from
    - publish fed - checkbox to publish fed (or not)
    - gloo commit SHA - (optional) gloo commit SHA to reference if unpublished OSS changes are desired
- run workflow
- wait for workflow to finish
- [scroll to bottom of workflow(this is an example link)](https://github.com/solo-io/solo-projects/actions/runs/4116181193), then see the Artifacts section for installation instructions

# Running Scale Tests
We have two scale tests for Gloo, run by the [resource scale test action](eks-resource-scale-runner.yaml) and the [fed scale test action](eks-fed-scale-test-runner.yaml).

During scale test runs, the created clusters can be found as part of our [AWS EKS Clusters](https://us-east-2.console.aws.amazon.com/eks/home?region=us-east-2#/clusters) prefixed with `cicd-`.
They are part of the [solo-eng-test-beds-vpc](https://us-east-2.console.aws.amazon.com/vpc/home?region=us-east-2#VpcDetails:VpcId=vpc-0ac45034c3b802489) VPC.

Our resource scale tests are currently configured to run a setup similar to T-Mobile's, and our Fed scale tests are configured to run a setup similar to SAP's.
The resources created can be found in our [/templates folder](../../ci/eks/templates) and are called in each test's respective script: [install-edge.sh](../../ci/eks/resource-scale.sh) for resource scale tests, and [fed-data-creation.sh](../../ci/eks/fed-data-creation.sh) for fed scale tests.

## Cleaning up scale tests
To clean up any remaining clusters after running scale tests, run the [EKS cleanup action](eks-clusters-cleanup.yaml).

**Note**: The `Environment ID` passed is treated as a prefix. For instance `edgeScaleTests-cl1` will delete the clusters `cl1`, `cl10`, `cl11`, etc. If deleting clusters created during scale tests, providing the `<base-name>` used when starting the workflow deletes all clusters created for it during the test.

## Error Resolution

### Debugging

#### Connecting to the EKS Cluster
On every scaling test workflow run, if the test fails, a `kubeconfig` artifact gets produced in the run's summary.
You can also instead run `aws eks update-kubeconfig --region <region> --name <cluster-name>` which stores the context in your kubeconfig and automatically sets it as the active context.

After connecting, running `glooctl check` is a good first step to see if anything stands out.

**Note**: You will need AWS credentials to connect to the EKS cluster.

#### Testing Scale Test Updates
To test updates made to the scale tests, you can update our [EKS Resource Scale Action](eks-resource-scale-runner.yaml) (or our [EKS Fed Scale Action](eks-fed-scale-test-runner.yaml)) to point to your branch.
By updating the `checkout` action to include `ref: <your-branch-name>`, then using that branch to run the workflow, the tests will run with any updates made to the tests.

**Note**: This method does not support testing local builds of Gloo, just updates to the [scale tests](../../ci/eks/resource-scale.sh).

### Common Errors

#### Clusters already present in EKS
When running scale tests, the build-environment job may fail with the following error:
```
The following clusters seem to be already present and running. Delete them and then try running again.
cicd-<base-name>-<something>:<region>
<...list of any other existing clusters>
```
To resolve this follow guidance on [cleaning up scale tests](#cleaning-up-scale-tests).

#### Tests running indefinitely / Getting 503 responses during requests
Another error observed when running the resource scale tests, is getting a response of `503` on the request to the VirtualService.

If this occurs, you can replicate the request using `curl` by doing the following:
1. [Connect to EKS cluster](#connecting-to-the-eks-cluster)
2. Run `curl -v -H "Host: <domain>" <glooctl proxy url>/status/200`
   * The `Host` is the domain we're trying to hit during the request which is logged in the GitHub Action run.
   * Running `glooctl proxy url` give us the accessible ELB we're trying to hit.

If the `503`s are occurring locally as well, make sure the upstream `host` we're hitting is active.

Our tests hit `httpbin.org` as the upstream, and it is possible that the domain is down. 
A possible workaround is to update references of `httpbin.org` to a host that does a similar function (returns a `200` when hitting `<host>/status/200`).
We have an [open issue](https://github.com/solo-io/solo-projects/issues/5280#issuecomment-1675173999) which proposes either hosting our own long-lived `httpbin` deployment or creating one dynamically during test runtime to limit this being an issue.

#### Pods Failing - With Event `FailedCreatePodSandbox`
This issue usually occurs with the following `Event` message: `networkPlugin cni failed to set up pod <pod> network: add cmd: failed to assign an IP address to container`.

One reason for the pods to fail to create can be due to our VPC's subnets running out of available IPs to use. 
All of our scale tests run under the same VPC, meaning they all share available IPs.

To resolve this issue:
- Run scale tests with fewer clusters.
  - This is more relevant for Fed scale tests.
- Delete old or unnecessary EKS clusters using the same VPC.
- Upgrade the VPC's subnet CIDR to allow for more assignable IPs.
  - Our VPC has 6 subnets - 3 IPv4, 3 IPv6 - of which the IPv4 subnets used by our scale tests have `/20` CIDR ranges. It may not be possible to reserve larger ranges; make any requests of this nature to Sundar and/or IT.
