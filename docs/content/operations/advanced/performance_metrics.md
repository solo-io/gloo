
Gloo Metrics

Gloo collects sync time metrics in the following pattern:
|Metric|Description|Aggregation|Tag|
|------|-----------|-----------|---|
|${DOMAIN}/sync/time_sec|The time taken for a given sync|Distribution|syncer_name|

and assocaited volume metrics
|Metric|Description|Aggregation|Tag|
|------|-----------|-----------|---|
|${DOMAIN}/emitter/resources_in|The number of resource lists received on open watch channels|Count|N/A|
|${DOMAIN}/emitter/snap_out|The number of snapshots out|Count|N/A|
|${DOMAIN}/emitter/snap_missed|The number of snapshots missed|Count|NamespaceKey, ResourceKey|

These metrics are collected for the following domains:
- translator.clusteringress.gloo.solo.io
- discovery.gloo.solo.io
- enterprise.gloo.solo.io
- api.gloosnapshot.gloo.solo.io
- status.ingress.solo.io
- translator.ingress.solo.io
- translator.knative.gloo.solo.io