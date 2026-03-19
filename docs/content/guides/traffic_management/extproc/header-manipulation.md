---
title: Header manipulation
weight: 40
description: Walk through an example for how to manipulate request headers by using an ExtProc server.
---

Set up an external processing (ExtProc) server that manipulates request headers for a sample app.

{{% notice note %}}
External processing is an Enterprise-only feature.
{{% /notice %}}

1. Before you begin, install [Gloo Gateway Enterprise]({{% versioned_link_path fromRoot="/installation/enterprise/" %}}) in your cluster.

2. Set up the ExtProc server. This example uses a prebuilt ExtProc server that manipulates request and response headers based on instructions that are sent in an `instructions` header.

   {{< tabs >}}
{{% tab %}}
```yaml
kubectl apply -f- <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ext-proc-grpc
spec:
  selector:
    matchLabels:
      app: ext-proc-grpc
  replicas: 1
  template:
    metadata:
      labels:
        app: ext-proc-grpc
    spec:
      containers:
        - name: ext-proc-grpc
          image: gcr.io/solo-test-236622/ext-proc-example-basic-sink:0.0.2
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 18080
---
apiVersion: v1
kind: Service
metadata:
  name: ext-proc-grpc
  labels:
    app: ext-proc-grpc
  annotations:
    gloo.solo.io/h2_service: "true"
spec:
  ports:
  - port: 4444
    targetPort: 18080
    protocol: TCP
  selector:
    app: ext-proc-grpc
EOF
```
{{% /tab %}}
   {{< /tabs >}}

   The `instructions` header must be provided as a JSON string in the following format:

   ```json
   {
     "addHeaders": {
       "header1": "value1",
       "header2": "value2"
     },
     "removeHeaders": [ "header3", "header4" ],
     }
   }
   ```

3. Verify that the ExtProc server is up and running.
   ```sh
   kubectl get pods
   ```

   Example output:
   ```
   NAME                             READY   STATUS    RESTARTS   AGE
   ext-proc-grpc-59d44ddf76-42q2x   1/1     Running   0          24m
   ```

4. Edit the default `Settings` custom resource to enable ExtProc in Gloo Gateway.
   ```
   kubectl edit settings default -n gloo-system
   ```

   Add the following ExtProc settings to the `spec` section. This example configures the standard `extProc` filter, which runs in the middle of the Envoy filter chain. Gloo Gateway also supports `extProcEarly` (runs early in the filter chain) and `extProcLate` (runs as the final filter before a request leaves Envoy). For more information, see [ExtProc filter variants]({{% versioned_link_path fromRoot="/guides/traffic_management/extproc/about/#extproc-filter-variants" %}}).

   ```yaml
   extProc:
     grpcService:
       extProcServerRef:
         name: default-ext-proc-grpc-4444
         namespace: gloo-system
     filterStage:
       stage: AuthZStage
       predicate: After
     failureModeAllow: false
     allowModeOverride: false
     processingMode:
       requestHeaderMode: SEND
       responseHeaderMode: SKIP
   ```

   |Setting|Description|
   |--|--|
   |`grpcService`| The configuration of the external processing server that you created earlier.|
   |`grpcService.exProcServerRef.name`| The name of the upstream that was created for the ExtProc server.|
   |`grpcService.exProcServerRef.namespace`| The namespace of the upstream that was created for the ExtProc server.|
   |`filterStage`|Where in the filter chain you want to apply the external processing. Applies to `extProcEarly` and `extProc`. Has no effect on `extProcLate`, which always runs as the final filter.|
   |`failureModeAllow`|Allow the ExtProc server to continue when an error is detected during external processing. If set to `true`, the ExtProc server continues. If set to `false`, external processing is stopped and an error is returned to the Envoy proxy. |
   |`allowModeOverride`|Allow the ExtProc server to override the processing mode settings that you set. Default value is `false`. |
   |`processingMode`|Decide how you want the ExtProc server to process request and response information. |
   |`processingMode.requestHeaderMode`|Send (`SEND`) or skip sending (`SKIP`) request header information to the ExtProc server. |
   |`processingMode.responseHeaderMode`|Send (`SEND`) or skip sending (`SKIP`) response header information to the ExtProc server. |

