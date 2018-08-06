# solo-kit
Tools for Building Declarative API, stateless Projects

## Major TODOS:

* pkg/api/v1/clients <- support / enforce resource versioning on updates / creates (be consistent with Kube)

steps:
XXX: snapshot + cache
xxx: split out SetStatus stuff for IsInputResource in code gen
XX4 - reporter

XXX - desired-state-achiever (syncer library)
 
xxx - support selectors/labels

5 - e2e tests

XXX 6 - callbacks for apiserver
xxx - something that works for 3rd party resources (configmaps, artifacts, etc)

7 - acl for apiserver

 post 3 weeks
- bootstrap
- installer
- CLI
- tests for consul and file


- caching plugin
- extending xds for rate limit, extauth


syncers:
- generic ADS server, returns the grpc server to which envoy, rate limit, etc register


gloo:
- plugins > eds > bootstrap for setup