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