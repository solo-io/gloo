# gloo-core
The Gloo Control Plane + Tools for Building Gloo-Based Projects

## Major TODOS:

* pkg/api/v1/clients <- support / enforce resource versioning on updates / creates (be consistent with Kube)

* Config Generator (Call it an inventory! :D ). Jut have the plugin
create a config watcher for all the resources it processes. user doesnt have to write
a new proto for that config (it's internal anyway)


steps:
1 - generate inventory container for resources, e.g.:
type Inventory struct {
  MockResourceList []*MockResource
  FakeResourceList []*FakeResource
}

2 - generate a config watcher for this config

3 - generate an event loop:
 - event loop should take other channels as parameters, as well as a sync funciton (for example we need
 an event loop to handle secrets, endpoints for gloo)