5. Deploy the `httpbin` sample app.
   {{< tabs >}}
{{% tab %}}
```yaml
kubectl apply -f- <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: httpbin
---
apiVersion: v1
kind: Service
metadata:
  name: httpbin
  labels:
    app: httpbin
spec:
  ports:
  - name: http
    port: 8000
    targetPort: 80
  selector:
    app: httpbin
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpbin
spec:
  replicas: 1
  selector:
    matchLabels:
      app: httpbin
      version: v1
  template:
    metadata:
      labels:
        app: httpbin
        version: v1
    spec:
      serviceAccountName: httpbin
      containers:
      - image: docker.io/kennethreitz/httpbin
        imagePullPolicy: IfNotPresent
        name: httpbin
        ports:
        - containerPort: 80
EOF
```
{{% /tab %}}
   {{< /tabs>}}

6. Verify that the httpbin pod is up an running.
   ```sh
   kubectl get pods | grep httpbin
   ```

7. Create a virtual service to expose the httpbin app on the gateway.
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: gateway.solo.io/v1
   kind: VirtualService
   metadata:
     name: vs
     namespace: gloo-system
   spec:
     virtualHost:
       domains:
       - '*'
       routes:
       - matchers:
         - prefix: /
         routeAction:
           single:
             upstream:
               name: default-httpbin-8000
               namespace: gloo-system
   EOF
   ```


8. Send a simple request to the httpbin app and make sure that you get back a 200 HTTP response code. The following request passes in two headers `header1` and `header2` that you see in your reponse.
   ```
   curl -vik $(glooctl proxy url)/get -H "header1: value1" -H "header2: value2"
   ```

   Example output:
   ```
   HTTP/1.1 200 OK
   ...
   {
     "args": {},
     "headers": {
       "Accept": "*/*",
       "Header1": "value1",
       "Header2": "value2",
       "Host": "example.com.com",
       "User-Agent": "curl/7.77.0",
       "X-Envoy-Expected-Rq-Timeout-Ms": "15000"
     },
     "origin": "10.0.11.109",
     "url": "http://example.com/get"
   }
   ```

9. Send another request to the httpbin app. This time, you pass along instructions for the ExtProc server in an `instruction` header. In this example, you use the ExtProc server to add the `header3` header, and to remove `header2`.
   ```sh
   curl -vik $(glooctl proxy url)/get -H "header1: value1" -H "header2: value2" -H 'instructions: {"addHeaders":{"header3":"value3","header4":"value4"},"removeHeaders":["instructions", "header2"]}'
   ```

   Example output:
   ```
   {
     "args": {},
     "headers": {
       "Accept": "*/*",
       "Header1": "value1",
       "Header3": "value3",
       "Header4": "value4",
       "Host": "example.com",
       "User-Agent": "curl/7.77.0",
       "X-Envoy-Expected-Rq-Timeout-Ms": "15000"
     },
     "origin": "10.0.11.109",
     "url": "http://example.com/get"
   }
   ```

10. Edit the `Settings` resource again and set `requestHeaderMode: SKIP`. This setting instructs the ExtProc filter to not send any request headers to the ExtProc server.
    ```sh
    kubectl edit settings default -n gloo-system
    ```

11. Send the same request to the httpbin app. Note that this time, `header2` is not removed and `header3` and `header4` are not added to the request, because the request headers are not sent to the ExtProc server.
    ```sh
    curl -vik $(glooctl proxy url)/get -H "header1: value1" -H "header2: value2" -H 'instructions: {"addHeaders":{"header3":"value3","header4":"value4"},"removeHeaders":["instructions", "header2"]}'
    ```

    Example output:
    ```
    {
     "args": {},
     "headers": {
       "Accept": "*/*",
       "Header1": "value1",
       "Header2": "value2",
       "Host": "example.com",
       "Instructions": "{\"addHeaders\":{\"header3\":\"value3\",\"header4\":\"value4\"},\"removeHeaders\":[\"instructions\", \"header2\"]}",
       "User-Agent": "curl/7.77.0",
       "X-Envoy-Expected-Rq-Timeout-Ms": "15000"
     },
     "origin": "10.0.11.109",
     "url": "http://example.com/get"
    }
    ```

## Use multiple extProc filter variants

You can configure `extProcEarly` and `extProcLate` alongside `extProc` to run multiple external processors at different stages of the filter chain. For example, you might want to debug your extProc server by logging requests at both the earliest and latest stages so that you can compare what changed in between. You can also use this setup to integrate with different extProc servers. 

1. Update the `default` Settings resource to configure all three extProc variants. In this example, all extProc variants use the same extProc server. However, you can configure a specific extProc server for each phase. The extProc stages are processed as follows: 
   * **extProcEarly**: Requests are modified before the `Fault` Envoy filter. The `Fault` filter is the first filter in the filter chain. For more information, see [Filter flow description]({{% versioned_link_path fromRoot="/introduction/traffic_filter/#filter-flow-description" %}}).
   * **extProc**: Requests are modified after the `AuthZ` Envoy filter. 
   * **extProcLate**: Requests are modified in the `upstream_http_filter` that is part of the Router phase. Note that although you must provide a `filterStage` setting, this setting is ignored as the `extProcLate` variant is always executed as part of the `upstream_http_filter` filter.
   ```sh
   kubectl edit settings default -n gloo-system
   ```

   Update your settings as follows: 
   ```yaml
   extProcEarly:
     grpcService:
       extProcServerRef:
         name: default-ext-proc-grpc-4444
         namespace: gloo-system
     filterStage:
       stage: FaulStage
       predicate: Before
     processingMode:
       requestHeaderMode: SEND
       responseHeaderMode: SEND
   extProc:
     grpcService:
       extProcServerRef:
         name: default-ext-proc-grpc-4444
         namespace: gloo-system
     filterStage: 
       stage: AuthZStage
       predicate: After
     failureModeAllow: false
     allowModeOverride: false
     processingMode:
       requestHeaderMode: SEND
       responseHeaderMode: SKIP
   extProcLate:
     # The filter stage is ignored as it is always executed in the upstream_http_filter. 
     filterStage:  
       stage: AuthZStage
       predicate: After
     grpcService:
       extProcServerRef:
         name: default-ext-proc-grpc-4444
         namespace: gloo-system
     processingMode:
       requestHeaderMode: SEND
       responseHeaderMode: SEND
   ```

2. Create a virtual service to expose the httpbin app on the gateway.
   ```yaml
   kubectl apply -f- <<EOF
   apiVersion: gateway.solo.io/v1
   kind: VirtualService
   metadata:
     name: vs
     namespace: gloo-system
   spec:
     virtualHost:
       domains:
       - '*'
       routes:
       - matchers:
         - prefix: /
         routeAction:
           single:
             upstream:
               name: default-httpbin-8000
               namespace: gloo-system
   EOF
   ```

3. Send a request to the httpbin app.
   ```sh
   curl -vik $(glooctl proxy url --name gateway-proxy)/get -H "header1: value1"
   ```

   Verify that you get back a 200 HTTP response code.
   ```
   HTTP/1.1 200 OK
   ```

4. Check the extProc server logs to verify that it received 3 processing requests (one for each stage). Because each variant sends request headers (`requestHeaderMode: SEND`), the server receives three separate gRPC processing calls per request. 
   ```sh
   kubectl logs -l app=ext-proc-grpc --tail=50
   ```

   Example output: 
   {{< highlight bash "hl_lines=2 5 8" >}}
   "Wed, 18 Mar 2026 19:57:41 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:207"	Starting gRPC server on port ":18080"
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:55"	Process
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:78"	Got RequestHeaders
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:182"	Sending ProcessingResponse
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:55"	Process
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:78"	Got RequestHeaders
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:182"	Sending ProcessingResponse
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:55"	Process
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:78"	Got RequestHeaders
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:182"	Sending ProcessingResponse
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:137"	Got ResponseHeaders
   "Thu, 19 Mar 2026 13:14:48 UTC: /Users/ben/Documents/solo-repos/ext-proc-examples/basic-sink/main.go:182"	Sending ProcessingResponse
   {{< /highlight >}}

   

## Cleanup

You can optionally remove the resources that you created as part of this guide.

```sh
kubectl delete deployment ext-proc-grpc httpbin
kubectl delete service ext-proc-grpc httpbin
kubectl delete servicaccount httpbin
kubectl delete virtualservice vs -n gloo-system
```







