# gloo-core
The Gloo Control Plane + Tools for Building Gloo-Based Projects

## Major TODOS:

* pkg/api/v1/clients <- support / enforce resource versioning on updates / creates (be consistent with Kube)

* Config Generator (Call it an inventory! :D ). Jut have the plugin
create a config watcher for all the resources it processes. user doesnt have to write
a new proto for that config (it's internal anyway)


steps:
XX4 - reporter

XXX - desired-state-achiever (syncer library)
 
xxx - support selectors/labels

5 - e2e tests

6 - callbacks/acl for apiserver


- something that works for 3rd party resources (configmaps, artifacts, etc)
 
 post 3 weeks
- bootstrap
- installer
- CLI
- tests for consul and file


- knative demo
- caching plugin
- framework
- extending xds for rate limit, extauth