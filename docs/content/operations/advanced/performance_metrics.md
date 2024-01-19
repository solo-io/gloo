
Gloo Metrics

Gloo collects sync time metrics in the following pattern:
|Metric|Description|Aggregation|Tag|
|------|-----------|-----------|---|
|${DOMAIN}/sync/time_sec|The time taken for a given sync|Distribution|syncer_name|

and associated volume metrics
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
- eds.gloo.solo.io (sync time only)


## A practical example

The view used for the sync time metrics is a distribution/histogram, but it is possible to measure a single event using the `_sync_time_sec_sum` metrics.
This view of the metric stores a cumumlative run time, so we can store the value, update the system, and then get the new value.

With a cluster running, start port-forwarding the Gloo metrics port:
```
kubectl port-forward -n gloo-system deploy/gloo 9091:9091
```

We can now curl `http://localhost:9091/metrics` and grep for the metrics we need. 
For example if we wanted to evaulate how long the gloosnapshot and discovery syncers were taking, we can run the following curl

```
curl http://localhost:9091/metrics | grep -E "api_gloosnapshot_gloo_solo_io_sync_time_sec_sum|eds_gloo_solo_io_sync_time_sec_sum"
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 15532    0 15532    0     0   470k      0 --:--:-- --:--:-- --:--:--  523k
api_gloosnapshot_gloo_solo_io_sync_time_sec_sum{syncer_name="gloosnapshot.ApiSyncers"} 3.285250713000001
eds_gloo_solo_io_sync_time_sec_sum{syncer_name="*discovery.syncer"} 1.5245133750000002
```

We can also grep for just `time_sec_sum` to see all of the available sync time metrics.

It is also possible to store the values in local variables:
```
export SNAPSHOT_TIME_START=`curl http://localhost:9091/metrics | grep api_gloosnapshot_gloo_solo_io_sync_time_sec_sum | sed 's/.* //'`
export EDS_TIME_START=`curl http://localhost:9091/metrics | grep eds_gloo_solo_io_sync_time_sec_sum | sed 's/.* //'` 
```

We can then apply a change:

```
kubectl apply -f https://raw.githubusercontent.com/solo-io/gloo/v1.14.x/example/petstore/petstore.yaml
```

and rerun the command:
```
curl http://localhost:9091/metrics | grep -E "api_gloosnapshot_gloo_solo_io_sync_time_sec_sum|eds_gloo_solo_io_sync_time_sec_sum" 
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 15523    0 15523    0     0   528k      0 --:--:-- --:--:-- --:--:--  659k
api_gloosnapshot_gloo_solo_io_sync_time_sec_sum{syncer_name="gloosnapshot.ApiSyncers"} 3.840931671
eds_gloo_solo_io_sync_time_sec_sum{syncer_name="*discovery.syncer"} 1.854135166
```

We can again store the values in environment variables, though 
```
export SNAPSHOT_TIME_END=`curl http://localhost:9091/metrics | grep api_gloosnapshot_gloo_solo_io_sync_time_sec_sum | sed 's/.* //'`
export EDS_TIME_END=`curl http://localhost:9091/metrics | grep eds_gloo_solo_io_sync_time_sec_sum | sed 's/.* //'`   
```

We can then subtract the values to see how long was spent in each action:
```
echo $(($SNAPSHOT_TIME_END - $SNAPSHOT_TIME_START)) 
0.55568095799999906

 echo $(($EDS_TIME_END - $EDS_TIME_START))
0.32962179099999989
```

In this case, when creating the petstore app we spent 0.56 API snapshot sync and 0.33s in the Discovery sync.
